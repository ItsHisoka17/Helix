package analyzer

import (
	"io/fs"
	"path/filepath"
	"regexp"

	"github.com/ItsHisoka17/Helix/internal/analyzer/types"
)

func IndexRepositoryFiles(path string) {
	fileRegex := regexp.MustCompile(`(?i)(^|\/)(server|main|app|index|bootstrap|startup|entry|api|gateway|router|http|grpc|worker|consumer|producer|queue|scheduler|cron|job|task|daemon|service|deploy|infra|k8s|kubernetes|helm|terraform|docker|compose|nginx|caddy|traefik|proxy|config|settings|environment|env|database|db|cache|redis|mongo|postgres|mysql|sqlite|prisma|typeorm|sequelize|gorm|auth|session|middleware|pipeline|storage|upload|s3|pubsub|event|webhook|socket|ws|graphql|rpc|lambda|function|ci|cd|workflow|actions?|jenkins|gitlab-ci)([-_.]?[a-z0-9]+)*\.(js|ts|jsx|tsx|mjs|cjs|go|py|java|kt|rs|rb|php|cs|scala|swift|lua|sh|bash|zsh|yaml|yml|json|toml|ini|conf|env|tf|tfvars|dockerfile)$`)
	ignoreRegex := regexp.MustCompile(`node_modules|.git|dist|build|coverage|vendor|target|bin`)
	var contextMap types.ContextMap = make(map[string]string)
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && ignoreRegex.Match([]byte(d.Name())) {
			return nil
		}
		if fileRegex.Match([]byte(d.Name())) {
			contextMap[d.Name()] = path
		}
		return nil
	})
}

func FilterContext(context types.ContextMap) {

}
