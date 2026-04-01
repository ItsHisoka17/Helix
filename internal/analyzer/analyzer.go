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
	RawSignals, err := ParseCodeContext(projectPath)
	if err != nil {
		return nil, err
	}
	var ports []int
	for _, signal := range RawSignals {
		if signal.Port != 0 {
			ports = append(ports, signal.Port)
		}
	}
	return &types.AnalysisResult{
		Framework:  "",
		Main:       pkg.Main,
		Scripts:    pkg.Scripts,
		Databases:  nil,
		Redis:      false,
		RawSignals: RawSignals,
		Ports:      ports,
		Confirm:    ConfirmFunc,
	}, nil
}

func ValidateAnalysis(Analysis *types.AnalysisResult) (*types.AnalysisResult, error) {
	signals := utils.MergeDuplicateSignals(Analysis.RawSignals)
	signals = utils.SortSignals(Analysis.RawSignals)
	scriptFileRgxCompiled := regexp.MustCompile(`([A-Za-z]|[0-9]|-|_)+\.(js|ts)`)
	var signalFileNames []string
	var portSignals []types.Signal
	for _, signal := range signals {
		signalFileNames = append(signalFileNames, signal.File)
	}
	if len(Analysis.Main) > 0 {
		ind, f := slices.BinarySearch(signalFileNames, Analysis.Main)
		if f {
			moved, _ := utils.MoveSignal(signals, slices.Index(signals, signals[ind]), ind)
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
					ind, f := slices.BinarySearch(signalFileNames, matches[0])
					if f {
						moved, _ := utils.MoveSignal(signals, slices.Index(signals, signals[ind]), ind)
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
			portSignals = append(portSignals, signal)
		}
	}
	if len(portSignals) > 1 {
		message, _ := Analysis.Confirm(portSignals)
		var selectedIndex int
		fmt.Print(message, "\n> ")
		_, err := fmt.Scanln(&selectedIndex)
		if err != nil {
			return nil, err
		}
		if selectedIndex > len(portSignals) {
			fmt.Printf("Invalid entry [%d]| Out of range", selectedIndex)
		}
		if len(portSignals[selectedIndex].File) > 0 {
			finalizedSignals, _ := utils.MoveSignal(portSignals, selectedIndex, 0)
			signals = finalizedSignals
			fmt.Printf("Selected Signal\nFile: [%s]\nPort: [%d]\nSignal_Type: [%s]", signals[0].File, signals[0].Port, signals[0].Type)
		} else {
			fmt.Printf("Invalid entry [%d] | Signal not found", selectedIndex)
		}
	}

	return &types.AnalysisResult{
		Framework:  string(signals[0].Type),
		RawSignals: signals,
	}, nil
}

func ConfirmFunc(signals []types.Signal) (string, error) {

	var message strings.Builder
	message.WriteString("Multiple signals detected | Confirmation required\nSignals:\n")
	for i, signal := range signals {
		format := "\n " + strconv.Itoa(i) + " File: " + string(signal.File) + "| " + "Signal: " + string(signal.Type)
		message.WriteString(format)
	}
	message.WriteString("\nInput the corresponding index for the correct file")
	return message.String(), nil
}
