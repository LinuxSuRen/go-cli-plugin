package config

import (
	"github.com/spf13/cobra"
)

// AppendPluginCmd create a command as root of config plugin
func AppendPluginCmd(root *cobra.Command, pluginOrg, pluginRepo string) {
	root.AddCommand(NewConfigPluginListCmd(pluginOrg, pluginRepo),
		NewConfigPluginFetchCmd(pluginOrg, pluginRepo),
		NewConfigPluginInstallCmd(pluginOrg, pluginRepo),
		NewConfigPluginUninstallCmd(pluginOrg, pluginRepo))
}
