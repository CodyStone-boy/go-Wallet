package add

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// AddCommand add subcommand
type AddCommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Signer
}

// Synopsis
func (c *AddCommand) Synopsis() string {
	return "add key for multisig address"
}

var (
	multisigSynopsis = "call `addmultisigaddress` which adds a P2SH multisig address to the wallet"
)

// Help
func (c *AddCommand) Help() string {
	return fmt.Sprintf(`Usage: sign add [Subcommands...]
Subcommands:
  multisig     %s
`, multisigSynopsis)
}

// RUn
func (c *AddCommand) Run(args []string) int {
	c.UI.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
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
