package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/command"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/keygen"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// keygen wallet as cold wallet
//  - generate key and seed for accounts
//  - create multisig address with full pubkey of auth accounts
//  - sing on unsigned transaction as first signature
//   (signature would not be completed if address is multisig)

// TODO: bitcoin functionalities
// - encrypt wallet itself by `encryptwallet` command
// - passphrase would be required when using secret key to sign unsigned transaction

var (
	walletType = wallet.WalletTypeKeyGen
	appName    = walletType.String()
	appVersion = "5.0.0"
)

func main() {
	// command line
	var (
		confPath        string
		accountConfPath string
		btcWallet       string
		coinTypeCode    string
		isHelp          bool
		isVersion       bool
		walleter        wallets.Keygener
	)
	flags := flag.NewFlagSet("main", flag.ContinueOnError)
	flags.StringVar(&confPath, "conf", "", "config file path")
	flags.StringVar(&coinTypeCode, "coin", "btc", "coin type code `btc`, `bch`, `eth`, `xrp`, `hyt`")
	flags.StringVar(&btcWallet, "wallet", "", "specify wallet.dat in bitcoin core")
	flags.BoolVar(&isVersion, "version", false, "show version")
	flags.BoolVar(&isHelp, "help", false, "show help")
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	// version
	if isVersion {
		fmt.Printf("%s v%s\n", appName, appVersion)
		os.Exit(0)
	}

	// validate coinTypeCode
	if !coin.IsCoinTypeCode(coinTypeCode) {
		log.Fatal("coin args is invalid. `btc`, `bch`, `eth`, `xrp`, `hyt` is allowed")
	}

	// set config path if environment variable is existing
	if confPath == "" {
		switch {
		case coinTypeCode == coin.BTC.String():
			confPath = os.Getenv("BTC_KEYGEN_WALLET_CONF")
		case coinTypeCode == coin.BCH.String():
			confPath = os.Getenv("BCH_KEYGEN_WALLET_CONF")
		case coin.IsETHGroup(coin.CoinTypeCode(coinTypeCode)):
			confPath = os.Getenv("ETH_KEYGEN_WALLET_CONF")
		case coinTypeCode == coin.XRP.String():
			confPath = os.Getenv("XRP_KEYGEN_WALLET_CONF")
		}
	}
	// account conf path for account settings
	if accountConfPath == "" {
		switch coinTypeCode {
		case coin.BTC.String():
			accountConfPath = os.Getenv("BTC_ACCOUNT_CONF")
		case coin.BCH.String():
			accountConfPath = os.Getenv("BCH_ACCOUNT_CONF")
		}
	}

	// help
	if !isHelp && len(os.Args) > 1 {
		// base config
		conf, err := config.NewWallet(confPath, walletType, coin.CoinTypeCode(coinTypeCode))
		if err != nil {
			log.Fatal(err)
		}
		// account config
		accountConf := &account.AccountRoot{}
		if accountConfPath != "" {
			accountConf, err = account.NewAccount(accountConfPath)
			if err != nil {
				log.Fatal(err)
			}
		}

		// override config
		conf.CoinTypeCode = coin.CoinTypeCode(coinTypeCode)

		// override conf.Bitcoin.Host
		if btcWallet != "" {
			conf.Bitcoin.Host = fmt.Sprintf("%s/wallet/%s", conf.Bitcoin.Host, btcWallet)
			log.Println("conf.Bitcoin.Host:", conf.Bitcoin.Host)
		}

		// create wallet
		regi := NewRegistry(conf, accountConf, walletType)
		walleter = regi.NewKeygener()
	}
	defer func() {
		walleter.Done()
	}()

	// sub command
	args := flags.Args()
	cmds := keygen.WalletSubCommands(walleter, appVersion)
	cl := command.CreateSubCommand(appName, appVersion, args, cmds)
	cl.HelpFunc = command.HelpFunc(cl.Name)

	flags.Usage = func() { fmt.Println(cl.HelpFunc(cl.Commands)) }

	code, err := cl.Run()
	if err != nil {
		log.Printf("fail to call Run() %s command: %v", appName, err)
	}
	os.Exit(code)
}
