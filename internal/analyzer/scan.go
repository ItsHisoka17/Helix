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

type PackageJson struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func ParsePackageJson(projectpath string) (*PackageJson, error) {
	path := projectpath + "/package.json"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pkg PackageJson

	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}
	return &pkg, nil
}

func ParseCodeContext(projectpath string) ([]string, []int, error) {
	var RawSignals []string
	var Ports []int
	err := filepath.WalkDir(projectpath, func(path string, d fs.DirEntry, err error) error {
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
				RawSignals = append(RawSignals, "express_usage")
			}
			if strings.Contains(rawText, "createClient(") {
				RawSignals = append(RawSignals, "redis_usage")
			}
			if strings.Contains(rawText, "listen") {
				Ports = parsePort(rawText)
			}
		}
		return nil
	})
	return RawSignals, Ports, err
}

func parsePort(s string) (Ports []int) {
	if rgx := regexp.MustCompile(`/.*\(([0-9]{2,4})(,|\)|\s).*/gm`); rgx.MatchString(s) {
		matches := rgx.FindAllString(s, 5)
		for _, m := range matches {
			if rgx := regexp.MustCompile(`/^[0-9]{2,4}$/gm`); rgx.MatchString(m) {
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
