package ast

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Zac-Garby/lang/token"
)

const treeIndent = 2

func in(indent int) string {
	return strings.Repeat(" ", treeIndent*indent)
}

func prefix(indent int, name string) string {
	str := in(indent)

	if name != "" {
		str += name + " ‣ "
	}

	return str
}

// Tree returns a tree representation of a node
func Tree(node Node, indent int, name string) string {
	val := reflect.ValueOf(node)

	typeName := fmt.Sprintf("%T", node)

	if name != "" {
		name = fmt.Sprintf("%s (%s)", name, typeName)
	} else {
		name = fmt.Sprintf("(%s)", typeName)
	}

	str := prefix(indent, name) + val.Type().Name()

	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.Type().NumField(); i++ {
		f := val.Field(i)
		if !f.CanInterface() {
			continue
		}

		field := f.Interface()
		label := val.Type().Field(i).Name

		if label == "Value" {
			label = ""
		}

		if _, ok := field.(*Nil); ok {
			str += "\n" + prefix(indent+1, label) + "<nil>"
			continue
		}

		switch n := field.(type) {
		case token.Token:
		case Node:
			str += "\n" + Tree(n, indent+1, label)

		case map[Statement]Statement:
			nodes := make(map[Node]Node)

			for key, value := range n {
				nodes[key] = value
			}

			str += "\n" + makeDictTree(indent+1, nodes, label)

		case map[Expression]Expression:
			nodes := make(map[Node]Node)

			for key, value := range n {
				nodes[key] = value
			}

			str += "\n" + makeDictTree(indent+1, nodes, label)

		case []Statement:
			var nodes []Node

			for _, item := range n {
				nodes = append(nodes, item.(Node))
			}

			str += "\n" + makeListTree(indent+1, nodes, label)

		case []Expression:
			var nodes []Node

			for _, item := range n {
				nodes = append(nodes, item.(Node))
			}

			str += "\n" + makeListTree(indent+1, nodes, label)

		case []MatchBranch:
			str := prefix(indent, name) + "["

			if len(n) == 0 {
				return str + "]"
			}

			for _, branch := range n {
				str += fmt.Sprintf("%s\n%s\n%s\n",
					in(indent),
					Tree(branch.Condition, indent+1, "cond"),
					Tree(branch.Body, indent+1, "body"),
				)
			}

			return str + "\n" + in(indent) + "]"

		default:
			str += "\n" + fmt.Sprintf("%s%s", prefix(indent+1, label), fmt.Sprintf("%v", n))
		}
	}

	return str
}

func makeListTree(indent int, nodes []Node, name string) string {
	str := prefix(indent, name) + "["

	if len(nodes) == 0 {
		return str + "]"
	}

	for _, node := range nodes {
		str += "\n" + Tree(node, indent+1, "")
	}

	return str + "\n" + in(indent) + "]"
}

func makeDictTree(indent int, pairs map[Node]Node, name string) string {
	str := prefix(indent, name) + "["

	if len(pairs) == 0 {
		return str + ":]"
	}

	for key, value := range pairs {
		str += fmt.Sprintf("%s\n%s\n%s\n",
			in(indent),
			Tree(key, indent+1, "key"),
			Tree(value, indent+1, "value"),
		)
	}

	return str + in(indent) + "]"
}
