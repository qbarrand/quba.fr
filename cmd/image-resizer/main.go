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

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/qbarrand/quba.fr/internal/metadata"
)

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

const (
	formatJPEG = "jpg"
	formatWebp = "webp"
)

type processor struct {
	mdb *metadata.DB
}

type formatWithExt struct {
	ext  string
	name string
}

var formats = []formatWithExt{
	{ext: "jpg", name: formatJPEG},
	{ext: "webp", name: formatWebp},
}

type params struct {
	baseName string
	format   formatWithExt
	height   int
	outDir   string
	ref      *vips.ImageRef
	width    int
}

func (r *processor) resize(ctx context.Context, p *params) error {
	var (
		dimension    string
		err          error
		shouldResize bool
	)

	if p.height != 0 {
		dimension = fmt.Sprintf("_h%d", p.height)
		shouldResize = true
	} else if p.width != 0 {
		dimension = fmt.Sprintf("_w%d", p.width)
		shouldResize = true
	}

	if shouldResize {
		var scale float64 = 1

		if p.width != 0 {
			scale = float64(p.width) / float64(p.ref.Width())
		} else if p.height != 0 {
			scale = float64(p.height) / float64(p.ref.Height())
		}

		if scale <= 1 {
			if err = p.ref.Resize(scale, vips.KernelAuto); err != nil {
				return fmt.Errorf("error resizing %s to h=%d w=%d: %v", p.baseName, p.height, p.width, err)
			}
		}
	}

	var b []byte

	switch p.format.name {
	case formatJPEG:
		ep := vips.NewJpegExportParams()
		ep.StripMetadata = true
		b, _, err = p.ref.ExportJpeg(ep)
	case formatWebp:
		ep := vips.NewWebpExportParams()
		ep.ReductionEffort = 6
		ep.StripMetadata = true
		b, _, err = p.ref.ExportWebp(ep)
	}

	if err != nil {
		return fmt.Errorf("could not export as %s: %v", p.format, err)
	}

	hash := fnv.New32a()

	if _, err = hash.Write(b); err != nil {
		return fmt.Errorf("could not write image bytes to the hasher: %v", err)
	}

	ext := filepath.Ext(p.baseName)

	dstFilename := fmt.Sprintf(
		"%s%s_%s_%s.%s",
		p.baseName[:len(p.baseName)-len(ext)],
		dimension,
		p.format.name,
		hex.EncodeToString(hash.Sum(nil)),
		p.format.ext,
	)

	if err = os.WriteFile(filepath.Join(p.outDir, dstFilename), b, 0644); err != nil {
		return fmt.Errorf("could not write image %s: %v", dstFilename, err)
	}

	if err = r.mdb.AddWebImage(ctx, dstFilename, p.baseName, p.width, p.height, p.format.name); err != nil {
		return fmt.Errorf("error adding the image %s to the database: %v", dstFilename, err)
	}

	return nil
}

func main() {
	var (
		breakpointsFile string
		concurrency     int
		imgInDir        string
		outDir          string
	)

	flag.StringVar(&imgInDir, "img-in-dir", "img-src", "directory in which source images are stored")
	flag.StringVar(&outDir, "img-out-dir", "img-out", "directory in which images are generated")
	flag.StringVar(&breakpointsFile, "bp-file", "config/breakpoints.json", "file in which breakpoints are defined")
	flag.IntVar(&concurrency, "concurrency", runtime.NumCPU(), "the VIPS concurrency level")

	flag.Parse()

	bp, err := readBreakpointsFromFile(breakpointsFile)
	if err != nil {
		log.Fatalf("could not read breakpoints: %v", err)
	}

	log.Printf("Using up to %d VIPS threads", concurrency)

	vips.Startup(&vips.Config{ConcurrencyLevel: concurrency})
	defer vips.Shutdown()

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

	rs := &processor{mdb: mdb}

	for baseName := range meta {
		log.Printf("Processing %s", baseName)

		fullPath := filepath.Join(imgInDir, baseName)

		origImg, err := vips.LoadImageFromFile(fullPath, vips.NewImportParams())
		if err != nil {
			log.Fatalf("Could not load image %s", fullPath)
		}

		for _, format := range formats {
			var copyImg *vips.ImageRef

			// First, one job that does not resize images to use native resolution
			copyImg, err = origImg.Copy()
			if err != nil {
				log.Fatalf("Could not copy image: %v", err)
			}

			p := params{
				baseName: baseName,
				format:   format,
				outDir:   outDir,
				ref:      copyImg,
			}

			if err = rs.resize(ctx, &p); err != nil {
				log.Fatalf("Could not write the image with native resolution :%v", err)
			}

			for _, h := range bp.Heights {
				p.ref, err = origImg.Copy()
				if err != nil {
					log.Fatalf("Could not copy image: %v", err)
				}

				p.height = h

				if err = rs.resize(ctx, &p); err != nil {
					log.Fatalf("Could not write the image with height=%d resolution :%v", h, err)
				}
			}

			p.height = 0

			for _, w := range bp.Widths {
				p.ref, err = origImg.Copy()
				if err != nil {
					log.Fatalf("Could not copy image: %v", err)
				}

				p.width = w

				if err = rs.resize(ctx, &p); err != nil {
					log.Fatalf("Could not write the image with width=%d resolution :%v", w, err)
				}
			}
		}
	}

	if err != nil {
		log.Fatalf("Processing error: %v", err)
	}
}
