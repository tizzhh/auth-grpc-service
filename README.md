# Auth 

Auth gRPC service

## Usage 💡

- clone the repository
- cd sso
- `make`
- `./auth_service`

## Requirements

- Go 1.22.5
- config in the format of local.yaml in config/ directory. Path to config should be set in CONFIG_PATH env variable

## Used packages / tools / stack

- gRPC.
- JWT based authentication.
- Sqlite3.
- Testing with [Mockery](https://github.com/vektra/mockery).
- [Migrations](https://github.com/golang-migrate/migrate).
- The usage of slog as the centralized logger.

## Available methods:

- Register - register new user.
- Login - login and get a JWT token.
- IsAdmin - check whether the user is an admin

##  Database design

#### url

| Column Name    | Datatype  | Not Null | Primary Key |
|----------------|-----------|----------|-------------|
| id             | INT      | ✅        | ✅           |
| email          | TEXT      | ✅        |             |
| pa          | TEXT      | ✅        |             |
| url         | TEXT      | ✅        |             |

#### users

| Column Name    | Datatype  | Not Null | Primary Key |
|----------------|-----------|----------|-------------|
| id             | INT      | ✅        | ✅           |
| email          | CHAR      | ✅        |             |
| pass_hash          | BLOB      | ✅        |             |
| is_admin          | CHAR      | ✅        |             |

#### apps

| Column Name    | Datatype  | Not Null | Primary Key |
|----------------|-----------|----------|-------------|
| id             | BIGINT      | ✅        | ✅           |
| name          | TEXT      | ✅        |             |
| secret         | TEXT      | ✅        |             |