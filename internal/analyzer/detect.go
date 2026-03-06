package analyzer

func DetectDependencies(pkg *PackageJson) (framework string, dbs []string, redis bool) {
	deps := mergeDeps(pkg.Dependencies, pkg.DevDependencies)
	for dep := range deps {
		switch dep {
		case "express":
			{
				framework = "express"
			}

		case "fastify":
			{
				framework = "fastify"
			}

		case "@nestjs/core":
			{
				framework = "nest"
			}
		}
		switch dep {
		case "pg":
			{
				dbs = append(dbs, "postgres")
			}

		case "mongoose":
			{
				dbs = append(dbs, "mongodb")
			}

		case "mysql2":
			{
				dbs = append(dbs, "mysql")
			}

		case "prisma":
			{
				dbs = append(dbs, "prisma")
			}
		}
		if dep == "redis" || dep == "ioredis" {
			redis = true
		}
	}

	return
}

func mergeDeps(a, b map[string]string) map[string]string {
	list := make(map[string]string)
	for k, v := range a {
		list[k] = v
	}

	for k, v := range b {
		list[k] = v
	}

	return list
}
