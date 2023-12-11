package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func usage() {
	fmt.Println("help....")
	os.Exit(0)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		usage()
	}
	cmd := os.Args[1]

	switch cmd {
	case "run":
		if _, err := os.Stat("cmd/main.go"); err != nil {
			fmt.Println("not in slick app root: cmd/main.go not found")
			os.Exit(1)
		}
		if err := exec.Command("templ", "generate").Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		exec.Command("go", "run", "cmd/main.go").Run()
	case "install":
		if err := installProject(); err != nil {
			fmt.Println(err)
		}
	case "new":
		if len(os.Args) != 3 {
			usage()
		}
		name := os.Args[2]
		if err := generateProject(name); err != nil {
			fmt.Println(err)
		}
	}
}

func generateProject(name string) error {
	fmt.Println("creating new slick project:", name)
	if err := os.Mkdir(name, os.ModePerm); err != nil {
		return err
	}

	folders := []string{"model", "handler", "view", "cmd", "public"}

	for _, folder := range folders {
		if err := os.Mkdir(name+"/"+folder, os.ModePerm); err != nil {
			return err
		}
	}

	if err := os.WriteFile(name+"/go.mod", writeGoModContents(name), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(name+"/.env", writeEnvFileContents(), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(name+"/public/app.css", []byte(""), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(name+"/cmd/main.go", writeMainContents(name), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(name+"/handler/hello.go", writeHandlerContent(name), os.ModePerm); err != nil {
		return err
	}

	if err := os.Mkdir(name+"/view/hello", os.ModePerm); err != nil {
		return err
	}
	if err := os.Mkdir(name+"/view/layout", os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(name+"/view/layout/base.templ", writeBaseLayoutContent(), os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(name+"/view/hello/hello.templ", writeViewContent(name), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func installProject() error {
	start := time.Now()
	fmt.Println("installing project...")
	if err := exec.Command("go", "get", "github.com/anthdm/slick@latest").Run(); err != nil {
		return err
	}
	if err := exec.Command("go", "get", "github.com/a-h/templ").Run(); err != nil {
		return err
	}
	if err := exec.Command("templ", "generate").Run(); err != nil {
		return err
	}
	fmt.Printf("done installing project in %v\n", time.Since(start))
	return nil
}

func writeEnvFileContents() []byte {
	return []byte(`
SLICK_HTTP_LISTEN_ADDR=:3000

SLICK_SQL_DB_NAME=
SLICK_SQL_DB_USER=
SLICK_SQL_DB_PASSWORD=
SLICK_SQL_DB_HOST=
SLICK_SQL_DB_PORT=
`)
}

func writeMainContents(mod string) []byte {
	c := fmt.Sprintf(`
package main

import (
	"github.com/anthdm/slick"
	"%s/handler"
)

func main() {
	app := slick.New()
	app.Get("/", handler.HandleHelloIndex)
	app.Start()
}
`, mod)
	return []byte(c)
}

func writeGoModContents(mod string) []byte {
	buf := strings.Builder{}
	buf.WriteString("module " + mod)
	buf.WriteString("\n")
	buf.WriteString("\n")
	buf.WriteString("go 1.21.0")
	return []byte(buf.String())
}

func writeViewContent(mod string) []byte {
	c := fmt.Sprintf(`
package hello

import (
	"%s/view/layout"
)

templ Index() {
	@layout.Base() {
		<h1>hello there</h1>
	}
}
`, mod)

	return []byte(c)
}

func writeHandlerContent(mod string) []byte {
	c := fmt.Sprintf(`
package handler

import (
	"github.com/anthdm/slick"
	"%s/view/hello"
)

func HandleHelloIndex(c *slick.Context) error {
	return c.Render(hello.Index())
}
`, mod)

	return []byte(c)
}

func writeBaseLayoutContent() []byte {
	c := `
package layout

templ Base() {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>Slick Application</title>

			<script src="https://unpkg.com/htmx.org@1.9.9" defer></script>
		</head>
		<body>
			{ children... }
		</body>
	</html>
}
`
	return []byte(c)
}
