# backend

[![Deploy backend](https://github.com/EmilyOng/tusk-manager-backend/actions/workflows/deploy-backend.yml/badge.svg?branch=main)](https://github.com/EmilyOng/tusk-manager-backend/actions/workflows/deploy-backend.yml)

Backend API: https://tusk-manager-backend.onrender.com

## Getting Started

Clone the repository.

```shell
git clone https://github.com/EmilyOng/tusk-manager-backend.git
```

### Prerequisites
1. go (version: [1.17](https://go.dev/doc/go1.17))
2. [cosmtrek/air](https://github.com/cosmtrek/air)
3. [tscriptify](https://github.com/tkrajina/typescriptify-golang-structs)
4. [Postgres](https://www.postgresql.org/download/)
5. [Postgres CLI Tools](http://postgresapp.com/documentation/cli-tools.html)

### Folder Structure
```
.
├── constants   # Stores constants used in the application
├── db          # Handles database connection and migration
├── handlers    # Handles incoming API requests
├── models      # Contains structs definitions for gorm
├── router      # Defines routes for the application
├── services    # Data access layer handling business logic that interfaces
├                 between the contollers layer and database systems
├── types       # Defines reusable types in the application
├── utils       # Contains utility functions
├── views       # Defines payloads/responses structs for the API resources
├                 shared with frontend
└── main.go     # Entry point to the application
```

### Setting up your environment
- (in `.env`) `AUTH_SECRET_KEY`: Requires any string
- (in `.env`) `DATABASE_URL`: The URL is obtained from Render's PostgreSQL deployment.

### Developing the application

Notably, the application uses the following tools:
- [Gorm](https://gorm.io/), to provide an ORM to abstract database interaction
- [Gin](https://github.com/gin-gonic/gin), as a HTTP web framework

Finally, there are several quick shortcuts to make development easier.

- Start the application: `make start`
  - The application uses [cosmtrek/air](https://github.com/cosmtrek/air) to provide live reload utility. Now, you can make changes to the files and the application will auto-reload.
- Generate types: `make generate-types`
  - This command generates TypeScript interfaces based on the Golang structs provided in [views](views) to ensure parity of types.

### Infrastructure

This application is hosted on [Render](https://render.com/), which can be found at https://tusk-manager-backend.onrender.com.

Extra work is required to define the environment variables in Render ([guide](https://render.com/docs/configure-environment-variables)).

**Provisioning Render Postgres**

> [Render Postgres](https://render.com/docs/databases) is a managed Postgres database service provided directly by Render.

### Deployment
Deployment is automatically handled by Github Actions when you push to the `main` branch of the repository.

- [**Github Actions**] [Deploy to Render](.github/workflows/deploy-backend.yml):
  - The deployment builds a Docker image and pushes to Render. More details on the Docker image at [Dockerfile](/backend/Dockerfile).
