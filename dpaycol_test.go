package main

import (
	"testing"
)

var testStats stats

func init() {
	testStats = stats{
		Ak:    "None",
		Am:    0,
		Jn:    "None",
		Jkid:  "None",
		IsEnd: false,
		RC:    -1,
	}
}

func TestCheckMonth(t *testing.T) {
	month := testStats.checkMonth()
	if month != false {
		t.Errorf("CheckMonth was incorrect, got: %v, want: %v", month, false)
	}

	testStats.Am = 01
	month = testStats.checkMonth()
	if month != true {
		t.Errorf("CheckMonth was incorrect, got: %v, want: %v", month, true)
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
