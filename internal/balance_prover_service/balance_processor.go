package balance_prover_service

type BalanceProofWithPublicInputs struct {
	// Proof is a base64 encoded string with public inputs
	Proof        string
	PublicInputs *BalancePublicInputs
}

type SpentProofWithPublicInputs struct {
	Proof        string
	PublicInputs *SpentPublicInputs
}
