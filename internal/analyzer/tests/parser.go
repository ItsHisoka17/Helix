package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/ItsHisoka17/Helix/internal/analyzer"
)

type Fixtures struct {
	path         string
	main         string
	port         int
	signal_types []string
}

func TestParsers() {
	path, err := os.Getwd()
	if err == nil {
		data, err1 := os.ReadFile(path + "/internal/analyzer/tests/fixtures/fixtures.json")
		if err1 == nil {
			var fixtures Fixtures
			err2 := json.Unmarshal(data, &fixtures)
			if err2 == nil {
				fPath := path
				fPath += "/internal/analyzer/tests/fixtures"
				pkg, parseError := analyzer.ParsePackageJson(fPath)
				if parseError != nil {
					e := fmt.Errorf("--- Test#1 Failed\nTest Path: %s\n%s\n", path, parseError.Error())
					if e != nil {
						fmt.Print(e.Error())
					}
				}
				if pkg.Main != fixtures.main || len(pkg.Scripts) < 1 || len(pkg.Dependencies) < 1 {
					e0 := fmt.Errorf("--- Test#1 Failed\nTest Path: %s\nTest Result:\n%+v\n", path, pkg)
					if e0 != nil {
						fmt.Print(e0.Error())
					}
				} else {
					fmt.Printf("--- Test#1 Passed\nTest Path: %s\nTest Result:\n%+v\n", path, pkg)
				}
				signals, parseError1 := analyzer.ParseCodeContext(fPath)
				if parseError1 != nil {
					e1 := fmt.Errorf("--- Test#2 Failed\nTest Path: %s\n%s\n", path, parseError1.Error())
					if e1 != nil {
						fmt.Print(e1)
					}
				}
				var port int
				var signalMatch bool = false
				for _, signal := range signals {
					if signal.Port > 0 {
						port = signal.Port
					}
					if !signalMatch {
						if slices.Contains(fixtures.signal_types, string(signal.Type)) {
							signalMatch = true
						}
					}
				}
				if port <= 0 {
					e2 := fmt.Errorf("--- Test#2 Failed\nTest Path: %s\nPort not found\nTest Result:\n%+v\n", path, signals)
					if e2 != nil {
						fmt.Print(e2.Error())
					}
				}
				if !signalMatch {
					e3 := fmt.Errorf("--- Test#3 Failed\nTest Path: %s\nSignal not found\nTest Result:\n%+v\n%+v\n", path, signals, signalMatch)
					if e3 != nil {
						fmt.Print(e3)
					}
				}
			} else {
				fmt.Print(err2.Error())
			}
		} else {
			fmt.Print(err1.Error())
		}
	} else {
		fmt.Print(err.Error())
	}
}

func main() {
	TestParsers()
}
