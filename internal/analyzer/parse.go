package analyzer

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ItsHisoka17/Helix/internal/analyzer/types"
	sitter "github.com/smacker/go-tree-sitter"
)

func ParsePackageJson(projectpath string) (*types.PackageJSON, error) {
	path := projectpath + "/package.json"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pkg types.PackageJSON

	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}
	return &pkg, nil
}

func ParseCodeContext(projectPath string) ([]types.Signal, error) {
	var RawSignals []types.Signal
	err := filepath.WalkDir(projectPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "node_modules" {
				return filepath.SkipAll
			}
			return nil
		}
		if strings.HasSuffix(path, ".ts") || strings.HasSuffix(path, ".js") {
			context, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			node, err := GetNode(context)
			if err == nil {
				_, err := d.Info()
				if err == nil {
					file := d.Name()
					var Detectors [3]types.DetectFunc = [3]types.DetectFunc{DetectPorts, DetectExpress, DetectRedis}
					for _, Detect := range Detectors {
						RawSignals = append(
							RawSignals,
							Detect(node, context, file)...,
						)
					}
				}
			}

		}
		return nil
	})
	return RawSignals, err
}

func DetectPorts(root *sitter.Node, source []byte, file string) []types.Signal {
	query := `
	(call_expression
	  function: (member_expression
	    property: (property_identifier) @method)
	  arguments: (arguments
	    (number) @port))
	`
	matches, _ := RunQuery(root, source, query)

	var Signals []types.Signal

	for _, m := range matches {
		method := m.Captures["method"].Content(source)
		if method == "listen" {
			strPort := m.Captures["port"].Content(source)
			port, err := strconv.Atoi(strPort)
			if err != nil {
				continue
			}
			Signals = append(Signals, types.Signal{
				File:       file,
				Type:       "port_detected",
				Port:       port,
				Confidence: 0.95,
			})
		}
	}
	return Signals

}

func DetectExpress(node *sitter.Node, source []byte, file string) []types.Signal {
	query := `
	(call_expression
	function: (identifier) @fn)
	`
	matches, _ := RunQuery(node, source, query)
	var Signals []types.Signal
	for _, m := range matches {
		fn := m.Captures["fn"].Content(source)
		if fn == "express" {
			Signals = append(Signals, types.Signal{
				File:       file,
				Type:       "express_detected",
				Confidence: 0.9,
			})
		}
	}
	return Signals
}

func DetectRedis(node *sitter.Node, source []byte, file string) []types.Signal {
	query := `
	(call_expression
	function: (member_expression
	property: (property_identifier) @method))`

	matches, _ := RunQuery(node, source, query)
	var Signals []types.Signal
	for _, m := range matches {
		method := m.Captures["method"].Content(source)
		if method == "createClient" {
			Signals = append(Signals, types.Signal{
				File:       file,
				Type:       "redis_detected",
				Confidence: 0.85,
			})
		}
	}
	return Signals
}
