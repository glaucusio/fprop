package fprop_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/glaucusio/fprop"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v2"
)

func TestTreeAdd(t *testing.T) {
	cases := map[string]string{
		"testdata/scylla.yaml/kv.yaml":      "testdata/scylla.yaml/tree.yaml",
		"testdata/rules.yaml/kv.yaml":       "testdata/rules.yaml/tree.yaml",
		"testdata/rule_config.yaml/kv.yaml": "testdata/rule_config.yaml/tree.yaml",
	}

	for test, golden := range cases {
		var kv map[string]struct {
			Value    interface{}
			Property byte
		}

		p, err := ioutil.ReadFile(test)
		if err != nil {
			t.Fatalf("ReadFile()=%s", err)
		}

		if err := yaml.Unmarshal(p, &kv); err != nil {
			t.Fatalf("Unmarshal()=%s", err)
		}

		got := make(fprop.Tree)

		for path, v := range kv {
			got.Add(path, v.Value, v.Property)
		}

		if *updateGolden {
			p, err := yaml.Marshal(got)
			if err != nil {
				t.Fatalf("Unmarshal()=%s", err)
			}

			if err := ioutil.WriteFile(golden, p, 0644); err != nil {
				t.Fatalf("WriteFile()=%s", err)
			}

			continue
		}

		p, err = ioutil.ReadFile(golden)
		if err != nil {
			t.Fatalf("Unmarshal()=%s", err)
		}

		var want fprop.Tree

		if err := yaml.Unmarshal(p, &want); err != nil {
			t.Fatalf("Unmarshal()=%s", err)
		}

		if !cmp.Equal(got, want) {
			t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
		}
	}
}

func TestTreeCopy(t *testing.T) {

}

func TestTreeMerge(t *testing.T) {

}

func TestTreeVisit(t *testing.T) {
	cases := map[string]string{
		"testdata/scylla.yaml/tree.yaml":      "testdata/scylla.yaml/kv.yaml",
		"testdata/rules.yaml/tree.yaml":       "testdata/rules.yaml/kv.yaml",
		"testdata/rule_config.yaml/tree.yaml": "testdata/rule_config.yaml/kv.yaml",
	}

	for test, golden := range cases {
		var tree fprop.Tree

		p, err := ioutil.ReadFile(test)
		if err != nil {
			t.Fatalf("ReadFile()=%s", err)
		}

		if err := yaml.Unmarshal(p, &tree); err != nil {
			t.Fatalf("Unmarshal()=%s", err)
		}

		got := make(map[string]struct {
			Value    interface{}
			Property byte
		})

		err = tree.Visit(func(path string, value interface{}, prop byte) error {
			v, ok := got[path]
			if ok {
				return fmt.Errorf("duplicated visit: %q", path)
			}

			v.Value = value
			v.Property = prop

			got[path] = v

			return nil
		})

		if err != nil {
			t.Fatalf("Visit()=%s", err)
		}

		if *updateGolden {
			p, err := yaml.Marshal(got)
			if err != nil {
				t.Fatalf("Unmarshal()=%s", err)
			}

			if err := ioutil.WriteFile(golden, p, 0644); err != nil {
				t.Fatalf("WriteFile()=%s", err)
			}

			continue
		}

		p, err = ioutil.ReadFile(golden)
		if err != nil {
			t.Fatalf("Unmarshal()=%s", err)
		}

		var want map[string]interface{}

		if err := yaml.Unmarshal(p, &want); err != nil {
			t.Fatalf("Unmarshal()=%s", err)
		}

		if ok, diff := equal(got, want); !ok {
			t.Fatalf("got != want:\n:%s", diff)
		}
	}
}
