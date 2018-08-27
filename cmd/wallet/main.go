package main

import (
	"log"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/chaincfg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/service"
	"github.com/jessevdk/go-flags"
)

//こちらはHotwallet、ただし、Watch Only Walletとしての機能を実装していく。
//ネットワークへの接続はGCP上のBitcoin Core
//Watch Only Walletとしてのセットアップが必要
// - Cold Wallet側から生成したPublic Key をMultisigアドレス変換後、`importaddress xxxxx`でimportする
//   これがかなり時間がかかる。。。実運用ではどうすべきか。rescanしなくても最初はOKかと

//TODO:coldwallet側(非ネットワーク環境)側の機能と明確に分ける
//TODO:オフラインで可能機能と、不可能な機能の切り分けが必要
//TODO:ウォレットの定期バックアップ機能 + import機能
//TODO:coldウォレットへのデータ移行機能が必要なはず
//TODO:multisigの実装
//TODO:生成したkeyの暗号化処理のpkgが必要になるはず
//TODO:入金時にMultisigでの送金は不要な気がする

// Options コマンドラインオプション
type Options struct {
	//Configパス
	ConfPath string `short:"c" long:"conf" default:"./data/toml/watch_config.toml" description:"Path for configuration toml file"`
	//実行される機能
	Mode uint8 `short:"m" long:"mode" description:"Mode i.e.Functionality"`
	//txファイルパス
	ImportFile string `short:"i" long:"import" default:"" description:"import file path for hex"`
	//調整fee
	Fee float64 `short:"f" long:"fee" default:"" description:"adjustment fee"`
	//Debugモード
	Debug bool `short:"d" long:"debug" description:"for only development use"`
}

var (
	opts      Options
	chainConf *chaincfg.Params
)

func init() {
	if _, err := flags.Parse(&opts); err != nil {
		panic(err)
	}
}

func main() {
	//initialSettings()
	wallet, err := service.InitialSettings(opts.ConfPath)
	if err != nil {
		logger.Fatal(err)
	}
	defer wallet.Done()

	if opts.Debug {
		//debug用 機能確認
		debugForCheck(wallet)
	} else {
		//switch mode
		switchFunction(wallet)
	}
}

// 実運用上利用するもののみ、こちらに定義する
func switchFunction(wallet *service.Wallet) {
	//TODO:ここから呼び出すべきはService系のみに統一したい
	switch opts.Mode {
	case 1:
		logger.Info("Run: 入金処理検知")
		//実際には署名処理は手動なので、ユーザーの任意のタイミングで走らせたほうがいい。
		//入金検知 + 未署名トランザクション作成
		hex, fileName, err := wallet.DetectReceivedCoin(opts.Fee)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		if hex == "" {
			logger.Info("No utxo")
			return
		}
		logger.Infof("[hex]: %s\n[fileName]: %s", hex, fileName)

	case 2:
		logger.Info("Run:出金のための未署名トランザクション作成")
		hex, fileName, err := wallet.CreateUnsignedTransactionForPayment(opts.Fee)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		if hex == "" {
			logger.Info("No utxo")
			return
		}
		logger.Infof("[hex]: %s, \n[fileName]: %s", hex, fileName)

	case 3:
		logger.Info("Run: ファイルから署名済みtxを送信する")
		// 1.GPSにupload(web管理画面から行う??)
		// 2.Uploadされたtransactionファイルから、送信する？
		if opts.ImportFile == "" {
			logger.Fatal("file path is required as argument file when running")
		}
		// フルパスを指定する
		txID, err := wallet.SendFromFile(opts.ImportFile)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("[Done]送信までDONE!! txID: %s", txID)

	case 10:
		logger.Info("Run: 送信済ステータスのトランザクションを監視する")
		err := wallet.UpdateStatus()
		if err != nil {
			logger.Fatalf("%+v", err)
		}

	case 20:
		logger.Info("Run: [Debug用]入金から送金までの一連の流れを確認")

		//入金検知 + 未署名トランザクション作成
		logger.Info("[1]Run: 入金検知")
		hex, fileName, err := wallet.DetectReceivedCoin(opts.Fee)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		if hex == "" {
			logger.Info("No utxo")
			return
		}
		logger.Infof("[hex]: %s\n[fileName]: %s", hex, fileName)

		//署名(本来はColdWalletの機能)
		logger.Info("\n[2]Run: 署名")
		hexTx, isSigned, generatedFileName, err := wallet.SignatureFromFile(fileName)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("[hex]: %s\n[署名完了]: %t\n[fileName]: %s", hexTx, isSigned, generatedFileName)

		//送信
		logger.Info("\n[3]Run: 送信")
		txID, err := wallet.SendFromFile(generatedFileName)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("[Done]送信までDONE!! txID: %s", txID)

		//一連の署名から送信までの流れをチェック
		//[WIF] cUW7ZSF9WX7FUTeHkuw5L9Rj26V5Kz8yCkYjZamyvATTwsu7KUCi - [Pub Address] muVSWToBoNWusjLCbxcQNBWTmPjioRLpaA
		//hash, tx, err := wallet.BTC.SequentialTransaction(hex)
		//if err != nil {
		//	log.Fatalf("%+v", err)
		//}
		////tx.MsgTx()
		//log.Printf("[Debug] 送信までDONE!! %s, %v", hash.String(), tx)

	case 21:
		logger.Info("Run: [Debug用]出金から送金までの一連の流れを確認")

		//出金準備
		logger.Info("[1]Run:出金のための未署名トランザクション作成")
		hex, fileName, err := wallet.CreateUnsignedTransactionForPayment(opts.Fee)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		if hex == "" {
			logger.Info("No utxo")
			return
		}
		logger.Infof("[hex]: %s, \n[fileName]: %s", hex, fileName)

		//署名(本来はColdWalletの機能)
		logger.Info("\n[2]Run: 署名")
		hexTx, isSigned, generatedFileName, err := wallet.SignatureFromFile(fileName)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("[hex]: %s\n[署名完了]: %t\n[fileName]: %s", hexTx, isSigned, generatedFileName)

		//送信
		logger.Info("\n[3]Run: 送信")
		txID, err := wallet.SendFromFile(generatedFileName)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("[Done]送信までDONE!! txID: %s", txID)

	default:
		logger.Info("該当Mode無し")
	}

}

