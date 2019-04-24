package aggregator

type PhoutEntry struct {
	// phout fields as is
	Time          float64
	Tag           string
	IntervalReal  int
	ConnectTime   int
	SendTime      int
	Latency       int
	ReceiveTime   int
	IntervalEvent int
	SizeOut       int
	SizeIn        int
	NetCode       int
	ProtoCode     int
	// custom fields
	ReceiveTs int // calculated field via Time + IntervalReal
}

type AggregatedSecond struct {
	Ts         int
	Overall    AggregatedChunk
	Tagged     map[string]AggregatedChunk
	CountedRps int
}

type AggregatedChunk struct {
	IntervalReal  *AggIntervalReal
	ConnectTime   *AggConnectTime
	SendTime      *AggSendTime
	Latency       *AggLatency
	ReceiveTime   *AggReceiveTime
	IntervalEvent *AggIntervalEvent
	SizeOut       *AggSizeOut
	SizeIn        *AggSizeIn
	NetCode       *AggNetCode
	ProtoCode     *AggProtoCode
}

type QuantilesPack struct {
	q50  float64
	q75  float64
	q80  float64
	q85  float64
	q90  float64
	q95  float64
	q98  float64
	q99  float64
	q100 float64
}

type Hist struct {
	dividers []float64
	x        []float64
}

type AggIntervalReal struct {
	Q     QuantilesPack
	Min   float64
	Max   float64
	Total float64
	Len   int
	Hist  Hist
}

type AggConnectTime struct {
	Min   float64
	Max   float64
	Total float64
	Len   int
}

type AggSendTime struct {
	Min   float64
	Max   float64
	Total float64
	Len   int
}

type AggLatency struct {
	Min   float64
	Max   float64
	Total float64
	Len   int
}

type AggReceiveTime struct {
	Min   float64
	Max   float64
	Total float64
	Len   int
}

type AggIntervalEvent struct {
	Min   float64
	Max   float64
	Total float64
	Len   int
}

type AggSizeOut struct {
	Min   float64
	Max   float64
	Total float64
	Len   int
}

type AggSizeIn struct {
	Min   float64
	Max   float64
	Total float64
	Len   int
}

type AggNetCode struct {
	Count int
}

type AggProtoCode struct {
	Count int
}
