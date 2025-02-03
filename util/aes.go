package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
)

const (
	saltSize   = 16
	keySize    = 32
	nonceSize  = 12
	iterations = 100000
)

func GenerateRandomBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, bytes)
	return bytes, err
}

func DeriveKey(password, salt []byte) []byte {
	return pbkdf2Key(password, salt, iterations, keySize, sha256.New)
}

// AES_GCM_Encrypt encrypts plaintext with password using AES-GCM
func AES_GCM_Encrypt(password, plaintext []byte) ([]byte, error) {
	// Generate random salt and derive key using PBKDF2
	salt, err := GenerateRandomBytes(saltSize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %v", err)
	}
	key := DeriveKey(password, salt)

	// Generate random nonce for AES-GCM
	nonce, err := GenerateRandomBytes(nonceSize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}

	// Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %v", err)
	}

	// Encrypt and append nonce and salt
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// Concatenate salt, nonce, and ciphertext
	buffer := bytes.NewBuffer(salt)
	buffer.Write(nonce)
	buffer.Write(ciphertext)

	return buffer.Bytes(), nil
}

// AES_GCM_Decrypt decrypts ciphertext with password using AES-GCM
func AES_GCM_Decrypt(password, ciphertext []byte) ([]byte, error) {
	// Extract salt, nonce, and actual ciphertext
	salt := ciphertext[:saltSize]
	nonce := ciphertext[saltSize : saltSize+nonceSize]
	encMessage := ciphertext[saltSize+nonceSize:]

	// Derive key from password and salt
	key := DeriveKey(password, salt)

	// Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %v", err)
	}

	// Decrypt the ciphertext
	plaintext, err := aesgcm.Open(nil, nonce, encMessage, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %v", err)
	}

	return plaintext, nil
}

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Using a higher iteration count will increase the cost of an exhaustive
// search but will also make derivation proportionally slower.
func pbkdf2Key(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}
