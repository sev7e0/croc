package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

type encryption struct {
	key        []byte
	passphrase []byte
	salt       []byte
}

// New generates a new encryption, using the supplied passphrase and
// an optional supplied salt.
func New(passphrase []byte, salt []byte) (e encryption, err error) {
	e.passphrase = passphrase
	if salt == nil {
		e.salt = make([]byte, 8)
		// http://www.ietf.org/rfc/rfc2898.txt
		// Salt.
		rand.Read(e.salt)
	} else {
		e.salt = salt
	}
	e.key = pbkdf2.Key([]byte(passphrase), e.salt, 100, 32, sha256.New)
	return
}

func (e encryption) Salt() []byte {
	return e.salt
}

// Encrypt will generate an encryption, prefixed with the IV
func (e encryption) Encrypt(plaintext []byte) []byte {
	// generate a random iv each time
	// http://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38d.pdf
	// Section 8.2
	ivBytes := make([]byte, 12)
	rand.Read(ivBytes)
	b, _ := aes.NewCipher(e.key)
	aesgcm, _ := cipher.NewGCM(b)
	encrypted := aesgcm.Seal(nil, ivBytes, plaintext, nil)
	return append(ivBytes, encrypted...)
}

// Decrypt an encryption
func (e encryption) Decrypt(encrypted []byte) (plaintext []byte, err error) {
	b, _ := aes.NewCipher(e.key)
	aesgcm, _ := cipher.NewGCM(b)
	plaintext, err = aesgcm.Open(nil, encrypted[:12], encrypted[12:], nil)
	return
}
