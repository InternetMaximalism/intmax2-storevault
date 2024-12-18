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

func (p *pgx) CreateBackupDeposit(
	recipient, depositHash, encryptedDeposit string,
	blockNumber int64,
) (*mDBApp.BackupDeposit, error) {
	const query = `
	    INSERT INTO backup_deposits
        (id, recipient, deposit_double_hash, encrypted_deposit, block_number, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
	`

	id := uuid.New().String()
	createdAt := time.Now().UTC()

	err := p.createBackupEntry(query, id, recipient, depositHash, encryptedDeposit, blockNumber, createdAt)
	if err != nil {
		return nil, err
	}

	return p.GetBackupDeposit([]string{"id"}, []interface{}{id})
}

func (p *pgx) GetBackupDepositByID(id string) (*mDBApp.BackupDeposit, error) {
	const (
		q = `
        SELECT id, recipient, deposit_double_hash, encrypted_deposit, block_number, created_at
        FROM backup_deposits
        WHERE id = $1 `
	)

	var b models.BackupDeposit
	err := errPgx.Err(p.queryRow(p.ctx, q, id).
		Scan(
			&b.ID,
			&b.Recipient,
			&b.DepositDoubleHash,
			&b.EncryptedDeposit,
			&b.BlockNumber,
			&b.CreatedAt,
		))
	if err != nil {
		return nil, err
	}
	deposit := p.backupDepositToDBApp(&b)
	return &deposit, nil
}

func (p *pgx) GetBackupDeposit(conditions []string, values []interface{}) (*mDBApp.BackupDeposit, error) {
	const baseQuery = `
        SELECT id, recipient, deposit_double_hash, encrypted_deposit, block_number, created_at 
        FROM backup_deposits 
        WHERE %s`

	whereClause := make([]string, len(conditions))
	for i, condition := range conditions {
		whereClause[i] = fmt.Sprintf("%s = $%d", condition, i+1)
	}

	query := fmt.Sprintf(baseQuery, strings.Join(whereClause, " AND "))

	var b models.BackupDeposit
	err := errPgx.Err(p.queryRow(p.ctx, query, values...).
		Scan(
			&b.ID,
			&b.Recipient,
			&b.DepositDoubleHash,
			&b.EncryptedDeposit,
			&b.BlockNumber,
			&b.CreatedAt,
		))
	if err != nil {
		return nil, err
	}

	deposit := p.backupDepositToDBApp(&b)
	return &deposit, nil
}

func (p *pgx) GetBackupDeposits(condition string, value interface{}) ([]*mDBApp.BackupDeposit, error) {
	const baseQuery = `
        SELECT id, recipient, deposit_double_hash, encrypted_deposit, block_number, created_at
        FROM backup_deposits
        WHERE %s = $1
    `
	query := fmt.Sprintf(baseQuery, condition)
	var deposits []*mDBApp.BackupDeposit
	err := p.getBackupEntries(query, value, func(rows *sql.Rows) error {
		var b models.BackupDeposit
		err := rows.Scan(&b.ID, &b.Recipient, &b.DepositDoubleHash, &b.EncryptedDeposit, &b.BlockNumber, &b.CreatedAt)
		if err != nil {
			return err
		}
		deposit := p.backupDepositToDBApp(&b)
		deposits = append(deposits, &deposit)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return deposits, nil
}

func (p *pgx) GetBackupDepositsByRecipient(
	recipient string,
	pagination mDBApp.PaginationOfListOfBackupDepositsInput,
	sorting mFL.Sorting, orderBy mFL.OrderBy,
	filters mFL.FiltersList,
) (
	paginator *mDBApp.PaginationOfListOfBackupDeposits,
	listDBApp mDBApp.ListOfBackupDeposit,
	err error,
) {
	var (
		q = `
SELECT id ,recipient ,deposit_double_hash ,encrypted_deposit ,block_number ,created_at
FROM backup_deposits
WHERE recipient = @recipient %s
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
	wParams["recipient"] = recipient

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

	var list models.ListOfBackupDeposit
	for rows.Next() {
		var b models.BackupDeposit
		err = rows.Scan(
			&b.ID,
			&b.Recipient,
			&b.DepositDoubleHash,
			&b.EncryptedDeposit,
			&b.BlockNumber,
			&b.CreatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		list = append(list, b)
	}

	listDBApp = make(mDBApp.ListOfBackupDeposit, len(list))
	if revers {
		for key := len(list) - 1; key >= 0; key-- {
			listDBApp[len(list)-1-key] = p.backupDepositToDBApp(&list[key])
		}
	} else {
		for key := range list {
			listDBApp[key] = p.backupDepositToDBApp(&list[key])
		}
	}

	paginator = &mDBApp.PaginationOfListOfBackupDeposits{
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

		paginator.Cursor = &mDBApp.CursorListOfBackupDeposits{
			Prev: &mDBApp.CursorBaseOfListOfBackupDeposits{
				ID: list[startV].ID,
			},
			Next: &mDBApp.CursorBaseOfListOfBackupDeposits{
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

func (p *pgx) backupDepositToDBApp(b *models.BackupDeposit) mDBApp.BackupDeposit {
	return mDBApp.BackupDeposit{
		ID:                b.ID,
		Recipient:         b.Recipient,
		DepositDoubleHash: b.DepositDoubleHash.String,
		EncryptedDeposit:  b.EncryptedDeposit,
		BlockNumber:       b.BlockNumber,
		CreatedAt:         b.CreatedAt,
	}
}
