package main

/**
 * This file defines the server struct which servers as something to attach the
 * router, templates, handlers, and setonotes services to rather than using them
 * as package-level variables. Note that this file is not its own package. It is
 * still within the `main` package.
 *
 * This design is based on this article:
 * https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-
 *     37c208122831)
 *
 * The template handling code is based on this article:
 * https://hackernoon.com/golang-template-2-template-composition-and-how-to-
 *     organize-template-files-4cb40bcdf8f6
 */

import (
    "log"
    "regexp"
    "path/filepath"
    "net/http"
    "html/template"

    "github.com/setonotes/pkg/user"
    "github.com/setonotes/pkg/page"

    "github.com/gorilla/mux"
    "github.com/oxtoacart/bpool"
)

type userService interface {
    GetByID(userID int) (*user.User, error)
    GetByUsername(username string) (*user.User, error)
    Create(username, email, password, authSalt, encryptionSalt,
        mainKeyEncrypted string) (*user.User, error)
    TrackActivity(userID int, path string) error
    CheckBetaTesterWhitelist(username string) (bool, error)
}

type authService interface {
    CheckAuthStatusCookie(r *http.Request) (int, bool, error)
    CheckAuthStatusBearer(r *http.Request) (int, bool, error)
    InitUserSession(w http.ResponseWriter, r *http.Request, u *user.User,
        password []byte) error
    EndUserSession(w http.ResponseWriter, r *http.Request, userID int) error
    CheckPassHash(passwordHash, password []byte) (bool, error)
}

type permissionService interface {
    GetPageTitles(u *user.User) (map[int][]byte, error)
    SavePage(p *page.Page, u *user.User) (int, error)
    LoadAndDecryptPage(pageID int, u *user.User) (*page.Page, error)
    DeletePage(pageID, userID int) error
    GetPageAndKey(userID, pageID int) (*page.Page, []byte, error)
}

type server struct {
    router            *mux.Router
    templates         map[string]*template.Template
    bufpool           *bpool.BufferPool // used for template rendering

    userService       userService
    authService       authService
    permissionService permissionService

    validPath         *regexp.Regexp
}

/**
 * Creates a new server instance given an authentication service
 *
 * Go's built-in net/http.ServeMux is used as the default router. In a perfectly
 * hexagonal architecture, we would define a more general router interface, but
 * this is okay for now
*/
func newServer(u userService, a authService, p permissionService) *server {
    s := &server{
        router:            mux.NewRouter(),
        userService:       u,
        authService:       a,
        permissionService: p,
    }

    log.Println("loading templates...")
    err := s.loadTemplates()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("templates loaded successfully")

    log.Println("defining routes...")
    s.routes()
    log.Println("routes defined successfully")

    return s
}

/**
 * Defines all routes for the site
 */
func (s *server) routes() {
    //s.router.Host("localhost")
    s.router.HandleFunc("/",          s.homePageHandler)
    s.router.HandleFunc("/api/salts", s.saltsHandler)

    // page API
    s.router.HandleFunc("/api/pages/{id}", s.getPage).Methods("GET")
    s.router.HandleFunc("/api/pages",      s.savePage).Methods("POST")
    s.router.HandleFunc("/api/pages/{id}", s.savePage).Methods("POST")
    s.router.HandleFunc("/api/pages/{id}", s.deletePage).Methods("DELETE")

    s.router.HandleFunc("/signup/",  s.signupHandler)
    s.router.HandleFunc("/signin/",  s.signinHandler)
    s.router.HandleFunc("/signout/", s.makeHandler(s.signoutHandler))
    s.router.HandleFunc("/save/",    s.makeHandler(s.saveHandler))
    s.router.HandleFunc("/edit/",    s.makeHandler(s.editHandler))
    s.router.HandleFunc("/delete/",  s.makeHandler(s.deleteHandler))

    // serve static files from the /static/ directory
    s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
        http.FileServer(http.Dir("static"))))

    s.validPath = regexp.MustCompile(
        "^/(new|save|edit|delete|signout)/([0-9]*)$")
}

/**
 * Check that the input describes a valid route
 */
func (s *server) checkValidPath(path string) bool {
    m := s.validPath.FindStringSubmatch(path)
    if m == nil {
        return false
    }
    return true
}

const mainTmpl = `{{define "main"}} {{template "base" .}} {{end}}`

/**
 * Load all templates from `templates/` and `templates/layout/` while respecting
 * nesting
 */
func (s *server) loadTemplates() error {
    includePath := "templates/"
    layoutPath  := "templates/layout/"

    if s.templates == nil {
        s.templates = make(map[string]*template.Template)
    }

    layoutFiles, err := filepath.Glob(layoutPath + "*.tmpl")
    if err != nil {
        log.Println("failed to get included templates")
        return err
    }

    includeFiles, err := filepath.Glob(includePath + "*.tmpl")
    if err != nil {
        log.Println("failed to get layout templates")
        return err
    }

    mainTemplate := template.New("main")
    mainTemplate, err = mainTemplate.Parse(mainTmpl)
    if err != nil {
        log.Println("failed to parse main template")
        return err
    }

    for _, file := range includeFiles {
        fileName := filepath.Base(file)
        files := append(layoutFiles, file)
        s.templates[fileName], err = mainTemplate.Clone()
        if err != nil {
            return err
        }
        s.templates[fileName] = template.Must(
            s.templates[fileName].ParseFiles(files...))
    }

    s.bufpool = bpool.NewBufferPool(64)
    return nil
}

/**
 * Render HTML template
 */
func (s *server) renderTemplate(w http.ResponseWriter, name string,
    data interface{}) {

    tmpl, ok := s.templates[name]
    if !ok {
        log.Printf("failed to get template with name <%v>", name)
        http.Error(w, "missing template", http.StatusInternalServerError)
    }

    buf := s.bufpool.Get()
    defer s.bufpool.Put(buf)

    err := tmpl.Execute(buf, data)
    if err != nil {
        log.Printf("failed to execute template with name <%v>", name)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    buf.WriteTo(w)
}

