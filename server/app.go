package main

import (
	"github.com/andviro/noodle"
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
)

type Application struct {
	Port string
	*bolt.DB
	*httprouter.Router
}

func (app *Application) Run() error {
	return http.ListenAndServe(app.Port, app)
}

func H(h noodle.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(context.TODO(), "params", p)
		h(ctx, w, r)
	}
}

func AppFactory(port string) *Application {
	bdb, err := bolt.Open("data.bdb", 0600, nil)
	if err != nil {
		panic(err)
	}
	app := &Application{
		Port:   port,
		Router: httprouter.New(),
		DB:     bdb,
	}

	tickets := Tickets{}
	tickets.Init(app)

	root := noodle.Default()

	app.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, "site/index.html")
	})
	app.ServeFiles("/static/*filepath", http.Dir("site/static"))
	app.GET("/tickets", H(root.Then(tickets.List)))
	app.POST("/tickets", H(root.Then(tickets.Create)))
	app.GET("/tickets/:id", H(root.Then(tickets.View)))
	app.GET("/tickets/:id/delete", H(root.Then(tickets.Delete)))

	return app
}
