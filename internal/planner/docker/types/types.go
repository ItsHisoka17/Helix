package types

type Main struct {
	File string
	Path string
	Port int
}

type DockerCompile struct {
	PortMap  map[int]string
	Commands map[string]string
	Main     Main
}

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
