package quote

import "time"

type pDate time.Time

type Period struct {
	begin, end pDate
}

// pDate

func (pd *pDate) String() string {
	return time.Time(*pd).Format("2006-01-02")
}

func (pd *pDate) Unix() int64 {
	return time.Time(*pd).Unix()
}

func (pd *pDate) isZero() bool {
	return time.Time(*pd).IsZero()
}

func (pd *pDate) Set(dt time.Time) {
	if isToday(dt) {
		dt = time.Now()
	} else {
		y, m, d := dt.Date()
		dt = time.Date(y, m, d, 11, 0, 0, 0, dt.Location())
	}

	*pd = pDate(dt)
}

// period

func NewPeriod(begin, end time.Time) Period {
	return Period{pDate(begin), pDate(end)}
}

func (p *Period) isSingleDay() bool {
	bs, es := p.begin.String(), p.end.String()
	return bs == es
}

// utils

func isToday(d time.Time) bool {
	return d.Format(time.DateOnly) == time.Now().Format(time.DateOnly)
}
