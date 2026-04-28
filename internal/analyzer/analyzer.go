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

func Analyze(projectPath string) (*types.AnalysisResult, *types.PackageJSON, error) {
	pkg, err := ParsePackageJson(projectPath)
	if err != nil {
		return nil, nil, err
	}
	RawSignals, err := ParseCodeContext(projectPath)
	if err != nil {
		return nil, nil, err
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
	}, pkg, nil
}

func ValidateAnalysis(projectPath string, test bool) (*types.AnalysisResult, error) {
	Analysis, pkg, err := Analyze(projectPath)
	if err != nil {
		fmt.Printf("Error during analysis [AnalyzerError]\n%s", err.Error())
		return nil, err
	}
	signals := filterSignals(Analysis.RawSignals, pkg)
	scriptFileRgxCompiled := regexp.MustCompile(`([A-Za-z]|[0-9]|-|_)+\.(js|ts)`)
	var signalFileNames []string
	var framework string
	var redis bool

	for _, signal := range signals {
		if len(framework) < 1 {
			if signal.Type != types.SignalPort && signal.Type != types.RedisType {
				framework = strings.Split(string(signal.Type), "_")[0]
			}
		}
		signalFileNames = append(signalFileNames, signal.File)
		if !redis && signal.Type == types.RedisType {
			redis = true
		}
	}
	if len(Analysis.Main) > 0 {
		ind := slices.Index(signalFileNames, Analysis.Main)
		if ind != -1 {
			moved, _ := utils.MoveSignal(signals, ind, 0)
			if moved != nil {
				signals = moved
			}
		}
	}
	if Analysis.Scripts != nil {
		for _, cmd := range types.MainCmds {
			if len(Analysis.Scripts[cmd]) > 0 {
				matches := scriptFileRgxCompiled.FindStringSubmatch(Analysis.Scripts[cmd])
				if len(matches) > 0 && len(matches[0]) > 0 {
					ind := slices.Index(signalFileNames, matches[0])
					if ind != -1 {
						moved, _ := utils.MoveSignal(signals, ind, 0)
						if moved != nil {
							signals = moved
						}
					}
				}
			}
		}
	}
	if len(signals) > 1 {
		confirmedSignals, _ := Analysis.Confirm(signals, projectPath, test)
		if confirmedSignals != nil {
			signals = confirmedSignals
		}
	}

	return &types.AnalysisResult{
		Framework:  framework,
		RawSignals: signals,
		Main:       Analysis.Main,
		Scripts:    Analysis.Scripts,
		Databases:  Analysis.Databases,
		Ports:      Analysis.Ports,
		Redis:      true,
	}, nil
}

func ConfirmFunc(signals []types.Signal, projectPath string, test bool) ([]types.Signal, error) {
	if test {
		return nil, nil
	}
	var message strings.Builder
	message.WriteString("Multiple signals detected | Confirmation required\nSignals:\n")
	for i, signal := range signals {
		if signal.Type == types.RedisType {
			continue
		}
		format := "\n " + strconv.Itoa(i) + " File: " + string(signal.File) + "| " + "Signal: " + string(signal.Type)
		message.WriteString(format)
	}
	message.WriteString("\nInput the corresponding index for the correct file")
	var selectedIndex int
	fmt.Print(message.String(), "\n> ")
	_, err := fmt.Scanln(&selectedIndex)
	if err != nil {
		return nil, err
	}

	if selectedIndex >= len(signals) {
		fmt.Printf("Invalid entry [%d]| Out of range\n", selectedIndex)
		ValidateAnalysis(projectPath, test)
		return nil, nil
	}
	if len(signals[selectedIndex].File) > 0 {
		finalizedSignals, _ := utils.MoveSignal(signals, selectedIndex, 0)
		signals = finalizedSignals
		fmt.Printf("Selected Signal\nFile: [%s]\nPort: [%d]\nSignal_Type: [%s]\n", signals[0].File, signals[0].Port, signals[0].Type)
	} else {
		fmt.Printf("Invalid entry [%d] | Signal not found\n", selectedIndex)
	}
	return signals, nil
}

func filterSignals(signals []types.Signal, pkg *types.PackageJSON) []types.Signal {
	signals = utils.MergeSignals(signals)
	signals = utils.AssignConfidence(signals, pkg)
	return utils.SortSignals(signals)
}
