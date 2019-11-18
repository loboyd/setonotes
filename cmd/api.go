package main

import (
    "log"
    "net/http"
    "encoding/json"

    //"github.com/setonotes/pkg/user"
)

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

