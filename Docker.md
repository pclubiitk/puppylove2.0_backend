# How to run PuppyLove in Docker Container

First of all, format the `.env` file as below.

```.env
host = localhost
port = 5432
password = "dbpass"
dbName = postgres
user = postgres
CfgAdminPass = pass

AdminId = adminID
AdminPass = adminPASS

# Make sure both are different
UserjwtSigningKey = "something"
HeartjwtSigningKey = "something2"
domain = localhost


EMAILID="email@iitk.ac.in"
EMAILPASS="emailpass"

```

After doing above changes. Run the below commands in your terminal.

```bash
// Build docker image
docker build -t puppylove .

// Run the docker image
docker run -p 8080:8080 --network host puppylove
