package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Generates a random array of bytes (using [crypto/rand]) of the given length.
// Panics if the length is less than or equal to 0.
//
// [crypto/rand]: https://pkg.go.dev/crypto/rand#example-Int
func GenerateSalt(length int) []byte {
	if length <= 0 {
		panic(fmt.Sprintf("invalid salt length - must be greater than 0 (was: %d)", length))
	}
	salt := make([]byte, length)
	byteMax := big.NewInt(256)
	for i := range length {
		// According to Go's own docs, rand.Int can't return an error if rand.Reader is used:
		// https://pkg.go.dev/crypto/rand#example-Int
		v, _ := rand.Int(rand.Reader, byteMax)
		salt[i] = byte(v.Int64())
	}

	return salt
}

// Performs SHA256 hash of the given salt + input.
// Returns an error if hash.Write returns an error, but this generally shouldn't
// be possible. See [implementation].
//
// [implementation]: https://cs.opensource.google/go/go/+/refs/tags/go1.25.6:src/crypto/internal/fips140/sha256/sha256.go;l=172
func Hash(salt, input []byte) ([]byte, error) {
	hash := sha256.New()
	data := append(salt, input...)
	if _, err := hash.Write(data); err != nil {
		return nil, fmt.Errorf("failed to initialize hash: %w", err)
	}
	return hash.Sum(nil), nil
}

// Decodes a hex-encoded string to a byte slice. Panics if invalid.
func MustDecodeHex(hstr string) []byte {
	bs, err := hex.DecodeString(hstr)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode required hex string: %v", err))
	}
	return bs
}

// Performs SHA256 hash of the given salt + input.
// Panics if an error is encountered, but generally that shouldn't be possible.
func MustHash(salt, input []byte) []byte {
	bs, err := Hash(salt, input)
	if err != nil {
		panic(fmt.Sprintf("failed to hash necessary input: %v", err))
	}
	return bs
}
