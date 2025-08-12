package plugin

import (
	"testing"
	"github.com/urfave/cli/v2"
)

// TestPluginImplementation verifies the interface can be implemented
type TestPlugin struct{}

func (t *TestPlugin) Name() string {
	return "test-plugin"
}

func (t *TestPlugin) Version() string {
	return "1.0.0"
}

func (t *TestPlugin) DescribeCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "describe",
			Usage: "Describe resources",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}
}

func (t *TestPlugin) GetCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "get",
			Usage: "Get resources",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}
}

func (t *TestPlugin) RegistrationCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "register",
			Usage: "Register resources",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}
}

func (t *TestPlugin) RunCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "run",
			Usage: "Run resources",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}
}

func (t *TestPlugin) RemoveCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "remove",
			Usage: "Remove resources",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}
}

func TestInterfaceCompiles(t *testing.T) {
	var _ EigenRuntimePlugin = (*TestPlugin)(nil)
	
	plugin := &TestPlugin{}
	
	// Verify all methods return expected types
	if name := plugin.Name(); name != "test-plugin" {
		t.Errorf("Expected name 'test-plugin', got %s", name)
	}
	
	if version := plugin.Version(); version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %s", version)
	}
	
	// Verify command getters return non-nil slices
	if cmds := plugin.DescribeCommands(); cmds == nil || len(cmds) == 0 {
		t.Error("DescribeCommands should return non-empty command slice")
	}
	
	if cmds := plugin.GetCommands(); cmds == nil || len(cmds) == 0 {
		t.Error("GetCommands should return non-empty command slice")
	}
	
	if cmds := plugin.RegistrationCommands(); cmds == nil || len(cmds) == 0 {
		t.Error("RegistrationCommands should return non-empty command slice")
	}
	
	if cmds := plugin.RunCommands(); cmds == nil || len(cmds) == 0 {
		t.Error("RunCommands should return non-empty command slice")
	}
	
	if cmds := plugin.RemoveCommands(); cmds == nil || len(cmds) == 0 {
		t.Error("RemoveCommands should return non-empty command slice")
	}
}