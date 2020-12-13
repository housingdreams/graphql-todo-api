//+build mage

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/shurcooL/vfsgen"
)

type Frontend mg.Namespace

func (Frontend) Install() error {
	return sh.Run("yarn", "--cwd", "frontend", "install")
}

func (Frontend) Build() error {
	return sh.RunV("yarn", "--cwd", "frontend", "build")
}

type Backend mg.Namespace

func (Backend) GenMigrations() error {
	if _, err := os.Stat("internal/migrations"); os.IsNotExist(err) {
		os.Mkdir("internal/migrations/", 0755)
	}
	var fs http.FileSystem = http.Dir("migrations")
	err := vfsgen.Generate(fs, vfsgen.Options{
		Filename:     "internal/migrations/migrations_generated.go",
		PackageName:  "migrations",
		VariableName: "Migrations",
	})
	if err != nil {
		panic(err)
	}
	return nil
}

func isIn(check string, deps []string) bool {
	for _, dep := range deps {
		if check == dep {
			return true
		}
	}
	return false
}

func pyAndPip() (pythonExe string, pipExe string) {
	pythonExe = "python3"
	pipExe = "pip3"

	platform := runtime.GOOS
	if platform == "windows" {
		pythonExe = "python"
		pipExe = "pip"
	}

	return pythonExe, pipExe
}

func checkDependencies() error {
	deps := []string{"django", "pytz", "sqlparse", "asgiref"}
	installedDeps := []string{}
	_, pipExe := pyAndPip()

	out, err := exec.Command(pipExe, "list").Output()
	if err != nil {
		return err
	}
	strOut := string(out)
	scanner := bufio.NewScanner(strings.NewReader(strOut))
	scanner.Split(bufio.ScanLines)
	i := 0
	for scanner.Scan() {
		i++
		if i == 1 || i == 2 { // first 2 line are just header of dependencies
			continue
		}
		packageName := strings.Split(scanner.Text(), " ")[0]
		packageName = strings.TrimSpace(packageName)
		packageName = strings.ToLower(packageName)
		if isIn(packageName, deps) {
			installedDeps = append(installedDeps, packageName)
		}
	}

	for _, dep := range deps {
		installed := false
		for _, installedDep := range installedDeps {
			if installedDep == dep {
				installed = true
			}
		}
		if !installed {
			log.Println(dep, "not installed. Install it with pip first")
		}
	}

	return nil
}

func (Backend) SqlMigrate() error {
	err := checkDependencies()
	if err != nil {
		return err
	}
	// pythonExe, _ := pyAndPip()
	// out, err := exec.Command("python3", "tododb/manage.py", "sqlmigrate", "user", "0001").Output()
	// fmt.Println(string(out))
	// return err
	files, err := ioutil.ReadDir("./websocket")
	if err != nil {
		return err
	}

	for _, f := range files {
		// if f.IsDir() {
		// 	continue
		// }
		if er := os.Remove(filepath.Join("websocket", f.Name())); er != nil {
			return er
		}
	}
	return nil
}

func (Backend) GenFrontEnd() error {
	if _, err := os.Stat("internal/frontend"); os.IsNotExist(err) {
		os.Mkdir("internal/frontend", 0755)
	}
	var fs http.FileSystem = http.Dir("frontend/build")
	err := vfsgen.Generate(fs, vfsgen.Options{
		Filename:     "internal/frontend/frontend_generated.go",
		PackageName:  "frontend",
		VariableName: "Frontend",
	})
	if err != nil {
		panic(err)
	}
	return nil
}

func (Backend) Build() error {
	fmt.Println("compiling binary dist/todo")
	return sh.Run("go", "build", "-tags", "prod", "-o", "dist/todo", "cmd/todo/main.go")
}

func (Backend) Schema() error {
	files, err := ioutil.ReadDir("internal/graph/schema/")
	if err != nil {
		panic(err)
	}

	var schemaBuf strings.Builder

	for _, file := range files {
		filename := "internal/graph/schema/" + file.Name()
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		content, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		trimContent := strings.TrimSpace(string(content))
		fmt.Fprintln(&schemaBuf, trimContent)
	}

	err = ioutil.WriteFile("internal/graph/schema.graphql", []byte(schemaBuf.String()), os.FileMode(0755))
	if err != nil {
		panic(err)
	}
	return sh.Run("go", "generate", "./internal/graph")
	// return nil
}

func Build() {
	mg.SerialDeps(Frontend.Build, Backend.GenMigrations, Backend.GenFrontEnd, Backend.Build)
}

func Install() {
	mg.SerialDeps(Frontend.Install)
}
