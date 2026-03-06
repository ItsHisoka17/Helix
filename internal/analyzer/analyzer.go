package analyzer

func Analyze(projectpath string) (*AnalysisResult, error) {
	deps, err := ParsePackageJson(projectpath)
	if err != nil {
		return nil, err
	}
	framework, dbs, redis := DetectDependencies(deps)
	signals, ports, err := ParseCodeContext(projectpath)
	if err != nil {
		return nil, err
	}

	return &AnalysisResult{
		Framework:  framework,
		Databases:  dbs,
		Redis:      redis,
		RawSignals: signals,
		Ports:      ports,
	}, nil
}
