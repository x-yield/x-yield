package aggregator

import (
	"strconv"
	"strings"

	"github.com/kniren/gota/series"
)

var (
	histDividers = HistDividers()
)

func NewPhoutEntry(chunk *string) (entry PhoutEntry) {
	// TODO(netort): optimize casting/parsing
	records := strings.SplitN(*chunk, "\t", 12)
	entry.Time, _ = strconv.ParseFloat(records[0], 10)
	entry.Tag = records[1]
	entry.IntervalReal, _ = strconv.Atoi(records[2])
	entry.ConnectTime, _ = strconv.Atoi(records[3])
	entry.SendTime, _ = strconv.Atoi(records[4])
	entry.Latency, _ = strconv.Atoi(records[5])
	entry.ReceiveTime, _ = strconv.Atoi(records[6])
	entry.IntervalEvent, _ = strconv.Atoi(records[7])
	entry.SizeOut, _ = strconv.Atoi(records[8])
	entry.SizeIn, _ = strconv.Atoi(records[9])
	entry.NetCode, _ = strconv.Atoi(records[10])
	entry.ProtoCode, _ = strconv.Atoi(records[11])
	entry.ReceiveTs = int(entry.Time + float64(entry.IntervalReal)/1e6)
	return
}

func GetOldestDataChunk(pendingData map[int][]PhoutEntry) (minKey int) {
	keys := make([]int, 0, len(pendingData))
	for k := range pendingData {
		keys = append(keys, k)
	}
	minKey = int(series.Ints(keys).Min())
	return
}