# JWT Go 

Example application that implements JWT based authentication. 
To run this application, run the Go binary:

```sh
go run ./...
```

Now, using any HTTP client with support for cookies (vscode ext humao.rest-client), or your web browser) make a sign-in request with the appropriate credentials:

```
POST http://localhost:8000/signin

{"username":"user1","password":"password1"}
```

You can now try hitting the welcome route from the same client to get the welcome message:

```
GET http://localhost:8000/welcome
```

Hit the refresh route, and then inspect the clients cookies to see the new value of the `token` cookie, ensuring that a new token is not issued until enough time has elapsed (a new token will only be issued if the old token is within 30 seconds of expiry. Otherwise, return a bad request status)

```
POST http://localhost:8000/refresh
```
