package handler

import (
	"reflect"
	"testing"
)

func TestParseFocus(t *testing.T) {
	ok := []struct {
		raw  string
		want []string
	}{
		{"", nil},
		{"  ", nil},
		{"dsa-refresher", []string{"dsa-refresher"}},
		{"system-design, dsa-refresher", []string{"system-design", "dsa-refresher"}},
		{"a,b,c", []string{"a", "b", "c"}},
	}
	for _, c := range ok {
		got, valid := parseFocus(c.raw)
		if !valid || !reflect.DeepEqual(got, c.want) {
			t.Errorf("parseFocus(%q) = %v, %v; want %v, true", c.raw, got, valid, c.want)
		}
	}
	for _, raw := range []string{"a,b,c,d", "Bad_Slug", "a,,b", "a,a", "x;drop"} {
		if _, valid := parseFocus(raw); valid {
			t.Errorf("parseFocus(%q) should be rejected", raw)
		}
	}
}
