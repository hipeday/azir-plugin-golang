package command

import "github.com/hipeday/azir-plugin-golang/pkg/plugin"

var (
	Registry plugin.CommandRegistry
)

func init() {
	Registry = *plugin.NewCommandRegistry()
}
