package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

//TODO:参考に(中国語のサイト)
//https://www.haowuliaoa.com/article/info/11350.html

// FundRawTransactionResult fundrawtransactionをcallしたresponseの型
type FundRawTransactionResult struct {
	Hex       string `json:"hex"`
	Fee       int64  `json:"fee"`
	Changepos int64  `json:"changepos"`
}

//Result:
//	{
//		"hex":       "value", (string)  The resulting raw transaction (hex-encoded string)
//		"fee":       n,         (numeric) Fee in BTC the resulting transaction pays
//		"changepos": n          (numeric) The position of the added change output, or -1
//	}

// GetTransactionByTxID txIDからトランザクション詳細を取得する
func (b *Bitcoin) GetTransactionByTxID(txID string) (*btcjson.GetTransactionResult, error) {
	// Transaction詳細を取得(必要な情報があるかどうか不明)
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, errors.Errorf("chainhash.NewHashFromStr(%s): error: %v", txID, err)
	}
	txResult, err := b.client.GetTransaction(hash)
	if err != nil {
		return nil, errors.Errorf("GetTransaction(%s): error: %v", hash, err)
	}

	return txResult, nil
}

// GetTxOutByTxID TxOutを指定したトランザクションIDから取得する
func (b *Bitcoin) GetTxOutByTxID(txID string, index uint32) (*btcjson.GetTxOutResult, error) {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, errors.Errorf("chainhash.NewHashFromStr(%s): error: %v", txID, err)
	}

	// Gettxout
	// txHash *chainhash.Hash, index uint32, mempool bool
	//return b.Client.GetTxOut(hash, 0, false)
	txOutResult, err := b.client.GetTxOut(hash, index, false)
	if err != nil {
		return nil, errors.Errorf("GetTxOut(%s, %d, false): error: %v", hash, index, err)
	}

	return txOutResult, nil
	//log.Printf("TxOut: %v\n", txOut): Output
	//grok.Value(txOut)
	//value *GetTxOutResult = {
	//	BestBlock string = "00000000000000a080461b99935872934fa35bc705453f9f2ad7712b2228e849" 64
	//	Confirmations int64 = 473
	//	Value float64 = 0.65
	//	ScriptPubKey ScriptPubKeyResult = {
	//		Asm string = "OP_HASH160 b24f4d8c8237c73a7299b6e790b309d477bb509c OP_EQUAL" 60
	//		Hex string = "a914b24f4d8c8237c73a7299b6e790b309d477bb509c87" 46
	//		ReqSigs int32 = 1
	//		Type string = "scripthash" 10
	//		Addresses []string = [
	//			0 string = "2N9W3GVam33jQc5FbkLKwMqH7RkvkYK7xvz" 35
	//		]
	//	}
	//	Coinbase bool = false
	//}
}

// ToHex 16進数のstringに変換する
func (b *Bitcoin) ToHex(tx *wire.MsgTx) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return "", errors.Errorf("tx.Serialize(): error: %v", err)
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

//ToMsgTx 16進数のstringから、wire.MsgTxに変換する
func (b *Bitcoin) ToMsgTx(txHex string) (*wire.MsgTx, error) {
	byteHex, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, errors.Errorf("hex.DecodeString(): error: %v", err)
	}

	var msgTx wire.MsgTx
	if err := msgTx.Deserialize(bytes.NewReader(byteHex)); err != nil {
		return nil, err
	}
	//return btcutil.NewTx(&msgTx), nil
	return &msgTx, nil
}

// DecodeRawTransaction Hex stringをデコードして、Rawのトランザクションデータに変換する
func (b *Bitcoin) DecodeRawTransaction(txHex string) (*btcjson.TxRawResult, error) {
	byteHex, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, errors.Errorf("hex.DecodeString(): error: %v", err)
	}
	resTx, err := b.client.DecodeRawTransaction(byteHex)
	if err != nil {
		return nil, errors.Errorf("client.DecodeRawTransaction(): error: %v", err)
	}

	return resTx, nil
}

