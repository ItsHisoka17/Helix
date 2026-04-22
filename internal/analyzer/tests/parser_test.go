package tests

import (
	"encoding/json"
	"os"
	"slices"
	"testing"

	"github.com/ItsHisoka17/Helix/internal/analyzer"
	testUtils "github.com/ItsHisoka17/Helix/internal/analyzer/tests/utils"
	"github.com/ItsHisoka17/Helix/internal/analyzer/utils"
	globalUtils "github.com/ItsHisoka17/Helix/internal/global/utils"
)

func TestParser(t *testing.T) {
	path, err := os.Getwd()
	if utils.HandleError(err) {
		return
	}
	data, err1 := os.ReadFile(globalUtils.JoinPath(path, FileMap["dir"], FileMap["json"]))
	if utils.HandleError(err1) {
		return
	}
	var fixtures Fixtures
	err2 := json.Unmarshal(data, &fixtures)
	if utils.HandleError(err2) {
		return
	}
	pathF := globalUtils.JoinPath(path, FileMap["dir"])
	pkg, parseError := analyzer.ParsePackageJson(pathF)
	if parseError != nil {
		testUtils.FormatError(t, 1, Tests[1], path, parseError.Error())
	}
	if pkg.Main != fixtures.Main || len(pkg.Scripts) < 1 || len(pkg.Dependencies) < 1 {

		testUtils.FormatError(t, 1, Tests[1], path, *pkg)
	} else {
		testUtils.FormatSuccess(1, Tests[1], path, *pkg)
	}
	signals, parseError1 := analyzer.ParseCodeContext(pathF)
	if parseError1 != nil {
		testUtils.FormatError(t, 1, Tests[1], path, parseError1.Error())
	}
	var ports []int
	var signalMatch bool = false
	for _, signal := range signals {
		if signal.Port > 0 {
			ports = append(ports, signal.Port)
		}
		if !signalMatch {
			if slices.Contains(fixtures.Signal_types, string(signal.Type)) {
				signalMatch = true
			}
		}
	}
	if len(ports) <= 0 {
		testUtils.FormatError(t, 2, Tests[2], path, signals)
	} else {
		testUtils.FormatSuccess(2, Tests[2], path, ports)
	}
	if !signalMatch {
		testUtils.FormatError(t, 3, path, Tests[3], signals, "\n", signalMatch)
	}
	if signalMatch && len(ports) > 0 {
		testUtils.FormatSuccess(3, Tests[3], path, signals, signalMatch)
	}
}
