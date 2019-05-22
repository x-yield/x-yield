package autostop

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ParseCriteria - receives list of autostop criteria from config
func ParseCriteria(criteria []string) (cc []*Criterion, err error) {

	for _, criterion := range criteria {
		ok, err := regexp.Match("^[a-z]+(.+)$", []byte(criterion))
		if err != nil {
			break
		}
		if !ok {
			err = errors.New("invalid autostop format")
			break
		}

		params := strings.Split(criterion, "(")
		c := strings.TrimSpace(params[0])
		p := strings.TrimRight(params[1], ")")

		switch c {
		case "time":
			var avgTimeCriterion Criterion
			avgTimeCriterion = &AvgTimeCriterion{}
			err = avgTimeCriterion.Parse(p)
			if err != nil {
				break
			}
			cc = append(cc, &avgTimeCriterion)
		case "net":
			var netCodeCriterion Criterion
			netCodeCriterion = &NetCodeCriterion{}
			err = netCodeCriterion.Parse(p)
			if err != nil {
				break
			}
			cc = append(cc, &netCodeCriterion)
		case "http":
			var httpCodeCriterion Criterion
			httpCodeCriterion = &HttpCodeCriterion{}
			err = httpCodeCriterion.Parse(p)
			if err != nil {
				break
			}
			cc = append(cc, &httpCodeCriterion)
		case "quantile":
			var quantileCriterion Criterion
			quantileCriterion = &QuantileCriterion{}
			err = quantileCriterion.Parse(p)
			if err != nil {
				break
			}
			cc = append(cc, &quantileCriterion)
		case "limit":
			var timeLimitCriterion Criterion
			timeLimitCriterion = &TimeLimitCriterion{}
			err = timeLimitCriterion.Parse(p)
			if err != nil {
				break
			}
			cc = append(cc, &timeLimitCriterion)
		default:
			err = errors.New(fmt.Sprintf("unknown autostop type %s", c))
		}
	}
	return
}
