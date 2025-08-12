package plugin

import (
	"github.com/urfave/cli/v2"
)

type EigenRuntimePlugin interface {
	Name() string
	Version() string
	DescribeCommands() []*cli.Command
	GetCommands() []*cli.Command
	RunCommands() []*cli.Command
	RemoveCommands() []*cli.Command
}
