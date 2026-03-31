package analyzer

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/ItsHisoka17/Helix/internal/analyzer/types"
	"github.com/ItsHisoka17/Helix/internal/analyzer/utils"
)

func Analyze(projectPath string) (*types.AnalysisResult, error) {
	pkg, err := ParsePackageJson(projectPath)
	if err != nil {
		return nil, err
	}
	frameworks, dbs, redis := DetectDependencies(pkg)
	codeSignals, err := ParseCodeContext(projectPath)
	if err != nil {
		return nil, err
	}
	var ports []int
	for _, signal := range codeSignals.RawSignals {
		if signal.Port != 0 {
			ports = append(ports, signal.Port)
		}
	}
	return &types.AnalysisResult{
		Framework:  frameworks,
		Main:       pkg.Main,
		Scripts:    pkg.Scripts,
		Databases:  dbs,
		Redis:      redis,
		RawSignals: codeSignals.RawSignals,
		Ports:      ports,
		Confirm:    ConfirmFunc,
	}, nil
}

func ValidateAnalysis(Analysis *types.AnalysisResult) (*types.AnalysisResult, error) {
	signals := utils.MergeDuplicateSignals(Analysis.RawSignals)
	signals = utils.SortSignals(Analysis.RawSignals)
	scriptFileRgxCompiled := regexp.MustCompile(`([A-Za-z]|[0-9]|-|_)+\.(js|ts)`)
	var SignalFileNames []string
	var PortSignals []types.Signal
	for _, signal := range signals {
		SignalFileNames = append(SignalFileNames, signal.File)
	}
	if len(Analysis.Main) > 0 {
		ind, f := slices.BinarySearch(SignalFileNames, Analysis.Main)
		if f {
			moved, _ := utils.MoveSignal(signals, signals[ind], ind)
			if moved != nil {
				signals = moved
			}
		}
	}
	if Analysis.Scripts != nil {
		var mainCmds [2]string = [2]string{"run", "start"}
		for _, cmd := range mainCmds {
			if len(Analysis.Scripts[cmd]) > 0 {
				matches := scriptFileRgxCompiled.FindStringSubmatch(Analysis.Scripts[cmd])
				if len(matches) > 0 && len(matches[0]) > 0 {
					ind, f := slices.BinarySearch(SignalFileNames, matches[0])
					if f {
						moved, _ := utils.MoveSignal(signals, signals[ind], ind)
						if moved != nil {
							signals = moved
						}
					}
				}
			}
		}
	}
	for _, signal := range signals {
		if signal.Port > 0 {
			PortSignals = append(PortSignals, signal)
		}
	}
	if len(PortSignals) > 1 {
		message, _ := Analysis.Confirm(PortSignals)
		var selectedIndex int
		fmt.Print(message, "\n> ")
		_, err := fmt.Scanln(&selectedIndex)
		if err != nil {
			return nil, err
		}
	}

	return &types.AnalysisResult{
		RawSignals: signals,
	}, nil
}

func ConfirmFunc(signals []types.Signal) (string, error) {

	var message strings.Builder
	message.WriteString("Multiple signals detected | Confirmation required\nSignals:\n")
	for _, signal := range signals {
		format := "\n " + strconv.Itoa(signal.Rank) + " File: " + string(signal.File) + "| " + "Signal: " + string(signal.Type)
		message.WriteString(format)
	}
	message.WriteString("\nInput the corresponding index for the correct file")
	return message.String(), nil
}
