# backend

[![Deploy backend](https://github.com/EmilyOng/cvwo-backend/actions/workflows/deploy-backend.yml/badge.svg?branch=main)](https://github.com/EmilyOng/cvwo-backend/actions/workflows/deploy-backend.yml)

Backend API: https://tusk-manager-backend.herokuapp.com/

## Getting Started

Clone the repository.

```shell
git clone https://github.com/EmilyOng/cvwo-backend.git
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
├── controllers # Handles incoming requests
├── db          # Handles database connection
├── models      # Contains structs definitions
├── services    # Contains API logics
├── utils       # Contains useful utility functions
└── main.go     # Entry point to the application
```

### Setting up your environment
- (in `.env`) `AUTH_SECRET_KEY`: Requires any string
- (in `.env`) `DATABASE_URL`: The url is obtained by executing the command `heroku config:get DATABASE_URL -a tusk-manager-backend`. Refer to this [article](https://devcenter.heroku.com/articles/connecting-to-heroku-postgres-databases-from-outside-of-heroku) for more information.

### Developing the application

Notably, the application uses the following tools:
- [Gorm](https://gorm.io/), to provide an ORM to abstract database interaction
- [Gin](https://github.com/gin-gonic/gin), as a HTTP web framework

Finally, there are several quick shortcuts to make development easier.

- Start the application: `make start`
  - The application uses [cosmtrek/air](https://github.com/cosmtrek/air) to provide live reload utility. Now, you can make changes to the files and the application will auto-reload.
- Generate types: `make generate-types`
  - This command generates TypeScript interfaces based on the Golang structs provided in [models](models) to ensure parity of types.

### Infrastructure

This application is hosted on [Heroku](https://www.heroku.com/), which can be found at https://tusk-manager-backend.herokuapp.com/.

**Provisioning Heroku Postgres**

> [Heroku Postgres](https://elements.heroku.com/addons/heroku-postgresql) is a managed SQL database service provided directly by Heroku.

```shell
heroku addons:create heroku-postgresql:hobby-dev -a tusk-manager-backend
```

**Accessing Heroku Postgres**

> Note: Make sure that [Postgres CLI Tools](http://postgresapp.com/documentation/cli-tools.html) is installed.

```shell
heroku pg:psql -a tusk-manager-backend
```

### Deployment
Deployment is automatically handled by Github Actions when you push to the `main` branch of the repository.

- [**Github Actions**] [Deploy to Heroku](.github/workflows/deploy-backend.yml): `make push/heroku`
  - The deployment builds a Docker image and pushes to Heroku. More details on the Docker image at [Dockerfile](/backend/Dockerfile).
