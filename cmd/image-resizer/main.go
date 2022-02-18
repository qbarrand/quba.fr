package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/qbarrand/quba.fr/internal/imgpro"
	"github.com/qbarrand/quba.fr/internal/metadata"
)

var processorMap = map[string]imgpro.Processor{
	"vips": &imgpro.VipsProcessor{},
}

type breakpoints struct {
	Heights []int `json:"heights"`
	Widths  []int `json:"widths"`
}

func readBreakpointsFromFile(path string) (*breakpoints, error) {
	bp := &breakpoints{}

	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not read the breakpoints file: %v", err)
	}
	defer fd.Close()

	return bp, json.NewDecoder(fd).Decode(bp)
}

type processor struct {
	p imgpro.Processor
}

func (p *processor) getResizedBytes(ctx context.Context, inBytes []byte, f imgpro.Format, h, w int) ([]byte, error) {
	handler, err := p.p.HandlerFromBytes(inBytes)
	if err != nil {
		return nil, fmt.Errorf("could not create a new handler: %v", err)
	}
	defer handler.Destroy()

	if err = handler.StripMetadata(); err != nil {
		return nil, fmt.Errorf("could not strip metadata: %v", err)
	}

	if err = handler.SetFormat(f); err != nil {
		return nil, fmt.Errorf("could not set the %s format: %v", f, err)
	}

	if err = handler.Resize(ctx, w, h); err != nil {
		return nil, fmt.Errorf("could not resize to %dx%d: %v", w, h, err)
	}

	return handler.Bytes()
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

type outVars struct {
	filename  string
	format    string
	height    int
	imageName string
	width     int
}

func (j *job) run(ctx context.Context) (*outVars, error) {
	ext := filepath.Ext(j.baseName)

	dimension := ""

	if j.height != 0 {
		dimension = fmt.Sprintf("_h%d", j.height)
	} else if j.width != 0 {
		dimension = fmt.Sprintf("_w%d", j.width)
	}

	b, err := j.proc.getResizedBytes(ctx, j.srcBytes, j.format.f, j.height, j.width)
	if err != nil {
		return nil, fmt.Errorf("could not get the resized bytes: %v", err)
	}

	hash := fnv.New32a()

	if _, err = hash.Write(b); err != nil {
		return nil, fmt.Errorf("could not write image bytes to the hasher: %v", err)
	}

	dstFilename := fmt.Sprintf(
		"%s%s_%s_%s.%s",
		j.baseName[:len(j.baseName)-len(ext)],
		dimension,
		j.format.name,
		hex.EncodeToString(hash.Sum(nil)),
		j.format.ext,
	)

	ov := &outVars{
		filename:  dstFilename,
		format:    j.format.name,
		height:    j.height,
		imageName: j.baseName,
		width:     j.width,
	}

	return ov, os.WriteFile(filepath.Join(j.outDir, dstFilename), b, 0644)
}

func main() {
	var (
		breakpointsFile string
		imgInDir        string
		outDir          string
		proc            processor
		workers         int
	)

	flag.StringVar(&imgInDir, "img-in-dir", "img-src", "directory in which source img-src are stored")
	flag.StringVar(&outDir, "img-out-dir", "img-src", "directory in which img-src are generated")
	flag.StringVar(&breakpointsFile, "bp-file", "config/breakpoints.json", "file in which breakpoints are defined")
	flag.IntVar(&workers, "parallel", runtime.NumCPU(), "the number of parallel goroutines in the processing pool")
	flag.Var(&proc, "processor", "the image processor to use to prepare img-src")

	flag.Parse()

	if proc.p == nil {
		log.Fatal("A processor must be defined.")
	}

	bp, err := readBreakpointsFromFile(breakpointsFile)
	if err != nil {
		log.Fatalf("could not read breakpoints: %v", err)
	}

	log.Printf("Using processor %q", &proc)

	if err = proc.p.Init(workers); err != nil {
		log.Fatalf("Could not initialize processor: %v", err)
	}
	defer proc.p.Destroy()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inMetadataPath := filepath.Join(imgInDir, "metadata.json")

	meta, err := ReadFromFile(inMetadataPath)
	if err != nil {
		log.Fatalf("Could not read the metadata file %s: %v", inMetadataPath, err)
	}

	mdbPath := filepath.Join(outDir, "metadata.db")

	if err = os.MkdirAll(outDir, 0744|os.ModeDir); err != nil {
		log.Fatalf("Could not create the destination directory: %v", err)
	}

	if err = os.Remove(mdbPath); err != nil && !os.IsNotExist(err) {
		log.Fatalf("could not remove the existing database (%s exists): %v", mdbPath, err)
	}

	mdb, err := metadata.OpenDB(mdbPath, false)
	if err != nil {
		log.Fatalf("Could not create the database: %v", err)
	}
	defer mdb.Close()

	if err = mdb.Init(ctx); err != nil {
		log.Fatalf("Could not initialize the database: %v", err)
	}

	// Add image metadata
	for k, v := range meta {
		if err = mdb.AddImage(ctx, k, v.Date, v.Location, v.MainColor); err != nil {
			log.Fatalf("Could not add %s to the metadata DB: %v", k, err)
		}
	}

	var wg sync.WaitGroup

	jobs := make(chan *job)

	log.Printf("Using %d goroutines", workers)

	// create workers
	for i := 0; i < workers; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case j := <-jobs:
					ov, err := j.run(ctx)
					if err != nil {
						cancel()
						log.Fatalf("Error running a job for %s: %v", j.baseName, err)
					}

					if err = mdb.AddWebImage(ctx, ov.filename, ov.imageName, ov.width, ov.height, ov.format); err != nil {
						cancel()
						log.Fatalf("Error adding the image to the database: %v", err)
					}

					wg.Done()
				}
			}
		}()
	}

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
				cancel()
				log.Fatalf("Could not read %s: %v", fullPath, err)
			}

			for _, format := range formats {
				// First, one job that does not resize images to use native resolution
				j := &job{
					baseName: baseName,
					format:   format,
					outDir:   outDir,
					proc:     proc,
					srcBytes: b,
				}

				select {
				case jobs <- j:
					wg.Add(1)
				case <-ctx.Done():
					return
				}

				for _, h := range bp.Heights {
					j := &job{
						baseName: baseName,
						format:   format,
						height:   h,
						outDir:   outDir,
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

				for _, w := range bp.Widths {
					j := &job{
						baseName: baseName,
						format:   format,
						width:    w,
						outDir:   outDir,
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
	}()

	go func() {
		wg.Wait()
		cancel()
	}()

	<-ctx.Done()

	if err != nil {
		log.Fatalf("Processing error: %v", err)
	}
}
