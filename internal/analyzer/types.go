package analyzer

type SignalType string

const (
	ExpressType SignalType = "express_usage"
	RedisType   SignalType = "redis_usage"
	SignalPort  SignalType = "port_detected"
)

type Signal struct {
	Type SignalType
	File string
	Port int
}

type AST struct {
}

type PackageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

type CodeSignals struct {
	RawSignals []Signal
}

type AnalysisResult struct {
	ProjectPath  string
	Dependencies map[string]string
	Framework    string
	Databases    []string
	Redis        bool
	RawSignals   []Signal
	Ports        []int
}
