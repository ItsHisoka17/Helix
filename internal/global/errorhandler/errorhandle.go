package errorhandler

import (
	"fmt"
	"runtime"
)

func HandleError(message string, errors ...error) (ErrorFound bool) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "Unknown"
	}
	if len(message) == 0 {
		message = "Unknown Error Occured"
	}
	for _, err := range errors {
		if err != nil {
			ErrorFound = true
			message = fmt.Sprintf("%s | %s - On line: %d\n  - %s", file, message, line, err.Error())

			fmt.Printf("%s%s%s", "\033[31m", message, "\033[0m")
		}
	}
	return ErrorFound
}
