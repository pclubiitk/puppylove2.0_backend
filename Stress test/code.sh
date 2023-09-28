#/bin/sh

curl --location 'localhost:8080/session/login' \
--header 'Content-Type: application/json' \
--header '' \
--data '{
    "_id":"210667",
    "passHash":"aaaa"
}'

