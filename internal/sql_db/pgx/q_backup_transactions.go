package pgx

import (
	"database/sql"
	"errors"
	"fmt"
	errPgx "intmax2-store-vault/internal/sql_db/pgx/errors"
	"intmax2-store-vault/internal/sql_db/pgx/models"
	"intmax2-store-vault/internal/sql_filter"
	mFL "intmax2-store-vault/internal/sql_filter/models"
	mDBApp "intmax2-store-vault/pkg/sql_db/db_app/models"
	errorsDB "intmax2-store-vault/pkg/sql_db/errors"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
	libPGX "github.com/jackc/pgx/v5"
)

func (p *pgx) CreateBackupTransaction(
	sender, encryptedTxHash, encryptedTx, signature string,
	blockNumber int64,
) (*mDBApp.BackupTransaction, error) {
	const query = `
	    INSERT INTO backup_transactions
        (id, sender, tx_double_hash, encrypted_tx, block_number, signature, created_at, encoding_version)
        VALUES ($1, $2, $3, $4, $5, $6, $7, 1)
	`

	id := uuid.New().String()
	createdAt := time.Now().UTC()

	err := p.createBackupEntry(query,
		id, sender, encryptedTxHash, encryptedTx, blockNumber, signature, createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup transaction: %w", err)
	}

	return p.GetBackupTransaction("id", id)
}

func (p *pgx) GetBackupTransaction(condition, value string) (*mDBApp.BackupTransaction, error) {
	const baseQuery = `
        SELECT id, sender, tx_double_hash, encrypted_tx, encoding_version, block_number, signature, created_at
        FROM backup_transactions
        WHERE %s = $1
    `
	query := fmt.Sprintf(baseQuery, condition)

	var b models.BackupTransaction
	err := errPgx.Err(p.queryRow(p.ctx, query, value).
		Scan(
			&b.ID,
			&b.Sender,
			&b.TxDoubleHash,
			&b.EncryptedTx,
			&b.EncodingVersion,
			&b.BlockNumber,
			&b.Signature,
			&b.CreatedAt,
		))
	if err != nil {
		return nil, fmt.Errorf("failed to get backup transaction: %w", err)
	}
	transaction := p.backupTransactionToDBApp(&b)
	return &transaction, nil
}

func (p *pgx) GetBackupTransactionByID(id string) (*mDBApp.BackupTransaction, error) {
	const (
		q = `
        SELECT id, sender, tx_double_hash, encrypted_tx, encoding_version, block_number, signature, created_at
        FROM backup_transactions
        WHERE id = $1 `
	)

	var b models.BackupTransaction
	err := errPgx.Err(p.queryRow(p.ctx, q, id).
		Scan(
			&b.ID,
			&b.Sender,
			&b.TxDoubleHash,
			&b.EncryptedTx,
			&b.EncodingVersion,
			&b.BlockNumber,
			&b.Signature,
			&b.CreatedAt,
		))
	if err != nil {
		return nil, err
	}
	transaction := p.backupTransactionToDBApp(&b)
	return &transaction, nil
}

