package utils

import (
	"errors"
	"sort"
	"strings"

	"github.com/ItsHisoka17/Helix/internal/analyzer/types"
)

func SortSignals(signals []types.Signal) []types.Signal {
	sort.Slice(signals, func(i, j int) bool {
		return signals[i].Confidence > signals[j].Confidence
	})
	return signals
}

func MoveSignal(signals []types.Signal, from int, index int) ([]types.Signal, error) {
	if index >= len(signals) || index < 0 || from >= len(signals) || from < 0 {
		return signals, errors.New("Invalid index | out of range")
	}
	if from == index {
		return signals, nil
	}
	target := signals[from]
	if from < index {
		for i := from; i < index; i++ {
			signals[i] = signals[i+1]
		}
	} else {
		for i := from; i > index; i-- {
			signals[i] = signals[i-1]
		}
	}
	signals[index] = target
	return signals, nil
}

func MergeSignals(signals []types.Signal) []types.Signal {
	var seenMap map[string]types.Signal = make(map[string]types.Signal)
	var keep []types.Signal
	var merged []types.Signal
	for _, signal := range signals {
		seen, exists := seenMap[signal.File]
		if exists {
			keepSignal := signal.Port < 1 && (signal.Type != types.SignalPort && seen.Type == types.SignalPort)
			keepSeen := seen.Port < 1 && (seen.Type != types.SignalPort && signal.Type == types.SignalPort)
			if keepSignal {
				seenMap[signal.File] = types.Signal{
					File:       signal.File,
					Port:       seen.Port,
					Type:       signal.Type,
					Confidence: signal.Confidence,
				}
				continue
			}
			if keepSeen {
				seenMap[signal.File] = types.Signal{
					File:       signal.File,
					Port:       signal.Port,
					Type:       seen.Type,
					Confidence: signal.Confidence,
				}
				continue
			}
			if !keepSignal && !keepSeen {
				keep = append(keep, signal)
				continue
			}
		} else {
			seenMap[signal.File] = signal
		}
	}
	for _, s := range seenMap {
		merged = append(merged, s)
	}
	merged = append(merged, keep...)
	return merged
}

func AssignConfidence(signals []types.Signal, pkg *types.PackageJSON) []types.Signal {
	type score struct {
		f bool
		s float64
	}
	for i, signal := range signals {
		var scoreMap map[string]score = make(map[string]score)
		scoreMap["f"] = score{f: signal.Framework, s: 0.3}
		scoreMap["m"] = score{f: (signal.File == pkg.Main), s: 0.5}
		scoreMap["p"] = score{f: (signal.Port > 0), s: 0.2}
		var isScript bool
		for _, cmd := range types.MainCmds {
			script, exists := pkg.Scripts[cmd]
			if exists {
				if strings.Contains(script, signal.File) {
					isScript = true
				}
				break
			}
		}
		scoreMap["isScript"] = score{f: isScript, s: 0.5}
		for _, v := range scoreMap {
			if v.f {
				signals[i].Confidence = signals[i].Confidence + v.s
			}
		}
	}
	return signals
}
