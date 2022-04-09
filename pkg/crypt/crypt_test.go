package crypt

import (
	"testing"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

var testingBadKey = "aje3Phaemoogh8poo6xah4EijeukeuNga"
var testingGoodKey = "ulief3Nea4biochiwe0wieh1Phech8ay"

var testingContent = ""

func TestAESGCM(t *testing.T) {
	_, err := NewAESGCM(testingBadKey)
	assert.Error(t, err)

	a, err := NewAESGCM(testingGoodKey)
	assert.NoError(t, err)

	faker := faker.New()
	content := faker.Lorem().Text(1024)

	ciphertext, err := a.Encrypt([]byte(content))
	assert.NoError(t, err)

	plaintext, err := a.Decrypt(ciphertext)
	assert.NoError(t, err)

	stringPlaintext := string(plaintext)
	assert.Equal(t, content, stringPlaintext)

}
