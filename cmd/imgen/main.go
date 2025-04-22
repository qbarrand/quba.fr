package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"hash/fnv"
	"iter"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sync"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func slogFatal(err error, msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

type runtimeConfig struct {
	inDir  string
	outDir string
}

type InputImage struct {
	ImageMetadata

	filename string
}

type orientation string

const (
	orientationLandscape orientation = "landscape"
	orientationPortrait  orientation = "portrait"

	overflowPercent = 130
)

type Source struct {
	Filename string `json:"filename"`
	Length   uint   `json:"length"`
}

type mimeType = string

type OrientationTree map[orientation][]Source

func NewOrientationTree() OrientationTree {
	return OrientationTree{
		orientationLandscape: make([]Source, 0),
		orientationPortrait:  make([]Source, 0),
	}
}

type MimeTree struct {
	Avif OrientationTree `json:"avif"`
	Webp OrientationTree `json:"webp"`
	Jpeg OrientationTree `json:"jpeg"`
}

type Picture struct {
	Date      string   `json:"date"`
	Location  string   `json:"location"`
	MainColor string   `json:"main_color"`
	Tree      MimeTree `json:"tree"`
}

func NewPicture(date, location, mainColor string) *Picture {
	return &Picture{
		Date:      date,
		Location:  location,
		MainColor: mainColor,
		Tree: MimeTree{
			Avif: NewOrientationTree(),
			Webp: NewOrientationTree(),
			Jpeg: NewOrientationTree(),
		},
	}
}

func (p *Picture) AddSource(or orientation, mt mimeType, length uint, filename string) {
	s := Source{
		Filename: filename,
		Length:   length,
	}

	switch mt {
	case avif.mimeType:
		p.Tree.Avif[or] = append(p.Tree.Avif[or], s)
	case webp.mimeType:
		p.Tree.Webp[or] = append(p.Tree.Webp[or], s)
	case jpeg.mimeType:
		p.Tree.Jpeg[or] = append(p.Tree.Jpeg[or], s)
	}
}

type formatWithExt struct {
	ext        string
	magickName string
	mimeType   mimeType
}

var (
	heights = []uint{480, 640, 800, 1000, 1200, 2000}
	widths  = []uint{576, 768, 992, 1200, 1400, 1600, 2000}

	avif = formatWithExt{ext: "avif", magickName: "avif", mimeType: "image/avif"}
	jpeg = formatWithExt{ext: "jpg", magickName: "jpg", mimeType: "image/jpeg"}
	webp = formatWithExt{ext: "webp", magickName: "webp", mimeType: "image/webp"}

	formats = []formatWithExt{avif, jpeg, webp}
)

func prepareImage(mw *imagick.MagickWand, wp, hp *uint, f formatWithExt) ([]byte, error) {
	var height, width uint

	switch {
	case wp != nil && hp != nil:
		return nil, errors.New("cannot set both width and height")
	case hp != nil:
		ratio := float32(*hp) / float32(mw.GetImageHeight())
		width = uint(float32(mw.GetImageWidth()) * ratio)
		height = *hp
	case wp != nil:
		ratio := float32(*wp) / float32(mw.GetImageWidth())
		height = uint(float32(mw.GetImageHeight()) * ratio)
		width = *wp
	}

	if width > 0 && height > 0 {
		if err := mw.ResizeImage(width, height, imagick.FILTER_LANCZOS); err != nil {
			return nil, fmt.Errorf("could not resize image: %v", err)
		}
	}

	if err := mw.SetImageFormat(f.magickName); err != nil {
		return nil, fmt.Errorf("could not set image format: %v", err)
	}

	if err := mw.SetCompressionQuality(80); err != nil {
		return nil, fmt.Errorf("could not set image compression quality: %v", err)
	}

	b, err := mw.GetImageBlob()
	if err != nil {
		return nil, fmt.Errorf("could not get image blob: %v", err)
	}

	hash := fnv.New32a()

	if _, err := hash.Write(b); err != nil {
		return nil, fmt.Errorf("could not write image blob: %v", err)
	}

	return hash.Sum(nil), nil
}

