package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

//create subcommand
type CreateCommand struct {
	Name        string
	Version     string
	SynopsisExp string
	UI          cli.Ui
	Wallet      wallet.Keygener
}

func (c *CreateCommand) Synopsis() string {
	return c.SynopsisExp
}

var (
	keySynopsis      = "create one key for debug use"
	hdkeySynopsis    = "create HD key"
	seedSynopsis     = "create seed"
	multisigSynopsis = "multisig address for debug use"
)

func (c *CreateCommand) Help() string {
	return fmt.Sprintf(`Usage: keygen create [Subcommands...]
Subcommands:
  key       %s
  hdkey     %s
  seed      %s
  multisig  %s
`, keySynopsis, keySynopsis, seedSynopsis, multisigSynopsis)
}

func (c *CreateCommand) Run(args []string) int {
	c.UI.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"key": func() (cli.Command, error) {
			return &HDKeyCommand{
				name:     "key",
				synopsis: keySynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"hdkey": func() (cli.Command, error) {
			return &HDKeyCommand{
				name:     "hdkey",
				synopsis: hdkeySynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"seed": func() (cli.Command, error) {
			return &SeedCommand{
				name:     "seed",
				synopsis: seedSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"multisig": func() (cli.Command, error) {
			return &MultisigCommand{
				name:     "multisig",
				synopsis: multisigSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
	}
	cl := command.CreateSubCommand(c.Name, c.Version, args, cmds)

	code, err := cl.Run()
	if err != nil {
		c.UI.Error(fmt.Sprintf("fail to call Run() subcommand of %s: %v", c.Name, err))
	}
	return code
}
