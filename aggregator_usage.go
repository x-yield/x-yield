package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"x_yield/aggregator"
)


const (
	maximumQueueLength = 3
)

func Aggregate(key int, pendingData map[int][]aggregator.PhoutEntry) (result aggregator.AggregatedSecond) {
	overallReadyToGo := pendingData[key]

	taggedMap := make(map[string]aggregator.AggregatedChunk)
	// TODO(netort): goroutines
	result = aggregator.AggregatedSecond{
		Ts:     key,
		Tagged: taggedMap,
	}

	// tagged aggregates
	pendingTagged := make(map[string][]aggregator.PhoutEntry)

	for _, taggedEntry := range overallReadyToGo {
		pendingTagged[taggedEntry.Tag] = append(pendingTagged[taggedEntry.Tag], taggedEntry)
	}

	for _, taggedReadyToGo := range pendingTagged {
		result.Tagged[taggedReadyToGo[0].Tag] = aggregator.NewAggregatedChunk(taggedReadyToGo)
	}

	// overall aggregated
	result.Overall = aggregator.NewAggregatedChunk(overallReadyToGo)
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

	pendingData := make(map[int][]aggregator.PhoutEntry)
	var pendingAggregates []aggregator.AggregatedSecond

	for scanner.Scan() {
		chunk := scanner.Text()

		entry := aggregator.NewPhoutEntry(&chunk)
		// TODO(netort): implement groupby into go-gota dataframe
		// groupby receive timestamp
		pendingData[entry.ReceiveTs] = append(pendingData[entry.ReceiveTs], entry)

		// if length of prepared data > maximumQueueLength then pop the ready-to-go values out
		if len(pendingData) > maximumQueueLength {
			// get ready-to-go (minimum ts) key
			key := aggregator.GetOldestDataChunk(pendingData)

			aggregated := Aggregate(key, pendingData)
			delete(pendingData, key)
			pendingAggregates = append(pendingAggregates, aggregated)
		}
	}

	// process the rest of data
	for key, _ := range pendingData {
		aggregated := Aggregate(key, pendingData)
		delete(pendingData, key)
		pendingAggregates = append(pendingAggregates, aggregated)
	}

	log.Printf("Done its work w/ finished seconds: %v", len(pendingAggregates))
}
