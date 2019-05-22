package autostop

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Criterion interface {
	Parse(params string) error
	Check(second []string) bool
}

type AvgTimeCriterion struct {
	Since             int64
	Progress          int32
	ThresholdDuration int32
	ThresholdValue    int32
	Tag               string
}

type NetCodeCriterion struct {
	Since             int64
	Progress          int32
	ThresholdDuration int32
	ThresholdValue    int32
	Tag               string
	Codes              *regexp.Regexp
	ThresholdValueType string // count || percent
}

type HttpCodeCriterion struct {
	Since             int64
	Progress          int32
	ThresholdDuration int32
	ThresholdValue    int32
	Tag               string
	Codes              *regexp.Regexp
	ThresholdValueType string // count || percent
}

type QuantileCriterion struct {
	Since             int64
	Progress          int32
	ThresholdDuration int32
	ThresholdValue    int32
	Tag               string
	Quantile int32
}

type TimeLimitCriterion struct {
	Since             int64
	Progress          int32
	ThresholdDuration int32
	ThresholdValue    int32
	Tag               string
}

// expect something like "1s500ms, 30s" or "50,15"
func (c *AvgTimeCriterion) Parse(params string) error {
	var paramsList []string
	for _, p := range strings.Split(params, ",") {
		paramsList = append(paramsList, strings.TrimSpace(p))
	}
	if len(paramsList) < 2 || len(paramsList) > 3 {
		return errors.New("invalid autostop format: time autostop must have 2 or 3 params")
	}

	// ThresholdValue
	thresholdValue, err := ParseMilliseconds(paramsList[0])
	if err != nil {
		return err
	}
	c.ThresholdValue = thresholdValue

	// ThresholdDuration
	thresholdDuration, err := ParseDuration(paramsList[1])
	if err != nil {
		return err
	}
	c.ThresholdDuration = thresholdDuration

	// Tag
	if len(paramsList) == 3 {
		c.Tag = paramsList[2]
	}

	return nil
}

func (c AvgTimeCriterion) Check(second []string) bool {
	return false
}

// expect something like "404,10,15, sometag" or "5xx, 10%, 1m"
func (c *NetCodeCriterion) Parse(params string) error {
	var paramsList []string
	for _, p := range strings.Split(params, ",") {
		paramsList = append(paramsList, strings.TrimSpace(p))
	}
	if len(paramsList) < 3 || len(paramsList) > 4 {
		return errors.New("invalid net autostop format: net autostop must have 3 or 4 params")
	}

	// Codes
	ok, err := regexp.Match("^[0-9x]+$", []byte(paramsList[0]))
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("invalid net autostop net code format")
	}
	codes, err := regexp.Compile(fmt.Sprintf("^%s$", strings.Replace(paramsList[0], "x", "[0-9]", -1)))
	if err != nil {
		return err
	}
	c.Codes = codes

	// ThresholdValue and ThresholdValueType
	ok, err = regexp.Match("^[0-9]+%?$", []byte(paramsList[1]))
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("invalid net autostop threshold format")
	}
	thresholdValue, err := strconv.Atoi(strings.TrimRight(paramsList[1], "%"))
	if err != nil {
		return err
	}
	c.ThresholdValue = int32(thresholdValue)
	if paramsList[1][len(paramsList[1])-1] == '%' {
		c.ThresholdValueType = "percent"
	} else {
		c.ThresholdValueType = "count"
	}

	// ThresholdDuration
	thresholdDuration, err := ParseDuration(paramsList[2])
	if err != nil {
		return err
	}
	c.ThresholdDuration = thresholdDuration

	// Tag
	if len(paramsList) == 4 {
		c.Tag = paramsList[3]
	}

	return nil
}

func (c NetCodeCriterion) Check(second []string) bool {
	return false
}

// expect something like "404,10,15, sometag" or "5xx, 10%, 1m"
func (c *HttpCodeCriterion) Parse(params string) error {
	var paramsList []string
	for _, p := range strings.Split(params, ",") {
		paramsList = append(paramsList, strings.TrimSpace(p))
	}
	if len(paramsList) < 3 || len(paramsList) > 4 {
		return errors.New("invalid http autostop format: http autostop must have 3 or 4 params")
	}

	// Codes
	ok, err := regexp.Match("^[0-9x]+$", []byte(paramsList[0]))
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("invalid http autostop http code format")
	}
	codes, err := regexp.Compile(fmt.Sprintf("^%s$", strings.Replace(paramsList[0], "x", "[0-9]", -1)))
	if err != nil {
		return err
	}
	c.Codes = codes

	// ThresholdValue and ThresholdValueType
	ok, err = regexp.Match("^[0-9]+%?$", []byte(paramsList[1]))
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("invalid http autostop threshold format")
	}
	thresholdValue, err := strconv.Atoi(strings.TrimRight(paramsList[1], "%"))
	if err != nil {
		return err
	}
	c.ThresholdValue = int32(thresholdValue)
	if paramsList[1][len(paramsList[1])-1] == '%' {
		c.ThresholdValueType = "percent"
	} else {
		c.ThresholdValueType = "count"
	}

	// ThresholdDuration
	thresholdDuration, err := ParseDuration(paramsList[2])
	if err != nil {
		return err
	}
	c.ThresholdDuration = thresholdDuration

	// Tag
	if len(paramsList) == 4 {
		c.Tag = paramsList[3]
	}

	return nil
}

func (c HttpCodeCriterion) Check(second []string) bool {
	return false
}

// expect something like "95,100ms,10s"
func (c *QuantileCriterion) Parse(params string) error {
	var paramsList []string
	for _, p := range strings.Split(params, ",") {
		paramsList = append(paramsList, strings.TrimSpace(p))
	}
	if len(paramsList) < 3 || len(paramsList) > 4 {
		return errors.New("invalid quantile autostop format: quantile autostop must have 3 or 4 params")
	}

	// Quantile
	q, err := strconv.Atoi(paramsList[0])
	if err != nil {
		return err
	}
	if q < 0 || q > 100 {
		return errors.New("invalid quantile autostop: quantile must be between 0 and 100")
	}
	c.Quantile = int32(q)

	// ThresholdValue
	thresholdValue, err := ParseMilliseconds(paramsList[1])
	if err != nil {
		return err
	}
	c.ThresholdValue = thresholdValue

	// ThresholdDuration
	thresholdDuration, err := ParseDuration(paramsList[2])
	if err != nil {
		return err
	}
	c.ThresholdDuration = thresholdDuration

	// Tag
	if len(paramsList) == 4 {
		c.Tag = paramsList[3]
	}

	return nil
}

func (c QuantileCriterion) Check(second []string) bool {
	return false
}

func (c *TimeLimitCriterion) Parse(params string) error {
	limit, err := ParseDuration(strings.TrimSpace(params))
	if err != nil {
		return err
	}
	c.ThresholdValue = int32(limit)

	return nil
}

func (c TimeLimitCriterion) Check(second []string) bool {
	return false
}

