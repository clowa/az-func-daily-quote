# Overview

This is my first Azure function. It's supposed to be a simple api that returns a quote from [quotable.io](https://api.quotable.io)

## Deployment

1. Complie the binary

   ```pwsh
   $env:GOOS="linux"; $env:GOARCH="amd64"; go build -o ./bin/main ./src/main.go
   ```

2. Deploy the function via zip bundle.

   ```pwsh
   func azure functionapp publish azure-func-demo-func
   ```

For super easy deployment, you can use [taskfile](https://taskfile.dev/#/installation) to deploy the function. Than it's a simple `task deploy` ✅
