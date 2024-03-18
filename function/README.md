# Overview

This is my first Azure function. It's supposed to be a simple api that returns a quote from [quotable.io](https://api.quotable.io). It uses cosmos db to store the quotes of each day and act like a cache.

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
