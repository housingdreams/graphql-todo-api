** TODO API GRAPHQL **

NOTE:
    prerequisites: 
    1) golang >= 1.14
    2) nodejs >= 12
    3) npm >= 6

1) create postgresql database
2) create file similar to file inside conf folder
3) run: go run cmd/mage/main.go build
4) run: go run dist/todo migrate
5) run: go run dist/todo web

Then go to http://localhost:5555
