package fprop_test

import (
	"flag"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"gopkg.in/yaml.v2"
)

var updateGolden = flag.Bool("update-golden", false, "Update golden files")

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func equal(v1, v2 interface{}) (bool, string) {
	p1, err := yaml.Marshal(v1)
	if err != nil {
		panic("unexpected error: " + err.Error())
	}

	p2, err := yaml.Marshal(v2)
	if err != nil {
		panic("unexpected error: " + err.Error())
	}

	var w1, w2 interface{}

	if err := yaml.Unmarshal(p1, &w1); err != nil {
		panic("unexpected error: " + err.Error())
	}

	if err := yaml.Unmarshal(p2, &w2); err != nil {
		panic("unexpected error: " + err.Error())
	}

	if !cmp.Equal(w1, w2) {
		return false, cmp.Diff(w1, w1)
	}

	return true, ""
}
