package main

import (
	"encoding/json"
	"fmt"
	"github.com/andviro/example-service/api"
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"time"
)

type Tickets struct {
	app    *Application
	bucket []byte
}

func (t *Tickets) Init(app *Application) {
	t.app = app
	t.bucket = []byte("tickets")
	err := t.app.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(t.bucket)
		return err
	})
	if err != nil {
		panic(err)
	}
}

func (t *Tickets) List(c context.Context, w http.ResponseWriter, r *http.Request) error {
	var res = make([]api.Ticket, 0)

	_ = t.app.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.bucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var t api.Ticket
			json.Unmarshal(v, &t)
			res = append(res, t)
		}
		return nil
	})
	return json.NewEncoder(w).Encode(res)
}

func (t *Tickets) View(c context.Context, w http.ResponseWriter, r *http.Request) error {
	var res api.Ticket
	id := c.Value("params").(httprouter.Params).ByName("id")

	_ = t.app.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.bucket)
		return json.Unmarshal(b.Get([]byte(id)), &res)
	})
	return json.NewEncoder(w).Encode(res)
}

func (t *Tickets) Create(c context.Context, w http.ResponseWriter, r *http.Request) error {
	var res api.Ticket
	place, err := strconv.Atoi(r.FormValue("place"))
	if err != nil {
		return err
	}
	res.Id = uuid.NewV4()
	res.Place = place
	res.Seance.Film = r.FormValue("film")
	res.Seance.Date = time.Now()

	_ = t.app.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.bucket)
		data, _ := json.Marshal(res)
		return b.Put(res.Id.Bytes(), data)
	})
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(res)
}

func (t *Tickets) Delete(c context.Context, w http.ResponseWriter, r *http.Request) error {
	id := c.Value("params").(httprouter.Params).ByName("id")
	uid, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	_ = t.app.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.bucket)
		return b.Delete(uid.Bytes())
	})
	if err != nil {
		return err
	}
	fmt.Fprintln(w, "Deleted", id)
	return nil
}
