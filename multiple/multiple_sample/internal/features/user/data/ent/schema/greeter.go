/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Greeter holds the schema definition for the Greeter entity.
type Greeter struct {
	ent.Schema
}

// Fields of the Greeter.
func (Greeter) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty(),
	}
}

// Edges of the Greeter.
func (Greeter) Edges() []ent.Edge {
	return nil
}
