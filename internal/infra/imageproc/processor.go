package imageproc

import (
	"bytes"
	"image"
	"image/png"
	"math"

	"github.com/disintegration/imaging"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

type imageProcImpl struct{}

func New() *imageProcImpl {
	return &imageProcImpl{}
}

func (p *imageProcImpl) targetFormat(currentFormat string, allowedFormats []Format, defaultFormat Format) Format {
	for _, format := range allowedFormats {
		if string(format) == currentFormat {
			return format
		}
	}

	return defaultFormat
}

func (p *imageProcImpl) convert(
	src image.Image,
	format Format,
	currentFileDataLen int,
	options ...option,
) ([]byte, *appErrors.AppError) {
	var buf bytes.Buffer

	cfg := defaultOptions()
	for _, option := range options {
		option(cfg)
	}

	switch format {
	case FormatJPEG:
		for _, q := range cfg.jpegQualities {
			buf.Reset()
			if err := imaging.Encode(&buf, src, imaging.JPEG, imaging.JPEGQuality(q)); err != nil {
				return nil, ErrCantConvertImage.WithWrap(err)
			}
			if buf.Len() <= currentFileDataLen {
				break
			}
		}
	case FormatPNG:
		if err := imaging.Encode(&buf, src, imaging.PNG, imaging.PNGCompressionLevel(png.BestCompression)); err != nil {
			return nil, ErrCantConvertImage.WithWrap(err)
		}
	default:
		return nil, ErrCantConvertImage.WithDetail("format", false, string(format))
	}

	return buf.Bytes(), nil
}

func (p *imageProcImpl) IsOpenable(fileData []byte) bool {
	_, err := imaging.Decode(bytes.NewReader(fileData))
	return err == nil
}

func (p *imageProcImpl) ScaleAndCrop(
	fileData []byte,
	width int,
	height int,
	options ...option,
) ([]byte, *appErrors.AppError) {
	cfg := defaultOptions()
	for _, option := range options {
		option(cfg)
	}

	_, format, err := image.DecodeConfig(bytes.NewReader(fileData))
	if err != nil {
		return nil, ErrInvalidImageFile.WithWrap(err)
	}

	srcImg, err := imaging.Decode(bytes.NewReader(fileData))
	if err != nil {
		return nil, ErrInvalidImageFile.WithWrap(err)
	}

	b := srcImg.Bounds()
	w, h := b.Dx(), b.Dy()

	if w < width || h < height {
		if !cfg.useUpscale {
			return nil, ErrInvalidImageSize.
				WithDetail("minWidth", false, width).
				WithDetail("minHeight", false, height)
		}

		scale := math.Max(float64(width)/float64(w), float64(height)/float64(h))
		newW := int(math.Ceil(float64(w) * scale))
		newH := int(math.Ceil(float64(h) * scale))

		srcImg = imaging.Resize(srcImg, newW, newH, imaging.Lanczos)
	}

	resizedImg := imaging.Fill(srcImg, width, height, imaging.Center, imaging.Lanczos)

	targetFormat := p.targetFormat(string(format), cfg.allowedFormats, cfg.defaultFormat)

	result, convErr := p.convert(resizedImg, targetFormat, len(fileData), WithJpegQualities(cfg.jpegQualities...))
	if convErr != nil {
		return nil, convErr
	}

	return result, nil
}

func (p *imageProcImpl) DownscaleIfLarger(
	fileData []byte,
	maxWidth int,
	maxHeight int,
	options ...option,
) ([]byte, *appErrors.AppError) {
	if maxWidth <= 0 || maxHeight <= 0 {
		return nil, appErrors.ErrBadRequest
	}

	cfg := defaultOptions()
	for _, opt := range options {
		opt(cfg)
	}

	_, format, err := image.DecodeConfig(bytes.NewReader(fileData))
	if err != nil {
		return nil, ErrInvalidImageFile.WithWrap(err)
	}

	srcImg, err := imaging.Decode(bytes.NewReader(fileData))
	if err != nil {
		return nil, ErrInvalidImageFile.WithWrap(err)
	}

	b := srcImg.Bounds()
	w, h := b.Dx(), b.Dy()

	var processed image.Image
	if w > maxWidth || h > maxHeight {
		processed = imaging.Fit(srcImg, maxWidth, maxHeight, imaging.Lanczos)
	} else {
		processed = srcImg
	}

	targetFormat := p.targetFormat(string(format), cfg.allowedFormats, cfg.defaultFormat)

	result, convErr := p.convert(processed, targetFormat, len(fileData), WithJpegQualities(cfg.jpegQualities...))
	if convErr != nil {
		return nil, convErr
	}

	return result, nil
}
