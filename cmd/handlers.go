package main

import (
    "log"
    "strconv"
    "net/http"
    "html/template"

    "github.com/setonotes/pkg/page"

    "github.com/blackfriday"
    "github.com/microcosm-cc/bluemonday"
)

/**
 * Make handler for protected routes; check user auth, check valid path string,
 * get page ID from URL and call proper handler
 */
func (s *server) makeHandler(fn func (http.ResponseWriter, *http.Request,
    int, int, bool)) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {
        // check authentication status
        log.Println("checking user auth status...")
        userID, authorized, err := s.authService.CheckUserAuthStatus(r)
        if err != nil {
            log.Println("failed to check user auth status")
            http.NotFound(w, r)
            return
        }

        // check for valid path string
        m := s.validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }

        // get ID from URL or set ID to 0
        pageID, err := strconv.Atoi(m[2])
        if err != nil {
            pageID = 0
        }

        // TODO: find a nice way to refresh session token with each
        // auth-only action

        // track each auth-only HTTP request -- this function is in database.go
        // current time is stored with userID and URL path
        // remove this if statement to start tracking all requests
        if authorized {
            err = s.userService.TrackActivity(userID, r.URL.Path)
            if err != nil {
                log.Println("failed to track user activity; continuing...")
            }
        }

        fn(w, r, pageID, userID, authorized)
    }
}

/**
 * Make handler for API routes; check client auth, check valid path string,
 * call proper handler
 */
func (s *server) makeAPIHandler(fn func (http.ResponseWriter,
    *http.Request)) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {

        // TODO SECURITY-CRITICAL: Replace this with actual auth-checking
        authorized := true;

        // check for valid path string
        if r.URL.Path != "/api/messaging" {
            http.NotFound(w, r)
            return
        }

        // This is inside an `if` just to be explicit
        if authorized {
            fn(w, r)
        }
    }
}

