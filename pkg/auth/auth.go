package auth

import (
    "log"
    "time"
    "strconv"
    "net/http"
    "crypto/sha256"
    "strings"
    "errors"

    "github.com/setonotes/pkg/user"

    "github.com/satori/go.uuid"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/crypto/pbkdf2"
)

type Cache interface {
    Set(key, value interface{}) error
    SetEx(key, value interface{}, lifetime int) error
    GetInt(key interface{}) (int, error)
    GetString(key interface{}) (string, error)
    Delete(key interface{}) error
}

type Service struct {
    sessionCache Cache
}

func NewService(sessionCache Cache) *Service {
    return &Service{sessionCache: sessionCache}
}

/**
 * Initialize a user session by storing a session token in the session cache,
 * storing a cookie on the user's browser and storing the password-generated key
 * in the cache
 */
func (s *Service) InitUserSession(w http.ResponseWriter, r *http.Request,
    u *user.User, password []byte) error {

    // create cache session token
    log.Println("creating new UUID session token...")
    sessionTokenTmp, err := uuid.NewV4()
    if err != nil {
        log.Println("failed to create new UUID session token")
        return err
    }
    log.Println("successfully created new UUID session token")
    log.Println("storing session token in cache with expiration 1 day...")
    sessionToken := sessionTokenTmp.String()
    err = s.sessionCache.SetEx(sessionToken, u.ID, 86400) // 86400s == 1 day
    if err != nil {
        log.Println("failed to store session token in cache")
        return err
    }
    log.Println("successfully stored session token in cache")

    // create session cookie on user's browser
    log.Println("setting cookie on user's brower...")
    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    sessionToken,
        Expires:  time.Now().Add(86400 * time.Second),
        Path:     "/",
        //HttpOnly: true,
    })

    // generate password-generated key
    log.Println("generating key from password...")
    key, err := s.generateKeyFromPassword([]byte(password), u.EncryptionSalt)
    if err != nil {
        log.Println("failed to generate key from password")
        return err
    }
    log.Println("successfully generated key from password; storing in cache...")

    // set key in session cache
    err = s.sessionCache.SetEx("pgkey_"+strconv.Itoa(u.ID), key, 86400)
    if err != nil {
        log.Println("failed to store key in cache")
        return err
    }
    log.Println("successfully stored key in cache")

    // increment user session count
    log.Println("incrementing user's session count")
    cacheString := "n_sessions_" + strconv.Itoa(u.ID)
    numberUserSessions, err := s.sessionCache.GetInt(cacheString)
    if err != nil {
        log.Println("count not found; starting new count at 1...")
        err2 := s.sessionCache.Set(cacheString, 1)
        if err2 != nil {
            // TODO: what action needs to be taken here?
            log.Println("failed to start new count")
        }
        log.Println("successfully started count")
    } else {
        log.Println("count found; incrementing count...")
        err2 := s.sessionCache.Set(cacheString, numberUserSessions + 1)
        if err2 != nil {
            // TODO: what action needs to be taken here?
            log.Println("failed to increment count")
        }
        log.Println("successfully incremented count")
    }

    return nil
}

/**
 * End the user session by removing their session token from the cache,
 * decrementing their session count (used for knowing when it is okay to remove
 * their password-generated key from the cache), and overwrite the cookie on
 * their browser with an immediately-expiring cookie
 */
