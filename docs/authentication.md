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
If user was successfully registered: HTTP 200
```