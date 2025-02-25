package quote_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/guneyin/gobist/quote"
	"github.com/guneyin/gobist/store"
)

var (
	symbols = []string{"TUPRS", "BIMAS", "KCHOL"}
)

func TestCrumb(t *testing.T) {
	crumbStr, err := quote.SetYahooCrumb()
	if err != nil {
		t.Fatal(err)
	}
	if crumbStr == "Too Many Requests" {
		t.Fatal(crumbStr)
	} else {
		t.Logf("crumb: %s", crumbStr)
	}
}

func Test_GetQuoteWithHistory(t *testing.T) {
	ctx := context.Background()
	api := quote.NewAPI(store.NewMemoryStore())

	d1, _ := time.Parse(time.DateOnly, "2024-10-06")
	d2, _ := time.Parse(time.DateOnly, "2024-10-13")

	q, err := api.GetQuoteList(ctx, symbols, quote.WithPeriod(quote.NewPeriod(d1, d2)))
	assertError(t, err)
	assertNotNil(t, q)

	if q != nil {
		t.Log(q.ToJSON())
	}
}

func assertError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error occurred [%v]", err)
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	if isNil(v) {
		t.Errorf("[%v] was expected to be non-nil", v)
	}
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	kind := rv.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && rv.IsNil() {
		return true
	}

	return false
}
