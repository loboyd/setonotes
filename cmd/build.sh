# build the go server
go build \
main.go \
server.go \
handlers.go \
user_auth.go \
api.go

# build and arrange EasyMDE if it exists
if [ -d ../easy-markdown-editor ]; then
    # if there aren't any static folders, then make them
    if [ ! -d ../cmd/static/js ]; then
        mkdir -p ../cmd/static/js
    fi
    if [ ! -d ../cmd/static/css ]; then
        mkdir -p ../cmd/static/css
    fi

    # copy the source files into the EasyMDE repo
    cd ../easy-markdown-editor
    cp -r ../editor-src/js/easymde.js ./src/js/easymde.js
    cp -r ../editor-src/css/easymde.css ./src/css/easymde.css

    # build EasyMDE and copy the minified files to the right places
    gulp
    cp -r ./dist/easymde.min.js ../cmd/static/js/easymde.min.js
    cp -r ./dist/easymde.min.css ../cmd/static/css/easymde.min.css
fi
