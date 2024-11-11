# Clerk Integration with OpenAPI-Generated Server using oapi-codegen

This repository provides a straightforward example of integrating [Clerk](https://clerk.dev/) for authentication into a server generated from an OpenAPI specification using [oapi-codegen](https://github.com/deepmap/oapi-codegen).

## Features

- **Authentication with Clerk:** Seamlessly incorporate Clerk’s robust authentication system into your API server.
- **OpenAPI Specification:** Utilize OpenAPI to define your API endpoints, ensuring clear and standardized documentation.
- **Code Generation with oapi-codegen:** Automatically generate type-safe Go server code from your OpenAPI specs, accelerating development.
- **Easy Setup:** Follow simple steps to get your authenticated server up and running quickly.

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/yun-jay/clerk-echo-oapi-middleware.git
cd clerk-echo-oapi-middleware
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Configure Clerk

- Sign up for a Clerk account and obtain your API keys.
- Set up environment variables as per the instructions in the Clerk documentation.
- Enable the organization feature
- Create the permissions org:things:r, org:things:w
- Create an organization and one user and assign it to the organization

### 4. Generate Server Code

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
go generate ./...
```

### 5. Run the Server

```bash
go run main.go
```