func (p *pgx) GetBackupTransactions(condition string, value interface{}) ([]*mDBApp.BackupTransaction, error) {
	const baseQuery = `
        SELECT id, sender, tx_double_hash, encrypted_tx, encoding_version, block_number, signature, created_at
        FROM backup_transactions
        WHERE %s = $1
`
	query := fmt.Sprintf(baseQuery, condition)
	var transactions []*mDBApp.BackupTransaction
	err := p.getBackupEntries(query, value, func(rows *sql.Rows) error {
		var b models.BackupTransaction
		err := rows.Scan(&b.ID, &b.Sender, &b.TxDoubleHash, &b.EncryptedTx, &b.EncodingVersion, &b.BlockNumber, &b.Signature, &b.CreatedAt)
		if err != nil {
			return err
		}
		transaction := p.backupTransactionToDBApp(&b)
		transactions = append(transactions, &transaction)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (p *pgx) GetBackupTransactionsBySender(
	sender string,
	pagination mDBApp.PaginationOfListOfBackupTransactionsInput,
	sorting mFL.Sorting, orderBy mFL.OrderBy,
	filters mFL.FiltersList,
) (
	paginator *mDBApp.PaginationOfListOfBackupTransactions,
	listDBApp mDBApp.ListOfBackupTransaction,
	err error,
) {
	var (
		q = `
SELECT id, sender, tx_double_hash, encrypted_tx, encoding_version, block_number, signature, created_at
FROM backup_transactions
WHERE sender = @sender %s
`
	)

	sorting = mFL.Sorting(strings.TrimSpace(string(sorting)))
	if sorting == "" {
		sorting = mFL.SortingDESC
	}

	var (
		cursor       string
		orderByValue string
	)
	switch orderBy {
	case mFL.DateCreate:
		const createdAtKey = "created_at"
		orderByValue = createdAtKey
		if pagination.Cursor != nil {
			cursor = time.Unix(0, pagination.Cursor.SortingValue.Int64()).UTC().Format(time.RFC3339Nano)
		}
	default:
		orderBy = mFL.DateCreate
		const startedAtKey = "created_at"
		orderByValue = startedAtKey
		if pagination.Cursor != nil {
			cursor = time.Unix(0, pagination.Cursor.SortingValue.Int64()).UTC().Format(time.RFC3339Nano)
		}
	}

	wParams := make(libPGX.NamedArgs)
	wParams["sender"] = sender

	var where string
	if len(filters) > 0 {
		var (
			fc     sql_filter.SQLFilter
			params map[string]interface{}
		)
		where, params = fc.FilterDataToWhereQuery(filters)
		where = fc.PrepareWhereString(where, false)
		where = fmt.Sprintf("AND (%s)", where)

		for pKey := range params {
			wParams[pKey] = params[pKey]
		}
	}

	var revers bool
	if pagination.Cursor != nil {
		rID := pagination.Cursor.ID
		cond := mFL.LessSymbol
		if sorting == mFL.SortingDESC && pagination.Direction == mFL.DirectionNext {
			cond = mFL.LessSymbol
		} else if sorting == mFL.SortingDESC && pagination.Direction == mFL.DirectionPrev {
			sorting = mFL.SortingASC
			cond = mFL.MoreSymbol
			revers = true
		} else if sorting == mFL.SortingASC && pagination.Direction == mFL.DirectionNext {
			cond = mFL.MoreSymbol
		} else if sorting == mFL.SortingASC && pagination.Direction == mFL.DirectionPrev {
			sorting = mFL.SortingDESC
			cond = mFL.LessSymbol
			revers = true
		}
		where += fmt.Sprintf(
			"AND ((%s, id) %s ('%s', '%s'))",
			orderByValue, cond, cursor, rID)
	}

	q += fmt.Sprintf(" ORDER BY %s %s, id %s", orderByValue, sorting, sorting)

	q += fmt.Sprintf(" FETCH FIRST %d ROWS ONLY ", pagination.Offset)

	var rows *sql.Rows
	rows, err = p.query(p.ctx, fmt.Sprintf(q, where), wParams)
	if err != nil {
		err = errPgx.Err(err)
		if errors.Is(err, errorsDB.ErrNotFound) {
			return nil, nil, nil
		}

		return nil, nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var list models.ListOfBackupTransaction
	for rows.Next() {
		var b models.BackupTransaction
		err = rows.Scan(
			&b.ID,
			&b.Sender,
			&b.TxDoubleHash,
			&b.EncryptedTx,
			&b.EncodingVersion,
			&b.BlockNumber,
			&b.Signature,
			&b.CreatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		list = append(list, b)
	}

	listDBApp = make(mDBApp.ListOfBackupTransaction, len(list))
	if revers {
		for key := len(list) - 1; key >= 0; key-- {
			listDBApp[len(list)-1-key] = p.backupTransactionToDBApp(&list[key])
		}
	} else {
		for key := range list {
			listDBApp[key] = p.backupTransactionToDBApp(&list[key])
		}
	}

	paginator = &mDBApp.PaginationOfListOfBackupTransactions{
		Offset: pagination.Offset,
	}

	if list != nil {
		const (
			int0Key = 0
			int1Key = 1
		)

		startV := int0Key
		endV := len(list) - int1Key
		if revers {
			startV = len(list) - int1Key
			endV = int0Key
		}

		paginator.Cursor = &mDBApp.CursorListOfBackupTransactions{
			Prev: &mDBApp.CursorBaseOfListOfBackupTransactions{
				ID: list[startV].ID,
			},
			Next: &mDBApp.CursorBaseOfListOfBackupTransactions{
				ID: list[endV].ID,
			},
		}

		switch orderBy { // nolint:gocritic
		case mFL.DateCreate:
			paginator.Cursor.Prev.SortingValue = new(big.Int).SetInt64(list[startV].CreatedAt.UTC().UnixNano())
			paginator.Cursor.Next.SortingValue = new(big.Int).SetInt64(list[endV].CreatedAt.UTC().UnixNano())
		}
	}

	return paginator, listDBApp, nil
}

func (p *pgx) backupTransactionToDBApp(b *models.BackupTransaction) mDBApp.BackupTransaction {
	return mDBApp.BackupTransaction{
		ID:              b.ID,
		Sender:          b.Sender,
		TxDoubleHash:    b.TxDoubleHash.String,
		EncryptedTx:     b.EncryptedTx,
		EncodingVersion: b.EncodingVersion,
		BlockNumber:     b.BlockNumber,
		Signature:       b.Signature,
		CreatedAt:       b.CreatedAt,
	}
}
