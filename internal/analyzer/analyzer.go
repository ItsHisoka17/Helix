package analyzer

func Analyze(projectPath string) (*AnalysisResult, error) {
	deps, err := ParsePackageJson(projectPath)
	if err != nil {
		return nil, err
	}
	framework, dbs, redis := DetectDependencies(deps)
	codeSignals, err := ParseCodeContext(projectPath)
	if err != nil {
		return nil, err
	}
	var ports []int
	for _, signal := range codeSignals.RawSignals {
		ports = append(ports, signal.Port)
	}
	return &AnalysisResult{
		Framework:  framework,
		Databases:  dbs,
		Redis:      redis,
		RawSignals: codeSignals.RawSignals,
		Ports:      ports,
	}, nil
}
