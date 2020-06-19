package db

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/command"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// DBCommand db subcommand
type DBCommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Watcher
}

// Synopsis is explanation for this subcommand
func (c *DBCommand) Synopsis() string {
	return "Database functionality"
}

var (
	createSynopsis = "create table"
)

// Help returns usage for this subcommand
func (c *DBCommand) Help() string {
	return fmt.Sprintf(`Usage: wallet db [Subcommands...]
Subcommands:
  create  %s
`, createSynopsis)
}

// Run executes this subcommand
func (c *DBCommand) Run(args []string) int {
	c.UI.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"create": func() (cli.Command, error) {
			return &CreateCommand{
				name:     "create",
				synopsis: createSynopsis,
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