package main

import (
    "log"
    "net/http"
    "encoding/json"
    b64 "encoding/base64"
    "strconv"

    "github.com/setonotes/pkg/page"

    "github.com/gorilla/mux"
)

/**
 * Handle the POST /pages/{id} API call
 *
 * The request data should contain `key`, `title`, and `body` fields of the
 * page encoded as base-64 strings:
 *  {
 *    "page": {
 *      "key": "some-base64-string",
 *      "title": "another-base64-string",
 *      "body": "final-base64-string",
 *    }
 *  }
 */
func (s *server) createPage(w http.ResponseWriter, r *http.Request) {
    log.Println("Inside the pages handler...")

    // check authentication with bearer token
    userID, authorized, err := s.authService.CheckAuthStatusBearer(r)
    if err != nil || !authorized {
        log.Println("unauthorized")
        http.Error(w, http.StatusText(http.StatusUnauthorized),
            http.StatusUnauthorized)
        return
    }
    log.Println("successfully confirmed authorization")

    // get user
    u, err := s.userService.GetByID(userID)
    if err != nil {
        log.Println("failed to get user by ID for POST api/pages")
        http.Error(w, http.StatusText(http.StatusInternalServerError),
            http.StatusInternalServerError)
        return
    }
    log.Println("successfully retrieved user", u.ID)

    type data struct {
        Key string   `json:"key"`
        Title string `json:"title"`
        Body string  `json:"body"`
    }

    // get request body
    decoder := json.NewDecoder(r.Body)
    var d data
    err = decoder.Decode(&d)
    if err != nil {
        log.Printf("error decoding JSON request body: %v", err)
        http.Error(w, http.StatusText(http.StatusBadRequest),
            http.StatusBadRequest)
        return;
    }

    key,   err1 := b64.StdEncoding.DecodeString(d.Key)
    title, err2 := b64.StdEncoding.DecodeString(d.Title)
    body,  err3 := b64.StdEncoding.DecodeString(d.Body)
    if err1 != nil || err2 != nil || err3 != nil {
        log.Printf("error decoding base64 values in JSON  body: %v, %v, %v",
            err1, err2, err3)
        http.Error(w, http.StatusText(http.StatusInternalServerError),
            http.StatusInternalServerError)
        return;
    }

    // build the page
    p := &page.Page{
        ID: 0,
        Title: title,
        Body:  body,
        OwnerID: userID ,
    }

    // create the page
    pageID, err := s.permissionService.CreatePage(p, u, key)
    if err != nil {
        http.Error(w, http.StatusText(http.StatusInternalServerError),
            http.StatusInternalServerError)
        return
    }
    log.Println("successfully saved page", pageID)

    // send JSON packet containing newly-created ID
    new_id := struct {
        ID int `json:"id"`
    }{
        pageID,
    }

    log.Println("Responding with JSON packet containing new page ID...")

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(new_id)
}

/* TODO: finish this
*/
func (s *server) updatePage(w http.ResponseWriter, r *http.Request) {
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

    // check authentication with bearer token
    userID, authorized, err := s.authService.CheckAuthStatusBearer(r)
    if err != nil || !authorized {
        log.Println("unauthorized")
        http.Error(w, http.StatusText(http.StatusUnauthorized),
            http.StatusUnauthorized)
        return
    }

    // get URL variables
    pageID, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        log.Printf("error converting page ID to string: %v", err)
        http.Error(w, http.StatusText(http.StatusInternalServerError),
            http.StatusInternalServerError)
        return
    }
    log.Println("page ID: ", pageID)

    // get page and key
    p, key, err := s.permissionService.GetPageAndKey(pageID, userID)
    if err != nil {
        log.Printf("error getting page and key: %v", err)
        http.Error(w, http.StatusText(http.StatusInternalServerError),
            http.StatusInternalServerError)
        return;
    }

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
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(salts)
}

