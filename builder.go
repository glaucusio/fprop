package fprop

var DefaultBuilder = Object

var Object = &Builder{}

var JSON = &Builder{}

var YAML = &Builder{}

type Builder struct {
}

func (b *Builder) Make(obj interface{}) Tree {
	return nil
}

func (b *Builder) Merge(t, u Tree) Tree {
	tmerge := func(path string, v interface{}, f byte) error {
		switch {
		case (f & FlagDelete) != 0:
			t.Del(path)
		case (f & FlagMerge) != 0:
			merge := func(_ string, t Tree, w interface{}, g byte) (Tree, interface{}, byte) {
				return t, v, f // todo
			}

			t.MapAt(path, merge)
		default:
			t.Add(path, v, 0)
		}
		return nil
	}

	if err := u.Visit(tmerge); err != nil {
		panic("unexpected error: " + err.Error())
	}

	return t
}

func (b *Builder) Unmarshal(t Tree, v interface{}) error {
	return nil
}

func Make(obj interface{}) Tree {
	return DefaultBuilder.Make(obj)
}

func Merge(t, u Tree) Tree {
	return DefaultBuilder.Merge(t, u)
}

func Unmarshal(t Tree, v interface{}) error {
	return DefaultBuilder.Unmarshal(t, v)
}