func mainColor(mw *imagick.MagickWand) (string, error) {
	if err := mw.SetDepth(8); err != nil {
		return "", fmt.Errorf("could not set image depth: %v", err)
	}

	if err := mw.QuantizeImage(1, imagick.COLORSPACE_UNDEFINED, 0, imagick.DITHER_METHOD_FLOYD_STEINBERG, false); err != nil {
		return "", fmt.Errorf("error quantizing image: %v", err)
	}

	colors, pws := mw.GetImageHistogram()
	if colors != 1 {
		return "", fmt.Errorf("expected 1 color, got %d", colors)
	}

	red := uint(pws[0].GetRed() * 256)
	green := uint(pws[0].GetGreen() * 256)
	blue := uint(pws[0].GetBlue() * 256)

	return fmt.Sprintf("#%X%X%X", red, green, blue), nil
}

type pool struct {
	errs    chan error
	inputs  chan *InputImage
	rc      *runtimeConfig
	res     chan *Picture
	wg      *sync.WaitGroup
	workers int
}

func newPool(rc *runtimeConfig, workers int) *pool {
	var wg sync.WaitGroup
	wg.Add(workers)

	return &pool{
		errs:    make(chan error),
		inputs:  make(chan *InputImage),
		rc:      rc,
		res:     make(chan *Picture),
		workers: workers,
		wg:      &wg,
	}
}

func (p *pool) StartWorkers(ctx context.Context) {
	for i := range p.workers {
		go p.worker(ctx, i)
	}

	go func() {
		p.wg.Wait()

		slog.Debug("Closing output channels")
		close(p.res)
	}()
}

func (p *pool) Errs() <-chan error {
	return p.errs
}

func (p *pool) Inputs() chan<- *InputImage {
	return p.inputs
}

func (p *pool) Outputs() <-chan *Picture {
	return p.res
}

func (p *pool) processImage(mw *imagick.MagickWand, ii *InputImage) (*Picture, error) {
	slog.Info("Processing image", "filename", ii.filename)

	ext := filepath.Ext(filepath.Base(ii.filename))
	imageNameWithoutExt := ii.filename[:len(ii.filename)-len(ext)]

	imagePath := filepath.Join(p.rc.inDir, ii.filename)

	mw.Clear()

	if err := mw.ReadImage(imagePath); err != nil {
		return nil, fmt.Errorf("could not read image %s: %v", imagePath, err)
	}

	if err := mw.StripImage(); err != nil {
		return nil, fmt.Errorf("could not strip image: %v", err)
	}

	slog.Info("Getting main color", "image", imagePath)

	mwMainColor := mw.Clone()

	mainColor, err := mainColor(mwMainColor)
	if err != nil {
		return nil, fmt.Errorf("could not get main color for image %s: %v", imagePath, err)
	}

	pic := NewPicture(ii.Date, ii.Location, mainColor)

	for _, w := range widths {
		ofWidth := (w * overflowPercent) / 100

		for _, format := range formats {
			mw := mw.Clone()

			hashBytes, err := prepareImage(mw, &ofWidth, nil, format)
			if err != nil {
				return nil, fmt.Errorf("could not prepare image %s: %v", imagePath, err)
			}

			outFileName := fmt.Sprintf(
				"%s_%dw_%s.%s",
				imageNameWithoutExt,
				ofWidth,
				hex.EncodeToString(hashBytes),
				format.ext)

			outPath := filepath.Join(p.rc.outDir, outFileName)

			slog.Info("Writing image", "path", outPath)

			if err = mw.WriteImage(outPath); err != nil {
				return nil, fmt.Errorf("could not write image %s: %v", outPath, err)
			}

			pic.AddSource(orientationLandscape, format.mimeType, w, outFileName)
		}
	}

	for _, h := range heights {
		for _, format := range formats {
			mw := mw.Clone()

			hashBytes, err := prepareImage(mw, nil, &h, format)
			if err != nil {
				return nil, fmt.Errorf("could not prepare image %s: %v", imagePath, err)
			}

			outFileName := fmt.Sprintf(
				"%s_%dh_%s.%s",
				imageNameWithoutExt,
				h,
				hex.EncodeToString(hashBytes),
				format.ext)

			outPath := filepath.Join(p.rc.outDir, outFileName)

			slog.Info("Writing image", "path", outPath)

			if err = mw.WriteImage(outPath); err != nil {
				return nil, fmt.Errorf("could not write image %s: %v", outPath, err)
			}

			pic.AddSource(orientationPortrait, format.mimeType, h, outFileName)
		}
	}

	// native resolution
	for _, format := range formats {
		hashBytes, err := prepareImage(mw, nil, nil, format)
		if err != nil {
			return nil, fmt.Errorf("could not prepare image in native resolution %s: %v", imagePath, err)
		}

		outFileName := fmt.Sprintf("%s_%s.%s",
			imageNameWithoutExt,
			hex.EncodeToString(hashBytes),
			format.ext)

		outPath := filepath.Join(p.rc.outDir, outFileName)

		slog.Info("Writing image", "path", outPath)

		if err = mw.WriteImage(outPath); err != nil {
			return nil, fmt.Errorf("could not write image %s: %v", outPath, err)
		}

		pic.AddSource(orientationPortrait, format.mimeType, 0, outFileName)
		pic.AddSource(orientationLandscape, format.mimeType, 0, outFileName)
	}

	return pic, nil
}

