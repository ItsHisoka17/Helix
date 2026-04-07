package types

type Docker struct {
	Compile        bool
	Environment    string
	ImageName      string
	ImageCount     string
	ContainerName  string
	ContainerPorts map[int]int
	Scripts        []string
	EntryPoint     string
	Build          bool
	BuildCommand   string
	BuildDir       string
}
