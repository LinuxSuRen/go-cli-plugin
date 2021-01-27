package config

import (
	"github.com/linuxsuren/go-cli-plugin/pkg"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// NewConfigPluginUninstallCmd create a command to uninstall a plugin
func NewConfigPluginUninstallCmd(pluginOrg, pluginRepo string) (cmd *cobra.Command) {
	jcliPluginUninstallCmd := jcliPluginUninstallCmd{
		PluginOrg:      pluginOrg,
		PluginRepo:     pluginRepo,
		PluginRepoName: pluginRepo,
	}

	cmd = &cobra.Command{
		Use:               "uninstall",
		Short:             "Remove a plugin",
		Long:              "Remove a plugin",
		Args:              cobra.MinimumNArgs(1),
		RunE:              jcliPluginUninstallCmd.RunE,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) (strings []string, directive cobra.ShellCompDirective) {
			return ValidPluginNames(cmd, args, toComplete, pluginOrg, pluginRepo)
		},
	}
	return
}

// RunE is the main entry point of command
func (c *jcliPluginUninstallCmd) RunE(cmd *cobra.Command, args []string) (err error) {
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	name := args[0]
	cachedMetadataFile := pkg.GetJCLIPluginPath(userHome, c.PluginRepoName, name, false)

	var data []byte
	if data, err = ioutil.ReadFile(cachedMetadataFile); err == nil {
		plugin := &pkg.Plugin{}
		if err = yaml.Unmarshal(data, plugin); err == nil {
			mainFile := pkg.GetJCLIPluginPath(userHome, c.PluginRepoName, plugin.Main, true)

			os.Remove(cachedMetadataFile)
			os.Remove(mainFile)
		}
	} else if os.IsNotExist(err) {
		err = nil
		cmd.Printf("plugin \"%s\" does not exists\n", name)
	}
	return
}
