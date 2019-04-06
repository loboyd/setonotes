package postgres

/**
 * This file contains permission-related repository functions
 */

import (
    "log"
)

/**
 * Creates a page permission row in the database
 */
func (r *Repository) CreatePagePermission(userID, pageID int, isOwner,
    canEdit bool, userEncryptedPageKey []byte) error {

    // create entry in `page_permissions`
    log.Println("creating new page permission row in DB...")
    psqlStmt := `
        INSERT INTO page_permissions (user_id, page_id, is_owner,
            can_edit, user_encrypted_page_key)
        VALUES ($1, $2, $3, $4, $5)`
    _, err := r.DB.Exec(psqlStmt, userID, pageID, isOwner, canEdit,
        userEncryptedPageKey)
    if err != nil {
        log.Println("failed to create new page permissions row in DB")
        return err
    }

    return nil
}

/**
 * Get user-encrypted page key given userID, pageID
 */
func (r *Repository) GetUserEncryptedPageKey(userID,
    pageID int) ([]byte, error) {

    psqlStmt := `
        SELECT user_encrypted_page_key
        FROM page_permissions
        WHERE user_id=$1 AND page_id=$2`
    var key []byte
    err := r.DB.QueryRow(psqlStmt, userID, pageID).Scan(&key)
    if err != nil {
        log.Printf("failed to get user-%v's page-%v key from DB: %v", userID,
            pageID, err)
        return nil, err
    }

    return key, nil
}

/**
 * Check userID can edit pageID
 */
func (r *Repository) CheckUserCanEditPage(userID, pageID int) (bool, error) {
    psqlStmt := `
        SELECT can_edit FROM page_permissions
        WHERE user_id=$1 AND page_id=$2`
    var canEdit bool
    err := r.DB.QueryRow(psqlStmt, userID, pageID).Scan(&canEdit)
    if err != nil {
        log.Printf("failed to check user-%v, page-%v read permission", userID,
            pageID)
        return false, err
    }

    return canEdit, nil
}
