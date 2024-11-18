package store_vault_server

import (
	"context"
	"fmt"
	"intmax2-store-vault/internal/open_telemetry"
	node "intmax2-store-vault/internal/pb/gen/store_vault_service/node"
	postBackupBalance "intmax2-store-vault/internal/use_cases/post_backup_balance"
	"intmax2-store-vault/pkg/grpc_server/utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *StoreVaultServer) BackupBalance(
	ctx context.Context,
	req *node.BackupBalanceRequest,
) (*node.BackupBalanceResponse, error) {
	resp := node.BackupBalanceResponse{}

	const (
		hName      = "Handler BackupBalance"
		requestKey = "request"
	)

	spanCtx, span := open_telemetry.Tracer().Start(ctx, hName,
		trace.WithAttributes(
			attribute.String(requestKey, req.String()),
		))
	defer span.End()

	input := postBackupBalance.UCPostBackupBalanceInput{
		User:                  req.User,
		EncryptedBalanceProof: req.EncryptedBalanceProof,
		EncryptedBalanceData:  req.EncryptedBalanceData,
		EncryptedTxs:          req.EncryptedTxs,
		EncryptedTransfers:    req.EncryptedTransfers,
		EncryptedDeposits:     req.EncryptedDeposits,
		Signature:             req.Signature,
		BlockNumber:           uint32(req.BlockNumber),
	}

	err := input.Valid()
	if err != nil {
		open_telemetry.MarkSpanError(spanCtx, err)
		return &resp, utils.BadRequest(spanCtx, err)
	}

	var info postBackupBalance.UCPostBackupBalance
	err = s.dbApp.Exec(spanCtx, &info, func(d interface{}, in interface{}) (err error) {
		q, _ := d.(SQLDriverApp)

		var result *postBackupBalance.UCPostBackupBalance
		result, err = s.commands.PostBackupBalance(s.config, s.log, q).Do(spanCtx, &input)
		if err != nil {
			open_telemetry.MarkSpanError(spanCtx, err)
			const msg = "failed to post backup balance: %w"
			return fmt.Errorf(msg, err)
		}

		if v, ok := in.(*postBackupBalance.UCPostBackupBalance); ok {
			v.ID = result.ID
			v.UserAddress = result.UserAddress
			v.EncryptedBalanceProof = result.EncryptedBalanceProof
			v.EncryptedBalanceData = result.EncryptedBalanceData
			v.EncryptedTxs = result.EncryptedTxs
			v.EncryptedTransfers = result.EncryptedTransfers
			v.EncryptedDeposits = result.EncryptedDeposits
			v.BlockNumber = result.BlockNumber
			v.Signature = result.Signature
			v.CreatedAt = result.CreatedAt
		} else {
			const msg = "failed to convert of the backup balance"
			err = fmt.Errorf(msg)
			open_telemetry.MarkSpanError(spanCtx, err)
			return err
		}

		return nil
	})
	if err != nil {
		const msg = "failed to post backup balance with DB App: %+v"
		return &resp, utils.Internal(spanCtx, s.log, msg, err)
	}

	resp.Success = true
	resp.Data = &node.BackupBalanceResponse_Data{
		Balance: &node.BackupBalanceResponse_Data_Balance{
			Id:                    info.ID,
			UserAddress:           info.UserAddress,
			EncryptedBalanceProof: info.EncryptedBalanceProof,
			EncryptedBalanceData:  info.EncryptedBalanceData,
			EncryptedTxs:          info.EncryptedTxs,
			EncryptedTransfers:    info.EncryptedTransfers,
			EncryptedDeposits:     info.EncryptedDeposits,
			BlockNumber:           info.BlockNumber,
			Signature:             info.Signature,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: info.CreatedAt.Unix(),
				Nanos:   int32(info.CreatedAt.Nanosecond()),
			},
		},
	}

	return &resp, utils.OK(spanCtx)
}
