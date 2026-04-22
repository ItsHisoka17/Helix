package types

import sitter "github.com/smacker/go-tree-sitter"

type SignalType string

const (
	ExpressType SignalType = "express_usage"
	RedisType   SignalType = "redis_usage"
	SignalPort  SignalType = "port_detected"
)

var MainCmds []string = []string{"start", "run"}

var MainFiles []string = []string{"server", "app", "main", "index"}

type Signal struct {
	Type       SignalType
	File       string
	Port       int
	Confidence float64
	Framework  bool
}

type PackageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Main            string            `json:"main"`
	Scripts         map[string]string `json:"scripts"`
}

type Confirm func([]Signal, string) ([]Signal, error)

type AnalysisResult struct {
	ProjectPath  string
	Dependencies map[string]string
	Main         string
	Scripts      map[string]string
	Framework    string
	Databases    []string
	Redis        bool
	RawSignals   []Signal
	Ports        []int
	Confirm      Confirm
}

type QueryMatch struct {
	Captures map[string]*sitter.Node
}

type DetectFunc func(*sitter.Node, []byte, string) []Signal
