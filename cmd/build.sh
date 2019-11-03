# build the go server
go build \
main.go \
server.go \
handlers.go \
user_auth.go

# build and arrange EasyMDE if it exists
if [ -d ../easy-markdown-editor ]; then
    cd ../easy-markdown-editor
    cp ../editor-src/js/easymde.js ./src/js/easymde.js
    gulp
    cp ./dist/easymde.min.js ../cmd/static/js/easymde.min.js
fi
