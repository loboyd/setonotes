package main

import (
    "flag"
    "log"
    "net/http"

    "github.com/setonotes/pkg/config"
    repo "github.com/setonotes/pkg/storage/postgres"
    cache "github.com/setonotes/pkg/cache/redis"
    "github.com/setonotes/pkg/user"
    "github.com/setonotes/pkg/page"
    "github.com/setonotes/pkg/encryption"
    "github.com/setonotes/pkg/auth"
    "github.com/setonotes/pkg/permission"
)

/**
 * Redirect all HTTP traffic to HTTPS for SeCuRiTy
 */
func httpsRedirect(w http.ResponseWriter, r *http.Request) {
    http.Redirect(
        w, r,
        "https://" + r.Host + r.URL.String(),
        http.StatusMovedPermanently,
    )
}

func main() {
    // define command line flags
    localFlag := flag.Bool("local", false,
        "Usage: ./<setonotes main> -local")

    log.Println("starting setonotes main...")
    flag.Parse()

    // get configuration settings
    log.Println("getting configuration from JSON file...")
    conf, err := config.New("../config.json")
    if err != nil {
        log.Fatalf("failed to get configuration settings from file: %v", err)
    }
    log.Println("successfully got configuration settings")

    // create new repository
    log.Println("creating new repository...")
    repository, err := repo.New(conf)
    if err != nil {
        log.Fatalln("failed to create new repository")
    }
    log.Println("successfully created new repository")

    // create a session cache
    log.Println("creating new session cache...")
    sessionCache, err := cache.New()
    if err != nil {
        log.Fatalln("failed to create new cache")
    }
    log.Println("successfully created new session cache")

    // create new encryption service
    log.Println("creating new encryption servic...")
    encryptionService := encryption.NewService(sessionCache)
    log.Println("successfully created new encryption service")

    // create new auth service
    log.Println("creating new authentication service...")
    authService := auth.NewService(sessionCache)
    log.Println("successfully created new authentication service")

    // initialize user service
    log.Println("creating new user service...")
    userService := user.NewService(repository, encryptionService, authService)
    log.Println("successfully created new user service")

    // initialize page service
    log.Println("creating new page service...")
    pageService := page.NewService(repository)
    log.Println("successfully created new page service")

    // initialize permission service
    log.Println("creating new permission service...")
    permissionService := permission.NewService(repository, encryptionService,
        userService, pageService)
    log.Println("successfully created new permission service")

    // initialize server (defined in `server.go`)
    server := newServer(userService, authService, permissionService)

    // find proper CA-certificates and keys for HTTPS
    var tlsCertPath string
    var tlsKeyPath string
    if *localFlag {
        tlsCertPath = `local_https/localhost.crt`
        tlsKeyPath  = `local_https/localhost.key`
        log.Println("using local self-signed ca-certificate")
    } else {
        tlsCertPath = `/etc/letsencrypt/live/setonotes.com/fullchain.pem`
        tlsKeyPath  = `/etc/letsencrypt/live/setonotes.com/privkey.pem`
        log.Println("using letsencrypt ca-certificate")
    }

    log.Println("listening on ports :80 and :443...")

    // redirect all HTTP traffic to HTTPS
    go http.ListenAndServe(":80", http.HandlerFunc(httpsRedirect))
    log.Fatal(http.ListenAndServeTLS(":443", tlsCertPath, tlsKeyPath,
        server.router))
}

