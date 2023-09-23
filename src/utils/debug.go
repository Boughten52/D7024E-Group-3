package utils

import (
	"fmt"
)

var severityThreshold int = 10
var errorPrefix string = "ERROR: "

// Logs message if severity is above the threshold (10 = ERROR)
func Log(severity int, format string, a ...any) {
	if severity > severityThreshold {
		fmt.Printf(format, a...)
	}
}

func LogError(format string, a ...any) {
	if 10 > severityThreshold {
		fmt.Printf(errorPrefix+format, a...)
	}
}
