package nanogo

import (
	"fmt"
	"golang.org/x/crypto/blake2b"
	"regexp"
	"strings"
)

// AddressIsValid checks if the wallet address is valid,
// address: the wallet address to check,
// returns true if the wallet address is valid, false otherwise.
func AddressIsValid(address string) bool {
	prefixes := []string{"nano", "xrb"}

	if address == "" {
		return false
	}

	validPrefix := false
	var prefix string

	for _, p := range prefixes {
		if strings.HasPrefix(address, p) {
			if len(address) == 61+len(p) {
				validPrefix = true
				prefix = p
				break
			}
		}
	}

	if !validPrefix {
		return false
	}

	pattern := fmt.Sprintf(`^(%s)_[13]{1}[13456789abcdefghijkmnopqrstuwxyz]{59}$`, prefix)
	re := regexp.MustCompile(pattern)

	if !re.MatchString(address) {
		return false
	}

	pubKey, err := AddressToPublicKey(address)

	if err != nil {
		return false
	}

	withoutPrefix := strings.ReplaceAll(address, prefix+"_", "")
	origSum := withoutPrefix[len(withoutPrefix)-8 : len(withoutPrefix)]
	h, err := blake2b.New(5, nil)

	if err != nil {
		return false
	}

	h.Write(pubKey[:])

	sumBytes := h.Sum(nil)
	sum := base32Encode(revertBytes(sumBytes[:5]))

	if sum != origSum {
		return false
	}

	return true
}
