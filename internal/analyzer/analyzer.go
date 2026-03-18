package analyzer

func Analyze(projectPath string) (*AnalysisResult, error) {
	pkg, err := ParsePackageJson(projectPath)
	if err != nil {
		return nil, err
	}
	frameworks, dbs, redis := DetectDependencies(pkg)
	codeSignals, err := ParseCodeContext(projectPath)
	if err != nil {
		return nil, err
	}
	var ports []int
	for _, signal := range codeSignals.RawSignals {
		if signal.Port != 0 {
			ports = append(ports, signal.Port)
		}
	}
	return &AnalysisResult{
		Framework:  frameworks,
		Main:       pkg.Main,
		Scripts:    pkg.Scripts,
		Databases:  dbs,
		Redis:      redis,
		RawSignals: codeSignals.RawSignals,
		Ports:      ports,
		Confirm:    nil,
	}, nil
}

func ValidateAnalysis(Analysis *AnalysisResult) (*AnalysisResult, error) {
	return nil, nil
}
