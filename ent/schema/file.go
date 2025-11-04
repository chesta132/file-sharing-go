package schema

import (
	"file-sharing/internal/lib/crypto"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type File struct {
	ent.Schema
}

func (File) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("file_size"),
		field.String("file_name"),
		field.String("mime"),

		field.String("password").Optional().Nillable().Sensitive(),
		field.Int("max_downloads").Optional().Nillable(),

		field.String("id").DefaultFunc(func() string {
			return uuid.New().String()
		}).Unique(),
		field.String("token").DefaultFunc(crypto.CreateToken),
		field.Time("expires_at").Default(func() time.Time {
			return time.Now().AddDate(0, 0, 7)
		}),
		field.Int("download_count").Default(0),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (File) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("token"),
	}
}
