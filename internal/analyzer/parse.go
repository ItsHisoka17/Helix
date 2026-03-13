package analyzer

import (
	"context"
	"strconv"

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

func ParsePort(node *sitter.Node, file []byte) (*Signal, error) {

	if node.Type() == "call_expression" {
		function := node.ChildByFieldName("function")
		if function != nil {
			if property := function.ChildByFieldName("property"); property != nil && property.Content(file) == "listen" {
				args := node.ChildByFieldName("arguments")
				portNode := args.NamedChild(0)
				if portNode != nil {
					portS := portNode.Content(file)
					port, err := strconv.Atoi(portS)
					if err != nil {
						return nil, err
					}
					return &Signal{
							Type: "port_detected",
							Port: port,
						},
						nil
				}
			}
		}
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		signal, err := ParsePort(node.Child(i), file)
		if signal != nil || err != nil {
			return signal, err
		}
	}
	return nil, nil
}
