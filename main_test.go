package main

import (
	"testing"
)

var testStats stats

func init() {
	testStats = stats{
		Ak:    "None",
		Am:    "None",
		Jn:    "None",
		Jkid:  "None",
		IsEnd: false,
		RC:    -1,
	}
}

func TestCheckMonth(t *testing.T) {

	tests := map[string]struct {
		input string
		want  bool
	}{
		"correct year and month < 10":       {"209903", true},
		"correct year and month > 9":        {"209911", true},
		"correct year and wrong month > 12": {"209913", false},
		"correct year and wrong month > 20": {"209922", false},
		"correct year and wrong month = 0":  {"209900", false},
	}

	for name, tc := range tests {
		testStats.Am = tc.input
		got := testStats.checkMonth()
		if got != tc.want {
			t.Fatalf("%s: expected: %v, got: %v", name, tc.want, got)
		}
	}
}

func TestCheckRequired(t *testing.T) {
	required := testStats.checkRequired()
	if required != false {
		t.Errorf("CheckRequired was incorrect, got: %v, want: %v", required, false)
	}
}

func TestCheckEnd(t *testing.T) {
	end := testStats.checkEnd()
	if end != true {
		t.Errorf("CheckEnd was incorrect, got: %v, want: %v", end, true)
	}

	testStats.RC = 0
	end = testStats.checkEnd()
	if end != false {
		t.Errorf("CheckEnd was incorrect, got: %v, want: %v", end, false)
	}

	testStats.IsEnd = true
	end = testStats.checkEnd()
	if end != true {
		t.Errorf("CheckEnd was incorrect, got: %v, want: %v", end, true)
	}

	testStats.RC = -1
	end = testStats.checkEnd()
	if end != false {
		t.Errorf("CheckEnd was incorrect, got: %v, want: %v", end, false)
	}
}
