# products-api

Simple api to build and create an api using golang and a psql database. Built using this [tutorial](https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql)

## Setup

To run the product you will need set up and download postgres

https://www.postgresql.org/

Once postgres in installed set up a user and database and add them to .env file in the pattern

.env file
```
TEST_DB_USERNAME="testUserName"
TEST_DB_PASSWORD="testPassword"
TEST_DB_NAME="testDatabaseName"
APP_DB_USERNAME="appUserName"
APP_DB_PASSWORD="appPassword"
APP_DB_NAME="appDatabaseName"
```

## Test

Run the below command in the terminal in the project directory

```
go test -v
```


