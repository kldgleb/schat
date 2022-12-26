package schat

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var key []byte

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("errror while reading env: %s", err.Error())
	}
	k := os.Getenv("KEY")
	key = []byte(k)
}

// cipher key
func encrypt(key []byte, text []byte) []byte {
	c, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(c)

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	result := gcm.Seal(nonce, nonce, text, nil)

	return result
}
func decrypt(key []byte, ciphertext []byte) string {
	c, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(c)

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		panic("ciphertext size is less than nonceSize")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)

	return string(plaintext)
}
