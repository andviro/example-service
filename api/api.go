package api

import (
	"github.com/satori/go.uuid"
	"time"
)

type Ticket struct {
	Id     uuid.UUID `json:"id"`
	Seance Seance    `json:"seance"`
	Place  int       `json:"place"`
}

type Seance struct {
	Film string    `json:"film"`
	Date time.Time `json:"date"`
}
