package analyzer

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

func GetNode(file []byte) (*sitter.Node, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(javascript.GetLanguage())
	tree, err := parser.ParseCtx(context.Background(), nil, file)
	node := tree.RootNode()
	if err != nil {
		return nil, err
	}
	return node, err
}
