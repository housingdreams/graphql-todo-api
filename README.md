## TODO API GRAPHQL

```
cp .env.sample .env
```

Create postgresql database

```
docker build -t postgres-db . && docker run -it -p 5432:5432 postgres-db
```

You can connect using dbeaver (`brew install --cask dbeaver-community`) or whatever GUI

```
Host: localhost
Port: 5432
User: postgres
Password: docker
Datbase: docker
```

```
go run cmd/todo/main.go password password
// $2a$10$xmfJsn8UBGbQITsmZYXeX.dvLLLPO0ISai/c8QuFefrtgAaKx0VdO
```

Create the first user

```
INSERT INTO user_account (user_id, first_name, last_name, username, email, password_hash) values
('3909fd4e-e8bd-4306-a66f-7e1c734b9cd2', 'Admin', 'istrator', 'admin', 'admin@example.com', '...');
```

start the go server

```
go run cmd/todo/main.go web
```

Then go to http://localhost:8888


Then generate a JWT for them

```
go run cmd/todo/main.go token 3909fd4e-e8bd-4306-a66f-7e1c734b9cd2
```

Add it to the http headers on the graphql page

```
{
  "Authorization": "Bearer eyJhbGci.......OiJIUzI1N"
}
```
