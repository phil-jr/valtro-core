package util

import (
	// "fmt"
	"log"
	"resources/types"
	"time"
)

func EvenlyBucketMetrics(metrics []types.Metric, c int) ([]types.EvenMetric, error) {
	if len(metrics) == 0 {
		return nil, nil
	}

	var buckets []types.EvenMetric
	historicMinutes := 24 * 60
	numBuckets := historicMinutes / c
	startTime := time.Now().Add(time.Duration(-historicMinutes) * time.Minute)
	// evaluatedTime := startTime

	uniqueMetricType := make(map[string][]types.Metric)
	for _, m := range metrics {
		uniqueMetricType[m.Name] = append(uniqueMetricType[m.Name], m)
	}


	var timeArray []time.Time
	evaluatedTime := startTime
	for i:=0; i<=numBuckets; i++ {
		timeArray = append(timeArray, evaluatedTime)
		evaluatedTime = evaluatedTime.Add(time.Duration(c) * time.Minute)
	}
	log.Printf("timeMap: %v", timeArray)


	// for key, value := range uniqueMetricType {
	// 	for _, m := range value {
	// 		metricStartTime := m.Timestamp
	// 		bucketValue := (m.Value / float64(m.Aggregate)) * float64(c)



	// 		if (evaluatedTime.After(metricStartTime)|| evaluatedTime.Equal(metricStartTime)) &&
	// 		   (evaluatedTime.Before(metricEndTime) || evaluatedTime.Equal(metricEndTime)) {

	// 		}

	// 		// 	buckets = append(buckets, types.EvenMetric{
	// 		// 		Name:  key,
	// 		// 		Start: curStart,
	// 		// 		End:   curEnd,
	// 		// 		Value: accValue,
	// 		// 		Unit:  unit,
	// 		// 	})

	// 		// }
	// 	}
	// }


	return buckets, nil
}
