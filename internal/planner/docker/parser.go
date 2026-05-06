package docker

import (
	"regexp"

	analyzerTypes "github.com/ItsHisoka17/Helix/internal/analyzer/types"
	types "github.com/ItsHisoka17/Helix/internal/planner/docker/types"
)

func ParseAnalysis(analysis *analyzerTypes.AnalysisResult) *types.DockerCompile {
	buildRgx := regexp.MustCompile(`(^tsc -b(\s(\w|[\/\.])+)?$)|(^(npm run build|npm install)(\s\-{1,2}(\w|[\/\.])+(\s(\w|[\/\.])+)?)?$)`)
	startRgx := regexp.MustCompile(`(^node \.$)|(^npm run start$)|(^node (\-{1,2}(\w|[\/\.])+\s)?(\w|[\/\.])+\.js$)`)
	var portMap map[int]string = map[int]string{}
	var main types.Main = types.Main{File: analysis.Main}
	var cmdMap map[string][]string = map[string][]string{}
	for _, signal := range analysis.RawSignals {
		if signal.Confidence > 0.85 && signal.Port > 0 {
			portMap[signal.Port] = signal.Path
		}
		if signal.File == main.File {
			main = types.Main{
				File: signal.File,
				Path: signal.Path,
				Port: signal.Port,
			}
		}
	}
	for _, script := range analysis.Scripts {
		if buildRgx.Match([]byte(script)) {
			cmdMap["build"] = append(cmdMap["build"], script)
			continue
		}
		if startRgx.Match([]byte(script)) {
			cmdMap["run"] = append(cmdMap["run"], script)
		}
	}
	return &types.DockerCompile{
		PortMap: portMap,
		Main:    main,
	}
}
