package get_backup_transactions_list

import (
	"errors"
	mFL "intmax2-store-vault/internal/sql_filter/models"
	"strconv"
	"time"

	"github.com/holiman/uint256"
	"github.com/prodadidb/go-validation"
	"github.com/prodadidb/go-validation/is"
)

// ErrValueInvalid error: value must be valid.
var ErrValueInvalid = errors.New("value must be valid")

func (input *UCGetBackupTransactionsListInput) Valid() error {
	const (
		int0Key      = 0
		int1Key      = 1
		intMinus1Key = -1
		int100Key    = 100
	)
	return validation.ValidateStruct(input,
		validation.Field(&input.Sender, validation.Required),
		validation.Field(&input.Pagination, validation.By(func(value interface{}) error {
			var isNil bool
			value, isNil = validation.Indirect(value)

			if isNil || validation.IsEmpty(value) {
				return nil
			}

			pagination, ok := value.(UCGetBackupTransactionsListPaginationInput)
			if !ok {
				return ErrValueInvalid
			}

			return validation.ValidateStruct(&pagination,
				validation.Field(&pagination.Direction, validation.In(mFL.DirectionPrev, mFL.DirectionNext)),
				validation.Field(&pagination.PerPage, validation.By(func(value interface{}) error {
					var v string
					v, ok = value.(string)
					if !ok {
						return ErrValueInvalid
					}

					perPage, err := strconv.Atoi(v)
					if err != nil {
						perPage = int100Key
					}

					if perPage == int0Key {
						perPage = intMinus1Key
					}

					err = validation.Min(int1Key).Validate(perPage)
					if err != nil {
						return err
					}

					err = validation.Max(int100Key).Validate(perPage)
					if err != nil {
						return err
					}

					input.Pagination.Offset = perPage

					return nil
				})),
				validation.Field(&pagination.Cursor, validation.By(func(value interface{}) error {
					value, isNil = validation.Indirect(value)
					if isNil || validation.IsEmpty(value) {
						return nil
					}

					var cursor UCGetBackupTransactionsListCursorBase
					cursor, ok = value.(UCGetBackupTransactionsListCursorBase)
					if !ok {
						return ErrValueInvalid
					}

					return validation.ValidateStruct(&cursor,
						validation.Field(&cursor.Uuid, validation.Required, is.UUID),
						validation.Field(&cursor.SortingValue, validation.By(func(value interface{}) error {
							var v string
							v, ok = value.(string)
							if !ok {
								return ErrValueInvalid
							}

							var test uint256.Int
							err := test.Scan(v)
							if err != nil {
								return ErrValueInvalid
							}

							input.Pagination.Cursor.ConvertSortingValue = test.ToBig()

							return nil
						})),
					)
				})),
			)
		})),
		validation.Field(&input.OrderBy, validation.In(
			mFL.DateCreate,
		)),
		validation.Field(&input.Sorting, validation.In(mFL.SortingASC, mFL.SortingDESC)),
		validation.Field(&input.Filters, validation.Each(validateFilter())),
	)
}

func validateFilter() validation.Rule {
	return validation.By(func(value interface{}) error {
		var isNil bool
		value, isNil = validation.Indirect(value)
		if isNil || validation.IsEmpty(value) {
			return ErrValueInvalid
		}

		f, ok := value.(mFL.Filter)
		if !ok {
			return ErrValueInvalid
		}

		return validation.ValidateStruct(&f,
			validation.Field(&f.Relation, validation.Required, validation.In(mFL.RelationAnd, mFL.RelationOr)),
			validation.Field(&f.DataField, validation.Required, validation.In(
				mFL.DataFieldBackupStartTime,
			)),
			validation.Field(&f.Condition,
				validation.When(f.DataField == mFL.DataFieldBackupStartTime, validation.In(
					mFL.ConditionGreaterThan, mFL.ConditionLessThan,
					mFL.ConditionGreaterThanOrEqualTo, mFL.ConditionLessThanOrEqualTo,
				)),
			),
			validation.Field(&f.Value,
				validation.When(f.DataField == mFL.DataFieldBackupStartTime,
					validation.Required,
					validation.By(func(value interface{}) (err error) {
						var v string
						v, ok = value.(string)
						if !ok {
							return ErrValueInvalid
						}

						_, err = time.Parse(time.RFC3339Nano, v)
						if err != nil {
							return ErrValueInvalid
						}

						return nil
					}),
				),
			),
		)
	})
}
