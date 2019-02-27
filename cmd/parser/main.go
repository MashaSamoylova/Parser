package main

import (
    "github.com/MashaSamoylova/Parser/pkg/parser"
)

func main() {
	system := parser.ParseSystem{
		Expression: "a + b",
		Grammar: []parser.Rule{
			{
				Term: "B",
				Alternatives: []string{
					"T + B",
					"T",
				},
			},
			{
				Term: "T",
				Alternatives: []string{
					"M",
					"M * T",
				},
			},
			{
				Term: "M",
				Alternatives: []string{
					"a",
					"b",
				},
			},
		},
	}
	system.Parse()
}
