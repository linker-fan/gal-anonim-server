# Authentication


## Register

### Endpoint
```
/users/register
Method: POST
```

Register endpoint registers a new user in the database.
### Example payload:
```
{
    "username": "hyperxpizza",
    "password1": "SomeTestPassword1#",
    "password2": "SomeTestPassword1#"
}
```

### Response:
```
If server can't unmarshal payload: HTTP 400
If username validation fails: HTTP 406
If username is already in the database: HTTP 409
If password validation fails: HTTP 406
If server fails: HTTP 500
If user was successfully registered: HTTP 201
```

## Login

### Endpoint
```
/users/login
METHOD: POST
```

Login endpoint logs user into the application and sets a http only cookie with a JWT token.

### Example payload
```
{
    "username": "hyperxpizza",
    "password1": "SomeTestPassword1#"
}
```

### Response
```
If server can't unmarshal payload: HTTP 400
If username does not exist in the database: HTTP 404
If password from the payload does not match the password hash from the database: HTTP 401
If server fails: HTTP 500
If everything is ok: HTTP 200
```

## Me

### Endpoint
```
/protected/me
METHOD: GET
```
Me Handler returns information contained in the jwt token.
Authentication required, cookie with the token must be set in order for the request to be successfull

### Response
```
If unauthenticated: HTTP 401
If server fails: HTTP 500
If everything OK: 
HTTP 200
{
    "id": 1,
    "username": hyperxpizza,
    "is_admin": false
}
```