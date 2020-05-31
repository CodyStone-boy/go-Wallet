package xrp

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	pb "github.com/hiromaily/ripple-lib-proto/pb/go/rippleapi"
)

// TxJSON is transaction json type
type TxJSON struct {
	TransactionType    string `json:"TransactionType"`
	Account            string `json:"Account"`
	Amount             string `json:"Amount"`
	Destination        string `json:"Destination"`
	Flags              uint64 `json:"Flags"`
	LastLedgerSequence uint64 `json:"LastLedgerSequence"`
	Fee                string `json:"Fee"`
	Sequence           uint64 `json:"Sequence"`
}

// PrepareTransaction calls PrepareTransaction API
func (r *Ripple) PrepareTransaction(senderAccount, receiverAccount string, amount float64) (*TxJSON, error) {

	ctx := context.Background()
	req := &pb.RequestPrepareTransaction{
		TxType:          pb.TX_PAYMENT,
		SenderAccount:   senderAccount,
		Amount:          amount,
		ReceiverAccount: receiverAccount,
		Instructions:    &pb.Instructions{MaxLedgerVersionOffset: 75},
	}

	//res: *pb.ResponsePrepareTransaction
	res, err := r.API.client.PrepareTransaction(ctx, req)

	if err != nil {
		return nil, errors.Wrap(err, "fail to call client.PrepareTransaction()")
	}
	r.logger.Debug("response",
		zap.String("TxJSON", res.TxJSON),
		zap.Any("Instructions", res.Instructions),
	)

	var txJSON TxJSON
	unquotedJSON, _ := strconv.Unquote(res.TxJSON)
	if err = json.Unmarshal([]byte(unquotedJSON), &txJSON); err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal(txJSON)")
	}

	return &txJSON, nil
}
