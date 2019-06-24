package main

import (
	"github.com/gostaticanalysis/unused"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(unused.Analyzer) }