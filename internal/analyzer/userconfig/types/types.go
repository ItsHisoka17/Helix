package types

type KubernetesConfig struct {
	UseKubernetes bool
	Version       string
	AppName       string
	MinReplicas   int
	MaxReplicas   int
	Ingress       bool
}

type DockerConfig struct {
	Dockerize     bool
	DockerCompose bool
}

type AWSConfig struct {
	UseTerraform string
}

type DefaulConfig struct {
	Docker     DockerConfig
	Kubernetes KubernetesConfig
	AWS        AWSConfig
}
