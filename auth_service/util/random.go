package util

import (
	crypto "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	t := time.Now().Unix()
	rand.New(rand.NewSource(t))
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(maxLength int) string {
	b := make([]byte, maxLength)
	//fmt.Printf("before: %s\n", string(b))
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	//fmt.Printf("%s\n", string(b))

	return string(b)
}

func RandomInteger(min, max int64) int64 {

	if min > max {
		panic("min cannot be greater than max")
	}

	// Generate a random number within the range [0, max-min]
	randomNumber := rand.Int63n(max - min + 1)

	// Add the minimum value to shift the range
	return randomNumber + min
}

func RandomOwner() string {
	return RandomString(7)
}

func RandomAmount() int64 {
	return RandomInteger(1, 100)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomOwner())
}

// RandomHex generates a random hex string of the specified byte length
func RandomHex(byteLength int) (string, error) {

	bytes := make([]byte, byteLength)

	if _, err := crypto.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
