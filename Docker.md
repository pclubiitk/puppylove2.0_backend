# How to run PuppyLove in Docker Container

First of all, format the `.env` file as below.

```.env
host = host.docker.internal
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
docker build . -t <your username>/absurd-agent

// Run the docker image
docker run -p 8000:8000 -d <your username>/absurd-agent
```
