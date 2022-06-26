package logging

import (
	"math"

	"k8s.io/utils/strings"
)

func SnipTo(data string, maxLen float64) string {
	suffix := ""
	dataLen := float64(len(data))
	if dataLen > maxLen {
		suffix = "..."
	}
	return strings.ShortenString(data, int(math.Min(dataLen, maxLen))) + suffix
}

func Ends(data string, partMaxLen float64) string {
	dataLen := float64(len(data))
	if dataLen <= partMaxLen*2 {
		return data
	}
	return SnipTo(data, partMaxLen) + data[int(dataLen-partMaxLen):]
}
