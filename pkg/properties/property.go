package properties

import "github.com/ideal-rucksack/workflow-glolang-plugin/pkg/logging"

type Property struct {
	InvokeId    string         `json:"invoke_id"`
	Logger      logging.Config `json:"logger"`
	Name        string         `json:"name"`
	Version     string         `json:"version"`
	Home        string         `json:"home"`
	Description string         `json:"description"`
	Logfile     string         `json:"logfile"`
	Cmd         CmdProperty    `json:"cmd"`
}

type CmdProperty struct {
	Linux   []string `json:"linux"`
	Windows []string `json:"windows"`
	Darwin  []string `json:"darwin"`
}
