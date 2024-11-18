package errors

import (
	"errors"
)

// ErrInvalidSequenceStr error: invalid sequence: invalid sequence.
const ErrInvalidSequenceStr = "invalid sequence: invalid sequence"

// Err520ScrollWebServerStr error: 520: Web server is returning an unknown error.
const Err520ScrollWebServerStr = "520: Web server is returning an unknown error"

// Err502ScrollWebServerStr error: 502: Bad gateway.
const Err502ScrollWebServerStr = "502: Bad gateway"

// Err502EthereumWevServerStr error: 502 Bad Gateway.
const Err502EthereumWevServerStr = "502 Bad Gateway"

// Err503EthereumWebServerStr error: 503 Service Unavailable.
const Err503EthereumWebServerStr = "503 Service Unavailable"

// ErrScrollChainIDInvalidStr error: the scroll chain ID must be equal: %s, %s.
const ErrScrollChainIDInvalidStr = "the scroll chain ID must be equal: %s, %s"

// ErrEthereumChainIDInvalidStr error: the ethereum chain ID must be equal: %s, %s.
const ErrEthereumChainIDInvalidStr = "the ethereum chain ID must be equal: %s, %s"

// ErrStdinProcessingFail error: stdin processing error occurred.
var ErrStdinProcessingFail = errors.New("stdin processing error occurred")

// ErrEthClientDialFail error: failed to dial ETH client.
var ErrEthClientDialFail = errors.New("failed to dial ETH client")

// ErrChainIDWithEthClientFail error: failed to get chain ID with ETH client.
var ErrChainIDWithEthClientFail = errors.New("failed to get chain ID with ETH client")
