package wallet

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/action"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
)

// UpdateTxStatus update transaction status
// - monitor transaction whose tx_type=3(TxTypeSent) in tx_payment/tx_receipt/tx_transfer
func (w *Wallet) UpdateTxStatus() error {
	//TODO: as possibility tx_type is not updated from `done`

	types := []action.ActionType{
		action.ActionTypeReceipt,
		action.ActionTypePayment,
		action.ActionTypeTransfer,
	}

	//1. update tx_type for TxTypeSent
	for _, actionType := range types {
		err := w.updateStatusTxTypeSent(actionType)
		if err != nil {
			return errors.Wrapf(err, "fail to call updateStatusForTxTypeSent() ActionType: %s", actionType)
		}
	}

	//2. update tx_type for TxTypeDone
	// - TODO: notification
	for _, actionType := range types {
		err := w.updateStatusTxTypeDone(actionType)
		if err != nil {
			return errors.Wrapf(err, "fail to call updateStatusForTxTypeDone() ActionType: %s", actionType)
		}
	}

	return nil
}

// update TxTypeSent to TxTypeDone if confirmation is 6 or more
func (w *Wallet) updateStatusTxTypeSent(actionType action.ActionType) error {
	// get records whose status is TxTypeSent
	hashes, err := w.repo.Tx().GetSentHashTx(actionType, tx.TxTypeSent)
	if err != nil {
		return errors.Wrapf(err, "fail to call txRepo.GetSentHashTx(TxTypeSent) ActionType: %s", actionType)
	}

	// get hash in detail and check confirmation
	// update txType if confirmation is 6 or more (or configured number
	for _, hash := range hashes {
		isDone, err := w.checkTxConfirmation(hash, actionType)
		if err != nil {
			w.logger.Error(
				"fail to call w.checkTransaction()",
				zap.String("actionType", actionType.String()),
				zap.String("hash", hash),
				zap.Error(err))
			continue
		}
		if isDone {
			// current confirmation meets 6 or more
			_, err = w.repo.Tx().UpdateTxTypeBySentHashTx(actionType, tx.TxTypeDone, hash)
			if err != nil {
				return errors.Wrapf(err, "fail to call repo.Tx().UpdateTxTypeBySentHashTx(tx.TxTypeDone) ActionType: %s", actionType)
			}
		}
	}
	return nil
}

func (w *Wallet) updateStatusTxTypeDone(actionType action.ActionType) error {
	// get records whose status is TxTypeDone
	hashes, err := w.repo.Tx().GetSentHashTx(actionType, tx.TxTypeDone)
	if err != nil {
		return errors.Wrapf(err, "fail to call txRepo.GetSentHashTx(TxTypeDone) ActionType: %s", actionType)
	}
	w.logger.Debug(
		"called repo.Tx().GetSentHashTx(TxTypeDone)",
		zap.String("actionType", actionType.String()),
		zap.Any("hashes", hashes))

	// notify tx get done
	for _, hash := range hashes {
		txID, err := w.notifyTxDone(hash, actionType)
		if err != nil {
			w.logger.Error(
				"fail to call w.notifyUsers()",
				zap.String("actionType", actionType.String()),
				zap.String("hash", hash),
				zap.Error(err))
			continue
		}
		// update is already done
		if txID == 0 {
			continue
		}

		// update tx_type to TxTypeNotified
		err = w.updateTxTypeNotified(txID, actionType)
		//TODO: even if update is failed, notification is done. so how to manage??
		if err != nil {
			w.logger.Error(
				"fail to call w.updateTxTypeNotified()",
				zap.String("actionType", actionType.String()),
				zap.String("hash", hash),
				zap.Error(err))
			continue
		}
	}
	return nil
}

// checkTxConfirmation check confirmation for hash tx
func (w *Wallet) checkTxConfirmation(hash string, actionType action.ActionType) (bool, error) {
	// get tx in detail by RPC `gettransaction`
	tran, err := w.btc.GetTransactionByTxID(hash)
	if err != nil {
		return false, errors.Wrapf(err, "fail to call btc.GetTransactionByTxID(): ActionType: %s, txID:%s", actionType, hash)
	}
	w.logger.Debug("confirmation detail",
		zap.String("actionType", actionType.String()),
		zap.Uint64("confirmation", tran.Confirmations))

	// check current confirmation
	if tran.Confirmations >= uint64(w.btc.ConfirmationBlock()) {
		// current confirmation meets 6 or more
		return true, nil
	}

	// not completed yet
	//TODO: what if confirmation doesn't proceed for a long time after signed tx is sent
	// - should it be canceled??
	// - then raise fee and should unsigned tx be re-created again??
	w.logger.Info("confirmation is not met yet",
		zap.Uint64("want", w.btc.ConfirmationBlock()),
		zap.Uint64("got", tran.Confirmations))

	return false, nil
}

