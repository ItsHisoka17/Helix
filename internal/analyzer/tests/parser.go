package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/ItsHisoka17/Helix/internal/analyzer"
	"github.com/ItsHisoka17/Helix/internal/analyzer/tests/utils"
)

type Fixtures struct {
	Path         string   `json:"path"`
	Main         string   `json:"main"`
	Port         int      `json:"port"`
	Signal_types []string `json:"signal_types"`
}

func TestParsers() {
	path, err := os.Getwd()
	if err == nil {
		data, err1 := os.ReadFile(path + "/fixtures/fixtures.json")
		if err1 == nil {
			var fixtures Fixtures
			err2 := json.Unmarshal(data, &fixtures)
			if err2 == nil {
				fPath := path
				fPath += "/fixtures"
				pkg, parseError := analyzer.ParsePackageJson(fPath)
				if parseError != nil {
					utils.FormatError(1, path, parseError.Error())
				}
				if pkg.Main != fixtures.Main || len(pkg.Scripts) < 1 || len(pkg.Dependencies) < 1 {
					utils.FormatError(1, Tests[1], path, *pkg)
				} else {
					utils.FormatSuccess(1, Tests[1], path, *pkg)
				}
				signals, parseError1 := analyzer.ParseCodeContext(fPath)
				if parseError1 != nil {
					utils.FormatError(1, Tests[1], path, parseError1.Error())
				}
				var port int
				var signalMatch bool = false
				for _, signal := range signals {
					if signal.Port > 0 {
						port = signal.Port
					}
					if !signalMatch {
						if slices.Contains(fixtures.Signal_types, string(signal.Type)) {
							signalMatch = true
						}
					}
				}
				if port <= 0 {
					utils.FormatError(2, Tests[2], path, signals)
				} else {
					utils.FormatSuccess(2, Tests[2], path, port)
				}
				if !signalMatch {
					utils.FormatError(3, path, Tests[3], signals, "\n", signalMatch)
				}
				if signalMatch && port > 0 {
					utils.FormatSuccess(3, Tests[3], path, signals, signalMatch)
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
