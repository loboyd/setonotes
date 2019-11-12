package main

/**
 * This file implements the handlers for `/signin/`, `/signup/`, and
 * `/siginout`. A lot of this functionality should probably be moved to the
 * `auth` package. These handlers should also be rewritten to use Go templates
 * rather than just serving up plain HTML files.
 */

import (
    "fmt"
    "log"
    "net/http"

    "github.com/setonotes/pkg/page"
)

/**
 * Handle user-signin by setting Redis cache token and session cookie
 * Expiration time: 24 hours
 */
func (s *server) signinHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("handling signin...")
    switch r.Method {
    case "GET":
        http.ServeFile(w, r, "templates/signin.html") // change to template?
    case "POST":
        if err := r.ParseForm(); err != nil {
            log.Printf("could not parse signin.html: %v\n", err)
            return
        }

        // get form values
        username := r.FormValue("username")
        password := r.FormValue("password")

        // get user
        log.Printf("getting user by username <%s>...", username)
        u, err := s.userService.GetByUsername(username)
        if err != nil {
            log.Printf("failed to get user with username <%s>: %v", username,
                err)
            return // something better should be done here
        }
        log.Println("succesfully got user by username")

        // check password
        auth, err := s.authService.CheckPassHash(
            u.PasswordHash, []byte(password))
        if err != nil || !auth{
            log.Printf("failed check-password-hash: %v", err)
            return // something better should be done here
        }

        // initialize user session
        log.Printf("initializing session for user-%v...", u.ID)
        err = s.authService.InitUserSession(w, r, u, []byte(password))
        if err != nil {
            // handle this error better ("those were the wrong credentials"?)
            log.Printf("failed to initialize session for user-%v", u.ID)
            return // do something better than this
        }
        log.Printf("successfully initialized session for user-%v", u.ID)

        // track user
        err = s.userService.TrackActivity(u.ID, r.URL.Path)
        if err != nil {
            log.Printf("failed to track user activity: %v", err)
        }

        log.Printf("signed user-%v in successfully\n", u.ID)
        http.Redirect(w, r, "/", http.StatusFound)

    default:
        http.Redirect(w, r, "/", http.StatusNotFound)
        // fmt.Fprintf(w, "Only GET and POST requests supported")
    }
}

/**
 * Creates new user and saves to database
 * TODO: Automatic login
 *       Redirect to authorized landing page
 */
func (s *server) signupHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("handling signup...")
    switch r.Method {
    case "GET":
        http.ServeFile(w, r, "templates/signup.html") // change to template?
    case "POST":
        if err := r.ParseForm(); err != nil {
            // fmt.Fprintf(w, "ParseForm() err: %v", err)
            // return
            panic(err)
        }

        username := r.FormValue("username")
        email    := r.FormValue("email")
        password := r.FormValue("password")
        authSalt := r.FormValue("auth_salt")
        encryptionSalt   := r.FormValue("encryption_salt")
        mainKeyEncrypted := r.FormValue("main_key_encrypted")

        log.Println("auth salt: ", authSalt)

        // right now, this is checking the username exists in the beta-testers
        // table, but eventually this will check against our imposed conditions
        // for valid usernames
        whitelisted, err := s.userService.CheckBetaTesterWhitelist(username)
        if err != nil {
            log.Printf("failed to check beta-tester whitelist: %v", err)
            return
        }
        if !whitelisted {
            fmt.Fprintf(w,
                "Only beta testers are allowed to sign-up right now.")
            return
        }
        // END

        //u, err := s.userService.Create(username, email, password)
        u, err := s.userService.Create(username, email, password, authSalt,
            encryptionSalt, mainKeyEncrypted)
        if err != nil {
            log.Println("failed to create new user: %v", err)
            return
        }

        // initialize user session
        log.Printf("initializing session for user-%v...", u.ID)
        err = s.authService.InitUserSession(w, r, u, []byte(password))
        if err != nil {
            // handle this error better ("those were the wrong credentials"?)
            log.Printf("failed to initialize session for user-%v", u.ID)
            return // do something better than this
        }
        log.Printf("successfully initialized session for user-%v", u.ID)

        // track user
        err = s.userService.TrackActivity(u.ID, r.URL.Path)
        if err != nil {
            log.Printf("failed to track user activity: %v", err)
        }
        // create reference pages that demonstrates Markdown
        err = s.createReferencePage(u.ID)
        if err != nil {
            log.Printf("failed to create reference page: %v", err)
            return
        }

        http.Redirect(w, r, "/", http.StatusFound)
    }
}

/**
 * Handle user sign-out by asking the auth service to end the session
 */
func (s *server) signoutHandler(w http.ResponseWriter, r *http.Request, _,
    userID int, authorized bool) {

    if authorized {
        err := s.authService.EndUserSession(w, r, userID)
        if err != nil {
            log.Println("failed to end user-%v session", userID)
        }
    }

    http.Redirect(w, r, "/", http.StatusFound)
}

/**
 * Create a reference page to demonstrate various features of Markdown
 * The actual data for the page should probably be stored in the database or in
 * some better way than hard-coding it. One consideration is how easily updated
 * this is (because it should be updated to reflect new features). Eventually
 * it would be nice if there was just a single instance of this page that was
 * read-only to all users.
 */
func (s *server) createReferencePage(userID int) error {
    p := &page.Page{
        ID: 0,
        Title: []byte("Reference Page (click me!)"),
        Body: []byte(
            "Welcome to setonotes! This page serves as a reference for the" +
            " various features available to you.\n\n" +

            "Click `edit` above to see how this page is formatted. The hope is" +
            " to add more awesome features with time.\n\n" +

            "Also, please tell us what is awful about setonotes so that we can" +
            " make it awesome instead! Email contact@setonotes.com.\n\n" +

            "# Header 1\n" +
            "## Header 2\n" +
            "### Header 3\n" +
            "#### Header 4\n" +
            "##### Header 5\n" +
            "###### Header 6\n" +

            "You've got **bold** and _italics,_ and **_both._**\n\n" +

            "You can link to other web pages: [Google](https://google.com).\n\n" +

            "You can include images.\n\n" +
            "![I'm just a happy little piece of alt-text.](https://i2.kym" +
            "-cdn.com/entries/icons/original/000/014/664/smallfctfrog.JPG)\n\n" +

            "* You\n" +
            "* can\n" +
            "* create\n" +
            "* unordered\n" +
            "1. or\n" +
            "2. ordered\n" +
            "3. lists.\n\n" +

            "`print(\"Inline code works.\")`\n\n" +

            "Fractions look nice: 1/2\n"),
    }

    // get user
    u, err := s.userService.GetByID(userID)
    if err != nil {
        log.Printf("failed to create reference page: %v", err)
        return err
    }

    _, err = s.permissionService.SavePage(p, u)
    if err != nil {
        log.Printf("failed to create reference page: %v", err)
        return err
    }

    return nil
}
