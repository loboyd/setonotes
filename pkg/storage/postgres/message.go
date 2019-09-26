package postgres

/**
 * This file contains message-related repository functions
 */

import (
    "log"

    //"github.com/setonotes/pkg/page"
    //"github.com/setonotes/pkg/user" // for current version number
)

/** TODO: Adapt this for messages
 * Given a page ID, check there exists a database row in the `pages` table with
 * that ID
func (r *Repository) CheckPageExists(pageID int) (bool, error) {
    pageExists := false
    psqlStmt := `
        SELECT EXISTS(
        SELECT 1 FROM pages
        WHERE id=$1)`
    err := r.DB.QueryRow(psqlStmt, pageID).Scan(&pageExists)
    if err != nil {
        return false, err
    }
    return pageExists, nil
}
 */

/**
 * Create new row in the `messages` table.
func (r *Repository) CreateMesssage(p *message.Message,
    authorID int) (int, error) {
 */
func (r *Repository) CreateMessage(from_user int, to_user int,
    body []byte) (int, error) {

    log.Println("creating row in `messages` table...")
    psqlStmt := `
        INSERT INTO messages (from_user, to_user, body)
        VALUES ($1, $2, $3)
        RETURNING id`
    messageID := 0
    err := r.DB.QueryRow(psqlStmt, from_user, to_user, body).Scan(&messageID)
    if err != nil {
        log.Println("failed to store message")
        return 0, err
    }
    return messageID, nil
}

/** TODO: Adapt this for messages
 * Update Title and Body of existing page
func (r *Repository) UpdatePage(p *page.Page) error {
    log.Printf("updating row for page-%v", p.ID)
    psqlStmt := `
        UPDATE pages
        SET title=$1, body=$2, version=$3
        WHERE id=$4`
    _, err := r.DB.Exec(psqlStmt, p.Title, p.Body, p.Version, p.ID)
    if err != nil {
        log.Println("failed to updated row for page-%v", p.ID)
        return err
    }
    log.Println("successfully updated row for page-%v", p.ID)
    return nil
}
 */

/** TODO: Adapt this for pages
 * Delete a page
func (r *Repository) DeletePage(pageID int) error {
    log.Printf("deleting page-%v row from pages table...", pageID)
    psqlStmt := `
        DELETE FROM pages
        WHERE id=$1`
    _, err := r.DB.Exec(psqlStmt, pageID)
    if err != nil {
        log.Printf("failed to delete page-%v row from database", pageID)
    }
    log.Printf("successfully deleted page-%v row from database", pageID)
    return nil
}
 */