func (s *Service) EndUserSession(w http.ResponseWriter, r *http.Request,
    userID int) error {

    // look for cookie on user's browser
    c, err := r.Cookie("session_token")
    if err != nil {
        if err == http.ErrNoCookie {
            http.Redirect(w, r, "/", http.StatusFound)
            return err
        }
        w.WriteHeader(http.StatusBadRequest)
    }
    sessionToken := c.Value

    // delete cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    sessionToken,
        MaxAge:   -1,
        Path:     "/",
        HttpOnly: true,
    })

    // look for token in Redis cache
    _, err = s.sessionCache.GetInt(sessionToken)
    if err != nil {
        http.Redirect(w, r, "/", http.StatusFound)
        return err
    }

    // delete user session from cache
    err = s.sessionCache.Delete(sessionToken)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return err
    }

    // delete user session from cache
    err = s.sessionCache.Delete(sessionToken)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return err
    }

    // get user session count
    redisString := "n_sessions_" + strconv.Itoa(userID)
    numberUserSessions, err := s.sessionCache.GetInt(redisString)
    if err != nil {
        log.Printf("failed to get cached session count for user-%v sign-out",
            userID)

        // delete the password-generated encryption key from the cache
        s.sessionCache.Delete("pgkey_"+strconv.Itoa(userID))
        // if err2 != nil {
            // TODO:
            // the case where the key is not found probably doesn't need to be
            // handled; what about other errors? There should be a log here
            // w.WriteHeader(http.StatusInternalServerError)
            // return
        // }
    } else if numberUserSessions <= 1 {
        log.Printf("session count for user-%v is now 0", userID)
        s.sessionCache.Set(redisString, 0)

        log.Printf("deleting password-generated key for user-%v...", userID)
        s.sessionCache.Delete("pgkey_"+strconv.Itoa(userID))
    } else {
        log.Printf("decrementing session count for user-%v...", userID)
        s.sessionCache.Set(redisString, numberUserSessions - 1)
    }

    return nil
}

/**
 * Check client's authentication status given a session token (UUID string)
 */
func (s *Service) checkSessionToken(sessionToken string) (int, bool, error) {
    log.Println("checking session token...")

    // look for token in cache
    log.Println("looking for session token in cache...")
    response, err := s.sessionCache.GetInt(sessionToken)
    if err != nil {
        log.Println("failed to get session token from cache")
        return 0, false, err
    }
    log.Println("successfully got session token from cache")

    // return true if no issues above
    return response, true, nil
}

/**
 * Check user's authentication status given *http.Request containing
 * `session_token` cookie
 */
func (s *Service) CheckAuthStatusCookie(r *http.Request) (int, bool, error) {
    log.Println("looking for session token cookie...")
    c, err := r.Cookie("session_token")
    if err != nil {
        log.Println("failed to find session token cookie")
        return 0, false, err
    }
    log.Println("successfully found session token cookie")
    sessionToken := c.Value

    return s.checkSessionToken(sessionToken)
}

/**
 * Check user's authenciation status given *http.Request containing bearer
 * token `Authorization` header
 */
func (s *Service) CheckAuthStatusBearer(r *http.Request) (int, bool, error) {
    // get bearer token
    log.Println("looking for session token header...")
    sessionToken := r.Header.Get("Authorization")
    splitToken := strings.Split(sessionToken, "Bearer")
    if len(splitToken) != 2 {
        return 0, false, errors.New("improper bearer token format");
    }
    log.Println("successfully found bearer token")
    sessionToken = strings.TrimSpace(splitToken[1])

    return s.checkSessionToken(sessionToken)
}

/**
 * Hash and salt a user's password using Bcrypt
 * see https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72
 */
func (s *Service) HashAndSalt(password []byte) ([]byte, error) {
    hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
    if err != nil {
        log.Printf("bcrypt hash+password comparison failure: %v", err)
        return nil, err
    }
    return hash, nil
}

/**
 * Check the password and hash match
 */
func (s *Service) CheckPassHash(hash, password []byte) (bool, error) {
    err := bcrypt.CompareHashAndPassword(hash, password)
    if err != nil {
        log.Printf("failed to compare hash and password: %v", err)
        return false, err
    }
    return true, nil // bcrypt will error upon failed comparison
}

/**
 * Get password-generated encryption key from cache for a given user-ID
 */
func (s *Service) GetPasswordGeneratedKey(userID int) ([]byte, error) {
    key, err := s.sessionCache.GetString("pgkey_"+strconv.Itoa(userID))
    if err != nil {
        return nil, err
    }
    return []byte(key), nil
}

/**
 * Generates a 128-bit encryption key given a user-specific salt and password
 * using PBKDF2
 */
func (s *Service) generateKeyFromPassword(password,
    salt []byte) ([]byte, error) {

    // use 3e5 iterations, 16 bytes == 128 bits
    // error is returned because it's not clear why pdkdf2 does not return an
    // error, and the error might be useful for forward-compatibility
    return pbkdf2.Key(password, salt, 3e5, 16, sha256.New), nil
}
