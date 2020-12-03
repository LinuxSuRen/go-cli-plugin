package config

import (
	"github.com/linuxsuren/go-cli-plugin/pkg"
	"io"
	"net/http"
)

type (
	configPluginListCmd struct {
		pkg.OutputOption
	}

	jcliPluginFetchCmd struct {
		PluginRepo string
		Reset      bool

		Username   string
		Password   string
		SSHKeyFile string

		output    io.Writer
		PluginOrg string
	}

	jcliPluginInstallCmd struct {
		RoundTripper http.RoundTripper
		ShowProgress bool

		output     io.Writer
		PluginOrg  string
		PluginRepo string
	}

	jcliPluginUninstallCmd struct {
	}

	jcliPluginUpdateCmd struct {
	}
)
