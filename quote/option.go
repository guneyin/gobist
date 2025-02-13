package quote

import "time"

type OptionFunc func(qf *fetcher)

type Option struct {
	period Period
}

func newDefaultOptions() *Option {
	opt := &Option{}
	opt.setPeriod(NewPeriod(time.Now(), time.Now()))

	return opt
}

func WithPeriod(p Period) OptionFunc {
	return func(qf *fetcher) {
		qf.opt.setPeriod(p)
	}
}

func (o *Option) setPeriod(p Period) {
	switch {
	case p.begin.isZero():
		dtToday := time.Now()
		o.period.begin.Set(dtToday)
		o.period.end.Set(dtToday)
	case p.end.isZero():
		o.period.begin = p.begin
		o.period.end = p.begin
	default:
		o.period.begin = p.begin
		o.period.end = p.end
	}
}
