package aggregator

import (
	"github.com/go-gota/gota/dataframe"
	"github.com/gonum/floats"
	"github.com/gonum/stat"
	"github.com/kniren/gota/series"
)

func NewAggIntervalReal(s series.Series, histDividers []float64) *AggIntervalReal {
	orderedSeries := SortFloatSeries(s)
	return &AggIntervalReal{
		*QuantilesPackForOrderedSeries(orderedSeries),
		s.Min(),
		s.Max(),
		floats.Sum(s.Float()),
		s.Len(),
		// TODO(netort): stat.Histogram has O(N^2), optimize search algo
		Hist{
			histDividers,
			stat.Histogram(nil, histDividers, orderedSeries, nil),
		},
	}
}

func NewAggConnectTime(s series.Series) *AggConnectTime {
	return &AggConnectTime{
		s.Min(),
		s.Max(),
		floats.Sum(s.Float()),
		s.Len(),
	}
}

func NewAggSendTime(s series.Series) *AggSendTime {
	return &AggSendTime{
		s.Min(),
		s.Max(),
		floats.Sum(s.Float()),
		s.Len(),
	}
}

func NewAggLatency(s series.Series) *AggLatency {
	return &AggLatency{
		s.Min(),
		s.Max(),
		floats.Sum(s.Float()),
		s.Len(),
	}
}

func NewAggReceiveTime(s series.Series) *AggReceiveTime {
	return &AggReceiveTime{
		s.Min(),
		s.Max(),
		floats.Sum(s.Float()),
		s.Len(),
	}
}

func NewAggIntervalEvent(s series.Series) *AggIntervalEvent {
	return &AggIntervalEvent{
		s.Min(),
		s.Max(),
		floats.Sum(s.Float()),
		s.Len(),
	}
}

func NewAggSizeOut(s series.Series) *AggSizeOut {
	return &AggSizeOut{
		s.Min(),
		s.Max(),
		floats.Sum(s.Float()),
		s.Len(),
	}
}

func NewAggSizeIn(s series.Series) *AggSizeIn {
	return &AggSizeIn{
		s.Min(),
		s.Max(),
		floats.Sum(s.Float()),
		s.Len(),
	}
}

func NewAggNetCode(s series.Series) *AggNetCode {
	return &AggNetCode{
		s.Len(),
	}
}

func NewAggProtoCode(s series.Series) *AggProtoCode {
	return &AggProtoCode{
		s.Len(),
	}
}

func HistDividers() (dividers []float64) {
	// 0 <= x < 5ms | 10µs accuracy
	for i := 0.0; i <= 4.9*1e3; i = i + 10 {
		dividers = append(dividers, i)
	}
	// 5ms <= x < 10ms | 100µs accuracy
	for i := 5.0 * 1e3; i <= 9.9*1e3; i = i + 100 {
		dividers = append(dividers, i)
	}
	// 10ms <= x < 500ms | 1ms accuracy
	for i := 10 * 1e3; i <= 499*1e3; i = i + 1*1e3 {
		dividers = append(dividers, i)
	}
	// 0.5s <= x < 2.995s | 5ms accuracy
	for i := 0.5 * 1e6; i <= 2.995*1e6; i = i + 5*1e3 {
		dividers = append(dividers, i)
	}
	// 3s <= x < 9.99s | 10ms accuracy
	for i := 3 * 1e6; i <= 9.999*1e6; i = i + 10*1e3 {
		dividers = append(dividers, i)
	}
	// 10s <= x < 29.95s | 50ms accuracy
	for i := 10 * 1e6; i <= 29.95*1e6; i = i + 50*1e3 {
		dividers = append(dividers, i)
	}
	// 30s <= x < 119.9s | 100ms accuracy
	for i := 30 * 1e6; i <= 119.9*1e6; i = i + 100*1e3 {
		dividers = append(dividers, i)
	}
	// 120s <= x < 300s | 1s accuracy
	for i := 120 * 1e6; i <= 300*1e6; i = i + 1*1e6 {
		dividers = append(dividers, i)
	}
	return
}

func NewAggregatedChunk(entries []PhoutEntry) (aggChunk AggregatedChunk) {
	// create dataframe, clean up the pending slice
	df := dataframe.LoadStructs(entries)

	// create series
	seriesIntervalReal := df.Col("IntervalReal")
	seriesConnectTime := df.Col("ConnectTime")
	seriesSendTime := df.Col("SendTime")
	seriesLatency := df.Col("Latency")
	seriesReceiveTime := df.Col("ReceiveTime")
	seriesIntervalEvent := df.Col("IntervalEvent")
	seriesSizeOut := df.Col("SizeOut")
	seriesSizeIn := df.Col("SizeIn")
	seriesNetCode := df.Col("NetCode")
	seriesProtoCode := df.Col("ProtoCode")

	// calculate aggregates
	aggChunk = AggregatedChunk{
		NewAggIntervalReal(seriesIntervalReal, histDividers),
		NewAggConnectTime(seriesConnectTime),
		NewAggSendTime(seriesSendTime),
		NewAggLatency(seriesLatency),
		NewAggReceiveTime(seriesReceiveTime),
		NewAggIntervalEvent(seriesIntervalEvent),
		NewAggSizeOut(seriesSizeOut),
		NewAggSizeIn(seriesSizeIn),
		NewAggNetCode(seriesNetCode),
		NewAggProtoCode(seriesProtoCode),
	}
	return
}
