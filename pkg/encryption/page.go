package encryption

/**
 * This file will contain all encryption functionality concerning a particular
 * page. This is the only code allowed to handle unencrypted page-keys (for now
 * at least -- there may be another file related to permissions, but it will
 * also be contained within this package
 */

import (
    "github.com/setonotes/pkg/user"
    "github.com/setonotes/pkg/page"
)

/**
 * Encrypt a page for a particular user --
 * Get the user-encrypted page key from the repo, decrypt it with
 * s.UserDecryptData(), and use it to encrypt the page
 */
func (s *Service) EncryptPage(p *page.Page, u *user.User,
    userEncryptedPageKey []byte) error {

    // decrypt page-key with user-decrypt function
    key, err := s.UserDecryptData(u, userEncryptedPageKey)
    if err != nil {
        return err
    }

    // encrypt the page Title
    title, err := s.EncryptData(p.Title, key)
    if err != nil {
        return err
    }

    // encrypt the page Body
    body, err := s.EncryptData(p.Body, key)
    if err != nil {
        return err
    }

    // assign encrypted values to page fields
    p.Title = title
    p.Body  = body
    return nil
}

/**
 * Decrypt a page for a particular user --
 * Get the user-encrypted page key from the repo, decrypt it with
 * s.UserDecryptData(), and use it to decrypt the page
 *
 * If either the Title or the Body of the page is an empty byteslice, it will be
 * left empty and decryption will not be attempted
 */
func (s *Service) DecryptPage(p *page.Page, u *user.User,
    userEncryptedPageKey []byte) error {

    // decrypt key with user-decrypt function
    key, err := s.UserDecryptData(u, userEncryptedPageKey)
    if err != nil {
        return err
    }

    var title []byte
    if len(p.Title) > 0 {
        title, err = s.DecryptData(p.Title, key)
        if err != nil {
            return err
        }
    }

    var body []byte
    if len(p.Body) > 0 {
        body, err = s.DecryptData(p.Body, key)
        if err != nil {
            return err
        }
    }

    // assign encrypted values to page fields
    p.Title = title
    p.Body  = body
    return nil
}