func (s *server) homePageHandler(w http.ResponseWriter, r *http.Request) {
    // check valid path
    if r.URL.Path != "/" {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    // check user-authorization status
    userID, authorized, _ := s.authService.CheckUserAuthStatus(r)
    if authorized {
        // show user-specific directory
        s.directoryHandler(w, r, userID, authorized)
    } else {
        // show landing page of site
        s.landingPageHandler(w, r, authorized)
    }
}

func (s *server) landingPageHandler(w http.ResponseWriter, r *http.Request,
    authorized bool) {

    http.ServeFile(w, r, "templates/landingpage.html")
}

func (s *server) directoryHandler(w http.ResponseWriter, r *http.Request,
    userID int, authorized bool) {

    // get user from ID
    u, err := s.userService.GetByID(userID)
    if err != nil {
        log.Println("failed to get user from ID for directory page")
        return // TODO: this should probably 404
    }

    // get page titles
    tmplMapBytes, err := s.permissionService.GetPageTitles(u)
    if err != nil {
        log.Printf("failed to get titles for user %v: %v", userID, err)
        return // TODO: this should also probably 404
    }

    // convert byteslice titles to strings
    tmplMap := make(map[int]string)
    for k, v := range tmplMapBytes {
        tmplMap[k] = string(v)
    }

    data := struct {
        Pages      map[int]string
        Navbar     bool
        Authorized bool
    }{
        tmplMap,
        true, // directory page always gets a navbar
        authorized,
    }

    s.renderTemplate(w, "directory.tmpl", data)
}

/**
 * Replace every '\r' in a byteslice with '\n'
 * Right now this is necessary to handle the text from the input
 * form on the edit page. If this issue isn't handled in a better
 * way eventually, this function should probably move to a different
 * file somewhere.
 */
func newlineDoctor(text []byte) []byte {
    for i, c := range text {
        if c == '\r' {
            text[i] = '\n'
        }
    }
    return text
}

func (s *server) viewHandler(w http.ResponseWriter, r *http.Request, pageID int,
    userID int, authorized bool) {

    // redirect visitors
    if !authorized {
        log.Println("authorized attempt to view /view/; redirecting...")
        // should this be status found?
        http.Redirect(w, r, "/", http.StatusNotFound)
    }

    // get user
    u, err := s.userService.GetByID(userID)
    if err != nil {
        log.Println("failed to get user by ID for view page")
        return // TODO: this should probably 404
    }

    p, err := s.permissionService.LoadAndDecryptPage(pageID, u)
    if err != nil {
        // don't do this because it's weird
        // http.Redirect(w, r, "/edit/"+string(pageID), http.StatusFound)
        log.Println("failed to decrypt page for view page")
        w.WriteHeader(http.StatusNotFound)
        return // TODO: this should probably 404
    }

    // there is probably a better way to handle this issue
    p.Body = newlineDoctor(p.Body)

    // use blackfriday Markdown processor to get HTML
    unsafeHTML := blackfriday.Run(p.Body)

    // use bluemonday HTML sanitizer to make HTML safe
    safeHTML := bluemonday.UGCPolicy().SanitizeBytes(unsafeHTML)

    // create a map to include markdown in template data
    md_tmpl := map[string]interface{} {
        "ID":    p.ID,
        "Title": string(p.Title),
        "Markdown": template.HTML(safeHTML),
    }

    // data for template
    data := struct {
        Page       map[string]interface{}
        Navbar     bool
        Authorized bool
    }{
        md_tmpl,
        true, // `/view/` always gets a navbar
        authorized,
    }

    s.renderTemplate(w, "view.tmpl", data)
}

func (s *server) editHandler(w http.ResponseWriter, r *http.Request, pageID int,
    userID int, authorized bool) {

    // redirect visitors
    if !authorized {
        log.Println("authorized attempt to view /view/")
        // should this be status found?
        http.Redirect(w, r, "/", http.StatusNotFound)
    }

    // get user
    u, err := s.userService.GetByID(userID)
    if err != nil {
        log.Println("failed to get user by ID for edit page")
        return // TODO: this should probably 404
    }

    p, err := s.permissionService.LoadAndDecryptPage(pageID, u)
    if err != nil {
        // create a new page
        // TODO: what if there is an unexpected error here?
        p = &page.Page{ID: pageID, Title: []byte("New Page"), Body: []byte("")}
    }

    // this struct is the same as a page.Page, but with a string title --
    // this is probably a temporary solution, because eventually we will have a
    // WYSIWYG editor and the titles will also be rendered in Markdown/LaTeX
    data := struct {
        ID        int
        Title     string
        Body      []byte
        Navbar     bool
        Authorized bool
    }{
        p.ID,
        string(p.Title),
        p.Body,
        true, // `/edit/` always gets a navbar
        authorized,
    }

    s.renderTemplate(w, "edit.tmpl", data)
}

func (s *server) saveHandler(w http.ResponseWriter, r *http.Request, pageID int,
    userID int, authorized bool) {

    // redirect visitors
    if !authorized {
        log.Println("authorized attempt to view /view/")
        // should this be status found?
        http.Redirect(w, r, "/", http.StatusNotFound)
    }

    // get user
    u, err := s.userService.GetByID(userID)
    if err != nil {
        log.Println("failed to get user by ID for /save/")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    title := r.FormValue("title")
    body := r.FormValue("body")
    // if pageID == 0, the value will not get used
    p := &page.Page{ID: pageID, Title: []byte(title), Body: []byte(body)}
    pageID, err = s.permissionService.SavePage(p, u)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/view/"+strconv.Itoa(pageID), http.StatusFound)
}

func (s *server) deleteHandler(w http.ResponseWriter, r *http.Request,
    pageID int, userID int, authorized bool) {

    // redirect visitors
    if !authorized {
        log.Println("authorized attempt to view /view/")
        // should this be status found?
        http.Redirect(w, r, "/", http.StatusNotFound)
    }

    err := s.permissionService.DeletePage(pageID, userID)
    if err != nil {
        log.Println("failed to delete page")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/", http.StatusFound)
}

func (s *server) sendMessageHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: Replace this with an actual handler
    if r.Method == "GET" {
        log.Println("received GET request!")
    } else if r.Method == "POST" {
        log.Println("received POST request!")
    } else {
        log.Println("didn't receive POST or GET request...")
    }
}
