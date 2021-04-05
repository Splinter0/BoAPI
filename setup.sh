#!/bin/bash

echo "Enter a name for your database"
read dbName
echo "Enter the username for the admin user of the db $dbName"
read username

rootPass=$(echo -n $(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '' | base64 </dev/urandom | head -c 36 ; echo '') | md5sum | awk '{print $1}')
userPass=$(echo -n $(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '' | base64 </dev/urandom | head -c 36 ; echo '') | md5sum | awk '{print $1}')

echo "Creating TLS Keys..."
openssl ecparam -genkey -name secp384r1 -out api/server.key
openssl req -new -x509 -sha256 -key api/server.key -out api/server.crt -days 3650

secretKey=$(echo -n $(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '' | base64 </dev/urandom | head -c 36 ; echo '') | md5sum | awk '{print $1}')

echo "MONGO_URI=mongodb://$username:$userPass@mongo:27017/$dbName" > api/.env
echo "SECRET_KEY=$secretKey" >> api/.env
sed -i "s/replace/$dbName/" api/database/db.go

sudo docker-compose up -d

sleep 5

echo "COPY THESE INTO THE MONGODB SHELL"
echo "1. use admin"
echo "2. db.createUser({user: \"root\", pwd: \"$rootPass\", roles:[\"root\"]});"
echo "3. use $dbName"
echo "4. db.createUser({user: \"$username\", pwd: \"$userPass\", roles:[{role: \"readWrite\", db: \"$dbName\"}]});"
echo "5. exit"
read 
mongo 

echo "COPY THESE INTO THE MONGODB SHELL"
echo "1. use $dbName"
echo "2. db.createCollection(\"users\")"
echo "3. exit"
read
mongo -u $username -p $userPass $dbName

echo "        command: [--auth]" >> docker-compose.yml
sudo docker-compose down
sudo docker-compose up