// 検証用
func debugForCheck(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		//[Debug用]入金検知処理後、lock解除を行う
		log.Print("Run: lockされたトランザクションの解除")
		err := wallet.BTC.UnlockAllUnspentTransaction()
		if err != nil {
			log.Fatalf("%+v", err)
		}
	case 2:
		//[Debug用]手数料算出
		log.Print("Run: 手数料算出 estimatesmartfee")
		feePerKb, err := wallet.BTC.EstimateSmartFee()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		log.Printf("Estimatesmartfee: %f\n", feePerKb)
	case 3:
		//[Debug用]ロギング
		log.Print("Run: ロギング logging")
		logData, err := wallet.BTC.Logging()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		//Debug
		grok.Value(logData)
	case 4:
		//[Debug用]getnetworkinfoの呼び出し
		log.Print("Run: INFO getnetworkinfo")
		infoData, err := wallet.BTC.GetNetworkInfo()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		//Debug
		grok.Value(infoData)
		log.Printf("%f", infoData.Relayfee)

	case 5:
		//[Debug用]ValidateAddress
		log.Print("Run: AddressのValidationチェック")
		err := wallet.BTC.ValidateAddress("2NFXSXxw8Fa6P6CSovkdjXE6UF4hupcTHtr")
		if err != nil {
			log.Fatalf("%+v", err)
		}
		err = wallet.BTC.ValidateAddress("4VHGkbQTGg2vN5P6yHZw3UJhmsBh9igsSos")
		if err == nil {
			log.Fatal("something is wrong")
		}
		log.Print("Done!")

	case 10:
		//[Debug用]hexから署名済みtxを送信する
		log.Print("Run: hexから署名済みtxを送信する")

		hex := "020000000001019dcbbda4e5233051f2bed587c1d48e8e17aa21c2c3012097899bda5097ce78e201000000232200208e1343e11e4def66d7102d9b0f36f019188118df5a5f30dacdd1008928b12f5fffffffff01042bbf070000000017a9148191d41a7415a6a1f6ee14337e039f50b949e80e870400483045022100f4975a5ea23e5799b1df65d699f85236b9d00bcda8da333731ffa508285d3c59022037285857821ee68cbe5f74239299170686b108ce44e724a9a280a3ef9291746901483045022100f94ce83946b4698b8dfbb7cb75eece12932c5097017e70e60d924aeae1ec829a02206e7b2437e9747a9c28a3a3d7291ea16db1d2f0a60482cdb8eca91c28c01aba790147522103d69e07dbf6da065e6fae1ef5761d029b9ff9143e75d579ffc439d47484044bed2103748797877523b8b36add26c9e0fb6a023f05083dd4056aedc658d2932df1eb6052ae00000000"
		hash, err := wallet.BTC.SendTransactionByHex(hex)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		log.Printf("[Debug] 送信までDONE!! %s", hash.String())

	case 11:
		//[Debug用]payment_requestテーブルの情報を初期化する
		log.Print("Run: payment_requestテーブルの情報を初期化する")
		_, err := wallet.DB.ResetAnyFlagOnPaymentRequestForTestOnly(nil, true)
		if err != nil {
			log.Fatalf("%+v", err)
		}

	default:
		log.Print("該当Mode無し")
	}

}
