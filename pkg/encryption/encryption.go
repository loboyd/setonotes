package encryption

import (
    "log"
    "io"
    "crypto/sha256"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"

    "golang.org/x/crypto/pbkdf2"
)

type CacheService interface {
    GetString(key interface{}) (string, error)
}

type Service struct{
    cacheService CacheService
}

func NewService(c CacheService) *Service {
    return &Service{c}
}

/**
 * Creates a new 128-bit symmetric encryption key
 *
 * TODO: SECURITY-SENSITIVE -- This should not be exported. This package should
 * be the only package with the privilege of creating (and handling) an
 * unencrypted key.
 */
func (s *Service) NewSymmetricKey() ([]byte, error) {
    key , err := getRandomBytes(16)
    if err != nil {
        return nil, err
    }
    return key, nil
}

/**
 * Creates a new 4096-bit RSA public key-pair
 * TODO: Implement this (the entire public-key system needs to be more carefully
 * designed and probably won't use RSA ultimately)
 */
func (s *Service) NewAssymetricKeyPair() ([]byte, []byte, error) {
    log.Println("attempted to create assymetric key-pair; not implemented")
    return nil, nil, nil
}

/**
 * Creates a new 128-bit salt
 */
func (s *Service) NewSalt() ([]byte, error) {
    salt, err := getRandomBytes(16)
    if err != nil {
        return nil, err
    }
    return salt, nil
}

/**
 * Read n bytes from crypto/rand` reader
 */
func getRandomBytes(n int) ([]byte, error) {
    data := make([]byte, n)
    if _, err := io.ReadFull(rand.Reader, data); err != nil {
        log.Printf("failed to read from crypto/rand reader: %v", err)
        return nil, err
    }
    return data, nil
}

/**
 * Generates a 128-bit encryption key given a user-specific salt and password
 * using PBKDF2
 *
 * TODO: SECURITY-SENSITIVE -- This should not be exported. This package should
 * be the only package with the privilege of generating (and handling) an
 * unencrypted key.
 */
func (s *Service) GenerateKeyFromPassword(password,
    salt []byte) ([]byte, error) {

    // use 3e5 iterations, 16 bytes == 128 bits
    // error is returned because it's not clear why pdkdf2 does not return an
    // error, and the error might be useful for forward-compatibility
    return pbkdf2.Key(password, salt, 3e5, 16, sha256.New), nil
}

/**
 * Encrypt byteslice of data using AES encryption
 * See https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-
 *     golang-application-crypto-packages/
 *
 * TODO: This function probably doesn't need to be exported because no caller
 * outside this package should be able to supply an unencrypted key anyway.
 */
func (s *Service) EncryptData(data, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        log.Printf("failed to create new AES cipher for encryption: %v", err)
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        log.Printf("failed to create new GCM for encryption: %v", err)
        return nil, err
    }

    nonce, err := getRandomBytes(gcm.NonceSize())
    if err != nil {
        return nil, err
    }

    // in some places the below line is this instead:
    // ciphertext := gcm.Seal(nil, nonce, data, nil)
    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    return ciphertext, nil
}

/**
 * Decrypt byteslice of data that was encrypted by encryptData()
 * See https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-
 *     golang-application-crypto-packages/
 *
 * TODO: This function probably doesn't need to be exported because no caller
 * outside this package should be able to supply an unencrypted key anyway.
 */
func (s *Service) DecryptData(data, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        log.Printf("failed to create new AES cipher for decrypttion: %v", err)
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        log.Printf("failed to create new GCM for decryption: %v", err)
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        log.Printf("failed to open GCM for decryption: %v", err)
        return nil, err
    }

    return plaintext, nil
}
