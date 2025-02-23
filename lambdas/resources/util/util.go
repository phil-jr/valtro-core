package util

import (
	// "fmt"
	"resources/types"
	"time"
)

func EvenlyBucketMetrics(metrics []types.Metric, c int) ([]types.EvenMetric, error) {
	if len(metrics) == 0 {
		return nil, nil
	}

	unit := ""
	bucketDur := time.Duration(c) * time.Minute
	var buckets []types.EvenMetric

	uniqueMetricType := make(map[string][]types.Metric)
	for _, m := range metrics {
		uniqueMetricType[m.Name] = append(uniqueMetricType[m.Name], m)
	}

	for key, value := range uniqueMetricType {

		unit = value[0].Unit
		firstTime := value[0].Timestamp
		curStart := firstTime
		curEnd := curStart.Add(bucketDur)
		var accValue float64
		var accDuration time.Duration

		for _, m := range value {

			rowDuration := time.Duration(m.Aggregate) * time.Minute
			valuePerMinute := m.Value / float64(m.Aggregate)
			remaining := rowDuration

			for remaining > 0 {
				// Calculate how much time remains in the current bucket.
				timeLeftInBucket := curEnd.Sub(curStart) - accDuration

				if remaining <= timeLeftInBucket {
					// If the remaining row duration fits in the current bucket, assign it.
					accValue += valuePerMinute * remaining.Minutes()
					accDuration += remaining
					remaining = 0
				} else {
					// Fill the current bucket with as much as possible.
					accValue += valuePerMinute * timeLeftInBucket.Minutes()
					accDuration += timeLeftInBucket
					remaining -= timeLeftInBucket
				}

				// Finalize the current bucket if it's full.
				if accDuration >= bucketDur {
					buckets = append(buckets, types.EvenMetric{
						Name:  key,
						Start: curStart,
						End:   curEnd,
						Value: accValue,
						Unit:  unit,
					})
					// Move to the next bucket.
					curStart = curEnd
					curEnd = curStart.Add(bucketDur)
					accValue = 0
					accDuration = 0
				}
			}
		}

		if accDuration > 0 {
			buckets = append(buckets, types.EvenMetric{
				Name:  key,
				Start: curStart,
				End:   curStart.Add(accDuration),
				Value: accValue,
				Unit:  unit,
			})
		}
	}

	return buckets, nil
}
