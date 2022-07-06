package pwgen

import (
	"crypto/rand"
	"math/big"
	"strings"
)

type PasswordGenerator interface {
	Generate(length int) (string, error)
}

type generator struct {
	alphabet string
}

var DefaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;':,./<>?"

// NewPasswordGenerator creates a new PasswordGenerator
func NewPasswordGenerator(alphabet string) *generator {
	return &generator{
		alphabet: alphabet,
	}
}

// Generate generates a password of the given length
func (g *generator) Generate(length int) (string, error) {
	var generated string
	characterSet := strings.Split(g.alphabet, "")
	max := big.NewInt(int64(len(characterSet)))
	for i := 0; i < length; i++ {
		val, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		generated += characterSet[val.Int64()]
	}
	return generated, nil
}
