# Overview

This is my first Azure function. It's supposed to be a simple api that returns a quote from [quotable.io](https://api.quotable.io). It uses cosmos db to store the quotes of each day and act like a cache.

## Run locally with docker

1. Run the app via docker compose.

   ```bash
   docker-compose up 
   ```

   Make sure to rebuild the image if you made changes to the code

   ```bash
   docker-compose up --build
   ```

2. Call the local API at [`http://localhost:8080/api/quote`](http://localhost:8080/api/quote).
3. The local mongodb is available at `mongodb://localhost:27017` and can be accessed via the user configured in the [`docker-compose.yaml`](./docker-compose.yaml).
4. When finished you can spin down the environment with

   ```bash
   docker-compose down
   ```

   __add `--volumes` to remove the volumes as well__

## Deployment

1. Complie the binary

   ```pwsh
   $env:GOOS="linux"; $env:GOARCH="amd64"; go build -o ./bin/main ./src/main.go
   ```

2. Deploy the function via zip bundle.

   ```pwsh
   func azure functionapp publish <FUNCTION_APP_NAME> --custom
   ```

For super easy deployment, you can use [taskfile](https://taskfile.dev/#/installation) to deploy the function. Than it's a simple `task deploy-all` ✅

## Improvement ideas

- [ ] Implement Azure CLI as well as Managed Identity authentication to allow running in an non Azure environment.
- [ ] Implement Cross Origin Resource Sharing (CORS) on database communication.
- [ ] Implement Cross Origin Resource Sharing (CORS) to allow the function to just be called from the frontend.

## Configuration

| Name                           | Required | Description                               | Example                                                                                                  |
| ------------------------------ | -------- | ----------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| `FUNCTIONS_CUSTOMHANDLER_PORT` | No       | Port number of the API to listen on       | `8080`                                                                                                   |
| `MONGODB_CONNECTION_STRING`    | Yes      | MongoDB connection string of the database | `mongodb://myDatabaseUser:D1fficultP%40ssw0rd@cluster0.example.mongodb.net/?retryWrites=true&w=majority` |
| `MONGODB_DATABASE`             | Yes      | MongoDB database name                     | `myDatabase`                                                                                             |
| `MONGODB_COLLECTION`           | Yes      | MongoDB collection name                   | `myCollection`                                                                                           |
