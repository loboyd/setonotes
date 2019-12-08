package user

/**
 * TODO: SECURITY-SENSITIVE -- Any of the code which handles unencrypted keys
 * should be re-implemented in the `encryption` package instead of here. This
 * mostly pertains to code in the `Create()` function. Ultimately, the following
 * functions should no longer be exported by `encryption`:
 *     NewSymmetricKey()
 *     GenerateKeyFromPassword()
 *     EncryptData()
 *     DecryptData()
 */

import (
    "log"
    "encoding/base64"
)

type User struct {
    ID                  int
    Username            string
    Email               string
    PasswordHash        []byte
    MainKeyEncrypted    []byte
    PrivateKeyEncrypted []byte // encryption service will handle marshaling
    PublicKey           []byte // same for public key
    EncryptionSalt      []byte
    AuthSalt            []byte
}

/**
 * The Repository interface contains all the user-relevant functions from the
 * storage repo
 */
type Repository interface {
    GetUserByID(id int) (*User, error)
    GetUserIDFromEmail(email string) (int, error) // returns userID
    GetUserIDFromUsername(username string) (int, error) // return userID
    CreateUser(u *User) (int, error) // returns userID
    TrackUserActivity(userID int, url string) error
    CheckBetaTesterWhitelist(username string) (bool, error)
}

/**
 * The EncryptService interface contains all the user-relevent encryption
 * functions
 */
type EncryptService interface {
    NewSymmetricKey() ([]byte, error)
    NewAssymetricKeyPair() ([]byte, []byte, error)
    NewSalt() ([]byte, error)
    GenerateKeyFromPassword(password, salt []byte) ([]byte, error)
    EncryptData(data, key []byte) ([]byte, error)
    DecryptData(data, key []byte) ([]byte, error)
    UserEncryptData(u *User, data []byte) ([]byte, error)
    UserDecryptData(u *User, data []byte) ([]byte, error)
}

/**
 * The AuthService interface contains all the auth-related functions
 */
type AuthService interface {
    HashAndSalt(password []byte) ([]byte, error)
}

type Service struct {
    repo       Repository
    encryption EncryptService
    auth       AuthService
}

/**
 * Creates a new User Service
 * The Service interface is implemented by the service struct with the receiver
 * functions defined below. Note the receiver takes a repository struct which
 * implements the Repository interface defined above
 */
func NewService(r Repository, e EncryptService, a AuthService) *Service {
    return &Service{
        repo: r,
        encryption: e,
        auth: a,
    }
}

/**
 * Given a user ID, returns the user
 */
func (s *Service) GetByID(userID int) (*User, error) {
    u, err := s.repo.GetUserByID(userID)
    if err != nil {
        return nil, err
    }
    return u, nil
}

/**
 * Given a user's username, returns the user
 */
func (s *Service) GetByUsername(username string) (*User, error) {
    userID, err := s.repo.GetUserIDFromUsername(username)
    if err != nil {
        return nil, err
    }

    u, err := s.repo.GetUserByID(userID) // returns nil user if err
    if err != nil {
        return nil, err
    }

    return u, nil
}

/**
 * Given a user's email, returns the user
 */
func (s *Service) GetByEmail(email string) (*User, error) {
    userID, err := s.repo.GetUserIDFromEmail(email)
    if err != nil {
        return nil, err
    }

    user, err := s.repo.GetUserByID(userID) // returns nil user if err
    return user, err
}

/**
 * Creates a new user in storage
 *
 * TODO: Decode `authSalt`, `encryptionSalt`, and `mainKeyEncrypted` from
 *   base-64 before casting to `[]byte`.
 *
 * TODO: This function should check a user doesn't yet exist and then build all
 * the shit necessary for a user given only username, email, and password
 * Also, check username, email, and password are all valid (actually, maybe this
 * should be done in the auth package with unexported functions)
 */
func (s *Service) Create(username, email, passwordStr, authSaltB64,
    encryptionSaltB64, mainKeyEncryptedB64 string) (*User, error) {

    // convert stuff from base-64 into byteslices
    authSalt, err         := base64.StdEncoding.DecodeString(authSaltB64)
    if err != nil {
        log.Printf("failed to decode authSalt from base-64: %v", err)
        return nil, err
    }
    encryptionSalt, err   := base64.StdEncoding.DecodeString(encryptionSaltB64)
    if err != nil {
        log.Printf("failed to decode encryptionSalt from base-64: %v", err)
        return nil, err
    }
    //mainKeyEncrypted, err := base64.StdEncoding.DecodeString(mainKeyEncryptedB64)
    _, err = base64.StdEncoding.DecodeString(mainKeyEncryptedB64)
    if err != nil {
        log.Printf("failed to decode mainKeyEncrypted from base-64: %v", err)
        return nil, err
    }

    // TODO: This needs to be changed such that the unencrypted main key is not
    // handled by any code outside of the encryption package
    // create main key
    mainKey, err := s.encryption.NewSymmetricKey()
    if err != nil {
        log.Printf("failed to create symmetric key: %v", err)
        return nil, err
    }

    // hash password
    password := []byte(passwordStr)
    passwordHash, err := s.auth.HashAndSalt(password)
    if err != nil {
        log.Printf("failed to hash and salt password: %v", err)
        return nil, err
    }

    // TODO: THIS FUNCTIONALITY NEEDS TO BE MOVED INTO THE ENCRYPTION PACKAGE
    // FOR SECURITY PURPOSES
    // create password-generated key
    passwordGeneratedKey, err := s.encryption.GenerateKeyFromPassword(password,
        encryptionSalt)
    if err != nil {
        log.Printf("failed to generate key from password: %v", err)
        return nil, err
    }

    // TODO: This will also be moved to the encryption package
    // encrypt main key
    mainKeyEncrypted, err := s.encryption.EncryptData(mainKey,
        passwordGeneratedKey)
    if err != nil {
        log.Printf("failed to encrypt main key: %v", err)
        return nil, err
    }

    u := &User{
        Username:            username,
        Email:               email,
        PasswordHash:        passwordHash,
        MainKeyEncrypted:    []byte(mainKeyEncrypted),
        //PrivateKeyEncrypted: privateKeyEncrypted,
        //PublicKey:           publicKey,
        EncryptionSalt:      []byte(encryptionSalt),
        AuthSalt:            []byte(authSalt),
    }
    userID, err := s.repo.CreateUser(u) // returns -1 userID if err
    u.ID = userID
    return u, err
}

/**
 * Stores userID, URL path, and timestamp
 */
func (s *Service) TrackActivity(userID int, url string) error {
    // timestamp is also tracked
    return s.repo.TrackUserActivity(userID, url)
}

/**
 * Check to see that username is on the beta-tester whitelist
 */
func(s *Service) CheckBetaTesterWhitelist(username string) (bool, error) {
    nameExists, err := s.repo.CheckBetaTesterWhitelist(username)
    return nameExists, err
}
