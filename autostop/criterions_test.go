package autostop

import (
	"regexp"
	"testing"
)

func TestAvgTimeCriterion_Parse(t *testing.T) {
	testSuccessData := map[string]*AvgTimeCriterion{
		" 1s501ms ,30s,sometag": {
			ThresholdDuration: 30,
			Tag: "sometag",
			ThresholdValue: 1501,
		},
		"3m ,10": {
			ThresholdDuration: 10,
			Tag: "",
			ThresholdValue: 180000,
		},
	}
	testFailData := []string{
		"", "4dh", "s1", "ppppp", "mmm", "3h141pm5", "3ms", "(3h141m5)", "()", " ",
	}
	for s, crit := range testSuccessData {
		c := &AvgTimeCriterion{}
		err := c.Parse(s)
		if err != nil {
			t.Errorf("failed to parse AvgTimeCriterion params from %s", s)
		}
		if c.ThresholdDuration != crit.ThresholdDuration || c.Tag != crit.Tag || c.ThresholdValue != crit.ThresholdValue {
			t.Errorf("wrong AvgTimeCriterion params got from %s: expected %v actual %v", s, crit, c)
		}
	}
	for _, s := range testFailData {
		c := &AvgTimeCriterion{}
		err := c.Parse(s)
		if err == nil {
			t.Errorf("did not fail on wrong params for AvgTimeCriterion %s", s)
		}
	}
}

func TestNetCodeCriterion_Parse(t *testing.T) {
	successNetCodes1, _ := regexp.Compile("^4[0-9][0-9]$")
	successNetCodes2, _ := regexp.Compile("^110$")
	testSuccessData := map[string]*NetCodeCriterion{
		" 4xx ,123, 30": {
			ThresholdDuration: 30,
			Tag: "",
			ThresholdValue: 123,
			ThresholdValueType: "count",
			Codes: successNetCodes1,
		},
		"110 ,10%,4d,     sometag": {
			ThresholdDuration: 4*24*3600,
			Tag: "sometag",
			ThresholdValue: 10,
			ThresholdValueType: "percent",
			Codes: successNetCodes2,
		},
	}
	testFailData := []string{
		" 4hx ,123, 30", "110 ,10%,4d     sometag", "", "4hx ,123, 30, hfhf, lsls", "ppppp", ".", " ",
	}
	for s, crit := range testSuccessData {
		c := &NetCodeCriterion{}
		err := c.Parse(s)
		if err != nil {
			t.Errorf("failed to parse NetCodeCriterion params from %s", s)
		}
		if c.ThresholdDuration != crit.ThresholdDuration || c.Tag != crit.Tag || c.ThresholdValue != crit.ThresholdValue || c.ThresholdValueType != crit.ThresholdValueType || c.Codes.String() != crit.Codes.String() {
			t.Errorf("wrong NetCodeCriterion params got from %s: expected %v actual %v", s, crit, c)
		}
	}
	for _, s := range testFailData {
		c := &NetCodeCriterion{}
		err := c.Parse(s)
		if err == nil {
			t.Errorf("did not fail on wrong params for NetCodeCriterion %s", s)
		}
	}
}

func TestHttpCodeCriterion_Parse(t *testing.T) {
	successNetCodes1, _ := regexp.Compile("^4[0-9][0-9]$")
	successNetCodes2, _ := regexp.Compile("^110$")
	testSuccessData := map[string]*HttpCodeCriterion{
		" 4xx ,123, 30": {
			ThresholdDuration: 30,
			Tag: "",
			ThresholdValue: 123,
			ThresholdValueType: "count",
			Codes: successNetCodes1,
		},
		"110 ,10%,4d,     sometag": {
			ThresholdDuration: 4*24*3600,
			Tag: "sometag",
			ThresholdValue: 10,
			ThresholdValueType: "percent",
			Codes: successNetCodes2,
		},
	}
	testFailData := []string{
		" 4hx ,123, 30", "110 ,10%,4d     sometag", "", "4hx ,123, 30, hfhf, lsls", "ppppp", ".", " ",
	}
	for s, crit := range testSuccessData {
		c := &HttpCodeCriterion{}
		err := c.Parse(s)
		if err != nil {
			t.Errorf("failed to parse HttpCodeCriterion params from %s", s)
		}
		if c.ThresholdDuration != crit.ThresholdDuration || c.Tag != crit.Tag || c.ThresholdValue != crit.ThresholdValue || c.ThresholdValueType != crit.ThresholdValueType || c.Codes.String() != crit.Codes.String() {
			t.Errorf("wrong HttpCodeCriterion params got from %s: expected %v actual %v", s, crit, c)
		}
	}
	for _, s := range testFailData {
		c := &HttpCodeCriterion{}
		err := c.Parse(s)
		if err == nil {
			t.Errorf("did not fail on wrong params for HttpCodeCriterion %s", s)
		}
	}
}

func TestQuantileCriterion_Parse(t *testing.T) {
	testSuccessData := map[string]*QuantileCriterion{
		" 13 ,123, 30": {
			ThresholdDuration: 30,
			Tag: "",
			ThresholdValue: 123,
			Quantile: 13,
		},
		"99 ,1s1,4d,     sometag": {
			ThresholdDuration: 4*24*3600,
			Tag: "sometag",
			ThresholdValue: 1001,
			Quantile: 99,
		},
	}
	testFailData := []string{
		" 4hx ,123, 30", "110 ,10%,4d     sometag", "", "4hx ,123, 30, hfhf, lsls", "ppppp", ".", " ",
	}
	for s, crit := range testSuccessData {
		c := &QuantileCriterion{}
		err := c.Parse(s)
		if err != nil {
			t.Errorf("failed to parse QuantileCriterion params from %s", s)
		}
		if c.ThresholdDuration != crit.ThresholdDuration || c.Tag != crit.Tag || c.ThresholdValue != crit.ThresholdValue || c.Quantile != crit.Quantile {
			t.Errorf("wrong QuantileCriterion params got from %s: expected %v actual %v", s, crit, c)
		}
	}
	for _, s := range testFailData {
		c := &QuantileCriterion{}
		err := c.Parse(s)
		if err == nil {
			t.Errorf("did not fail on wrong params for QuantileCriterion %s", s)
		}
	}
}

func TestTimeLimitCriterion_Parse(t *testing.T) {
	testSuccessData := map[string]*TimeLimitCriterion{
		"30s": {
			ThresholdValue: 30,
		},
		"3m      ": {
			ThresholdValue: 180,
		},
	}
	testFailData := []string{
		"30s,sometag", " 4hx ,123, 30", "110 ,10%,4d     sometag", "", "4hx ,123, 30, hfhf, lsls", "ppppp", ".", " ",
	}
	for s, crit := range testSuccessData {
		c := &TimeLimitCriterion{}
		err := c.Parse(s)
		if err != nil {
			t.Errorf("failed to parse TimeLimitCriterion params from %s", s)
		}
		if c.ThresholdValue != crit.ThresholdValue {
			t.Errorf("wrong TimeLimitCriterion params got from %s: expected %v actual %v", s, crit, c)
		}
	}
	for _, s := range testFailData {
		c := &TimeLimitCriterion{}
		err := c.Parse(s)
		if err == nil {
			t.Errorf("did not fail on wrong params for TimeLimitCriterion %s", s)
		}
	}
}