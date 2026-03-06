package analyzer

type AnalysisResult struct {
	ProjectPath  string
	Dependencies map[string]string
	Framework    string
	Databases    []string
	Redis        bool
	RawSignals   []string
	Ports        []int
}
