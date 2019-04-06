if [ -f to_server.zip ]; then
    rm to_server.zip
fi

cp main setonotes_main
zip to_server.zip -r setonotes_main templates
rm setonotes_main
