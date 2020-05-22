package watchsrv

import (
	"database/sql"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
)

// TxMonitor type
type TxMonitor struct {
	eth          ethgrp.Ethereumer
	logger       *zap.Logger
	dbConn       *sql.DB
	txRepo       watchrepo.ETHTxRepositorier
	txDetailRepo watchrepo.EthDetailTxRepositorier
	payReqRepo   watchrepo.PaymentRequestRepositorier
	confirmNum   uint64
	wtype        wallet.WalletType
}

// NewTxMonitor returns TxMonitor object
func NewTxMonitor(
	eth ethgrp.Ethereumer,
	logger *zap.Logger,
	dbConn *sql.DB,
	txRepo watchrepo.ETHTxRepositorier,
	txDetailRepo watchrepo.EthDetailTxRepositorier,
	payReqRepo watchrepo.PaymentRequestRepositorier,
	confirmNum uint64,
	wtype wallet.WalletType) *TxMonitor {

	return &TxMonitor{
		eth:          eth,
		logger:       logger,
		dbConn:       dbConn,
		txRepo:       txRepo,
		txDetailRepo: txDetailRepo,
		payReqRepo:   payReqRepo,
		confirmNum:   confirmNum,
		wtype:        wtype,
	}
}

// TODO: implementation

// UpdateTxStatus update transaction status
// - monitor transaction whose tx_type=3(TxTypeSent) in tx_payment/tx_deposit/tx_transfer
func (t *TxMonitor) UpdateTxStatus() error {
	//TODO: as possibility tx_type is not updated from `done`

	//types := []action.ActionType{
	//	action.ActionTypeDeposit,
	//	action.ActionTypePayment,
	//	action.ActionTypeTransfer,
	//}

	//1. update tx_type for TxTypeSent
	err := t.updateStatusTxTypeSent()
	if err != nil {
		return errors.Wrap(err, "fail to call updateStatusTxTypeSent()")
	}

	//2. update tx_type for TxTypeDone
	// - TODO: notification
	//for _, actionType := range types {
	//	err := t.updateStatusTxTypeDone(actionType)
	//	if err != nil {
	//		return errors.Wrapf(err, "fail to call updateStatusTxTypeDone() ActionType: %s", actionType)
	//	}
	//}

	return nil
}

// update TxTypeSent to TxTypeDone if confirmation is 6 or more
func (t *TxMonitor) updateStatusTxTypeSent() error {
	// get records whose status is TxTypeSent
	hashes, err := t.txDetailRepo.GetSentHashTx(tx.TxTypeSent)
	if err != nil {
		return errors.Wrap(err, "fail to call txRepo.GetSentHashTx(TxTypeSent)")
	}

	// get hash in detail and check confirmation
	for _, sentHash := range hashes {
		// check confirmation
		confirmNum, err := t.eth.GetConfirmation(sentHash)
		if err != nil {
			return errors.Wrapf(err, "fail to call eth.GetConfirmation() sentHash: %s", sentHash)
		}
		t.logger.Info("confirmation",
			zap.String("sentHash", sentHash),
			zap.Uint64("confirmation num", confirmNum))
		if confirmNum < t.confirmNum {
			continue
		}
		// update status
		t.txDetailRepo.UpdateTxTypeBySentHashTx(tx.TxTypeDone, sentHash)
	}
	return nil
}
