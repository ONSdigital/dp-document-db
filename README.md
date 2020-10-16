# dp-document-db-spike

POC for running dp CMD apps against AWS DocumentDB instead of MongoDB. 

Contains dockerfile/docker-compose files to run a modified instance of the `dp-recipe-api` against a test DocumentDB 
instance instead of Mongodb using a stubbed auth/identity API in place of Zebedee.

## Getting started

### Prerequisites
- Go >= 1.15
- Go Modules enabled.
- The `ons-web-development.pem` gpg.
- Guide assumes you have already set up a DoumentDB cluster & EC2 instance (see AWS documentation).

Get the code:
```bash
git clone git@github.com:ONSdigital/dp-document-db-spike.git
```

### Set up
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

Set the following environment vars have been set on the EC2 instance:

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

# The Zebedee stub URL to use (use this value - **"api-stub"** is the container name in the docker-compose.yml)
ZEBEDEE_URL=http://api-stub:8082

# Flag to configure the recipe API to use DocumentDB instead of standard Mongo.
MONGODB_IS_DOC_DB=true
```
Exit the SSH connection.

### Install on EC2 and run the demo

1) Package the POC binaries & docker config and SCP them on to the EC2 instance:
    ```
    make install
    ```

2) SSH onto the box:
    ```
    make ssh
    ```

3) Run the install script from the home dir:
    ```
    ./install.sh
    ```
    This will stop and clean up any existing running containers, build and start new containers using the latest 
    binaries version using docker-compose. To check the app logs run `docker logs -f <container_name>`. Exit SSH session.

4) Assuming everything has been installed & configured correctly you should now be able run the demo app. From the project root dir locally run:
   ```
   make example
   ```
   This will POST a new recipe request to the recipe API (which is now running against DocumentDB) & then send a GET request to retrieve it 
   by ID and pretty print the response JSON in the console. You should see output similar to:
   
   ```bash
    [doc-db-demo] 🦄  2020-10-16T15:32:37+01:00 posting new recipe
    [doc-db-demo] 🦄  2020-10-16T15:32:37+01:00 executing request to Recipe API
    [doc-db-demo] 🦄  2020-10-16T15:32:37+01:00 post recipe response status OK
    [doc-db-demo] 🦄  2020-10-16T15:32:37+01:00 create recipe completed successfully : ID: 2943f3c5-c3f1-4a9a-aa6e-14d21c33524c Alias CPIH
   
    [doc-db-demo] 🦄  2020-10-16T15:32:37+01:00 retrieving recipe from API
    [doc-db-demo] 🦄  2020-10-16T15:32:37+01:00 executing request to Recipe API
    [doc-db-demo] 🦄  2020-10-16T15:32:37+01:00 get recipe response status OK
    [doc-db-demo] 🦄  2020-10-16T15:32:37+01:00 get recipe completed successfully : ID: 2943f3c5-c3f1-4a9a-aa6e-14d21c33524c Alias CPIH
   
    [doc-db-demo] 🦄  2020-10-16T15:32:37+01:00 recipe response json:
   
   {
     "id": "2943f3c5-c3f1-4a9a-aa6e-14d21c33524c",
     "alias": "CPIH",
     "format": "v4",
     "files": [
       {
         "description": "CPIH v4"
       }
     ],
     "output_instances": [
       {
         "dataset_id": "cpih01",
         "editions": [
           "time-series"
         ],
         "title": "Consumer Prices Index including owner occupiers' housing costs (CPIH)",
         "code_lists": [
           {
             "id": "mmm-yy",
             "href": "http://localhost:22400/code-lists/mmm-yy",
             "name": "time"
           },
           {
             "id": "uk-only",
             "href": "http://localhost:22400/code-lists/uk-only",
             "name": "geography",
             "is_hierarchy": true
           },
           {
             "id": "cpih1dim1aggid",
             "href": "http://localhost:22400/code-lists/cpih1dim1aggid",
             "name": "aggregate",
             "is_hierarchy": true
           }
         ]
       }
     ]
   }
   
    [doc-db-demo] 🦄  2020-10-16T15:32:37+01:00 demo complete 🚀  🎉
   ```
End.
 
