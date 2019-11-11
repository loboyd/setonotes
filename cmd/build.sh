# build the go server
go build \
main.go \
server.go \
handlers.go \
user_auth.go

# build and arrange EasyMDE if it exists
if [ -d ../easy-markdown-editor ]; then
    cd ../easy-markdown-editor
    cp -r ../editor-src/js/easymde.js ./src/js/easymde.js
    cp -r ../editor-src/css/easymde.css ./src/css/easymde.css
    gulp
    cp -r ./dist/easymde.min.js ../cmd/static/js/easymde.min.js
    cp -r ./dist/easymde.min.css ../cmd/static/css/easymde.min.css
fi
