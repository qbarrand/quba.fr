package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/qbarrand/quba.fr/internal/imgpro"
	"github.com/qbarrand/quba.fr/internal/metadata"
)

const (
	ExtensionJPEG = ".jpg"
	ExtensionWEBP = ".webp"

	metadataFileName = "metadata.json"
)

var processorMap = map[string]imgpro.Processor{
	"vips": &imgpro.VipsProcessor{},
}

type breakpointsSlice struct {
	bp []int
}

func (bs *breakpointsSlice) Set(s string) error {
	elems := strings.Split(s, ",")

	for _, e := range elems {
		i, err := strconv.ParseInt(e, 10, 32)
		if err != nil {
			return err
		}

		bs.bp = append(bs.bp, int(i))
	}

	return nil
}

func (bs *breakpointsSlice) String() string {
	var sb strings.Builder

	if len(bs.bp) >= 1 {
		fmt.Fprintf(&sb, "%d", bs.bp[0])

		for i := 1; i < len(bs.bp); i++ {
			sb.WriteRune(',')
			fmt.Fprintf(&sb, "%d", bs.bp[i])
		}
	}

	return sb.String()
}

type processor struct {
	p imgpro.Processor
}

func (p *processor) processImage(ctx context.Context, inBytes []byte, f imgpro.Format, h, w int, outPath string) error {
	handler, err := p.p.HandlerFromBytes(inBytes)
	if err != nil {
		return fmt.Errorf("could not create a new handler: %v", err)
	}
	defer handler.Destroy()

	if err := handler.StripMetadata(); err != nil {
		return fmt.Errorf("could not strip metadata: %v", err)
	}

	if err := handler.SetFormat(f); err != nil {
		return fmt.Errorf("could not set the %s format: %v", f, err)
	}

	if err := handler.Resize(ctx, w, h); err != nil {
		return fmt.Errorf("could not resize to %dx%d: %v", w, h, err)
	}

	b, err := handler.Bytes()
	if err != nil {
		return fmt.Errorf("could not get the resulting bytes: %v", err)
	}

	return os.WriteFile(outPath, b, 0666)
}

func (p *processor) Set(s string) error {
	var ok bool

	if p.p, ok = processorMap[s]; !ok {
		return fmt.Errorf("%s: invalid processor", s)
	}

	return nil
}

func (p *processor) String() string {
	if p.p == nil {
		return "undefined"
	}

	return p.p.String()
}

type formatWithExt struct {
	ext  string
	f    imgpro.Format
	name string
}

var formats = []formatWithExt{
	{ext: "jpg", f: imgpro.JPEG, name: "jpg"},
	{ext: "webp", f: imgpro.Webp, name: "webp"},
}

type job struct {
	baseName string
	format   formatWithExt
	height   int
	outDir   string
	proc     processor
	srcBytes []byte
	width    int
}

func (j *job) run(ctx context.Context) error {
	ext := filepath.Ext(j.baseName)

	var dimension string

	if j.height != 0 {
		dimension = fmt.Sprintf("h%d", j.height)
	} else {
		dimension = fmt.Sprintf("w%d", j.width)
	}

	dstFilename := fmt.Sprintf(
		"%s_%s_%s.%s",
		j.baseName[:len(j.baseName)-len(ext)],
		dimension,
		j.format.name,
		j.format.ext,
	)

	dst := filepath.Join(j.outDir, dstFilename)

	log.Printf("Processing to %s", dst)

	return j.proc.processImage(ctx, j.srcBytes, j.format.f, j.height, j.width, dst)
}

type result struct {
	err error
	job *job
}

func main() {
	var (
		heightBreakpoints breakpointsSlice
		imgInDir          string
		imgOutDir         string
		proc              processor
		widthBreakpoints  breakpointsSlice
		workers           int
	)

	flag.Var(&heightBreakpoints, "height-breakpoints", "a comma-separated list of height breakpoints")
	flag.StringVar(&imgInDir, "img-in-dir", "img-src", "directory in which source img-src are stored")
	flag.StringVar(&imgOutDir, "img-out-dir", "img-src", "directory in which img-src are generated")
	flag.IntVar(&workers, "parallel", runtime.NumCPU(), "the number of parallel goroutines in the processing pool")
	flag.Var(&proc, "processor", "the image processor to use to prepare img-src")
	flag.Var(&widthBreakpoints, "width-breakpoints", "a comma-separated list of width breakpoints")

	flag.Parse()

	log.Printf("%v", heightBreakpoints.bp)
	log.Printf("%s", heightBreakpoints.String())

	if proc.p == nil {
		log.Fatal("A processor must be defined.")
	}

	log.Printf("Using processor %q", &proc)

	if err := proc.p.Init(workers); err != nil {
		log.Fatalf("Could not initialize processor: %v", err)
	}
	defer proc.p.Destroy()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inMetadataPath := filepath.Join(imgInDir, metadataFileName)

	meta, err := metadata.ReadFromFile(inMetadataPath)
	if err != nil {
		log.Fatalf("Could not read the metadata file %s: %v", inMetadataPath, err)
	}

	if err := os.MkdirAll(imgOutDir, 0744|os.ModeDir); err != nil {
		log.Fatalf("Could not create the destination directory: %v", err)
	}

	var wg sync.WaitGroup

	jobs := make(chan *job)
	results := make(chan *result)

	log.Printf("Using %d goroutines", workers)

	// create workers
	for i := 0; i < workers; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case j := <-jobs:
					if err := j.run(ctx); err != nil {
						results <- &result{err: err}
					}

					results <- &result{job: j}
					wg.Done()
				}
			}

		}()
	}

	// watch results
	go func() {
		for {
			select {
			case r := <-results:
				if err = r.err; err != nil {
					cancel()
				} else {
					// Update metadata
					// TODO
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Without this, the goroutine that runs wg.Wait() returns before even one
	// job has been sent on jobs.
	// Done() deferred to the end of the producing goroutine
	wg.Add(1)

	// Producer goroutine
	go func() {
		defer wg.Done()

		for baseName, _ := range meta {
			fullPath := filepath.Join(imgInDir, baseName)

			b, err := os.ReadFile(fullPath)
			if err != nil {
				log.Fatalf("Could not read %s: %v", fullPath, err)
			}

			for _, format := range formats {
				for _, h := range heightBreakpoints.bp {
					j := &job{
						baseName: baseName,
						format:   format,
						height:   h,
						outDir:   imgOutDir,
						proc:     proc,
						srcBytes: b,
					}

					select {
					case jobs <- j:
						wg.Add(1)
					case <-ctx.Done():
						return
					}
				}

				for _, w := range widthBreakpoints.bp {
					j := &job{
						baseName: baseName,
						format:   format,
						width:    w,
						outDir:   imgOutDir,
						proc:     proc,
						srcBytes: b,
					}

					select {
					case jobs <- j:
						wg.Add(1)
					case <-ctx.Done():
						return
					}
				}
			}
		}

		wg.Done()
	}()

	go func() {
		wg.Wait()
		cancel()
	}()

	<-ctx.Done()

	if err != nil {
		log.Fatalf("Processing error: %v", err)
	}

	outMetadataPath := filepath.Join(imgOutDir, metadataFileName)

	if err := meta.WriteToFile(outMetadataPath); err != nil {
		log.Fatalf("could not write the out metadata to %s: %v", outMetadataPath, err)
	}
}