func (p *pool) worker(ctx context.Context, i int) {
	defer p.wg.Done()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	logger := slog.With("worker", i)

	for {
		select {
		case ii, more := <-p.inputs:
			if !more {
				logger.Info("Input channel closed, exiting")
				return
			}

			logger.Info("Processing image", "image", ii.filename)
			pic, err := p.processImage(mw, ii)
			if err != nil {
				p.errs <- fmt.Errorf("worker %d: %v", i, err)
				return
			}

			logger.Debug("Submitting result", "image", ii.filename)
			p.res <- pic

			logger.Debug("Clearing MagickWand")
			mw.Clear()
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	var (
		inDir           string
		imagesToProcess iter.Seq[string]
		logLevel        slog.Level
		outDir          string
		workers         int
	)

	flag.StringVar(&inDir, "in-dir", "img-src", "directory in which source images are stored")
	flag.TextVar(&logLevel, "log-level", slog.LevelInfo, "log level")
	flag.StringVar(&outDir, "out-dir", "img-out", "directory in which images are generated")
	flag.IntVar(&workers, "workers", runtime.NumCPU(), "number of workers to spawn")

	flag.Parse()

	slog.SetLogLoggerLevel(logLevel)
	slog.Debug("Using log level", "level", logLevel)

	if workers < 1 {
		workers = 1
	}

	metadataFilePath := filepath.Join(inDir, "metadata.json")

	slog.Info("Reading metadata file", "path", metadataFilePath)

	m, err := ReadFromFile(metadataFilePath)
	if err != nil {
		slogFatal(err, "Error reading metadata file", "path", metadataFilePath)
	}

	if flag.NArg() > 0 {
		imagesToProcess = slices.Values(flag.Args())
	} else {
		imagesToProcess = maps.Keys(m)
	}

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	rc := runtimeConfig{
		inDir:  inDir,
		outDir: outDir,
	}

	var (
		ctx, cancel = context.WithCancel(context.Background())
		p           = newPool(&rc, workers)
		pictures    = make([]Picture, 0)
	)

	errs := p.Errs()
	inputs := p.Inputs()
	outputs := p.Outputs()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case err := <-errs:
				slog.Error("Error from worker", "err", err)
				cancel()
			case p := <-outputs:
				if p == nil {
					// channel is closed
					return
				}

				slog.Debug("got result", "date", p.Date, "location", p.Location, "mainColor", p.MainColor)
				pictures = append(pictures, *p)
			}
		}
	}()

	slog.Info("Starting workers", "count", workers)
	p.StartWorkers(ctx)

	for imageName := range imagesToProcess {
		slog.Debug("Submitting image", "name", imageName)

		ii := InputImage{
			ImageMetadata: ImageMetadata{
				Date:     m[imageName].Date,
				Location: m[imageName].Location,
			},
			filename: imageName,
		}

		inputs <- &ii
	}

	close(inputs)
	wg.Wait()

	backgroundsFilePath := filepath.Join(outDir, "backgrounds.json")

	fd, err := os.Create(backgroundsFilePath)
	if err != nil {
		slogFatal(err, "Error creating backgrounds file", "path", backgroundsFilePath)
	}
	defer func() {
		if cerr := fd.Close(); cerr != nil {
			slogFatal(cerr, "Error closing backgrounds file", "path", backgroundsFilePath)
		}
	}()

	encoder := json.NewEncoder(fd)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(pictures); err != nil {
		slogFatal(err, "Error encoding backgrounds file into JSON")
	}
}
