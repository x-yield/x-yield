package aggregator

import (
	"github.com/gonum/stat"
	"github.com/kniren/gota/series"
)

func QuantilesPackForOrderedSeries(orderedSeries []float64) *QuantilesPack {
	return &QuantilesPack{
		stat.Quantile(.50, stat.Empirical, orderedSeries, nil),
		stat.Quantile(.75, stat.Empirical, orderedSeries, nil),
		stat.Quantile(.80, stat.Empirical, orderedSeries, nil),
		stat.Quantile(.85, stat.Empirical, orderedSeries, nil),
		stat.Quantile(.90, stat.Empirical, orderedSeries, nil),
		stat.Quantile(.95, stat.Empirical, orderedSeries, nil),
		stat.Quantile(.98, stat.Empirical, orderedSeries, nil),
		stat.Quantile(.99, stat.Empirical, orderedSeries, nil),
		stat.Quantile(1, stat.Empirical, orderedSeries, nil),
	}
}

func SortFloatSeries(s series.Series) []float64 {
	return s.Subset(s.Order(false)).Float()
}
