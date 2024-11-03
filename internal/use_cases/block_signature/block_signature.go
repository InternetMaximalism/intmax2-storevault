package block_signature

import (
	"encoding/base64"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate mockgen -destination=../mocks/mock_block_signature.go -package=mocks -source=block_signature.go

type Plonky2Proof struct {
	PublicInputs []uint64 `json:"publicInputs"`
	Proof        []byte   `json:"proof"`
}

func (dst *Plonky2Proof) Set(src *Plonky2Proof) *Plonky2Proof {
	dst.PublicInputs = make([]uint64, len(src.PublicInputs))
	copy(dst.PublicInputs, src.PublicInputs)
	dst.Proof = make([]byte, len(src.Proof))
	copy(dst.Proof, src.Proof)

	return dst
}

type EnoughBalanceProofInput struct {
	PrevBalanceProof  *Plonky2Proof `json:"prevBalanceProof"`
	TransferStepProof *Plonky2Proof `json:"transferStepProof"`
}

func (dst *EnoughBalanceProofInput) Set(src *EnoughBalanceProofInput) *EnoughBalanceProofInput {
	dst.PrevBalanceProof = new(Plonky2Proof).Set(src.PrevBalanceProof)
	dst.TransferStepProof = new(Plonky2Proof).Set(src.TransferStepProof)

	return dst
}

type EnoughBalanceProofBodyInput struct {
	PrevBalanceProofBody  string `json:"prevBalanceProof"`
	TransferStepProofBody string `json:"transferStepProof"`
}

func (dst *EnoughBalanceProofBodyInput) Set(src *EnoughBalanceProofBodyInput) *EnoughBalanceProofBodyInput {
	dst.PrevBalanceProofBody = src.PrevBalanceProofBody
	dst.TransferStepProofBody = src.TransferStepProofBody

	return dst
}

func (dst *EnoughBalanceProofBodyInput) FromEnoughBalanceProofInput(src *EnoughBalanceProofInput) *EnoughBalanceProofBodyInput {
	dst.PrevBalanceProofBody = base64.StdEncoding.EncodeToString(src.PrevBalanceProof.Proof)
	dst.TransferStepProofBody = base64.StdEncoding.EncodeToString(src.TransferStepProof.Proof)

	return dst
}

func (proof *EnoughBalanceProofBodyInput) FromEnoughBalanceProofBody(input *EnoughBalanceProofBody) *EnoughBalanceProofBodyInput {
	proof.PrevBalanceProofBody = base64.StdEncoding.EncodeToString(input.PrevBalanceProofBody)
	proof.TransferStepProofBody = base64.StdEncoding.EncodeToString(input.TransferStepProofBody)

	return proof
}

type EnoughBalanceProofBody struct {
	PrevBalanceProofBody  []byte
	TransferStepProofBody []byte
}

func (proof *EnoughBalanceProofBodyInput) EnoughBalanceProofBody() (*EnoughBalanceProofBody, error) {
	prevBalanceProofBodyBytes, err := base64.StdEncoding.DecodeString(proof.PrevBalanceProofBody)
	if err != nil {
		return nil, err
	}

	transferStepProofBodyBytes, err := base64.StdEncoding.DecodeString(proof.TransferStepProofBody)
	if err != nil {
		return nil, err
	}

	return &EnoughBalanceProofBody{
		PrevBalanceProofBody:  prevBalanceProofBodyBytes,
		TransferStepProofBody: transferStepProofBodyBytes,
	}, nil
}

func (proof *EnoughBalanceProofBody) Hash() string {
	buf := []byte{}
	buf = append(buf, proof.PrevBalanceProofBody...)
	buf = append(buf, proof.TransferStepProofBody...)
	output := crypto.Keccak256(buf)

	return hexutil.Encode(output)
}
