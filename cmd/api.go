package main

import (
    "log"
    "net/http"
    "encoding/json"

    //"github.com/setonotes/pkg/user"
)

func (s *server) saltsHandler(w http.ResponseWriter, r *http.Request) {
    // check valid path
    if r.URL.Path != "/api/salts" {
        w.WriteHeader(http.StatusNotFound)
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
    }

    // get the user
    user, err := s.userService.GetByUsername(d.Username)
    if err != nil {
        log.Printf("error getting user: %v", err)
    }

    // grab salts for the response data
    salts := struct {
        Encryption_salt []byte `json:"encryption_salt"`
        Auth_salt       []byte `json:"auth_salt"`
    }{
        user.Salt,
        user.AuthSalt,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(salts)
}

