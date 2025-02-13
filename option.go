package gobist

import "time"

type QuoteOptionFunc func(qf *quoteFetcher)

type QuoteOption struct {
	period Period
}

func NewDefaultOptions() *QuoteOption {
	opt := &QuoteOption{}
	opt.setPeriod(NewPeriod(time.Now(), time.Now()))

	return opt
}

func WithPeriod(p Period) QuoteOptionFunc {
	return func(qf *quoteFetcher) {
		qf.opt.setPeriod(p)
	}
}

func (o *QuoteOption) setPeriod(p Period) {
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
