// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"origadmin/basic-layout/internal/mods/helloworld/dal/ent/greeter"
	"origadmin/basic-layout/internal/mods/helloworld/dal/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// GreeterUpdate is the builder for updating Greeter entities.
type GreeterUpdate struct {
	config
	hooks     []Hook
	mutation  *GreeterMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the GreeterUpdate builder.
func (gu *GreeterUpdate) Where(ps ...predicate.Greeter) *GreeterUpdate {
	gu.mutation.Where(ps...)
	return gu
}

// SetName sets the "name" field.
func (gu *GreeterUpdate) SetName(s string) *GreeterUpdate {
	gu.mutation.SetName(s)
	return gu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (gu *GreeterUpdate) SetNillableName(s *string) *GreeterUpdate {
	if s != nil {
		gu.SetName(*s)
	}
	return gu
}

// Mutation returns the GreeterMutation object of the builder.
func (gu *GreeterUpdate) Mutation() *GreeterMutation {
	return gu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (gu *GreeterUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, gu.sqlSave, gu.mutation, gu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (gu *GreeterUpdate) SaveX(ctx context.Context) int {
	affected, err := gu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (gu *GreeterUpdate) Exec(ctx context.Context) error {
	_, err := gu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gu *GreeterUpdate) ExecX(ctx context.Context) {
	if err := gu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (gu *GreeterUpdate) check() error {
	if v, ok := gu.mutation.Name(); ok {
		if err := greeter.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Greeter.name": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (gu *GreeterUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *GreeterUpdate {
	gu.modifiers = append(gu.modifiers, modifiers...)
	return gu
}

func (gu *GreeterUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := gu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(greeter.Table, greeter.Columns, sqlgraph.NewFieldSpec(greeter.FieldID, field.TypeInt))
	if ps := gu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := gu.mutation.Name(); ok {
		_spec.SetField(greeter.FieldName, field.TypeString, value)
	}
	_spec.AddModifiers(gu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, gu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{greeter.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	gu.mutation.done = true
	return n, nil
}

// GreeterUpdateOne is the builder for updating a single Greeter entity.
type GreeterUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *GreeterMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetName sets the "name" field.
func (guo *GreeterUpdateOne) SetName(s string) *GreeterUpdateOne {
	guo.mutation.SetName(s)
	return guo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (guo *GreeterUpdateOne) SetNillableName(s *string) *GreeterUpdateOne {
	if s != nil {
		guo.SetName(*s)
	}
	return guo
}

// Mutation returns the GreeterMutation object of the builder.
func (guo *GreeterUpdateOne) Mutation() *GreeterMutation {
	return guo.mutation
}

// Where appends a list predicates to the GreeterUpdate builder.
func (guo *GreeterUpdateOne) Where(ps ...predicate.Greeter) *GreeterUpdateOne {
	guo.mutation.Where(ps...)
	return guo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (guo *GreeterUpdateOne) Select(field string, fields ...string) *GreeterUpdateOne {
	guo.fields = append([]string{field}, fields...)
	return guo
}

// Save executes the query and returns the updated Greeter entity.
func (guo *GreeterUpdateOne) Save(ctx context.Context) (*Greeter, error) {
	return withHooks(ctx, guo.sqlSave, guo.mutation, guo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (guo *GreeterUpdateOne) SaveX(ctx context.Context) *Greeter {
	node, err := guo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (guo *GreeterUpdateOne) Exec(ctx context.Context) error {
	_, err := guo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (guo *GreeterUpdateOne) ExecX(ctx context.Context) {
	if err := guo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (guo *GreeterUpdateOne) check() error {
	if v, ok := guo.mutation.Name(); ok {
		if err := greeter.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Greeter.name": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (guo *GreeterUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *GreeterUpdateOne {
	guo.modifiers = append(guo.modifiers, modifiers...)
	return guo
}

func (guo *GreeterUpdateOne) sqlSave(ctx context.Context) (_node *Greeter, err error) {
	if err := guo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(greeter.Table, greeter.Columns, sqlgraph.NewFieldSpec(greeter.FieldID, field.TypeInt))
	id, ok := guo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Greeter.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := guo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, greeter.FieldID)
		for _, f := range fields {
			if !greeter.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != greeter.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := guo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := guo.mutation.Name(); ok {
		_spec.SetField(greeter.FieldName, field.TypeString, value)
	}
	_spec.AddModifiers(guo.modifiers...)
	_node = &Greeter{config: guo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, guo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{greeter.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	guo.mutation.done = true
	return _node, nil
}

// SetGreeter set the Greeter
func (gu *GreeterUpdate) SetGreeter(input *Greeter, fields ...string) *GreeterUpdate {
	m := gu.mutation
	_ = m.SetFields(input, fields...)
	return gu
}

// SetGreeter set the Greeter
func (guo *GreeterUpdateOne) SetGreeter(input *Greeter, fields ...string) *GreeterUpdateOne {
	m := guo.mutation
	_ = m.SetFields(input, fields...)
	return guo
}

// Omit allows the unselect one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func (guo *GreeterUpdateOne) Omit(fields ...string) *GreeterUpdateOne {
	omits := make(map[string]struct{}, len(fields))
	for i := range fields {
		omits[fields[i]] = struct{}{}
	}
	guo.fields = []string(nil)
	for _, col := range greeter.Columns {
		if _, ok := omits[col]; !ok {
			guo.fields = append(guo.fields, col)
		}
	}
	return guo
}
