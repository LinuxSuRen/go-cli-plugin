package config

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/linuxsuren/go-cli-plugin/pkg"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

// NewConfigPluginInstallCmd create a command for fetching plugin metadata
func NewConfigPluginInstallCmd(pluginOrg, pluginRepo string) (cmd *cobra.Command) {
	pluginInstallCmd := jcliPluginInstallCmd{
		PluginOrg:      pluginOrg,
		PluginRepo:     pluginRepo,
		PluginRepoName: pluginRepo,
	}

	cmd = &cobra.Command{
		Use:   "install",
		Short: "install a jcli plugin",
		Long:  "install a jcli plugin",
		Args:  cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) (strings []string, directive cobra.ShellCompDirective) {
			return ValidPluginNames(cmd, args, toComplete, pluginOrg, pluginRepo)
		},
		RunE: pluginInstallCmd.Run,
	}

	// add flags
	flags := cmd.Flags()
	flags.BoolVarP(&pluginInstallCmd.ShowProgress, "show-progress", "", true,
		"If you want to show the progress of download")
	return
}

// Run main entry point for plugin install command
func (c *jcliPluginInstallCmd) Run(cmd *cobra.Command, args []string) (err error) {
	name := args[0]
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	var data []byte
	pluginsMetadataFile := fmt.Sprintf("%s/.%s/plugins-repo/%s.yaml", userHome, c.PluginRepoName, name)
	if data, err = ioutil.ReadFile(pluginsMetadataFile); err == nil {
		plugin := pkg.Plugin{}
		if err = yaml.Unmarshal(data, &plugin); err == nil {
			err = c.download(plugin)
		}
	}

	if err == nil {
		cachedMetadataFile := pkg.GetJCLIPluginPath(userHome, c.PluginRepoName, name, true)
		err = ioutil.WriteFile(cachedMetadataFile, data, 0664)
	}
	return
}

func (c *jcliPluginInstallCmd) download(plu pkg.Plugin) (err error) {
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	link := c.getDownloadLink(plu)
	output := fmt.Sprintf("%s/.%s/plugins/%s.tar.gz", userHome, c.PluginOrg, plu.Main)

	downloader := pkg.HTTPDownloader{
		RoundTripper:   c.RoundTripper,
		TargetFilePath: output,
		URL:            link,
		ShowProgress:   c.ShowProgress,
	}
	if err = downloader.DownloadFile(); err == nil {
		err = c.extractFiles(plu, output)
	}
	return
}

func (c *jcliPluginInstallCmd) getDownloadLink(plu pkg.Plugin) (link string) {
	link = plu.DownloadLink
	if link == "" {
		operationSystem := runtime.GOOS
		arch := runtime.GOARCH
		link = fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s-%s-%s.tar.gz",
			c.PluginOrg, plu.Main, plu.Version, plu.Main, operationSystem, arch)
	}
	return
}

func (c *jcliPluginInstallCmd) extractFiles(plugin pkg.Plugin, tarFile string) (err error) {
	var f *os.File
	var gzf *gzip.Reader
	if f, err = os.Open(tarFile); err != nil {
		return
	}
	defer f.Close()

	if gzf, err = gzip.NewReader(f); err != nil {
		return
	}

	tarReader := tar.NewReader(gzf)
	var header *tar.Header
	for {
		if header, err = tarReader.Next(); err == io.EOF {
			err = nil
			break
		} else if err != nil {
			break
		}
		name := header.Name

		switch header.Typeflag {
		case tar.TypeReg:
			if name != plugin.Main {
				continue
			}
			var targetFile *os.File
			if targetFile, err = os.OpenFile(fmt.Sprintf("%s/%s", filepath.Dir(tarFile), name),
				os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode)); err != nil {
				break
			}
			if _, err = io.Copy(targetFile, tarReader); err != nil {
				break
			}
			targetFile.Close()
		}
	}
	return
}
