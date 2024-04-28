package nanogo

import (
	"encoding/binary"
	"encoding/hex"
	"filippo.io/edwards25519"
	"fmt"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/blake2b"
	"strings"
)

// SeedToPrivateKey converts a seed to a private key,
// seed: the seed to convert,
// index: the index of the private key to generate,
// returns the private key or an error.
func SeedToPrivateKey(seed string, index int) ([32]byte, error) {
	seedBytes, err := hex.DecodeString(seed)

	if err != nil {
		return [32]byte{}, fmt.Errorf("could not decode seed: %v", err)
	}

	if len(seedBytes) != 32 {
		return [32]byte{}, fmt.Errorf("seed length is not 32 bytes")
	}

	iBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(iBytes, uint32(index))
	comb := append(seedBytes, iBytes...)
	privKeyBytes := blake2b.Sum256(comb)

	return privKeyBytes, nil
}

// PrivateKeyToPublicKey converts a private key to a public key,
// privateKey: the private key to convert,
// returns the public key or an error.
func PrivateKeyToPublicKey(privateKey [32]byte) ([32]byte, error) {
	hashBytes := blake2b.Sum512(privateKey[:])
	scalar, err := new(edwards25519.Scalar).SetBytesWithClamping(hashBytes[:32])

	if err != nil {
		return [32]byte{}, err
	}

	pubKeyBytes := new(edwards25519.Point).ScalarBaseMult(scalar).Bytes()

	return [32]byte(pubKeyBytes), nil
}

// PublicKeyToAddress converts a public key to a wallet address,
// publicKey: the public key to convert,
// returns the address or an error.
func PublicKeyToAddress(publicKey [32]byte) (string, error) {
	b32PubKey := base32Encode(publicKey[:])
	h, err := blake2b.New(5, nil)

	if err != nil {
		return "", err
	}

	if _, err := h.Write(publicKey[:]); err != nil {
		return "", err
	}

	hashBytes := h.Sum(nil)
	b32Hash := base32Encode(revertBytes(hashBytes))

	address := "nano_" + strings.Repeat("1", 52-len(b32PubKey)) + b32PubKey + strings.Repeat("1", 8-len(b32Hash)) + b32Hash

	return address, nil
}

// AddressToPublicKey converts a wallet address to a public key,
// address: the address to get the public key from,
// returns the public key or an error.
func AddressToPublicKey(address string) ([32]byte, error) {
	if len(address) == 64 {
		bytes, err := base32Decode(address[4:56])
		return [32]byte(bytes), err
	} else if len(address) == 65 {
		bytes, err := base32Decode(address[5:57])
		return [32]byte(bytes), err
	}

	return [32]byte{}, fmt.Errorf("could not parse address (%s)", address)
}

// NanoToRaw converts a nano amount to raw,
// nano: the nano amount to convert,
// returns the raw or an error.
func NanoToRaw(nano string) (string, error) {
	nanoDec, err := decimal.NewFromString(nano)

	if err != nil {
		return "", fmt.Errorf("could not parse nano: %v", err)
	}

	rawPerNano, err := decimal.NewFromString("1000000000000000000000000000000")

	if err != nil {
		return "", fmt.Errorf("could not parse raw per nano: %v", err)
	}

	raw := nanoDec.Mul(rawPerNano)

	return raw.String(), nil
}

// RawToNano converts a raw to nano amount,
// raw: the raw to convert,
// returns the nano amount or an error.
func RawToNano(raw string) (string, error) {
	decimal.DivisionPrecision = 30
	rawDec, err := decimal.NewFromString(raw)

	if err != nil {
		return "", fmt.Errorf("could not parse raw: %v", err)
	}

	rawPerNano, err := decimal.NewFromString("1000000000000000000000000000000")

	if err != nil {
		return "", fmt.Errorf("could not parse raw per nano: %v", err)
	}

	nano := rawDec.Div(rawPerNano)

	return nano.String(), nil
}
