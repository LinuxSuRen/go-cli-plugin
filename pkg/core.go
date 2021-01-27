package pkg

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type Plugin struct {
	Use          string `yaml:"use"`
	Short        string
	Long         string
	Main         string
	Version      string
	DownloadLink string `yaml:"downloadLink"`
	Installed    bool
}

type PluginError struct {
	error
	code int
}

func FindPlugins(pluginOrg, pluginRepo string) (plugins []Plugin, err error) {
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	plugins = make([]Plugin, 0)
	pluginsDir := fmt.Sprintf("%s/.%s/plugins-repo/*.yaml", pluginRepo, userHome)
	//fmt.Println("start to parse plugin file from dir", pluginsDir)
	var files []string
	if files, err = filepath.Glob(pluginsDir); err == nil {
		for _, metaFile := range files {
			var data []byte
			plugin := Plugin{}
			if data, err = ioutil.ReadFile(metaFile); err == nil {
				if err = yaml.Unmarshal(data, &plugin); err != nil {
					fmt.Println(err)
				} else {
					if plugin.Main == "" {
						plugin.Main = fmt.Sprintf("jcli-%s-Plugin", plugin.Use)
					}

					if _, fileErr := os.Stat(GetJCLIPluginPath(userHome, pluginRepo, plugin.Main, true)); !os.IsNotExist(fileErr) {
						plugin.Installed = true
					}
					plugins = append(plugins, plugin)
				}
			} else {
				fmt.Println("failed to parse file", metaFile)
			}
		}
	}
	return
}

// LoadPlugins loads the plugins
func LoadPlugins(cmd *cobra.Command, pluginOrg, pluginRepo string) {
	var plugins []Plugin
	var err error
	if plugins, err = FindPlugins(pluginOrg, pluginRepo); err != nil {
		cmd.PrintErrln("Cannot load plugins successfully")
		return
	}
	//cmd.Println("found plugins, count", len(plugins), plugins)

	for _, plugin := range plugins {
		if !plugin.Installed {
			continue
		}

		// This function is used to setup the environment for the Plugin and then
		// call the executable specified by the parameter 'main'
		callPluginExecutable := func(cmd *cobra.Command, main string, argv []string, out io.Writer) error {
			env := os.Environ()

			prog := exec.Command(main, argv...)
			prog.Env = env
			prog.Stdin = os.Stdin
			prog.Stdout = out
			prog.Stderr = os.Stderr
			if err := prog.Run(); err != nil {
				if eerr, ok := err.(*exec.ExitError); ok {
					os.Stderr.Write(eerr.Stderr)
					status := eerr.Sys().(syscall.WaitStatus)
					return PluginError{
						error: errors.Errorf("Plugin %s exited with error", main),
						code:  status.ExitStatus(),
					}
				}
				return err
			}

			return nil
		}

		//cmd.Println("register Plugin name", Plugin.Use)
		c := &cobra.Command{
			Use:   plugin.Use,
			Short: plugin.Short,
			Long:  plugin.Long,
			Annotations: map[string]string{
				"main": plugin.Main,
			},
			RunE: func(cmd *cobra.Command, args []string) (err error) {
				var userHome string
				if userHome, err = homedir.Dir(); err != nil {
					return
				}

				pluginExec := GetJCLIPluginPath(userHome, pluginRepo, cmd.Annotations["main"], true)
				err = callPluginExecutable(cmd, pluginExec, args, cmd.OutOrStdout())
				return
			},
		}
		cmd.AddCommand(c)
	}
}

// GetJCLIPluginPath returns the path of a jcli plugin
func GetJCLIPluginPath(userHome, pluginRepoName, name string, binary bool) string {
	suffix := ".yaml"
	if binary {
		suffix = ""
	}
	return fmt.Sprintf("%s/.%s/plugins-repo/%s%s", userHome, pluginRepoName, name, suffix)
}
