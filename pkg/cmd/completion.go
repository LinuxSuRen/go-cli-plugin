package config

import (
	"github.com/linuxsuren/go-cli-plugin/pkg"
	"github.com/spf13/cobra"
	"strings"
)

// ValidPluginNames returns the valid plugin name list
func ValidPluginNames(cmd *cobra.Command, args []string, prefix, pluginOrg, pluginRepo string) (pluginNames []string, directive cobra.ShellCompDirective) {
	directive = cobra.ShellCompDirectiveNoFileComp
	if plugins, err := pkg.FindPlugins(pluginOrg, pluginRepo); err == nil {
		pluginNames = make([]string, 0)
		for i := range plugins {
			plugin := plugins[i]
			name := plugin.Use

			switch cmd.Use {
			case "install":
				if plugin.Installed {
					continue
				}
			case "uninstall":
				if !plugin.Installed {
					continue
				}
			}

			duplicated := false
			for j := range args {
				if name == args[j] {
					duplicated = true
					break
				}
			}

			if !duplicated && strings.HasPrefix(name, prefix) {
				pluginNames = append(pluginNames, name)
			}
		}
	}
	return
}

// NoFileCompletion avoid completion with files
func NoFileCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}
