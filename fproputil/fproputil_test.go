package fproputil_test

import (
	"flag"
	"os"
	"testing"
)

var updateGolden = flag.Bool("update-golden", false, "Update golden files")

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}
