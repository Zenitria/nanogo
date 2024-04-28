package nanogo

import (
	"encoding/hex"
	"fmt"
	"github.com/zenitria/nanogo/ed25519"
	"golang.org/x/crypto/blake2b"
	"math/big"
)

// Block is a block of the Nano blockchain,
// Type: the type of the block,
// Account: the account of the block,
// Previous: the account previous block hash,
// Representative: the representative wallet address,
// Balance: the new balance of the account in raw,
// Link: the link of the block,
// LinkAsAccount: the link as account,
// Signature: the signature of the block,
// Work: the work of the block.
type Block struct {
	Type           string `json:"type"`
	Account        string `json:"account"`
	Previous       string `json:"previous"`
	Representative string `json:"representative"`
	Balance        string `json:"balance"`
	Link           string `json:"link"`
	LinkAsAccount  string `json:"link_as_account"`
	Signature      string `json:"signature"`
	Work           string `json:"work"`
}

// Sign signs block,
// privateKey: private key with which you sign the block,
// returns an error.
func (b *Block) Sign(privateKey [32]byte) error {
	pubKey, err := PrivateKeyToPublicKey(privateKey)

	if err != nil {
		return err
	}

	hash, err := b.hashBytes()

	if err != nil {
		return err
	}

	sig, err := ed25519.Sign(pubKey, privateKey, hash)

	if err != nil {
		return err
	}

	b.Signature = fmt.Sprintf("%0128X", sig)

	return nil
}

// AddWork adds work to the block,
// work: the work to add.
func (b *Block) AddWork(work string) {
	b.Work = work
}

func (b *Block) hashBytes() ([]byte, error) {
	msg := make([]byte, 176)
	msg[31] = 0x6
	pubKey, err := AddressToPublicKey(b.Account)

	if err != nil {
		return nil, err
	}

	copy(msg[32:64], pubKey[:])

	prev, err := hex.DecodeString(b.Previous)

	if err != nil {
		return nil, err
	}

	copy(msg[64:96], prev)

	rep, err := AddressToPublicKey(b.Representative)

	if err != nil {
		return nil, err
	}

	copy(msg[96:128], rep[:])

	bal, ok := new(big.Int).SetString(b.Balance, 10)

	if !ok {
		return nil, fmt.Errorf("could not convert string to big int")
	}

	copy(msg[128:144], bal.FillBytes(make([]byte, 16)))

	link, err := hex.DecodeString(b.Link)

	if err != nil {
		return nil, err
	}

	copy(msg[144:176], link)

	hash := blake2b.Sum256(msg)

	return hash[:], nil
}
