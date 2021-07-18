package main

import (
	"github.com/owenoclee/checkexhaustive"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(checkexhaustive.Analyzer)
}
