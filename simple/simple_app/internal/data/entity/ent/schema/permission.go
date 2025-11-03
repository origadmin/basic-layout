package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Permission holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

// Fields of the Permission.
func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").
			Comment("The ID of the user.").
			Immutable().
			NotEmpty(),
		field.String("object").
			Comment("The object on which the action is performed.").
			Immutable().
			NotEmpty(),
		field.String("action").
			Comment("The action being performed on the object.").
			Immutable().
			NotEmpty(),
	}
}

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
	return nil
}

// Indexes of the Permission.
func (Permission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "object", "action").
			Unique(), // Ensure unique combination of user, object, and action
	}
}
