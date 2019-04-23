package aggregator

import "testing"

func BenchmarkAggregatorGota(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main()
	}
}
