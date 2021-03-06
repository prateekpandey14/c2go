package ast

import "fmt"

type UnaryOperator struct {
	Address  string
	Position string
	Type     string
	IsLvalue bool
	IsPrefix bool
	Operator string
	Children []Node
}

func parseUnaryOperator(line string) *UnaryOperator {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.*?)'
		(?P<lvalue> lvalue)?
		(?P<prefix> prefix)?
		(?P<postfix> postfix)?
		 '(?P<operator>.*?)'`,
		line,
	)

	return &UnaryOperator{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		IsLvalue: len(groups["lvalue"]) > 0,
		IsPrefix: len(groups["prefix"]) > 0,
		Operator: groups["operator"],
		Children: []Node{},
	}
}

func (n *UnaryOperator) render(ast *Ast) (string, string) {
	operator := n.Operator
	expr, exprType := renderExpression(ast, n.Children[0])

	if operator == "!" {
		if exprType == "bool" {
			return fmt.Sprintf("!(%s)", expr), exprType
		}

		ast.addImport("github.com/elliotchance/c2go/noarch")
		return fmt.Sprintf("%s(%s)", fmt.Sprintf("noarch.Not%s", ucfirst(exprType)), expr), exprType
	}

	if operator == "*" {
		if exprType == "const char *" {
			return fmt.Sprintf("%s[0]", expr), "char"
		}

		return fmt.Sprintf("*%s", expr), "int"
	}

	if operator == "++" {
		return fmt.Sprintf("%s += 1", expr), exprType
	}

	if operator == "~" {
		operator = "^"
	}

	return fmt.Sprintf("%s%s", operator, expr), exprType
}

func (n *UnaryOperator) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
