package unused_test

import (
	"testing"

	"github.com/gostaticanalysis/unused"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, unused.Analyzer, "a")
}
