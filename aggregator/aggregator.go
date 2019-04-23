package aggregator

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kniren/gota/series"

)

const (
	maximumQueueLength = 3
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

func main() {
	file, err := os.Open("/Users/ttorubarov/phout-site.log")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	pendingData := make(map[int][]PhoutEntry)

	var pendingAggregates []AggregatedSecond
	for scanner.Scan() {
		chunk := scanner.Text()
		entry := NewPhoutEntry(&chunk)

		// TODO(netort): implement groupby into go-gota dataframe
		// groupby receive timestamp
		pendingData[entry.ReceiveTs] = append(pendingData[entry.ReceiveTs], entry)

		// if length of prepared data > maximumQueueLength then pop the ready-to-go values out
		if len(pendingData) > maximumQueueLength {
			// get ready-to-go (minimum ts) key
			key := GetOldestDataChunk(pendingData)

			overallReadyToGo := pendingData[key]

			taggedMap := make(map[string]AggregatedChunk)
			// TODO(netort): goroutines
			aggregate := AggregatedSecond{
				Ts:     key,
				Tagged: taggedMap,
			}
			// tagged aggregates
			pendingTagged := make(map[string][]PhoutEntry)

			for _, taggedEntry := range overallReadyToGo {
				pendingTagged[taggedEntry.Tag] = append(pendingTagged[taggedEntry.Tag], taggedEntry)
			}

			for _, taggedReadyToGo := range pendingTagged {
				aggregate.Tagged[taggedReadyToGo[0].Tag] = NewAggregatedChunk(taggedReadyToGo)
			}

			// overall aggregated
			aggregate.Overall = NewAggregatedChunk(overallReadyToGo)

			delete(pendingData, key)

			pendingAggregates = append(pendingAggregates, aggregate)
		}
	}
	log.Printf("Done its work w/ finished seconds: %v", len(pendingAggregates))
}
