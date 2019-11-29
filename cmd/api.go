package main

import (
    "log"
    "net/http"
    "encoding/json"
    "strconv"

    //"github.com/setonotes/pkg/user"
    "github.com/setonotes/pkg/page"

    "github.com/gorilla/mux"
)

/* TODO: finish this
*/
func (s *server) savePage(w http.ResponseWriter, r *http.Request) {
    log.Println("Inside the pages handler...")
    // check authentication

    // check authorization(?)

    // get URL variables

    // return page as JSON
}

/**
 * Get a page and return as JSON
 */
func (s *server) getPage(w http.ResponseWriter, r *http.Request) {
    log.Println("Inside the pages handler...")

    // check authorization with bearer token
    userID, authorized, err := s.authService.CheckAuthStatusBearer(r)
    if err != nil {
        log.Printf("error checking bearer token: %v", err)
        return
    } else {
        if !authorized {
            log.Println("unauthorized")
            return
        }
    }

    // get URL variables
    pageID, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        log.Printf("error converting page ID to string: %v", err)
        return
    }
    log.Println("page ID: ", pageID)

    // get page and key
    p, key, err := s.permissionService.GetPageAndKey(pageID, userID)

    // make JSON packet with the response data
    response := struct {
        Page *page.Page `json:"page"`
        Key      []byte `json:"key"`
    }{
        p,
        key,
    }

    log.Println("Responding with JSON packet containing page and key...")

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

/**
 *
 */
func (s *server) deletePage(w http.ResponseWriter, r *http.Request) {
}

/**
 * Responds with a JSON packet containing the `auth_salt` and `encryption_salt`
 * for a particular user
 *
 * Expects a JSON packet with the key `username` in the request body
 */
func (s *server) saltsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Getting salts via API...")

    // check valid path
    if r.URL.Path != "/api/salts" {
        w.WriteHeader(http.StatusNotFound)
        log.Println("error: invalid path name")
        return
    }

    type data struct {
        Username string
    }

    // get request body
    decoder := json.NewDecoder(r.Body)
    var d data
    err := decoder.Decode(&d)
    if err != nil {
        log.Printf("error decoding JSON request body: %v", err)
        return;
    }

    // get the user
    user, err := s.userService.GetByUsername(d.Username)
    if err != nil {
        log.Printf("error getting user: %v", err)
        return;
    }

    // grab salts for the response data
    salts := struct {
        Encryption_salt []byte `json:"encryption_salt"`
        Auth_salt       []byte `json:"auth_salt"`
    }{
        user.EncryptionSalt,
        user.AuthSalt,
    }

    log.Println("Responding with JSON packet of salts...")

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(salts)
}

