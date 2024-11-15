package gobist

import (
	"reflect"
	"testing"
	"time"
)

var (
	symbols = []string{"TUPRS", "BIMAS", "KCHOL"}
)

func TestCrumb(t *testing.T) {
	crumbStr, err := setYahooCrumb()
	if err != nil {
		t.Fatal(err)
	}
	if crumbStr == "Too Many Requests" {
		t.Fatalf(crumbStr)
	} else {
		t.Logf("crumb: %s", crumbStr)
	}
}

func TestBist_GetQuote(t *testing.T) {
	bist := New()

	q, err := bist.GetQuote(symbols)
	assertError(t, err)
	assertNotNil(t, q)

	if q != nil {
		t.Logf(q.ToJson())
	}
}

func TestBist_GetQuoteWithHistory(t *testing.T) {
	bist := New()

	d1, _ := time.Parse(time.DateOnly, "2024-10-06")
	d2, _ := time.Parse(time.DateOnly, "2024-10-13")
	q, err := bist.GetQuoteWithHistory(symbols, d1, d2)
	assertError(t, err)
	assertNotNil(t, q)

	if q != nil {
		t.Logf(q.ToJson())
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
