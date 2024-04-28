package nanogo

import (
	"fmt"
	"math/big"
)

func revertBytes(in []byte) []byte {
	for i := 0; i < len(in)/2; i++ {
		in[i], in[len(in)-1-i] = in[len(in)-1-i], in[i]
	}

	return in
}

func base32Encode(in []byte) string {
	alph := []byte("13456789abcdefghijkmnopqrstuwxyz")
	z := big.NewInt(0)
	radix := big.NewInt(32)
	num := big.NewInt(0).SetBytes(in)
	o := make([]byte, 0)
	mod := new(big.Int)

	for num.Cmp(z) > 0 {
		num.DivMod(num, radix, mod)
		o = append(o, alph[mod.Int64()])
	}

	for i := 0; i < len(o)/2; i++ {
		o[i], o[len(o)-1-i] = o[len(o)-1-i], o[i]
	}

	return string(o)
}

func base32Decode(in string) ([]byte, error) {
	revAlph := map[rune]*big.Int{}
	revAlph['1'] = big.NewInt(0)
	revAlph['3'] = big.NewInt(1)
	revAlph['4'] = big.NewInt(2)
	revAlph['5'] = big.NewInt(3)
	revAlph['6'] = big.NewInt(4)
	revAlph['7'] = big.NewInt(5)
	revAlph['8'] = big.NewInt(6)
	revAlph['9'] = big.NewInt(7)
	revAlph['a'] = big.NewInt(8)
	revAlph['b'] = big.NewInt(9)
	revAlph['c'] = big.NewInt(10)
	revAlph['d'] = big.NewInt(11)
	revAlph['e'] = big.NewInt(12)
	revAlph['f'] = big.NewInt(13)
	revAlph['g'] = big.NewInt(14)
	revAlph['h'] = big.NewInt(15)
	revAlph['i'] = big.NewInt(16)
	revAlph['j'] = big.NewInt(17)
	revAlph['k'] = big.NewInt(18)
	revAlph['m'] = big.NewInt(19)
	revAlph['n'] = big.NewInt(20)
	revAlph['o'] = big.NewInt(21)
	revAlph['p'] = big.NewInt(22)
	revAlph['q'] = big.NewInt(23)
	revAlph['r'] = big.NewInt(24)
	revAlph['s'] = big.NewInt(25)
	revAlph['t'] = big.NewInt(26)
	revAlph['u'] = big.NewInt(27)
	revAlph['w'] = big.NewInt(28)
	revAlph['x'] = big.NewInt(29)
	revAlph['y'] = big.NewInt(30)
	revAlph['z'] = big.NewInt(31)
	o := big.NewInt(0)
	radix := big.NewInt(32)

	for _, r := range in {
		o.Mul(o, radix)
		val, ok := revAlph[r]

		if !ok {
			return []byte{}, fmt.Errorf("'%c' is no legal base32 character", r)
		}

		o.Add(o, val)
	}

	return o.Bytes(), nil
}
