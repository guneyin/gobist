package kap

type OptionFunc func(sc *scraper)

type Option struct {
	withShares bool
}

func WithShares() OptionFunc {
	return func(sc *scraper) {
		sc.opt.withShares = true
	}
}
