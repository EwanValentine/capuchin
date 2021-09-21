package query

// One day...
// import (
// 	"log"

// 	"github.com/pingcap/parser"
// 	"github.com/pingcap/parser/ast"
// 	_ "github.com/pingcap/parser/parser_driver"
// )

// func parse(sql string) (*ast.StmtNode, error) {
// 	p := parser.New()

// 	stmtNodes, _, err := p.Parse(sql, "", "")
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &stmtNodes[0], nil
// }

// // ParseQuery -
// func ParseQuery(sql string) (*Query, error) {
// 	astNode, err := parse("SELECT a, b FROM t")
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Println(astNode)

// 	return &Query{}, nil
// }
