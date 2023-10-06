## Quick Start

Copy .env.orders and fill it with your environment values.

Run unit tests with  `make unit-test``.
    - it out puts the coverage at the end of the test
There are no integration tests since it was not part of the requirement

### Start the server
#### Requirements
- A functioning postgres database with the credentials populated in the `.env.orders` file

#### Steps
1. Create a database and replace the database credentials in the .env.orders file.
2. Run `make run-api` in the terminal (in the orders directory the default port is `:5000`)
3. Some functionality in this service communicate with the `payment service`. Ensure that the payment service is up and healthy to test all the functionality of this api

### Code/File structure
All the code logic is written in the internal folder and its subdirectories
./internal
    ./platform/database folder with database setup functions (by default, PostgreSQL)
    ./platform/migrations folder with migration files (used with golang-migrate/migrate tool)

    ./model folder with the models that are used to interact with the database

    ./handler folder with the grpc implementation and the business logic

    ./repository folder with the fuctionality to persist data in the database

    ./config filder consist of all the code that sets up configurations for the api to run ie databases,env variabled etc

    ./common contains all the shared funtions

    ./grpc_clients contains the interface and code to interact synchronously with other service

./build folder contains the latest build of the code which is updated every time `make run-api` is run.

#### Note
Logs are printed as json objects to facilitate 3rd part analysis.
Default log level is Debug. This can be changed in the main folder

### Additional information specific to the test
- The repository layer is to be tested using integration tests
-