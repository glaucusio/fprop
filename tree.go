package fprop

import (
	"errors"
	stdpath "path"
	"sort"

	"github.com/glaucusio/fprop/internal/object"
)

const (
	FlagDelete = 1 << iota
	FlagMerge
)

type Tree map[string]struct {
	Value    interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	Children Tree        `json:"children,omitempty" yaml:"children,omitempty"`
	Flag     byte        `json:"flag,omitempty" yaml:"flag,omitempty"`
}

func (t Tree) Add(path string, v interface{}, flag byte) Tree {
	fn := func(_ string, t Tree, _ interface{}, _ byte) (Tree, interface{}, byte) {
		return t, v, flag
	}

	return t.MapAt(path, fn)
}

func (t Tree) Del(path string) {
	dir, file := stdpath.Split(stdpath.Clean(path))

	for _, name := range object.Split(dir) {
		node, ok := t[name]
		if !ok {
			return
		}
		t = node.Children
	}

	delete(t, file)
}

var _ = make(Tree).Copy()

func (t Tree) Copy() Tree {
	u := make(Tree)

	uadd := func(path string, v interface{}, prop byte) error {
		u.Add(path, v, prop)
		return nil
	}

	if err := t.Visit(uadd); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return u
}

func (t Tree) MapAt(path string, fn MapFunc) Tree {
	path = stdpath.Join("/", stdpath.Clean(path))
	dir, file := stdpath.Split(path)

	for _, name := range object.Split(dir) {
		node, ok := t[name]
		if !ok || node.Children == nil {
			node.Children = make(Tree)
			t[name] = node
		}
		t = node.Children
	}

	n := t[file]
	n.Children, n.Value, n.Flag = fn(path, n.Children, n.Value, n.Flag)
	t[file] = n

	return t
}

func (t Tree) Map(fn MapFunc) {
	t.mapeach("/", fn)
}

func (t Tree) mapeach(path string, fn MapFunc) {
	for _, k := range t.Keys() {
		path := stdpath.Join(path, k)

		if n := t[k]; len(n.Children) != 0 {
			n.Children.mapeach(path, fn)
		} else {
			n.Children, n.Value, n.Flag = fn(path, nil, n.Value, n.Flag)
			t[k] = n
		}
	}
}

func (t Tree) Keys() []string {
	keys := make([]string, 0, len(t))

	for k := range t {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func (t Tree) Object() interface{} {
	obj := make(map[string]interface{})

	for name, node := range t {
		if len(node.Children) != 0 {
			obj[name] = node.Children.Object()
		} else {
			obj[name] = node.Value
		}
	}

	if a := object.Slice(obj); len(a) != 0 {
		return a
	}

	return obj
}

func (t Tree) Merge(u Tree) Tree {
	return Merge(t, u)
}

func (t Tree) Unmarshal(v interface{}) error {
	return Unmarshal(t, v)
}

func (t Tree) Visit(fn VisitFunc) error {
	return t.visit("/", fn)
}

func (t Tree) visit(path string, fn VisitFunc) (err error) {
	for _, k := range t.Keys() {
		path := stdpath.Join(path, k)

		if len(t[k].Children) != 0 {
			err = t[k].Children.visit(path, fn)
		} else {
			if t[k].Value != nil || t[k].Flag != 0 { // todo: ensure it's needed
				err = fn(path, t[k].Value, t[k].Flag)
			}
		}

		switch {
		case errors.Is(err, SkipNode):
			return nil
		case err != nil:
			return err
		}
	}

	return nil
}
