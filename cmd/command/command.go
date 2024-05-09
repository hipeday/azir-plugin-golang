package command

import "github.com/ideal-rucksack/workflow-glolang-plugin/pkg/plugin"

var (
	Registry plugin.CommandRegistry
)

func init() {
	Registry = *plugin.NewCommandRegistry()
}
