package ed25519

import (
	"filippo.io/edwards25519"
	"golang.org/x/crypto/blake2b"
)

func Sign(pubKey, privKey [32]byte, msg []byte) ([]byte, error) {
	sig := make([]byte, 64)

	h, err := blake2b.New512(nil)

	if err != nil {
		return nil, err
	}

	h.Write(privKey[:])

	var dig, msgDig, hram [64]byte
	h.Sum(dig[:0])

	s1, err := new(edwards25519.Scalar).SetBytesWithClamping(dig[:32])

	if err != nil {
		return nil, err
	}

	h.Reset()
	h.Write(dig[32:])
	h.Write(msg)
	h.Sum(msgDig[:0])

	rr, err := new(edwards25519.Scalar).SetUniformBytes(msgDig[:])

	if err != nil {
		return nil, err
	}

	rp := new(edwards25519.Point).ScalarBaseMult(rr)
	enc := rp.Bytes()

	h.Reset()
	h.Write(enc[:])
	h.Write(pubKey[:])
	h.Write(msg)
	h.Sum(hram[:0])

	kr, err := new(edwards25519.Scalar).SetUniformBytes(hram[:])

	if err != nil {
		return nil, err
	}

	s2 := new(edwards25519.Scalar).MultiplyAdd(kr, s1, rr)

	copy(sig[:], enc[:])
	copy(sig[32:], s2.Bytes())

	return sig, nil
}
