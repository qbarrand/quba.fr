package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/qbarrand/quba.fr/internal/imgpro"
)

const (
	ExtensionJPEG = ".jpg"
	ExtensionWEBP = ".webp"
)

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

func (p *processor) Set(s string) error {
	switch s {
	case "magick":
		p.p = &imgpro.ImageMagickProcessor{}
	case "vips":
		p.p = &imgpro.VipsProcessor{}
	default:
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

func main() {
	var (
		heightBreakpoints breakpointsSlice
		imgInDir          string
		imgOutDir         string
		proc              processor
		widthBreakpoints  breakpointsSlice
	)

	flag.Var(&heightBreakpoints, "height-breakpoints", "a comma-separated list of height breakpoints")
	flag.StringVar(&imgInDir, "img-in-dir", "images", "directory in which source images are stored")
	flag.StringVar(&imgOutDir, "img-out-dir", "images-out", "directory in which images are generated")
	flag.Var(&proc, "processor", "the image processor to use to prepare images")
	flag.Var(&widthBreakpoints, "width-breakpoints", "a comma-separated list of width breakpoints")

	flag.Parse()

	log.Printf("%v", heightBreakpoints.bp)
	log.Printf("%s", heightBreakpoints.String())

	if proc.p == nil {
		log.Fatal("A processor must be defined.")
	}

	log.Printf("Using processor %q", &proc)
	defer proc.p.Destroy()

	filepath.WalkDir(imgInDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path != imgInDir && d.IsDir() {
			return filepath.SkipDir
		}

		log.Printf("Processing %s", path)

		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read %s: %v", path, err)
		}

		h, err := proc.p.HandlerFromBytes(b)
		if err != nil {
			return fmt.Errorf("could not create a new handler: %v", err)
		}
		defer h.Destroy()

		if err := h.StripMetadata(); err != nil {
			return fmt.Errorf("%s: could not strip metadata: %v", path, err)
		}

		// for

		return nil
	})
}
