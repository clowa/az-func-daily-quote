# function/src Directory

## Overview

The `function/src` directory contains the main source code for the Quotes Azure Function App. This directory is organized into sub-packages that implement the API endpoints, business logic, configuration, and integration with external services such as MongoDB and the Quotable API.

## Structure

- **main.go**  
  Entry point for the application. Sets up HTTP routes and loads configuration.

- **handlers/**  
  Contains HTTP handler functions the different API endpoints defined in the [openapi](../openapi.yaml) specification.

- **lib/**
  - **config/**  
    Loads and manages application configuration (e.g., database connection strings, API port).
  - **quote/**  
    Defines the `Quote` struct and methods for manipulating and storing quotes.
  - **quotable/**  
    Client and types for interacting with the external [Quotable API](https://github.com/lukePeavey/quotable).

- **tests/**  
  Contains acceptance/end-to-end tests for the API.
