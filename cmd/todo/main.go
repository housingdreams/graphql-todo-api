package main

import (
	"github.com/leminhson2398/todo-api/internal/commands"
	_ "github.com/lib/pq"
)

func main() {
	commands.Execute()
}
