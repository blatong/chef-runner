package openssh_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/openssh"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	expect := "OpenSSH (host: some-host)"
	actual := openssh.NewClient("some-host").String()
	assert.Equal(t, expect, actual)
}

var commandTests = []struct {
	client openssh.Client
	cmd    []string
}{
	{
		openssh.Client{
			Host: "some-host",
		},
		[]string{"ssh", "some-host", "uname -a"},
	},
	{
		openssh.Client{
			Host: "some-host",
			User: "some-user",
		},
		[]string{"ssh", "-l", "some-user", "some-host", "uname -a"},
	},
	{
		openssh.Client{
			Host: "some-host",
			Port: 1234,
		},
		[]string{"ssh", "-p", "1234", "some-host", "uname -a"},
	},
	{
		openssh.Client{
			Host:        "some-host",
			PrivateKeys: []string{"some-key", "another-key"},
		},
		[]string{"ssh", "-i", "some-key", "-i", "another-key",
			"some-host", "uname -a"},
	},
	{
		openssh.Client{
			Host: "some-host",
			Options: map[string]string{
				"SomeOption":    "yes",
				"AnotherOption": "no",
			},
		},
		[]string{"ssh", "-o", "AnotherOption=no", "-o", "SomeOption=yes",
			"some-host", "uname -a"},
	},
	{
		openssh.Client{
			Host:        "some-host",
			User:        "some-user",
			Port:        1234,
			PrivateKeys: []string{"some-key"},
			Options:     map[string]string{"SomeOption": "yes"},
		},
		[]string{"ssh", "-l", "some-user", "-p", "1234", "-i", "some-key",
			"-o", "SomeOption=yes", "some-host", "uname -a"},
	},
}

func TestCommand(t *testing.T) {
	for _, test := range commandTests {
		cmd, err := test.client.Command("uname -a")
		if assert.NoError(t, err) {
			assert.Equal(t, test.cmd, cmd)
		}
	}
}

func TestRunCommand_MissingCommand(t *testing.T) {
	err := openssh.Client{Host: "some-host"}.RunCommand("")
	assert.EqualError(t, err, "no command given")
}

func TestRunCommand_MissingHost(t *testing.T) {
	err := openssh.Client{}.RunCommand("uname -a")
	assert.EqualError(t, err, "no host given")
}
