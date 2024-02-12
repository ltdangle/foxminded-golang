package main

import "time"

const layoutIn = "3:04:05PM"
const layoutOut = "15:04:05"

func convertTime(timeStr string) (string, error) {
	t, err := time.Parse(layoutIn, timeStr)
	if err != nil {
		return "", err
	}
	return t.Format(layoutOut), nil
}
