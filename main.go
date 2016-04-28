package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gocraft/web"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var (
	BuildRevision string = "development"
)

var App = &cli.App{
	Name:        "grinder",
	HelpName:    "grinder",
	Usage:       "coffee dating",
	UsageText:   "",
	HideVersion: true,
	Writer:      os.Stdout,
	Action:      cli.ShowAppHelp,

	Commands: []cli.Command{
		{
			Name:   "server",
			Usage:  "launch server",
			Action: Server,

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "host",
					Value:  "127.0.0.1",
					Usage:  "http server address",
					EnvVar: "HOST",
				},
				cli.StringFlag{
					Name:   "port",
					Value:  "8080",
					Usage:  "http server port",
					EnvVar: "PORT",
				},
				cli.StringFlag{
					Name:   "database",
					Value:  "postgres://localhost/grinder?sslmode=disable",
					Usage:  "postgresql database dsn",
					EnvVar: "DATABASE_URL",
				},
			},
		},
		{
			Name:   "version",
			Usage:  "show version information",
			Action: Version,
		},
	},
}

func Server(c *cli.Context) {
	router := web.New(Grinder{})

	// open database connection
	db, err := sqlx.Open("postgres", c.String("database"))
	if err != nil {
		log.Fatalln("fatal: database:", err)
	}
	log.Println("info: database:", c.String("database"))

	// middleware to inject database connection
	router.Middleware(func(g *Grinder, w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
		g.DB = db.Unsafe()

		next(w, r)
	})

	// define unauthenticated routes
	router.Get("/claim", (*Grinder).Claims)
	router.Post("/claim/:user_id", (*Grinder).Claim)

	internal := router.Subrouter(Internal{}, "/internal")
	internal.Get("/reset", (*Internal).Reset)
	internal.Get("/match", (*Internal).Match)

	// define authenticated routes
	users := router.Subrouter(Users{}, "/user")
	users.Middleware((*Users).Auth)
	users.Get("/", (*Users).Index)
	users.Post("/available", (*Users).Available)

	matches := users.Subrouter(Matches{}, "/match")
	matches.Get("/", (*Matches).Index)
	matches.Post("/:user_id", (*Matches).Match)

	// create http listen addres
	addr := c.String("host") + ":" + c.String("port")
	log.Println("info: listen:", addr)

	// launch http server
	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalln("fatal: listen:", err)
	}
}

func Version(c *cli.Context) {
	fmt.Fprintln(c.App.Writer, "Build Revision:", BuildRevision)
}

func main() {
	App.Run(os.Args)
}
