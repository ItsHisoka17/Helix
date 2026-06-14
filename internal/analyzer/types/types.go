package types

type ContextMap map[string]string

type ReconcilationResult struct {
	FileCount int
	Context   ContextMap
	Review    bool
}

type InferencePacket struct {
	Context ContextMap
}
