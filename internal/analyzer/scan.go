package analyzer

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var rgx *regexp.Regexp = regexp.MustCompile(`.*\(([0-9]{2,4})(,|\)|\s).*`)
var rgxP *regexp.Regexp = regexp.MustCompile(`^[0-9]{2,5}$`)

func ParsePackageJson(projectpath string) (*PackageJSON, error) {
	path := projectpath + "/package.json"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pkg PackageJSON

	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}
	return &pkg, nil
}

func ParseCodeContext(projectPath string) (*CodeSignals, error) {
	var RawSignals []Signal
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
			rawText := string(context)
			if strings.Contains(rawText, "express") {
				RawSignals = append(RawSignals, Signal{
					Type: "express_usage",
					File: path,
				})
			}
			if strings.Contains(rawText, "createClient(") {
				RawSignals = append(RawSignals, Signal{
					Type: "redis_usage",
					File: path,
				})
			}
			ports := parsePort(rawText)
			for _, port := range ports {
				RawSignals = append(RawSignals, Signal{
					Type: SignalPort,
					File: path,
					Port: port,
				})
			}
		}
		return nil
	})
	return &CodeSignals{RawSignals: RawSignals}, err
}

func parsePort(s string) (Ports []int) {
	if rgx.MatchString(s) {
		matches := rgx.FindAllString(s, 5)
		for _, m := range matches {
			if rgxP.MatchString(m) {
				port, err := strconv.Atoi(m)
				if err != nil {
					continue
				}
				Ports = append(Ports, port)
			}
		}
	}
	return
}
