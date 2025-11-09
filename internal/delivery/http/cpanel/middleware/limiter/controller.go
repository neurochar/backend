package limiter

type Controller struct {
	mws map[string]*Middleware
}

func NewController() *Controller {
	return &Controller{
		mws: make(map[string]*Middleware),
	}
}

// DefaultName - default limiter with rate "1-S"
const DefaultName = "default"

const defaultRate = "1-S"

func (c *Controller) Get(name string, rateFormatted ...string) *Middleware {
	rateVal := defaultRate

	if name != DefaultName && len(rateFormatted) > 0 && rateFormatted[0] != "" {
		rateVal = rateFormatted[0]
	}

	mw, ok := c.mws[name]

	if !ok {
		mw = newMiddleware(rateVal)
		c.mws[name] = mw
	}

	return mw
}
