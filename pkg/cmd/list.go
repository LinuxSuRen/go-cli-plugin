package config

import (
	"github.com/linuxsuren/go-cli-plugin/pkg"
	"github.com/spf13/cobra"
)

// NewConfigPluginListCmd create a command for list all jcli plugins
func NewConfigPluginListCmd(pluginOrg, pluginRepo string) (cmd *cobra.Command) {
	configPluginListCmd := &configPluginListCmd{
		PluginOrg:      pluginOrg,
		PluginRepo:     pluginRepo,
		PluginRepoName: pluginRepo,
	}

	cmd = &cobra.Command{
		Use:               "list",
		Short:             "List all installed plugins",
		Long:              "List all installed plugins",
		RunE:              configPluginListCmd.RunE,
		ValidArgsFunction: NoFileCompletion,
	}

	configPluginListCmd.SetFlagWithHeaders(cmd, "Use,Version,Installed,DownloadLink")
	return
}

// RunE is the main entry point of config plugin list command
func (c *configPluginListCmd) RunE(cmd *cobra.Command, args []string) (err error) {
	c.Writer = cmd.OutOrStdout()
	var plugins []pkg.Plugin
	if plugins, err = pkg.FindPlugins(c.PluginOrg, c.PluginRepoName); err == nil {
		err = c.OutputV2(plugins)
	}
	return
}
