package fproputil_test

import (
	"io/ioutil"
	"testing"

	"github.com/glaucusio/fprop"
	"github.com/glaucusio/fprop/fproputil"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v2"
)

func TestMakeTree(t *testing.T) {
	cases := map[string]string{
		"testdata/scylla.yaml/object.yaml":      "testdata/scylla.yaml/tree.yaml",
		"testdata/rules.yaml/object.yaml":       "testdata/rules.yaml/tree.yaml",
		"testdata/rule_config.yaml/object.yaml": "testdata/rule_config.yaml/tree.yaml",
	}

	for test, golden := range cases {
		p, err := ioutil.ReadFile(test)
		if err != nil {
			t.Fatalf("ReadFile()=%s", err)
		}

		var obj map[string]interface{}

		if err := yaml.Unmarshal(p, &obj); err != nil {
			t.Fatalf("Unmarshal()=%s", err)
		}

		if *updateGolden {
			tree := fproputil.MakeTree(obj)

			p, err := yaml.Marshal(tree)
			if err != nil {
				t.Fatalf("Marshal()=%s", err)
			}

			if err := ioutil.WriteFile(golden, p, 0644); err != nil {
				t.Fatalf("WriteFile()=%s", err)
			}

			continue
		}

		p, err = ioutil.ReadFile(golden)
		if err != nil {
			t.Fatalf("ReadFile()=%s", err)
		}

		var got, want fprop.Tree

		got = fproputil.MakeTree(obj)

		if err := yaml.Unmarshal(p, &want); err != nil {
			t.Fatalf("Unmarshal()=%s", err)
		}

		if !cmp.Equal(got, want) {
			t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
		}

		if got, want := got.Object(), want.Object(); !cmp.Equal(got, want) {
			t.Fatalf("obj != want:\n%s", cmp.Diff(got, want))
		}

		p, err = yaml.Marshal(got.Object())
		if err != nil {
			t.Fatalf("Marshal()=%s", err)
		}

		q, err := yaml.Marshal(obj)
		if err != nil {
			t.Fatalf("Marshal()=%s", err)
		}

		var v, w interface{}

		if err := yaml.Unmarshal(p, &v); err != nil {
			t.Fatalf("Unmarshal()=%s", err)
		}

		if err := yaml.Unmarshal(q, &w); err != nil {
			t.Fatalf("Unmarshal()=%s", err)
		}

		if !cmp.Equal(v, w) {
			t.Fatalf("v != w:\n%s", cmp.Diff(v, w))
		}
	}
}
