package analyzer

import "maps"

func DetectDependencies(pkg *PackageJSON) (frameworks []string, dbs []string, redis bool) {
	deps := mergeDeps(pkg.Dependencies, pkg.DevDependencies)
	for dep := range deps {
		switch dep {
		case "express":
			{
				frameworks = append(frameworks, "express")
			}

		case "fastify":
			{
				frameworks = append(frameworks, "fastify")
			}

		case "@nestjs/core":
			{
				frameworks = append(frameworks, "nest")
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
	maps.Copy(list, a)

	maps.Copy(list, b)

	return list
}
