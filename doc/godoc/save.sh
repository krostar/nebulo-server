#!/bin/bash
set -u

DOC_DIR="./build/doc/godoc"
PKG="github.com/krostar/nebulo"
ADDRESS="localhost:53412"

# Run a godoc server which we will scrape. Clobber the GOPATH to include
# only our dependencies.
godoc -http="$ADDRESS" &
DOC_PID=$!

# Wait for the server to init
while :
do
	curl -s "http://$ADDRESS" 2>&1>/dev/null
	if [ $? -eq 0 ] # exit code is 0 if we connected
	then
		break
	fi
done

# Scrape the pkg directory for the API docs. Scrap lib for the CSS/JS. Ignore everything else.
wget -q -r -m -k -E -p -erobots=off --include-directories="/pkg,/lib,/src" --exclude-directories="*" "http://$ADDRESS/pkg/$PKG/../"

# Stop the godoc server
kill -9 $DOC_PID

# Delete the old directory or else mv will put the localhost dir into
# the DOC_DIR if it already exists.
rm -rf $DOC_DIR
mv $ADDRESS $DOC_DIR
