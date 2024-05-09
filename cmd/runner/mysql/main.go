package main

import (
	"github.com/ideal-rucksack/workflow-glolang-plugin/cmd/runner"
	_ "github.com/ideal-rucksack/workflow-glolang-plugin/cmd/runner/mysql/plugin"
)

func main() {
	runner.Run()
}