// notifyTxDone notify tx is sent and met specific confirmation number
func (w *Wallet) notifyTxDone(hash string, actionType action.ActionType) (int64, error) {

	var (
		txID int64
		err  error
	)

	switch actionType {
	case action.ActionTypeReceipt:
		// 1. get txID from hash
		txID, err = w.repo.Tx().GetTxIDBySentHash(actionType, hash)
		if err != nil {
			return 0, errors.Wrapf(err, "fail to call txRepo.GetTxIDBySentHash() ActionType: %s", actionType)
		}

		// 2. get txInputs
		txInputs, err := w.repo.TxInput().GetAllByTxID(txID)
		if err != nil {
			return 0, errors.Wrapf(err, "fail to call txInRepo.GetAllByTxID(%d) ActionType: %s", txID, actionType)
		}
		if len(txInputs) == 0 {
			w.logger.Debug("txInputs is not found in tx_input table",
				zap.Int64("tx_id", txID))
			return 0, nil
		}

		// 3. notify to given input_addresses tx is done
		// TODO:how to notify
		for _, input := range txInputs {
			w.logger.Debug("address in txInputs", zap.String("input.InputAddress", input.InputAddress))
		}
	case action.ActionTypePayment:
		// 1. get txID from hash
		txID, err = w.repo.Tx().GetTxIDBySentHash(actionType, hash)
		if err != nil {
			return 0, errors.Wrapf(err, "fail to call txRepo.GetTxIDBySentHash() ActionType: %s", actionType)
		}

		// 2. get info from payment_request table
		paymentUsers, err := w.repo.PayReq().GetAllByPaymentID(txID)
		if err != nil {
			return 0, errors.Wrapf(err, "fail to call repo.GetPaymentRequestByPaymentID(%d) ActionType: %s", txID, actionType)
		}
		if len(paymentUsers) == 0 {
			w.logger.Debug("payment user is not found",
				zap.Int64("tx_id", txID))
			return 0, nil
		}

		// 3. notify to given input_addresses tx is done
		// TODO:how to notify
		for _, user := range paymentUsers {
			w.logger.Debug("address in paymentUsers", zap.String("user.AddressFrom", user.SenderAddress))
		}
	case action.ActionTypeTransfer:
		//TODO: not implemented yet
		w.logger.Warn("action.ActionTypeTransfer is not implemented yet in notifyTxDone()")
		return 0, errors.New("action.ActionTypeTransfer is not implemented yet in notifyTxDone()")
	}

	return txID, nil
}

// update tx_type TxTypeNotified
func (w *Wallet) updateTxTypeNotified(id int64, actionType action.ActionType) error {
	switch actionType {
	case action.ActionTypeReceipt:
		_, err := w.repo.Tx().UpdateTxType(id, tx.TxTypeNotified)
		if err != nil {
			return errors.Wrapf(err, "fail to call repo.Tx().UpdateTxType(tx.TxTypeNotified) ActionType: %s", actionType)
		}
	case action.ActionTypePayment:
		dtx, err := w.repo.BeginTx()
		if err != nil {
			return errors.Wrapf(err, "fail to start transaction")
		}
		defer func() {
			if err != nil {
				dtx.Rollback()
			} else {
				dtx.Commit()
			}
		}()
		_, err = w.repo.Tx().UpdateTxType(id, tx.TxTypeNotified)
		if err != nil {
			return errors.Wrapf(err, "fail to call repo.Tx().UpdateTxType(tx.TxTypeNotified) ActionType: %s", actionType)
		}

		// update is_done=true in payment_request
		_, err = w.repo.PayReq().UpdateIsDone(id)
		if err != nil {
			return errors.Wrapf(err, "fail to call repo.UpdateIsDoneOnPaymentRequest() ActionType: %s", actionType)
		}
	case action.ActionTypeTransfer:
		//TODO: not implemented yet, it could be same to action.ActionTypeReceipt
		w.logger.Warn("action.ActionTypeTransfer is not implemented yet in updateTxTypeNotified()")
		return errors.New("action.ActionTypeTransfer is not implemented yet in updateTxTypeNotified()")
	}

	return nil
}
