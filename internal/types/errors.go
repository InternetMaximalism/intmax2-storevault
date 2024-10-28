package types

import "errors"

// ErrETHAddressInvalid error: the Ethereum address should be 20 bytes.
var ErrETHAddressInvalid = errors.New("the Ethereum address should be 20 bytes")

// ErrINTMAXAddressInvalid error: the INTMAX address should be 32 bytes.
var ErrINTMAXAddressInvalid = errors.New("the INTMAX address should be 32 bytes")

var ErrTokenTypeRequired = errors.New("token type is required")

var ErrInvalidETHArgs = errors.New("ETH operation requires additional arguments")

var ErrInvalidERC20Args = errors.New("ERC20 operation requires a token address")

var ErrInvalidERC721Args = errors.New("ERC721 operation requires a token address and token ID")

var ErrInvalidERC1155Args = errors.New("ERC1155 operation requires a token address and token ID")

var ErrInvalidTokenType = errors.New("invalid token type. Use 'eth', 'erc20', 'erc721', or 'erc1155'")
