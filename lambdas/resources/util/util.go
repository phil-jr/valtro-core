package util

import (
	"fmt"
	"resources/types"
	"time"
)

// evenlyBucketMetrics takes the raw metrics (assumed to be sorted by Timestamp ascending)
// and produces evenly spaced buckets of length "c" minutes using the provided Aggregate field.
func EvenlyBucketMetrics(metrics []types.Metric, c int) ([]types.EvenMetric, error) {
	if len(metrics) == 0 {
		return nil, nil
	}

	// desired bucket duration in minutes
	bucketDur := time.Duration(c) * time.Minute

	// Parse the timestamp of the first metric to initialize our bucket timeline.
	firstTime, err := time.Parse(time.RFC3339, metrics[0].Timestamp)
	if err != nil {
		return nil, fmt.Errorf("cannot parse timestamp %q: %v", metrics[0].Timestamp, err)
	}

	var buckets []types.EvenMetric
	curStart := firstTime
	curEnd := curStart.Add(bucketDur)
	var accValue float64
	var accDuration time.Duration

	// Assume all metrics share the same unit.
	unit := metrics[0].Unit

	// Loop through each metric.
	for _, m := range metrics {
		// Convert the provided Aggregate minutes into a duration.
		rowDuration := time.Duration(m.Aggregate) * time.Minute

		// Compute the distributed value per minute.
		valuePerMinute := m.Value / float64(m.Aggregate)

		// Distribute the metric's value over the rowDuration.
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

	// If a partial bucket remains, add it.
	if accDuration > 0 {
		buckets = append(buckets, types.EvenMetric{
			Start: curStart,
			End:   curStart.Add(accDuration),
			Value: accValue,
			Unit:  unit,
		})
	}

	return buckets, nil
}
