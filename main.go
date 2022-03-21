package main

import (
	"os"

	"github.com/smmr-software/mabel/internal/full"
	"github.com/smmr-software/mabel/internal/mini"
)

func main() {
	args := os.Args[1:]
	if len(args) == 1 {
		mini.Execute(&args[0])
	} else {
		full.Execute()
	}
}
