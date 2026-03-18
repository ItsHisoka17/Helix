package analyzer

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

func RunQuery(node *sitter.Node, source []byte, queryStr string) ([]QueryMatch, error) {
	var QueryMatches []QueryMatch
	query, err := sitter.NewQuery([]byte(queryStr), javascript.GetLanguage())
	if err != nil {
		return nil, err
	}
	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, node)
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		captures := make(map[string]*sitter.Node)
		for _, c := range match.Captures {
			name := query.CaptureNameForId(c.Index)
			captures[name] = c.Node
		}
		QueryMatches = append(QueryMatches,
			QueryMatch{
				Captures: captures,
			})
	}
	return QueryMatches, nil
}
