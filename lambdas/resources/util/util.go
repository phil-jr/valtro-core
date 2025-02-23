package util

import (
	// "fmt"
	// "log"
	"resources/types"
	"strings"
	"time"
)

// func EvenlyBucketMetrics(metrics []types.Metric, c int) ([]types.EvenMetric, error) {
// 	if len(metrics) == 0 {
// 		return nil, nil
// 	}

// 	var buckets []types.EvenMetric
// 	historicMinutes := 24 * 60
// 	numBuckets := historicMinutes / c
// 	startTime := time.Now().Add(time.Duration(-historicMinutes) * time.Minute)
// 	evaluatedTime := startTime

// 	uniqueMetricType := make(map[string][]types.Metric)
// 	for _, m := range metrics {
// 		uniqueMetricType[m.Name] = append(uniqueMetricType[m.Name], m)
// 	}

// 	var timeArray []time.Time
// 	for i := 0; i <= numBuckets; i++ {
// 		timeArray = append(timeArray, evaluatedTime)
// 		evaluatedTime = evaluatedTime.Add(time.Duration(c) * time.Minute)
// 	}
// 	log.Printf("timeMap: %v", timeArray)

// 	for key, value := range uniqueMetricType {
// 		for _, m := range value {
// 			metricStartTime := m.Timestamp
// 			bucketValue := (m.Value / float64(m.Aggregate)) * float64(c)

// 			for index, _ := range timeArray {

// 				if index+1 == len(timeArray) {
// 					continue
// 				}

// 				if (metricStartTime.After(timeArray[index]) || metricStartTime.Equal(timeArray[index])) &&
// 					(metricStartTime.Before(timeArray[index+1]) || metricStartTime.Equal(timeArray[index+1])) {
// 					buckets = append(buckets, types.EvenMetric{
// 						Name:  key,
// 						Start: timeArray[index],
// 						End:   timeArray[index+1],
// 						Value: bucketValue,
// 						Unit:  m.Unit,
// 					})
// 				}

// 				// buckets = append(buckets, types.EvenMetric{
// 				// 	Name:  key,
// 				// 	Start: curStart,
// 				// 	End:   curEnd,
// 				// 	Value: accValue,
// 				// 	Unit:  unit,
// 				// })

// 			}

// 		}
// 	}

// 	log.Printf("buckets: %v", buckets)

// 	return buckets, nil
// }


func EvenlyBucketMetrics(metrics []types.Metric, c int) ([]types.EvenMetric, error) {
	if len(metrics) == 0 {
		return nil, nil
	}

	// Define the overall lookback period (here: 24 hours)
	historicMinutes := 24 * 60
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
		Sum       float64  // Sum for sum metrics OR weighted sum for avg metrics.
		Weight    float64  // Total minutes contributing (for avg metrics)
		Unit      string   // The metric unit.
		IsAverage bool     // true if metric name has "avg" in it.
	}

	// We'll use a map keyed by metric name and bucket start time.
	// Structure:  metricName -> bucket start -> bucketData
	bucketMap := make(map[string]map[time.Time]*bucketData)

	// Helper to create or get the accumulator for a given metric and bucket.
	getBucketData := func(metricName string, t time.Time) *bucketData {
		if bucketMap[metricName] == nil {
			bucketMap[metricName] = make(map[time.Time]*bucketData)
		}
		if bucketMap[metricName][t] == nil {
			// flag as average if the name contains "avg" (case-insensitive).
			isAvg := strings.Contains(strings.ToLower(metricName), "avg")
			bucketMap[metricName][t] = &bucketData{
				Sum:       0,
				Weight:    0,
				Unit:      "",
				IsAverage: isAvg,
			}
		}
		return bucketMap[metricName][t]
	}

	// Process each metric and distribute its data across the buckets
	for _, m := range metrics {
		metricStart := m.Timestamp
		metricEnd := m.Timestamp.Add(time.Duration(m.Aggregate) * time.Minute)

		// Iterate over each bucket interval to check for overlap.
		// (If performance becomes an issue, consider searching only for buckets
		// that might overlap the metric.)
		for _, bucket := range bucketsIntervals {
			// Determine the overlap between the metric's time window and bucket.
			overlapStart := maxTime(metricStart, bucket.Start)
			overlapEnd := minTime(metricEnd, bucket.End)

			// If there is no overlap, skip.
			if !overlapStart.Before(overlapEnd) {
				continue
			}

			// Get the overlapping duration in minutes.
			overlapMinutes := overlapEnd.Sub(overlapStart).Minutes()

			// Initialize or get the corresponding bucket accumulator.
			bd := getBucketData(m.Name, bucket.Start)
			bd.Unit = m.Unit // It is assumed all values share the same unit.

			if !bd.IsAverage { // Sum metric.
				// Assume m.Value represents a total over m.Aggregate minutes.
				perMinuteVal := m.Value / float64(m.Aggregate)
				bd.Sum += perMinuteVal * overlapMinutes
			} else { // Average metric.
				// For avg, each metric m.Value is an average over (m.Aggregate) minutes,
				// so weighted contribution is m.Value * overlapMinutes.
				bd.Sum += m.Value * overlapMinutes
				bd.Weight += overlapMinutes
			}
		}
	}

	// Convert the map into a slice of EvenMetric.
	var result []types.EvenMetric
	for metricName, buckMap := range bucketMap {
		for bucketStart, bd := range buckMap {
			bucketEnd := bucketStart.Add(time.Duration(c) * time.Minute)
			bucketValue := bd.Sum
			if bd.IsAverage && bd.Weight > 0 {
				bucketValue = bd.Sum / bd.Weight
			}
			result = append(result, types.EvenMetric{
				Name:  metricName,
				Start: bucketStart,
				End:   bucketEnd,
				Value: bucketValue,
				Unit:  bd.Unit,
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
