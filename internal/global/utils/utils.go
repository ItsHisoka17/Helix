package utils

func JoinPath(paths ...string) (path string) {
	for i, p := range paths {
		if i == len(paths) {
			path += p
			continue
		}
		path += ("/" + p)
	}
	return path
}
