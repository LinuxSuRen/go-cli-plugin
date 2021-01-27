package config

import (
	cobra_ext "github.com/linuxsuren/cobra-extension"
	"io"
	"net/http"
)

type (
	configPluginListCmd struct {
		cobra_ext.OutputOption
	}

	jcliPluginFetchCmd struct {
		PluginRepo string
		Reset      bool

		Username   string
		Password   string
		SSHKeyFile string

		output         io.Writer
		PluginOrg      string
		PluginRepoName string
	}

	jcliPluginInstallCmd struct {
		RoundTripper http.RoundTripper
		ShowProgress bool

		output         io.Writer
		PluginOrg      string
		PluginRepo     string
		PluginRepoName string
	}

	jcliPluginUninstallCmd struct {
		RoundTripper http.RoundTripper
		ShowProgress bool

		output         io.Writer
		PluginOrg      string
		PluginRepo     string
		PluginRepoName string
	}

	jcliPluginUpdateCmd struct {
		RoundTripper http.RoundTripper
		ShowProgress bool

		output         io.Writer
		PluginOrg      string
		PluginRepo     string
		PluginRepoName string
	}
)
