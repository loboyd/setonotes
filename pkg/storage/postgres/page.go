package postgres

/**
 * This file contains page-related repository functions
 */

import (
    "log"

    "github.com/setonotes/pkg/page"
    "github.com/setonotes/pkg/user" // for current version number
)

/**
 * Given an page ID, return the page
 */
func (r *Repository) GetPageByID(pageID int) (*page.Page, error) {
    var (
        title     []byte
        body      []byte
        ownerID   int
        version   int
    )
    psqlStmt := `
        SELECT title, body, author_id, version
        FROM pages
        WHERE id=$1`
    log.Printf("getting page-%v from DB...", pageID)
    err := r.DB.QueryRow(psqlStmt, pageID).Scan(&title, &body, &ownerID,
        &version)
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
        Version: version,
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
func (r *Repository) CreatePage(p *page.Page, authorID int) (int, error) {
    log.Println("creating row in `pages` table...")
    psqlStmt := `
        INSERT INTO pages (title, body, author_id, version)
        VALUES ($1, $2, $3, $4)
        RETURNING id`
    pageID := 0
    err := r.DB.QueryRow(psqlStmt, p.Title, p.Body, authorID,
        user.CurrentVersion).Scan(&pageID)
    if err != nil {
        log.Println("failed to store page")
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

