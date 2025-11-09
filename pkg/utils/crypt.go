package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mxmrykov/polonium-auth/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func Hash(val string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(val), 7)
	return string(bytes), err
}

func CheckHash(val, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(val))
	return err == nil
}

func NewCert() string {
	id := []byte(uuid.New().String())
	sshDecoded := sha512.Sum512(id)
	return convertToSymbolPairs(string(sshDecoded[:]), " ")
}

func GetTOTPCodes(secret string, duration int) []string {
	timeStep := time.Now().Unix() / int64(duration)

	prev := generateTOTP([]byte(secret), timeStep-1)
	cur := generateTOTP([]byte(secret), timeStep)
	next := generateTOTP([]byte(secret), timeStep+1)

	return []string{prev, cur, next}
}

func GenerateRSAKeys(bits int) (*model.RSAKeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA keys: %w", err)
	}

	// Validate the key
	err = privateKey.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate RSA key: %w", err)
	}

	return &model.RSAKeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

func generateTOTP(secret []byte, timeInterval int64) string {
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(timeInterval))

	hmacHash := hmac.New(sha1.New, secret)
	hmacHash.Write(timeBytes)
	hash := hmacHash.Sum(nil)

	offset := hash[len(hash)-1] & 0xf
	binaryCode := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7fffffff

	return fmt.Sprintf("%06d", binaryCode%1000000)
}

func convertToSymbolPairs(s string, separator string) string {
	if len(s) == 0 {
		return ""
	}

	var result strings.Builder
	length := len(s)

	for i := 0; i < length; i += 2 {
		end := i + 2
		if end > length {
			end = length
		}

		pair := s[i:end]

		if i > 0 {
			result.WriteString(separator)
		}

		result.WriteString(pair)
	}

	return result.String()
}
