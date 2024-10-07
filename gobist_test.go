package gobist

import (
	"reflect"
	"testing"
	"time"
)

var (
	symbol = "TUPRS"
)

func TestCrumb(t *testing.T) {
	crumbStr := getYahooCrumb()
	if crumbStr == "Too Many Requests" {
		t.Fatalf(crumbStr)
	} else {
		t.Logf("crumb: %s", crumbStr)
	}
}

func TestBist_GetQuote(t *testing.T) {
	bist, _ := New()

	q, err := bist.GetQuote(symbol)
	assertError(t, err)
	assertNotNil(t, q)

	if q != nil {
		t.Logf("Symbol=%s Name=%s Price=%f", q.Symbol, q.Name, q.Price)
	}
}

func TestBist_GetQuoteWithHistory(t *testing.T) {
	bist, _ := New()

	q, err := bist.GetQuoteWithHistory(symbol, time.Now().Add(-1*24*time.Hour))
	assertError(t, err)
	assertNotNil(t, q)

	if q != nil {
		t.Logf("Symbol=%s Name=%s Current Price=%f History Price=%f", q.Symbol, q.Name, q.Price, q.History.Price)
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
