package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

// multisig subcommand
type MultisigCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallet.Keygener
}

func (c *MultisigCommand) Synopsis() string {
	return c.synopsis
}

func (c *MultisigCommand) Help() string {
	return `Usage: keygen key create multisig [options...]
Options:
  -table  target table name
`
}

func (c *MultisigCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	var (
		tableName string
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&tableName, "table", "", "table name of database")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	c.ui.Output(fmt.Sprintf("-table: %s", tableName))

	////validator
	//if tableName == "" {
	//	tableName = "payment_request"
	//	//c.ui.Error("table name option [-table] is required")
	//	//return 1
	//}
	//
	////create payment_request table
	//err := testdata.CreateInitialTestData(c.wallet.GetDB(), c.wallet.GetBTC())
	//if err != nil {
	//	c.ui.Error(fmt.Sprintf("fail to call testdata.CreateInitialTestData() %+v", err))
	//	return 1
	//}
	//c.ui.Info("Done!")

	return 0
}
