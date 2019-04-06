package encryption

/**
 * This file will contain all encryption functionality concerning a particular
 * user. This is the only code allowed to handle the cached password-generated
 * keys or unencrypted main-keys.
 */

import (
    "strconv"

    "github.com/setonotes/pkg/user"
)

/**
 * Get a user's password-generated key from the cache (notice that this function
 * is not exported -- ideally, the password-generated key should not be passed
 * outside this package for any reason)
 */
func (s *Service) getPasswordGeneratedKey(userID int) ([]byte, error) {
    key, err := s.cacheService.GetString("pgkey_"+strconv.Itoa(userID))
    if err != nil {
        return nil, err
    }
    return []byte(key), nil
}

/**
 * Generate a new symmetric key and encrypt with the user's main-key before
 * returning
 */
func (s *Service) NewUserEncryptedSymmetricKey(u *user.User) ([]byte, error) {
    // generate new symmetric key
    key, err := s.NewSymmetricKey()
    if err != nil {
        return nil, err
    }

    // user-encrypt
    keyEncrypted, err := s.UserEncryptData(u, key)
    if err != nil {
        return nil, err
    }

    return keyEncrypted, nil
}

/**
 * Encrypt data for a particular user --
 * This funciton gets the user's password-generated key from the cache, uses it
 * to decrypt their main-key, and uses the main-key to encrypt the data.
 */
func (s *Service) UserEncryptData(u *user.User, data []byte) ([]byte, error) {
    // get the user's password-generated key
    passwordGeneratedKey, err := s.getPasswordGeneratedKey(u.ID)
    if err != nil {
        return nil, err
    }

    // decrypt user's main-key
    mainKey, err := s.DecryptData(u.MainKeyEncrypted, passwordGeneratedKey)
    if err != nil {
        return nil, err
    }

    // encrypt data
    result, err := s.EncryptData(data, mainKey)
    if err != nil {
        return nil, err
    }

    return result, nil
}

/**
 * Decrypt data for a particular user --
 * This funciton gets the user's password-generated key from the cache, uses it
 * to decrypt their main-key, and uses the main-key to decrypt the data.
 */
func (s *Service) UserDecryptData(u *user.User, data []byte) ([]byte, error) {
    // get the user's password-generated key
    passwordGeneratedKey, err := s.getPasswordGeneratedKey(u.ID)
    if err != nil {
        return nil, err
    }

    // decrypt user's main-key
    mainKey, err := s.DecryptData(u.MainKeyEncrypted, passwordGeneratedKey)
    if err != nil {
        return nil, err
    }

    // decrypt data
    result, err := s.DecryptData(data, mainKey)
    if err != nil {
        return nil, err
    }

    return result, nil
}
