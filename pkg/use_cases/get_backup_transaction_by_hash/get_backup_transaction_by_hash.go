package get_backup_transaction_by_hash

import (
	"context"
	"intmax2-store-vault/configs"
	"intmax2-store-vault/internal/logger"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	service "intmax2-store-vault/internal/store_vault_service"
	getBackupTransactionByHash "intmax2-store-vault/internal/use_cases/get_backup_transaction_by_hash"

	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// uc describes use case
type uc struct {
	cfg *configs.Config
	log logger.Logger
	db  SQLDriverApp
}

func New(
	cfg *configs.Config,
	log logger.Logger,
	db SQLDriverApp,
) getBackupTransactionByHash.UseCaseGetBackupTransactionByHash {
	return &uc{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (u *uc) Do(
	ctx context.Context,
	input *getBackupTransactionByHash.UCGetBackupTransactionByHashInput,
) (*node.GetBackupTransactionByHashResponse_Data, error) {
	const (
		hName     = "UseCase GetBackupTransactionByHash"
		senderKey = "sender"
		txHashKey = "tx_hash"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName)
	defer span.End()

	if input == nil {
		open_telemetry.MarkSpanError(spanCtx, ErrUCGetBackupTransactionByHashInputEmpty)
		return nil, ErrUCGetBackupTransactionByHashInputEmpty
	}

	span.SetAttributes(
		attribute.String(txHashKey, input.TxHash),
		attribute.String(senderKey, input.Sender),
	)

	transaction, err := service.GetBackupTransactionByHash(ctx, u.cfg, u.log, u.db, input)
	if err != nil {
		return nil, err
	}

	data := node.GetBackupTransactionByHashResponse_Data{
		Transaction: &node.GetBackupTransactionByHashResponse_Transaction{
			Id:          transaction.ID,
			Sender:      transaction.Sender,
			Signature:   transaction.Signature,
			BlockNumber: uint64(transaction.BlockNumber),
			EncryptedTx: transaction.EncryptedTx,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: transaction.CreatedAt.Unix(),
				Nanos:   int32(transaction.CreatedAt.Nanosecond()),
			},
		},
	}

	return &data, nil
}
