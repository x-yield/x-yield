package autostop

import "testing"

func TestParseDuration(t *testing.T) {
	testSuccessData := map[string]int32{
		"1m1s": 61,
		"1m": 60,
		"4d5h": 5*3600+4*24*3600,
		"1s": 1,
		"3h141m5s": 19265,
		"55": 55,
		"3h141m5": 19265,
		"0": 0,
	}
	testFailData := []string{
		"", "4dh", "s1", "ppppp", "mmm", "3h141pm5", "3ms", "(3h141m5)", "()", " ",
	}
	for s, i := range testSuccessData {
		duration, err := ParseDuration(s)
		if err != nil {
			t.Errorf("failed to parse duration from %s", s)
		}
		if duration != i {
			t.Errorf("wrong duration %v got from %s", i, s)
		}
	}
	for _, s := range testFailData {
		_, err := ParseDuration(s)
		if err == nil {
			t.Errorf("did not fail on wrong duration string %s", s)
		}
	}
}

func TestParseMilliseconds(t *testing.T) {
	testSuccessData := map[string]int32{
		"1m1s": 61000,
		"1m": 60000,
		"4d5h": (5*3600+4*24*3600)*1000,
		"1s": 1000,
		"3h141m5s": 19265000,
		"55": 55,
		"3h141m5": 19260005,
		"0": 0,
		"1s500ms": 1500,
	}
	testFailData := []string{
		"", "4dh", "s1", "ppppp", "mmm", "3h141mss", "3msms", "(3h141m5)", "()", " ",
	}
	for s, i := range testSuccessData {
		duration, err := ParseMilliseconds(s)
		if err != nil {
			t.Errorf("failed to parse milliseconds from %s", s)
		}
		if duration != i {
			t.Errorf("wrong milliseconds %v got from %s", i, s)
		}
	}
	for _, s := range testFailData {
		_, err := ParseMilliseconds(s)
		if err == nil {
			t.Errorf("did not fail on wrong milliseconds string %s", s)
		}
	}
}