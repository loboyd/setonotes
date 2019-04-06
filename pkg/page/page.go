package page

/**
 * This is a pretty light package that just defines the page struct and allows
 * us to get it from a repository. Most of the page-specific functions (such as
 * creating, updating, deleting, etc) require a *user.User as well or at least a
 * user ID, so they live in the permission packge.
 */

import (
    "log"
)

type Page struct {
    ID      int
    Title   []byte
    Body    []byte
    OwnerID int
    Version int
}

type Repository interface {
    GetPageByID(id int) (*Page, error)
}

type Service struct {
    repo       Repository
}

/**
 * Creates a new Page Service
 */
func NewService(r Repository) *Service {
    return &Service{
        repo:       r,
    }
}

/**
 * Returns a pointer to a page given the page's ID
 */
func (s *Service) GetByID(id int) (*Page, error) {
    page, err := s.repo.GetPageByID(id)
    if err != nil {
		log.Println("failed to get page by ID from repository")
        return nil, err
    }
    return page, nil
}

