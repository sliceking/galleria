#!/bin/bash
# Change directory to our code that we plan to work from
cd "root/go/src/galleria"

echo "====Releasing galleria===="

echo "  Deleting local binary if it exists..."
rm galleria
echo "  Done."

echo "  Deleting existing code"
ssh root@galleria.slowterminal.com "rm -rf /root/go/src/galleria"
echo "  Code deleted successfully."

echo "  Uploading code..."
rsync -avr --exclude '.git/*' --exclude 'tmp/*' --exclude 'images/*' ./ root@galleria.slowterminal.com:/root/go/src/galleria
echo "  Code uploaded successfully."

# echo "  Go get deps..."
# ssh root@galleria.slowterminal.com "export GOPATH=/root/go; /usr/local/go/bin/go get golang.org/x/crypto/bcrypt"
# ssh root@galleria.slowterminal.com "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/mux"
# ssh root@galleria.slowterminal.com "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/schema"
# ssh root@galleria.slowterminal.com "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/lib/pq"
# ssh root@galleria.slowterminal.com "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/jinzhu/gorm"
# ssh root@galleria.slowterminal.com "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/csrf"

echo "  Building code on the remote server..."
ssh root@galleria.slowterminal.com "cd /root/go/src/galleria; /usr/local/go/bin/go build -o ./server; cp server /root/app/"
echo "  Code built successfully."

echo "  Moving assets..."
ssh root@galleria.slowterminal.com "cd /root/app; cp -R /root/go/src/galleria/assets ."
echo "  Assets moved successfully."

echo "  Moving views..."
ssh root@galleria.slowterminal.com "cd /root/app; cp -R /root/go/src/galleria/views ."
echo "  Views moved successfully."

echo "  Moving caddyfile..."
ssh root@galleria.slowterminal.com "cd /root/app; cp -R /root/go/src/galleria/Caddyfile ."
echo "  CaddyFile moved successfully."

echo "  Restarting service..."
ssh root@galleria.slowterminal.com "sudo service galleria restart"
echo "  Service restart successful."

echo "  Restarting Caddy..."
ssh root@galleria.slowterminal.com "sudo service caddy restart"
echo "  Caddy restart successful."

echo "==== Finished releasing galleria ===="