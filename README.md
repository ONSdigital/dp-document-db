# dp-document-db-spike

POC for running dp CMD apps against AWS DocumentDB instead of MongoDB. 

Contains dockerfile/docker-compose files to run a modified instance of the `dp-recipe-api` against a test DocumentDB 
instance instead of Mongodb using a stubbed auth/identity API in place of Zebedee.

## Getting started

### Prerequisites
- Go >= 1.15
- Go Modules enabled.
- The `ons-web-development.pem` gpg.
- Guide assumed you have already set up a DoumentDB cluster & EC2 instance.

Get the code:
```bash
git clone git@github.com:ONSdigital/dp-document-db-spike.git
```

Set the following environment vars locally:
```bash
# The user/address EC2 instance running the Docdb cluster.
DOC_DB_EC2_ADDR=<USER@ADDRESS>

# The IP address of the EC2 instance in the VPC of the Docdb cluster.
DOC_DB_POC_IP=<IP ADDRESS>
```
SSH on to the EC2 box:
```bash
make ssh
```

Ensure the following environment vars have been set on the EC2 instance:

```bash
# The DocumentDB instance URL:PORT - see AWS console.
MONGODB_HOST=<URL:PORT>

# The DocumentDB instance username - see AWS console.
MONGODB_USERNAME=<USERNAME>>

# The DocumentDB instance password - see AWS console.
MONGODB_PASSWORD=<THE_PASSWORD>

# The path the cert file need to connect to the DocumentDB instance - see AWS console.
MONGODB_CERT=rds-combined-ca-bundle.pem

# Enable human log format for the Go apps. 
HUMAN_LOG=1

# The Zebedee stub URL to use (use this value - stub-api is the container name in the docker-compose.yml)
ZEBEDEE_URL=http://stub-api:8082

# Flag to configure the recipe API to use DocumentDB instead of standard Mongo.
MONGODB_IS_DOC_DB=true
```
Exit the SSH connection.

Package up the POC binaries and Docker configurations, scp it on to the EC2 instance and the run install script to 
install and start up the POC.
```bash
make install
```
Assuming everything has been installed & configured correctly you should now be able to curl the GET recipes endpoint:
```
curl -XGET "http://${DOC_DB_EC2_ADDR}:22300/recipes" | jq .
```
If this returns an empty response you can try posting a new recipe:
```bash
curl -h "Content-Type: application/json" -d `{"count":1,"limit":1,"items":[{"id":"2943f3c5-c3f1-4a9a-aa6e-14d21c33524c","alias":"CPIH","format":"v4","files":[{"description":"CPIH v4"}],"output_instances":[{"dataset_id":"cpih01","editions":["time-series"],"title":"Consumer Prices Index including owner occupiers' housing costs (CPIH)","code_lists":[{"id":"mmm-yy","href":"http://localhost:22400/code-lists/mmm-yy","name":"time","is_hierarchy":false},{"id":"uk-only","href":"http://localhost:22400/code-lists/uk-only","name":"geography","is_hierarchy":true},{"id":"cpih1dim1aggid","href":"http://localhost:22400/code-lists/cpih1dim1aggid","name":"aggregate","is_hierarchy":true}]}]}],"total_count":1}` http://${DOC_DB_EC2_ADDR}:22300/recipes 

```


curl -d "@example-recipe.json" \
    -h "Content-Type: application/json" \
    -X POST http://54.246.67.239:22300/recipes
    
curl -d @example-recipe.json  -H "Content-Type: application/json" -H "Authorization: Bearer 7e0d1238-cf25-4239-adfb-7f1a460a0580" -XPOST http://54.246.67.239:22300/recipes
    

curl -d `{"alias":"CPIH","files":[{"description":"CPIH v4"}],"format":"v4","id":"2943f3c5-c3f1-4a9a-aa6e-14d21c33524c","output_instances":[{"code_lists":[{"href":"http://localhost:22400/code-lists/mmm-yy","id":"mmm-yy","is_hierarchy":false,"name":"time"},{"href":"http://localhost:22400/code-lists/uk-only","id":"uk-only","is_hierarchy":true,"name":"geography"},{"href":"http://localhost:22400/code-lists/cpih1dim1aggid","id":"cpih1dim1aggid","is_hierarchy":true,"name":"aggregate"}],"dataset_id":"cpih01","editions":["time-series"],"title":"Consumer Prices Index including owner occupiers' housing costs (CPIH)"}]}` -H "Content-Type: application/json" -H "Authorization: Bearer 7e0d1238-cf25-4239-adfb-7f1a460a0580" -XPOST http://54.246.67.239:22300/recipes


 
