package fproputil

import (
	"github.com/glaucusio/fprop"
	"github.com/glaucusio/fprop/internal/object"
)

type Builder struct {
	MapFunc func(interface{}) interface{}
}

func MakeTree(obj map[string]interface{}) fprop.Tree {
	root := make(fprop.Tree)

	queue := []struct {
		obj   map[string]interface{}
		nodes fprop.Tree
	}{{obj, root}}

	for i := queue[0]; len(queue) != 0; {
		i, queue = queue[0], queue[1:]

		for k, v := range i.obj {
			node, ok := i.nodes[k]
			if !ok {
				node.Children = make(fprop.Tree)
			}

			j := i
			j.nodes = node.Children

			if obj := object.Object(v); obj != nil {
				j.obj = obj
				queue = append(queue, j)
			} else {
				node.Value = v
				node.Children = nil
			}

			i.nodes[k] = node
		}
	}

	return root
}
