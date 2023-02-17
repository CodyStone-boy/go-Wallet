package watchsrv

import (
	"database/sql"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
)

// TxCreator is TxCreator interface
type TxCreator interface {
	CreateDepositTx() (string, string, error)
	CreatePaymentTx() (string, string, error)
	CreateTransferTx(sender, receiver account.AccountType, floatAmount float64) (string, string, error)
}

// TxCreate type
type TxCreate struct {
	xrp             xrpgrp.Rippler
	logger          *zap.Logger
	dbConn          *sql.DB
	addrRepo        watchrepo.AddressRepositorier
	txRepo          watchrepo.TxRepositorier
	txDetailRepo    watchrepo.XrpDetailTxRepositorier
	payReqRepo      watchrepo.PaymentRequestRepositorier
	txFileRepo      tx.FileRepositorier
	depositReceiver account.AccountType
	paymentSender   account.AccountType
	wtype           wallet.WalletType
}

// NewTxCreate returns TxCreate object
func NewTxCreate(
	xrp xrpgrp.Rippler,
	logger *zap.Logger,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	txRepo watchrepo.TxRepositorier,
	txDetailRepo watchrepo.XrpDetailTxRepositorier,
	payReqRepo watchrepo.PaymentRequestRepositorier,
	txFileRepo tx.FileRepositorier,
	depositReceiver account.AccountType,
	paymentSender account.AccountType,
	wtype wallet.WalletType,
) *TxCreate {
	return &TxCreate{
		xrp:             xrp,
		logger:          logger,
		dbConn:          dbConn,
		addrRepo:        addrRepo,
		txRepo:          txRepo,
		txDetailRepo:    txDetailRepo,
		payReqRepo:      payReqRepo,
		txFileRepo:      txFileRepo,
		depositReceiver: depositReceiver,
		paymentSender:   paymentSender,
		wtype:           wtype,
	}
}

// TODO: it can be commonized to ./pkg/wallet/service/eth/watchsrv/tx_creator.go
func (t *TxCreate) updateDB(
	targetAction action.ActionType,
	txDetailItems []*models.XRPDetailTX,
	paymentRequestIds []int64,
) (int64, error) {
	// start transaction
	dtx, err := t.dbConn.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "fail to start transaction")
	}
	defer func() {
		if err != nil {
			// nolint:errcheck
			dtx.Rollback()
		} else {
			// nolint:errcheck
			dtx.Commit()
		}
	}()

	// Insert tx
	txID, err := t.txRepo.InsertUnsignedTx(targetAction)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call txRepo.InsertUnsignedTx()")
	}
	// Insert to eth_detail_tx
	for idx := range txDetailItems {
		txDetailItems[idx].TXID = txID
	}
	if err = t.txDetailRepo.InsertBulk(txDetailItems); err != nil {
		return 0, errors.Wrap(err, "fail to call txDetailRepo.InsertBulk()")
	}

	if targetAction == action.ActionTypePayment {
		_, err = t.payReqRepo.UpdatePaymentID(txID, paymentRequestIds)
		if err != nil {
			return 0, errors.Wrap(err, "fail to call repo.PayReq().UpdatePaymentID(txID, paymentRequestIds)")
		}
	}
	return txID, nil
}

// generateHexFile generate file for hex txID and encoded previous addresses
// TODO: it can be commonized to ./pkg/wallet/service/eth/watchsrv/tx_creator.go
func (t *TxCreate) generateHexFile(actionType action.ActionType, senderAccount account.AccountType, txID int64, serializedTxs []string) (string, error) {
	// add senderAccount to first line
	serializedTxs = append([]string{senderAccount.String()}, serializedTxs...)

	// create file
	path := t.txFileRepo.CreateFilePath(actionType, tx.TxTypeUnsigned, txID, 0)
	generatedFileName, err := t.txFileRepo.WriteFileSlice(path, serializedTxs)
	if err != nil {
		return "", errors.Wrap(err, "fail to call txFileRepo.WriteFile()")
	}

	return generatedFileName, nil
}
