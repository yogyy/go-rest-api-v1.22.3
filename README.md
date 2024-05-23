# Learning Golang with Restful API

to run the Server use

```bash
go run *.go
```

```bash
http://127.0.0.1:8080
```

## Routes

using prefix `/api/v1/`

### Register

```json
POST /auth/register

{
  "firstName": "",
  "lastName": "",
  "password": "6 charactest long"
}
```

### Login

```json
POST /auth/login

{
  "number": number,
  "password": ""
}
```

### Get all accounts

```json
GET /account
```

### Get account by id

- This endpoint requires a valid JSON Web Token (JWT) for authorization.
- Each JWT token is unique and grants access only to the account associated with the token itself.

need jwt token, each token only available to get account by token

```json
GET /account/12345 HTTP/1.1
Authorization: Bearer <your_jwt_token_here>
```

### Delete account by id

```json
DELETE /account/{id}
```
