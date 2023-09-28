# Me-My Encryption

## How to setup for development

Clone the repo from the github link & use command `go mod tidy` to install all dependencies (make sure you've installed go).

Install(if not already) & Run the postgres server(refer [link](https://www.postgresql.org/download/)) on your device (depending the OS you're using).

Depending on the authentication details you've set up modify the `.env` file(refer envformat.txt) by filling in details of the postgres server.

> Format of the .env file -- Don't forget to disable SSL if not already.
```
host = localhost
port = 5432
password = 'password'
dbName = postgres
user = postgres
CfgAdminPass = something

AdminId = may_be_aleatoryfreak
AdminPass = you_can_just_guess_it

# Make sure both are different
UserjwtSigningKey = "something"
HeartjwtSigningKey = "something2"
``` 

Build the app by using `go mod build` & then run `./me-my_encryption` (the generated file).
OR you may directly run the main.go file as well.

Server should be up & listening at port 8080.


## About the Algorithm (I'll try eplaining better when I get time)

let the key to be sent here be 
k = my_roll_me_roll_rand (eg. `210667_21xxxx_&89h9hKJbx`)
```mermaid
sequenceDiagram
Me ->> Server: SHA(k1), enc(SHA(k1))<br/> SHA(k2), enc(SHA(k2))...
Note right of Server: On login
My-->>Server: Give me all enc 
Server-->>My: enc1, enc2, ...
My-->>Server: enc1 is mine
Server-->>My: How do I verify ?
My-->>Server: Here is SHA I got<br/> on decoding,
Note left of Server: SHA & enc match.<br/> So, assign the pair to user.
Note left of My: Now when "my" send his/her<br/>hearts, assigned pairs are also sent.<br/> Also, "my" gets a token for 10min.<br/>To send claimed heart again <br/>if not sent already. (quick fix for clash)
My->>Server: SHA(l1), enc(SHA(l1))<br/> SHA(l2), enc(SHA(l2))...<br/> + <br/> SHA(k1), enc(SHA(k1)) <br/>-enc with pubkey of receiver.
Server-->Me: Syncing data with Me...
Me->>Server: I got my heart back.
Server->>Me: How do I verify ?
Me->>Server: Here is my k(210667_21xxxx_&89h9hKJbx).<br/>Verify it with SHA.
Me->My: Matched
 ```
