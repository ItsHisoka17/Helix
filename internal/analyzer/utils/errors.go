package utils

import "fmt"

func HandleError(errors ...error) (found bool) {
	for _, err := range errors {
		if err != nil {
			e := fmt.Errorf("%s", err.Error())
			if e != nil {
				fmt.Print(e)
			}
			if !found {
				found = true
			}
		}
	}
	return found
}
