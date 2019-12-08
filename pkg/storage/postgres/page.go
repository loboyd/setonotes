package postgres

/**
 * This file contains page-related repository functions
 */

import (
    "log"

    "github.com/setonotes/pkg/page"
)

/**
 * Given an page ID, return the page
 */
func (r *Repository) GetPageByID(pageID int) (*page.Page, error) {
    var (
        title     []byte
        body      []byte
        ownerID   int
    )
    psqlStmt := `
        SELECT title, body, author_id
        FROM pages
        WHERE id=$1`
    log.Printf("getting page-%v from DB...", pageID)
    err := r.DB.QueryRow(psqlStmt, pageID).Scan(&title, &body, &ownerID)
    if err != nil {
        log.Printf("failed to get page-%v from DB", pageID)
        return nil, err
    }
    log.Println("successfully got page from DB")

    return &page.Page{
        ID:      pageID,
        Title:   title,
        Body:    body,
        OwnerID: ownerID,
    }, nil
}

/**
 * Given a page ID, check there exists a database row in the `pages` table with
 * that ID
 */
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

/**
 * Create new row in the `pages` table.
 *
 * This function returns the page ID because right now the ID is assinged by
 * Postgres. This is not ideal as the page ID-ing should not depend on any
 * particular database implementation, so this should ultimately be switched to
 * UUID or something else
 */
func (r *Repository) CreatePage(authorID int) (int, error) {
    log.Println("creating row in `pages` table...")
    psqlStmt := `
        INSERT INTO pages (author_id)
        VALUES ($1)
        RETURNING id`
    pageID := 0
    err := r.DB.QueryRow(psqlStmt, authorID).Scan(&pageID)
    if err != nil {
        log.Println("failed to create page")
        return 0, err
    }
    return pageID, nil
}

/**
 * Update Title and Body of existing page
 */
func (r *Repository) UpdatePage(p *page.Page) error {
    log.Printf("updating row for page-%v", p.ID)
    psqlStmt := `
        UPDATE pages
        SET title=$1, body=$2
        WHERE id=$3`
    _, err := r.DB.Exec(psqlStmt, p.Title, p.Body, p.ID)
    if err != nil {
        log.Println("failed to updated row for page-%v", p.ID)
        return err
    }
    log.Println("successfully updated row for page-%v", p.ID)
    return nil
}

/**
 * Delete a page
 */
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

