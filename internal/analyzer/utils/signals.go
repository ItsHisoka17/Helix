package utils

import (
	"errors"
	"slices"
	"sort"
	"strconv"

	"github.com/ItsHisoka17/Helix/internal/analyzer/types"
)

func SortSignals(signals []types.Signal) []types.Signal {
	sort.Slice(signals, func(i, j int) bool {
		return signals[i].Confidence > signals[j].Confidence
	})
	return signals
}

func MoveSignal(signals []types.Signal, signal types.Signal, index int) ([]types.Signal, error) {
	if index > len(signals) {
		return nil, errors.New("Invalid index " + strconv.Itoa(index) + " out of range")
	}
	target := signals[index]
	currnt := slices.Index(signals, signal)
	if currnt > 0 {
		signals[currnt] = target
		signals[index] = signal
	}
	signals[len(signals)+1] = target
	signals[index] = signal
	return signals, nil
}

func MergeDuplicateSignals(signals []types.Signal) []types.Signal {
	var new map[string]types.Signal = make(map[string]types.Signal)
	var fixed []types.Signal
	for _, signal := range signals {
		if len(new) > 0 {
			if len(new[signal.File].Type) > 0 {
				if signal.Port > 0 && new[signal.File].Port == 0 {
					new[signal.File] = types.Signal{
						Type:       new[signal.File].Type,
						File:       signal.File,
						Port:       signal.Port,
						Confidence: new[signal.File].Confidence,
					}
					continue
				}
			}
		}
		new[signal.File] = signal
	}
	for _, signal := range new {
		fixed = append(fixed, signal)
	}
	return fixed
}
