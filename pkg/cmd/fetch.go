package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io/ioutil"
	"os"
	"strings"
)

// NewConfigPluginFetchCmd create a command for fetching plugin metadata
func NewConfigPluginFetchCmd(pluginOrg, pluginRepo string) (cmd *cobra.Command) {
	pluginFetchCmd := &jcliPluginFetchCmd{
		PluginOrg:      pluginOrg,
		PluginRepoName: pluginRepo,
	}

	cmd = &cobra.Command{
		Use:   "fetch",
		Short: "fetch metadata of plugins",
		Long: fmt.Sprintf(`fetch metadata of plugins
The official metadata git repository is https://github.com/%s/%s,
but you can change it by giving a command parameter.`, pluginFetchCmd.PluginOrg, pluginFetchCmd.PluginRepo),
		ValidArgsFunction: NoFileCompletion,
		RunE:              pluginFetchCmd.Run,
	}

	// add flags
	flags := cmd.Flags()
	flags.StringVarP(&pluginFetchCmd.PluginRepo, "plugin-repo", "",
		fmt.Sprintf("https://github.com/%s/%s", pluginFetchCmd.PluginOrg, pluginFetchCmd.PluginRepoName),
		"The plugin git repository URL")
	flags.BoolVarP(&pluginFetchCmd.Reset, "reset", "", true,
		"If you want to reset the git local repo when pulling it")
	flags.StringVarP(&pluginFetchCmd.Username, "username", "u", "",
		"The username of git repository")
	flags.StringVarP(&pluginFetchCmd.Password, "password", "p", "",
		"The password of git repository")

	sshKeyFile := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
	flags.StringVarP(&pluginFetchCmd.SSHKeyFile, "ssh-key-file", "", sshKeyFile,
		"SSH key file")
	return
}

// Run is the main entry point of plugin fetch command
func (c *jcliPluginFetchCmd) Run(cmd *cobra.Command, args []string) (err error) {
	var userHome string
	if userHome, err = homedir.Dir(); err != nil {
		return
	}

	pluginRepoCache := fmt.Sprintf("%s/.%s/plugins-repo", c.PluginRepo, userHome)
	c.output = cmd.OutOrStdout()

	cmd.Printf("start to fetch metadata from '%s', cache to '%s'\n", c.PluginRepo, pluginRepoCache)
	var r *git.Repository
	if r, err = git.PlainOpen(pluginRepoCache); err == nil {
		var w *git.Worktree
		if w, err = r.Worktree(); err != nil {
			return
		}

		if c.Reset {
			if err = w.Reset(&git.ResetOptions{
				Mode: git.HardReset,
			}); err != nil {
				return
			}
		}

		err = w.Pull(c.getPullOptions())
		if err == git.NoErrAlreadyUpToDate {
			err = nil // consider it's ok
		}
	} else {
		cloneOptions := c.getCloneOptions()
		_, err = git.PlainClone(pluginRepoCache, false, cloneOptions)
	}
	return
}

func (c *jcliPluginFetchCmd) getCloneOptions() (cloneOptions *git.CloneOptions) {
	cloneOptions = &git.CloneOptions{
		URL:      c.PluginRepo,
		Progress: c.output,
		Auth:     c.getAuth(),
	}
	return
}

func (c *jcliPluginFetchCmd) getPullOptions() (pullOptions *git.PullOptions) {
	pullOptions = &git.PullOptions{
		RemoteName: "origin",
		Progress:   c.output,
		Auth:       c.getAuth(),
	}
	return
}

func (c *jcliPluginFetchCmd) getAuth() (auth transport.AuthMethod) {
	if c.Username != "" {
		auth = &githttp.BasicAuth{
			Username: c.Username,
			Password: c.Password,
		}
	}

	if strings.HasPrefix(c.PluginRepo, "git@") {
		if sshKey, err := ioutil.ReadFile(c.SSHKeyFile); err == nil {
			signer, _ := ssh.ParsePrivateKey(sshKey)
			auth = &gitssh.PublicKeys{User: "git", Signer: signer}
		}
	}
	return
}
