package pwgen

import (
	"testing"
)

var testingLengths = []int{8, 12, 16, 32, 100}

// testGenerate tests the Generate function
func testGenerate(t *testing.T, generator PasswordGenerator, length int) string {
	t.Helper()
	actual, err := generator.Generate(length)
	if err != nil {
		t.Errorf("Error generating password: %v", err)
	}
	if len(actual) != length {
		t.Errorf("Expected password of length %d, got %d", length, len(actual))
	}
	return actual
}

func Test_generator_Generate(t *testing.T) {
	t.Helper()

	var pwd string

	for _, i := range testingLengths {
		newPwd := testGenerate(t, NewPasswordGenerator(DefaultAlphabet), i)
		if newPwd == pwd {
			t.Errorf("Expected different password, got %s", pwd)
		}
		pwd = newPwd
	}

}
