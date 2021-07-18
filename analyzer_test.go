package checkexhaustive_test

import (
	"testing"

	"github.com/owenoclee/checkexhaustive"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, checkexhaustive.Analyzer, "full")
}
