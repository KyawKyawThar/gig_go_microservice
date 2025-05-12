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

	fmt.Printf("min: %d, max: %d\n", min, max)
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

func pow10(n int) int64 {

	res := int64(1)

	for i := 0; i < n; i++ {
		res *= 10
	}

	return res
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

// GenerateOTP generates an n-digit numeric OTP (4-8 digits recommended)
func GenerateOTP(length int) (string, error) {

	if length < 3 || length > 8 {

		return "", fmt.Errorf("OTP length must be between 4 and 8 digits")
	}

	mini := int64(pow10(length - 1)) // e.g., 1000 for length=4
	maxi := int64(pow10(length) - 1) // e.g., 9999 for length=4

	otp := RandomInteger(mini, maxi)
	
	return fmt.Sprintf("%0*d", length, otp), nil // Ensures leading zeros

}
