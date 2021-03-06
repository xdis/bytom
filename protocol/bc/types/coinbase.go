package types

// CoinbaseInput is record the coinbase message
type CoinbaseInput struct {
	Arbitrary []byte
}

// NewCoinbaseInput create a new coinbase input struct
func NewCoinbaseInput(arbitrary []byte) *TxInput {
	return &TxInput{
		AssetVersion: 1,
		TypedInput:   &CoinbaseInput{Arbitrary: arbitrary},
	}
}

// InputType is the interface function for return the input type
func (cb *CoinbaseInput) InputType() uint8 { return CoinbaseInputType }
