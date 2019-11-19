package postgres

/**
 * This file contains user-related repository functions
 */

import (
    "log"
    "time"

    "github.com/setonotes/pkg/user"
    "github.com/setonotes/pkg/page"
)

/**
 * Returns a user given the user's id
 */
func (r *Repository) GetUserByID(userID int) (*user.User, error) {
    var (
        username            string
        email               string
        passwordHash        []byte
        mainKeyEncrypted    []byte
        privateKeyEncrypted []byte
        publicKey           []byte
        salt                []byte
        authSalt            []byte
        version             int
    )

    // query database for user-fields
    psqlStmt := `
        SELECT
            username,
            email,
            password_hash,
            main_key_encrypted,
            private_key_encrypted,
            public_key,
            salt,
            auth_salt,
            version
        FROM users
        WHERE id=$1`
    err := r.DB.QueryRow(psqlStmt, userID).Scan(
        &username,
        &email,
        &passwordHash,
        &mainKeyEncrypted,
        &privateKeyEncrypted,
        &publicKey,
        &salt,
        &authSalt,
        &version,
    )
    if err != nil {
        log.Printf("failed to get user-%v from storage: %v", userID, err)
        return nil, err
    }

    return &user.User{
        ID:                  userID,
        Username:            username,
        Email:               email,
        PasswordHash:        passwordHash,
        MainKeyEncrypted:    mainKeyEncrypted,
        PrivateKeyEncrypted: privateKeyEncrypted,
        PublicKey:           publicKey,
        EncryptionSalt:      salt,
        AuthSalt:            authSalt,
        Version:             version,
    }, nil
}

/**
 * Returns userID corresponding to given username
 */
func (r *Repository) GetUserIDFromUsername(username string) (int, error) {
    psqlStmt := `
        SELECT id
        FROM users
        WHERE username=$1`
    var userID int
    err := r.DB.QueryRow(psqlStmt, username).Scan(&userID)
    if err != nil {
        return -1, err
    }
    return userID, err
}

/**
 * Returns userID corresponding to given email address
 */
func (r *Repository) GetUserIDFromEmail(email string) (int, error) {
    psqlStmt := `
        SELECT id
        FROM users
        WHERE email=$1`
    var userID int
    err := r.DB.QueryRow(psqlStmt, email).Scan(&userID)
    if err != nil {
        return -1, err
    }
    return userID, err
}

/**
 * Stores all user.User fields in a database row
 * This function assumes a user does not yet exist (this should be checked
 * by the caller)
 */
func (r *Repository) CreateUser(user *user.User) (int, error) {
    psqlStmt := `
        INSERT INTO users (
            username,
            email,
            password_hash,
            main_key_encrypted,
            /*private_key_encrypted,
              public_key,*/
            salt,
            auth_salt,
            version)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`
    var userID int
    err := r.DB.QueryRow(psqlStmt,
        user.Username,
        user.Email,
        user.PasswordHash,
        user.MainKeyEncrypted,
        //user.PrivateKeyEncrypted,
        //user.PublicKey,
        user.EncryptionSalt,
        user.AuthSalt,
        user.Version,
    ).Scan(&userID)
    if err != nil {
        log.Printf("failed to create row in `users`: %v", err)
        return -1, err
    }

    return userID, nil
}

/**
 * Tracks user ID, URL path and timestamp for each authorized HTTP request
 */
func (r *Repository) TrackUserActivity(userID int, url string) error {
    psqlStmt := `
        INSERT INTO user_activity (user_id, url, timestamp)
        VALUES ($1, $2, $3)`
    _, err := r.DB.Exec(psqlStmt, userID, url, time.Now())
    if err != nil {
        log.Printf("failed to write user activity")
        return err
    }
    return nil
}

/**
 * Get all (disembodied) pages for which userID has read-permission
 *
 * See https://www.calhoun.io/querying-for-multiple-records-with-gos-sql-
 * package/ for querying multiple records
 */
func (r *Repository) GetUserDisembodiedPages(userID int) ([]*page.Page, error) {
    psqlStmt := `
        SELECT id, title, version, author_id
        FROM pages JOIN page_permissions
        ON (pages.id=page_permissions.page_id)
        WHERE user_id=$1`
    log.Printf("getting page rows for user-%v from DB...", userID)
    rows, err := r.DB.Query(psqlStmt, userID)
    if err != nil {
        log.Printf("failed to get page rows for user-%v from DB", userID)
        return nil, err
    }
    log.Printf("successfully got page rows for user-%v from DB", userID)
    defer rows.Close()

    // loop over rows and create array of pages
    var pages = []*page.Page{}
    for rows.Next() {
        var (
            pageID         int
            titleEncrypted []byte
            ownerID        int
            version        int
        )
        log.Println("scanning row for page ID, title, owner ID, version...")
        err = rows.Scan(&pageID, &titleEncrypted, &version, &ownerID)
        if err != nil {
            log.Println("failed to get disembodied page from row")
            return nil, err
        }
        log.Println("successfully got disembodied page from row")

        // append page to array
        pages = append(pages, &page.Page{
            ID:      pageID,
            Title:   titleEncrypted,
            Body:    []byte(""),
            OwnerID: ownerID,
            Version: version,
        })
    }

    // get any errors encountered during iteration
    err = rows.Err()
    if err != nil {
        log.Println("error iterating over rows in database")
        return nil, err
    }

    return pages, nil
}

/**
 * Check for beta name in beta-tester username whitelist
 */
func (r *Repository) CheckBetaTesterWhitelist(username string) (bool, error) {
    nameExists := false
    psqlStmt := `
        SELECT EXISTS(SELECT 1 FROM beta_testers
        WHERE username=$1)`
    err := r.DB.QueryRow(psqlStmt, username).Scan(&nameExists)
    if err != nil {
        log.Println("failed to check beta-tester whitelist")
        return false, err
    }
    return nameExists, nil
}
