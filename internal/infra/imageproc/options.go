package imageproc

type option func(*options)

type Format string

const (
	FormatJPEG Format = "jpeg"
	FormatPNG  Format = "png"
)

type options struct {
	jpegQualities  []int
	useUpscale     bool
	defaultFormat  Format
	allowedFormats []Format
}

func defaultOptions() *options {
	return &options{
		jpegQualities: []int{99, 95, 90, 85, 80, 75},
		useUpscale:    true,
		defaultFormat: FormatJPEG,
		allowedFormats: []Format{
			FormatJPEG,
			FormatPNG,
		},
	}
}

func WithJpegQualities(value ...int) option {
	return func(o *options) {
		o.jpegQualities = value
	}
}

func WithNoUpscale() option {
	return func(o *options) {
		o.useUpscale = false
	}
}

func WithAllowedFormats(formats ...Format) option {
	return func(o *options) {
		o.allowedFormats = formats
	}
}

func WithDefaultFormat(format Format) option {
	return func(o *options) {
		o.defaultFormat = format
	}
}
