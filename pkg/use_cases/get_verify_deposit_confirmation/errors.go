package get_verify_deposit_confirmation

import "errors"

// ErrUCGetVerifyDepositConfirmationInputEmpty error: ucGetVerifyDepositConfirmationInput must not be empty.
var ErrUCGetVerifyDepositConfirmationInputEmpty = errors.New("ucGetVerifyDepositConfirmationInput must not be empty")

// ErrGetVerifyDepositConfirmationFail error: failed to get verify deposit confirmation.
var ErrGetVerifyDepositConfirmationFail = errors.New("failed to get verify deposit confirmation")