// GetRawTransactionByHex Hexからトランザクションを取得する
func (b *Bitcoin) GetRawTransactionByHex(strTxHash string) (*btcutil.Tx, error) {

	txHash, err := chainhash.NewHashFromStr(strTxHash)
	if err != nil {
		return nil, errors.Errorf("chainhash.NewHashFromStr(%s): error: %v", strTxHash, err)
	}

	tx, err := b.client.GetRawTransaction(txHash)
	if err != nil {
		return nil, errors.Errorf("GetRawTransaction(hash): error: %v", err)
	}
	//MsgTx()
	//tx.MsgTx()

	return tx, nil
}

// CreateRawTransaction Rawトランザクションを作成する
func (b *Bitcoin) CreateRawTransaction(sendAddr string, amount btcutil.Amount, inputs []btcjson.TransactionInput) (*wire.MsgTx, error) {
	sendAddrDecoded, err := btcutil.DecodeAddress(sendAddr, b.GetChainConf())
	//TODO:sendAddrの厳密なチェックがセキュリティ的に必要な場面もありそう
	//TODO:このタイミングでfeeの計算も必要っぽい
	//TODO:トランザクションのkbに応じて、手数料を算出
	//TODO:でもfeeのパラメータを入れるのは、sendrawtransaction
	if err != nil {
		return nil, errors.Errorf("btcutil.DecodeAddress(%s): error: %v", sendAddr, err)
	}

	//TODO: 手数料を考慮せず、全額送金しようとすると、SendRawTransaction()で、`min relay fee not met`
	//つまり、そのチェックロジックも入れたほうがいいかもしれない
	log.Printf("[Debug] amout:%d, %v", amount, amount) // 1.95 BTC => %v表示だと、単位まで表示される！

	outputs := make(map[btcutil.Address]btcutil.Amount)
	outputs[sendAddrDecoded] = amount //satoshi
	lockTime := int64(0)              //TODO:Raw locktime ここは何をいれるべき？
	msgTx, err := b.client.CreateRawTransaction(inputs, outputs, &lockTime)
	if err != nil {
		return nil, errors.Errorf("btcutil.CreateRawTransaction(): error: %v", err)
	}

	return msgTx, nil
}

// FundRawTransaction 送信したい金額に応じて、自動的にutxoを算出してくれる
// TODO:現時点で使う予定無し
func (b *Bitcoin) FundRawTransaction(hex string) (*FundRawTransactionResult, error) {
	//fundrawtransaction
	//https://bitcoincore.org/en/doc/0.16.2/rpc/rawtransactions/fundrawtransaction/
	//"{\"changePosition\":2}"

	//hex
	bHex, err := json.Marshal(hex)
	if err != nil {
		return nil, errors.Errorf("json.Marchal(hex): error: %v", err)
	}

	//fee rate
	feePerKb, err := b.EstimateSmartFee()
	if err != nil {
		return nil, errors.Errorf("EstimateSmartFee(): error: %v", err)
	}

	bFeeRate, err := json.Marshal(struct {
		FeeRate float64 `json:"feeRate"`
	}{
		FeeRate: feePerKb,
	})
	if err != nil {
		return nil, errors.Errorf("json.Marchal(feeRate): error: %v", err)
	}

	rawResult, err := b.client.RawRequest("fundrawtransaction", []json.RawMessage{bHex, bFeeRate})
	//rawResult, err := b.client.RawRequest("fundrawtransaction", []json.RawMessage{bHex})
	if err != nil {
		//FIXME: error: -4: Insufficient funds
		return nil, errors.Errorf("json.RawRequest(fundrawtransaction): error: %v", err)
	}

	fundRawTransactionResult := FundRawTransactionResult{}
	err = json.Unmarshal([]byte(rawResult), &fundRawTransactionResult)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %v", err)
	}

	log.Printf("[Debug]fundRawTransactionResult: %v: %s\n", fundRawTransactionResult, fundRawTransactionResult.Hex)

	return &fundRawTransactionResult, nil
}

