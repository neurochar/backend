package limiter

type Controller struct {
	mws map[string]*Limiter
}

func NewController() *Controller {
	return &Controller{
		mws: make(map[string]*Limiter),
	}
}

const DefaultName = "default"

const defaultRate = "1-S"

func (c *Controller) Get(name string, rateFormatted ...string) *Limiter {
	rateVal := defaultRate

	if name != DefaultName && len(rateFormatted) > 0 && rateFormatted[0] != "" {
		rateVal = rateFormatted[0]
	}

	mw, ok := c.mws[name]

	if !ok {
		mw = newLimiter(rateVal)
		c.mws[name] = mw
	}

	return mw
}
