package testUtils

import "fmt"

func FormatError(n int, t string, p string, result ...any) {
	err := fmt.Errorf("%s--- Test#%d [%s] Failed\n -- Test Path: %s\n -- Test Results:\n%+v\n%s\n", "\033[31m", n, t, p, result, "\033[0m")
	if err != nil {
		fmt.Print(err.Error())
	}
}

func FormatSuccess(n int, t string, p string, result ...any) {
	fmt.Printf("%s--- Test#%d [%s] Passed\n -- Test Path: %s\n -- Test Result:\n%+v\n%s\n", "\033[32m", n, t, p, result, "\033[0m")
}

func JoinPath(paths ...string) (path string) {
	for i, p := range paths {
		if i == len(paths) {
			path += p
			continue
		}
		path += ("/" + p)
	}
	return path
}