// SignRawTransactionByHex HexからRawトランザクションを生成し、署名する
// オフライン(coldwallet)での利用を想定
func (b *Bitcoin) SignRawTransactionByHex(hex string) (string, bool, error) {
	// Hexからトランザクションを取得
	msgTx, err := b.ToMsgTx(hex)
	if err != nil {
		return "", false, err
	}

	//署名
	signedTx, isSigned, err := b.SignRawTransaction(msgTx)
	if err != nil {
		return "", false, err
	}

	//Hexに変換
	hexTx, err := b.ToHex(signedTx)
	if err != nil {
		return "", false, errors.Errorf("w.Btc.ToHex(msgTx): error: %v", err)
	}

	//return signedTx, nil
	return hexTx, isSigned, nil
}

// SignRawTransaction Rawのトランザクションに署名する
func (b *Bitcoin) SignRawTransaction(tx *wire.MsgTx) (*wire.MsgTx, bool, error) {
	//TODO: It should be implemented on Cold Strage
	//この処理がHotwallet内で動くということは、重要な情報がwallet内に含まれてしまっているということでは？
	msgTx, isSigned, err := b.client.SignRawTransaction(tx)
	if err != nil {
		return nil, false, errors.Errorf("SignRawTransaction(): error: %v", err)
	}
	//if !isSigned {
	//	return nil, errors.New("SignRawTransaction() can not sign on given transaction")
	//}

	return msgTx, isSigned, nil
}

// SendRawTransaction Rawトランザクションを送信する
func (b *Bitcoin) SendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {
	hash, err := b.client.SendRawTransaction(tx, true)
	if err != nil {
		return nil, errors.Errorf("SendRawTransaction(): error: %v", err)
	}

	return hash, nil
}

// SendTransactionByHex 外部から渡されたバイト列からRawトランザクションを送信する
func (b *Bitcoin) SendTransactionByHex(hex string) (*chainhash.Hash, error) {
	// Hexからトランザクションを取得
	msgTx, err := b.ToMsgTx(hex)
	if err != nil {
		return nil, err
	}

	hash, err := b.client.SendRawTransaction(msgTx, true)
	if err != nil {
		return nil, errors.Errorf("SendRawTransaction(): error: %v", err)
	}

	return hash, nil
}

// SendTransactionByByte 外部から渡されたバイト列からRawトランザクションを送信する
func (b *Bitcoin) SendTransactionByByte(rawTx []byte) (*chainhash.Hash, error) {
	//TODO:渡された文字列は暗号化されていることを想定
	//TODO:ここで復号化の処理が必要

	wireTx := new(wire.MsgTx)
	r := bytes.NewBuffer(rawTx)

	if err := wireTx.Deserialize(r); err != nil {
		return nil, errors.Errorf("wireTx.Deserialize(): error: %v", err)
	}

	hash, err := b.client.SendRawTransaction(wireTx, true)
	if err != nil {
		return nil, errors.Errorf("SendRawTransaction(): error: %v", err)
	}

	return hash, nil
}

// SequentialTransaction 検証用: 一連の未署名トランザクション作成から送信までの流れ
//func (b *Bitcoin) SequentialTransaction(tx *wire.MsgTx) error {
func (b *Bitcoin) SequentialTransaction(hex string) (*chainhash.Hash, *btcutil.Tx, error) {
	// Hexからトランザクションを取得
	msgTx, err := b.ToMsgTx(hex)
	if err != nil {
		return nil, nil, err
	}

	//署名(オフライン)
	signedTx, isSigned, err := b.SignRawTransaction(msgTx)
	if err != nil {
		return nil, nil, err
	}
	if !isSigned {
		return nil, nil, errors.New("SignRawTransaction() can not sign on given transaction or multisig may be required")
	}

	//送金(オンライン)
	hash, err := b.SendRawTransaction(signedTx)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("[Debug] txID hash: %s", hash.String())

	//txを取得
	resTx, err := b.GetRawTransactionByHex(hash.String())
	if err != nil {
		return nil, nil, err
	}

	return hash, resTx, nil
}
