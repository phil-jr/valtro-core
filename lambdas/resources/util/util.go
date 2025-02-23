package util

import (
	// "fmt"
	// "log"
	"resources/types"
	"slices"
	"time"
)

func EvenlyBucketMetrics(metrics []types.Metric, c int) ([]types.EvenMetric, error) {
	if len(metrics) == 0 {
		return nil, nil
	}

	// Define the overall lookback period (here: 6 days)
	historicMinutes := 6 * 24 * 60
	numBuckets := historicMinutes / c
	bucketStartTime := time.Now().Add(-time.Duration(historicMinutes) * time.Minute)

	// Build a slice specifying bucket intervals.
	type bucketInterval struct {
		Start time.Time
		End   time.Time
	}
	bucketsIntervals := make([]bucketInterval, numBuckets)
	for i := 0; i < numBuckets; i++ {
		start := bucketStartTime.Add(time.Duration(i*c) * time.Minute)
		bucketsIntervals[i] = bucketInterval{
			Start: start,
			End:   start.Add(time.Duration(c) * time.Minute),
		}
	}

	// bucketData will hold our accumulator.
	// For average metrics we accumulate:
	//   - weighted sum  (Sum)
	//   - total contributing minutes (Weight)
	// For sum metrics, Weight is not used.
	type bucketData struct {
		Sum       float64 // Sum for sum metrics OR weighted sum for avg metrics.
		Weight    float64 // Total minutes contributing (for avg metrics)
		Unit      string  // The metric unit.
		IsAverage bool    // true if metric name has "avg" in it.
	}

	// We'll use a map keyed by metric name and bucket start time.
	// Structure: metricName -> bucket start -> bucketData
	bucketMap := make(map[string]map[time.Time]*bucketData)

	// Helper to create or get the accumulator for a given metric and bucket.
	getBucketData := func(metricName string, t time.Time) *bucketData {
		if bucketMap[metricName] == nil {
			bucketMap[metricName] = make(map[time.Time]*bucketData)
		}
		if bucketMap[metricName][t] == nil {
			// flag as average if the metric name is in the averagedMetrics list.
			averagedMetrics := []string{"duration", "avg_memory"}
			isAvg := slices.Contains(averagedMetrics, metricName)
			bucketMap[metricName][t] = &bucketData{
				Sum:       0,
				Weight:    0,
				Unit:      "",
				IsAverage: isAvg,
			}
		}
		return bucketMap[metricName][t]
	}

	// Process each metric and distribute its data across the buckets.
	for _, m := range metrics {
		metricStart := m.Timestamp
		metricEnd := m.Timestamp.Add(time.Duration(m.Aggregate) * time.Minute)

		// Iterate over each bucket interval to check for overlap.
		// (If performance becomes an issue, consider optimizing the lookup.)
		for _, bucket := range bucketsIntervals {
			// Determine the overlap between the metric window and the bucket.
			overlapStart := maxTime(metricStart, bucket.Start)
			overlapEnd := minTime(metricEnd, bucket.End)
			// If there is no overlap, skip.
			if !overlapStart.Before(overlapEnd) {
				continue
			}
			// Get the overlapping duration in minutes.
			overlapMinutes := overlapEnd.Sub(overlapStart).Minutes()

			// Get the bucket accumulator.
			bd := getBucketData(m.Name, bucket.Start)
			// Assume all values of a metric share the same unit.
			bd.Unit = m.Unit

			if !bd.IsAverage { // For sum metrics.
				// Assume m.Value represents the total over m.Aggregate minutes.
				perMinuteVal := m.Value / float64(m.Aggregate)
				bd.Sum += perMinuteVal * overlapMinutes
			} else { // For average metrics.
				// m.Value represents an average over m.Aggregate minutes.
				if m.Value != 0 {
					bd.Sum += m.Value * overlapMinutes
					bd.Weight += overlapMinutes
				}
			}
		}
	}

	// Build a list of unique metric names from the bucketMap.
	metricNames := make([]string, 0, len(bucketMap))
	for metricName := range bucketMap {
		metricNames = append(metricNames, metricName)
	}

	// Convert the accumulator into a slice of EvenMetric, while ensuring that we
	// output a bucket for every metric name even if no data was recorded in that interval.
	var result []types.EvenMetric
	for _, metricName := range metricNames {
		// For each bucket interval...
		for _, bucket := range bucketsIntervals {
			// Look for an existing entry for the current metric and bucket.
			bd, exists := bucketMap[metricName][bucket.Start]
			var bucketValue float64
			unit := ""
			if !exists {
				// Fill gap with a zero (or default) value.
				bucketValue = 0
			} else {
				if bd.IsAverage && bd.Weight > 0 {
					bucketValue = bd.Sum / bd.Weight
				} else {
					bucketValue = bd.Sum
				}
				unit = bd.Unit
			}

			result = append(result, types.EvenMetric{
				Name:  metricName,
				Start: bucket.Start,
				End:   bucket.End,
				Value: bucketValue,
				Unit:  unit,
			})
		}
	}

	return result, nil
}

// maxTime returns the later of two time.Time values.
func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// minTime returns the earlier of two time.Time values.
func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
