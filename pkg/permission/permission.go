package permission

/**
 * The purpose of this package is to facilitate the interaction between users
 * and pages. Any page functionality tied to a paritcular user or user
 * functionality tied to a particular page will live here.
 */

import (
    "log"
    "errors"

    "github.com/setonotes/pkg/user"
    "github.com/setonotes/pkg/page"
)

type Repository interface {
    CheckPageExists(pageID int) (bool, error)
    UpdatePage(p *page.Page) error
    CreatePage(userID int) (int, error) // returns pageID
    DeletePage(pageID int) error
    DeleteAllPagePermissions(pageID int) error
    GetUserEncryptedPageKey(userID, pageID int) ([]byte, error)
    GetUserDisembodiedPages(userID int) ([]*page.Page, error)
    CreatePagePermission(userID, pageID int, isOwner, canEdit bool,
        userEncryptedPageKey []byte) error
    CheckUserCanEditPage(userID, pageID int) (bool, error)
}

type EncryptionService interface {
    EncryptPage(p *page.Page, u *user.User, userEncryptedPageKey []byte) (error)
    DecryptPage(p *page.Page, u *user.User, userEncryptedPageKey []byte) (error)
    NewUserEncryptedSymmetricKey(u *user.User) ([]byte, error)
}

/**
 * Service holds interfaces for a repository and an encryption service. It also
 * holds pointers to user and pages services. Notice that these do not have to
 * use an interface, as this package is below the domain level and thus can
 * depend on domain-level packages.
 */
type Service struct {
    repo        Repository
    encryption  EncryptionService
    userService *user.Service
    pageService *page.Service
}

/**
 * Creates a new permission service
 */
func NewService(r Repository, e EncryptionService, u *user.Service,
    p *page.Service) *Service {

    return &Service {
        repo:        r,
        encryption:  e,
        userService: u,
        pageService: p,
    }
}

var ErrNotImplemented error = errors.New("not yet implemented")

/**
 * Gets a particular user's encrypted page-key
 */
func (s *Service) GetUserEncryptedPageKey(userID, pageID int) ([]byte, error) {
    key, err := s.repo.GetUserEncryptedPageKey(userID, pageID)
    return key, err
}

/**
 * Get all page titles for which the given user has read-permission
 * Returns a map from pageID to page Title
 * 
 * TODO: passing back owner info will be necessary for displaying owernship of
 * shared pages -- disembodied pages already contain the ownerID, but that info
 * is not passed back here.
 *
 * Page titles are returned encrypted from the database, and then decrypted with
 * the encryption service
 */
func (s *Service) GetPageTitles(u *user.User) (map[int][]byte, error) {
    // get titles from database
    // titles, err := s.repo.GetUserPageTitles(u.ID)
    pages, err := s.repo.GetUserDisembodiedPages(u.ID)
    if err != nil {
        return nil, err
    }

    // loop over titles and decrypt each
    titles := make(map[int][]byte)
    for _, p := range pages {
        log.Println("decrypting disembodied page...")
        err = s.UserDecryptPage(u, p)
        if err != nil {
            log.Println("failed to decrypt disembodied page")
            return nil, err
        }
        log.Println("successfully decrypted disembodied page")

        // add title to title map
        titles[p.ID] = p.Title
    }

    return titles, nil
}

/**
 * Gets a user's encrypted page key and encrypts a page
 */
func (s *Service) UserEncryptPage(u *user.User, p *page.Page) error {
    // get user-encrypted page key
    key, err := s.GetUserEncryptedPageKey(u.ID, p.ID)
    if err != nil {
        log.Printf("failed to get user-%v-encrypted page-%v key", u.ID, p.ID)
        return err
    }

    // encrypt page
    err = s.encryption.EncryptPage(p, u, key)
    if err != nil {
        log.Printf("failed to encrypt page-%v for user-%v", p.ID, u.ID)
        return err
    }

    return nil
}

/**
 * Gets a user's encrypted page key and decrypts a page
 */
func (s *Service) UserDecryptPage(u *user.User, p *page.Page) error {
    // get user-encrypted page key
    key, err := s.GetUserEncryptedPageKey(u.ID, p.ID)
    if err != nil {
        log.Printf("failed to get user-%v-encrypted page-%v key", u.ID, p.ID)
        return err
    }

    // decrypt page
    err = s.encryption.DecryptPage(p, u, key)
    if err != nil {
        log.Printf("failed to decrypt page-%v for user-%v", p.ID, u.ID)
        return err
    }

    return nil
}

/**
 * Check that userID is allowed to edit pageID
 */
