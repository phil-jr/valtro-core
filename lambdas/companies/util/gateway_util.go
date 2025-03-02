package util

import (
	"fmt"
	"time"
)

func ParseQueryTime(params map[string]string, key string, defaultTime time.Time) (time.Time, error) {
	if timeStr, ok := params[key]; ok && timeStr != "" {
		t, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return defaultTime, fmt.Errorf("invalid %s format: %w", key, err)
		}
		return t, nil
	}
	return defaultTime, nil
}

func GetMapValue(params map[string]string, key string) (string, error) {
	if v, ok := params[key]; ok && v != "" {
		return v, nil
	}
	return "", fmt.Errorf("missing path parameter: %s", key)
}
