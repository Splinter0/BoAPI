# BoAPI

A boilerplate for building dockerized APIs with Golang and MongoDB

## Setup

Run `setup.sh`, this will:
- Generate credentials for db
- Generate cert and key for TLS connection
- Suggest commands to add users to the DB
- Create SECRET_KEY for JWT tokens and put them in `api/.env`
- Add authentication requirement to the db container
- Restart and run the containers

This will bind the MongoDB container to port `27017` and the API to `443`

## Running it

After the first setup you may start the API with `sudo docker-compose up`.
You may also uncomment in `api/main.go`:

```go
	// Authorization required
	r.GET("/", views.Index)

	//security.CreateAdmin("splinter", "wow") <---

	// Start HTTPS
	err := http.ListenAndServeTLS(":443", CRT, KEY, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
```

To create an admin user in your `users` table, this is used to authentication on the API.

## Contribute

Feel free to contribute by improving the code as it can definelity be improved. Do not add 
features though, as this should be a boilerplate!