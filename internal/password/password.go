package password

import (
	"crypto/rand"
	"errors"
	"math"
	"math/big"
	"strings"
)

var ErrEntropyTooLow = errors.New(
	"the entropy of the requested password is too low; please request a password with a larger length",
)

func calculateEntropy(lenCharset, lenPwd int) error {
	minEntropy := 60 // generally acceptable amount of entropy

	// Entropy = log2(lenCharset ** lenPwd)
	entropy := int(math.Round(math.Log2(math.Pow(float64(lenCharset), float64(lenPwd)))))

	if entropy < minEntropy {
		return ErrEntropyTooLow
	}

	return nil
}

func Generate(lenPwd int) (string, error) {
	chSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcefghijklmopqrstuvwxyz0123456789!@#$%^&*()"
	lenChSet := len(chSet)

	if entropyErr := calculateEntropy(lenChSet, lenPwd); entropyErr != nil {
		return "", entropyErr
	}

	var b strings.Builder
	for i := 0; i <= lenPwd-1; i++ {
		rndIdx, rndErr := rand.Int(rand.Reader, big.NewInt(int64(lenChSet)))
		if rndErr != nil {
			// rand.Int returns an error if rand.Read returns 1. Let's discard the result
			// by decrementing i and trying again.
			i--
			continue
		}

		b.WriteByte(chSet[rndIdx.Int64()])
	}

	return b.String(), nil
}
