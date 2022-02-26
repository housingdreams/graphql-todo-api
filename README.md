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
go run cmd/mage/main.go backend:genMigrations backend:build
```

```
go run cmd/todo/main.go token
go run cmd/todo/main.go web
```

Then go to http://localhost:8000