func (s *Service) CheckUserCanEditPage(userID, pageID int) (bool, error) {
    canEdit, err := s.repo.CheckUserCanEditPage(userID, pageID)
    if err != nil {
        return false, err
    }
    return canEdit, nil
}

/**
 * Return page with user-encrypted page key
 */
func (s *Service) GetPageAndKey(pageID, userID int) (*page.Page, []byte,
    error) {

    log.Printf("getting page-%v for user-%v", pageID, userID)

    // get page
    p, err := s.pageService.GetByID(pageID)
    if err != nil {
        log.Println("failed to get page by ID")
        return nil, nil, err
    }

    // get key
    key, err := s.GetUserEncryptedPageKey(userID, pageID)
    if err != nil {
        log.Println("failed to get page-key")
        return nil, nil, err
    }

    return p, key, nil
}

/** TODO: deprecate this
 * Given a pageID and user, load, decrypt, and return a page
 */
func (s *Service) LoadAndDecryptPage(pageID int,
    u *user.User) (*page.Page, error) {

    log.Println("permissions: loading and decrypting page...")

    // get page
    p, err := s.pageService.GetByID(pageID)
    if err != nil {
        log.Println("failed to get page by ID")
        return nil, err
    }

    // decrypt page
    err = s.UserDecryptPage(u, p)
    if err != nil {
        log.Printf("failed to decrypt page-%v for user-%v", p.ID, u.ID)
        return nil, err
    }

    return p, nil
}

/**
 * Update page's Title and Body attributes in storage
 *
 * Returns page ID
 */
func (s *Service) UpdatePage(p *page.Page, u *user.User) error {
    // check the the given user has permission to update the given page
    // this should probably ultimately be handled by a `permission` package
    canEdit, err := s.CheckUserCanEditPage(u.ID, p.ID)
    if err != nil {
        log.Println("failed to check permission")
        return err
    }
    if !canEdit {
        log.Printf("user-%v cannot edit page-%v", u.ID, p.ID)
        return err
    }

    // store page
    log.Println("storing updated page-%v", p.ID)
    err = s.repo.UpdatePage(p)
    if err != nil {
        // should we decrypt the page in memory here?
        log.Printf("failed to update page-%v", p.ID)
        return err
    }
    log.Printf("succesfully stored updated page-%v", p.ID)

    return nil
}

/**
 * Given a sparse Page struct (containing only a title, body, and null ID),
 * generate remaining fields and store it. Also, create a new page permission
 *
 * The way this function works is a bit strange only because Postgres assigns
 * the pageID (this will ultimately be replaced by something like UUID) and it
 * is required to set the page permission. First, a new "empty" page is created
 * in order to generate a new ID. Then, the page is updated as if it is a normal
 * pre-existing page, but with the original Title and Body fields. The
 * encryption is performed by the update function.
 *
 * returns page ID
 */
func (s *Service) CreatePage(u *user.User, userEncryptedPageKey []byte) (int,
    error) {

    // Create page with empty byteslices for Title and Body. This allows us to
    // get a meaninful ID from the database (eventually this will be replaced
    // with UUID or something else). Further down we will update this page.
    log.Println("generating ID for new page...")
    pageID, err := s.repo.CreatePage(u.ID)
    if err != nil {
        log.Println("failed to generate ID")
        return 0, err
    }
    log.Println("successfully generated ID")

    // create new page permission and store user-encrypted page key
    // is-owner and can-edit flags are both set
    log.Println("creating new page permission...")
    s.repo.CreatePagePermission(u.ID, pageID, true, true,
        userEncryptedPageKey)
    log.Println("successfully created page permission")

    return pageID, nil
}

var ErrPermissionConflict = errors.New("permission conflict")

/** TODO: COMPLETE THIS
 * Delete a page after checking the user has delete-permission (is the owner of
 * the page)
 *
 * The page is built from the pageID within this function so that the caller
 * does not have to do it, as the caller is likely a handler that only has the
 * page ID at hand
 */
func (s *Service) DeletePage(pageID, userID int) (error) {
    // get the page
    p, err := s.pageService.GetByID(pageID)
    if err != nil {
        return err
    }

    // check user has delete-permission (is the owner)
    if p.OwnerID != userID {
        return ErrPermissionConflict
    }

    // delete permissions
    err = s.repo.DeleteAllPagePermissions(pageID)
    if err != nil {
        log.Printf("Unauthorized attempt by user-%v to delete page-%v permissions",
            userID, pageID)
        return err
    }

    // delete page
    err = s.repo.DeletePage(pageID)
    if err != nil {
        log.Printf("Unauthorized attempt by user-%v to delete page-%v", userID,
            pageID)
        return err
    }

    return nil
}
