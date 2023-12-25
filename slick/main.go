package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type File struct {
	Path    string
	Content []byte
}

func writeFileWithCheck(file File) error {
	if err := os.WriteFile(file.Path, file.Content, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func main() {
	cmd := NewCommand()
	cmd.Register(
		runProject,
		installProject,
		generateProject,
		generateModel,
		generateView,
		generateHandler,
	)

	cmd.Execute()
}

func runProject() *cobra.Command {
	return &cobra.Command{
		Use:     "run",
		Example: "slick run",
		Short:   "Run slick development server",
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := os.Stat("cmd/main.go"); err != nil {
				fmt.Println("not in slick app root: cmd/main.go not found")
				return
			}
			if err := exec.Command("templ", "generate").Run(); err != nil {
				fmt.Println(err)
				return
			}

			if err := exec.Command("go", "run", "cmd/main.go").Run(); err != nil {
				fmt.Println(err)
			}
		},
	}
}

func installProject() *cobra.Command {
	return &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Example: "slick install",
		Short:   "Install project's dependency",
		Run: func(cmd *cobra.Command, args []string) {
			start := time.Now()
			fmt.Println("installing project...")
			if err := exec.Command("go", "get", "github.com/anthdm/slick@latest").Run(); err != nil {
				fmt.Println(err)
				return
			}

			if err := exec.Command("go", "get", "github.com/a-h/templ").Run(); err != nil {
				fmt.Println(err)
				return
			}
			if err := exec.Command("templ", "generate").Run(); err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("done installing project in %v\n", time.Since(start))
		},
	}
}

func generateProject() *cobra.Command {
	return &cobra.Command{
		Use:     "new",
		Example: "slick new hello-world",
		Short:   "Create new slick project",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("invalid arguments")
				return
			}

			name := args[0]

			fmt.Println("creating new slick project:", name)
			if err := os.Mkdir(name, os.ModePerm); err != nil {
				fmt.Println(err)
				return
			}

			files := []File{
				// setup directory
				{Path: name + "/model", Content: nil},
				{Path: name + "/handler", Content: nil},
				{Path: name + "/view", Content: nil},
				{Path: name + "/cmd", Content: nil},
				{Path: name + "/public", Content: nil},
				{Path: name + "/view/hello", Content: nil},
				{Path: name + "/view/layout", Content: nil},

				// setup files
				{Path: name + "/go.mod", Content: writeGoModContents(name)},
				{Path: name + "/.air.toml", Content: writeAirTomlContents()},
				{Path: name + "/.env", Content: writeEnvFileContents()},
				{Path: name + "/.gitignore", Content: writeGitignore()},
				{Path: name + "/public/app.css", Content: []byte("")},
				{Path: name + "/cmd/main.go", Content: writeMainContents(name)},
				{Path: name + "/handler/hello.go", Content: writeHandlerContent(name)},
				{Path: name + "/view/layout/base.templ", Content: writeBaseLayoutContent()},
				{Path: name + "/view/hello/hello.templ", Content: writeViewContent(name)},
			}

			errors := []error{}
			for _, file := range files {
				if file.Content == nil {
					if err := os.Mkdir(file.Path, os.ModePerm); err != nil {
						errors = append(errors, err)
					}
				} else {
					if err := writeFileWithCheck(file); err != nil {
						errors = append(errors, err)
					}
				}
			}

			if len(errors) != 0 {
				fmt.Println("slick encountered errors during file initialization:", errors)
				return
			}
		},
	}
}

func generateModel() *cobra.Command {
	return &cobra.Command{
		Use:     "model",
		Example: "slick model user",
		Short:   "Generate new model",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO:
		},
	}
}

func generateView() *cobra.Command {
	return &cobra.Command{
		Use:     "view",
		Example: "slick view user",
		Short:   "Generate new view",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO:

		},
	}
}

func generateHandler() *cobra.Command {
	return &cobra.Command{
		Use:     "handler",
		Example: "slick handler home",
		Short:   "Generate new handler",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO:

		},
	}
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

func writeDBFileContents() []byte {
	c := `
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	conn *sql.DB
	once sync.Once
)

type config struct {
	Hostname string
	Username string
	Password string
	Port     string
	DBName   string
}

func (c *config) Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		c.Username,
		c.Password,
		c.Hostname,
		c.Port,
		c.DBName,
	)

	log.Println("Creating database connection using DSN:", dsn)

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error opening database connection:", err)
		return nil, err
	}

	log.Println("Pinging database to verify connection...")

	err = conn.Ping()
	if err != nil {
		log.Println("Error pinging database:", err)
		return nil, err
	}

	log.Println("Database connection established successfully.")

	return conn, nil
}

func OpenConnection() {
	log.Println("Attempting to open database connection...")

	mysql := config{
		Hostname: os.Getenv("SLICK_SQL_DB_HOST"),
		Username: os.Getenv("SLICK_SQL_DB_USER"),
		Password: os.Getenv("SLICK_SQL_DB_PASSWORD"),
		Port:     os.Getenv("SLICK_SQL_DB_PORT"),
		DBName:   os.Getenv("SLICK_SQL_DB_NAME"),
	}

	once.Do(func() {
		log.Println("Creating database connection...")

		db, err := mysql.Connect()
		if err != nil {
			log.Println("Error connecting to database:", err)
			return
		}

		log.Println("Storing database connection...")
		conn = db
	})

	log.Println("Database connection opened.")
}

func GetConnection() *sql.DB {
	if conn == nil {
		OpenConnection()
	}

	return conn
}

type Config struct {
	DB *sql.DB
}

func NewConfig() *Config {
	/**
	 * Open connection
	 */

	log.Println("Attempting to retrieve database connection...")
	db := GetConnection()
	log.Printf("Database connection retrieved: %v\n", db)

	if db == nil {
		log.Println("ERROR: Database connection failed.")
		return nil
	}

	config := Config{
		DB: db,
	}

	return &config
}

`
	return []byte(c)
}

func writeMainContents(mod string) []byte {
	c := fmt.Sprintf(`
package main

import (
	"log"
	"github.com/anthdm/slick"
	"%s/handler"
)

func main() {
	app := slick.New()
	app.Get("/", handler.HandleHelloIndex)
	log.Fatal(app.Start())
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

func writeAirTomlContents() []byte {
	c := `
root = "."
testdata_dir = "testdata"
tmp_dir = ".build"

[build]
  args_bin = []
  bin = "./.build/main"
  cmd = "templ generate && go build -o ./.build/main ./cmd"
  delay = 1000
  exclude_dir = ["assets", ".build", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go", "_templ.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html", "templ"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  post_cmd = []
  pre_cmd = []
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true

`
	return []byte(c)
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

func writeGitignore() []byte {
	c := `
bin
.build
.env
`
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
