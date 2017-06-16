package task

import (
	"testing"
	"Andariel/pkg/mongo"
)

func TestServer(t *testing.T) {
	sess := mongo.InitMongoSes("mongodb://127.0.0.1/github")

	m := MgoQueueEngine{
		Sess:   sess.MongoSession,
	}
	s := ServerOption{
		queueEngine:   &m,
		WorkerSize:   10,
	}

	server := NewServer(&s)

	server.Register(1, func(task *Task) error {
		return nil
	})

	server.Start()
}
