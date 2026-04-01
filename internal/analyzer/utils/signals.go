package utils

import (
	"errors"
	"sort"

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

func MergeDuplicateSignals(signals []types.Signal) []types.Signal {
	var new map[string]types.Signal = make(map[string]types.Signal)
	var fixed []types.Signal
	for _, signal := range signals {
		seen, exists := new[signal.File]
		if exists {
			if signal.Port > 0 && seen.Port < 1 {
				seen.Port = signal.Port
			}
			if seen.Confidence < signal.Confidence {
				seen.Confidence = signal.Confidence
			}
			if seen.Type == types.SignalPort && signal.Type != types.SignalPort {
				seen.Type = signal.Type
			}
			new[signal.File] = seen
			continue
		}
		new[signal.File] = signal
	}
	for _, signal := range new {
		fixed = append(fixed, signal)
	}
	return fixed
}
