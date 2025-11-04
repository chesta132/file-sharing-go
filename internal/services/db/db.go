package db

import (
	"context"
	"file-sharing/ent"
	"fmt"
	"log"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
)

func Connect(path string, autoCreateSchema bool) *ent.Client {
	client, err := ent.Open(dialect.SQLite, fmt.Sprintf("file:%v?cache=shared&_fk=1", path))
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	if autoCreateSchema {
		ctx := context.Background()
		if err := client.Schema.Create(ctx); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
	}

	return client
}
