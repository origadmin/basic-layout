package dal

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/basic-layout/internal/mods/helloworld/dto"
)

type greeterDao struct {
	db  *Database
	log *log.Helper
}

// NewGreeterDao .
func NewGreeterDao(db *Database, logger log.Logger) dto.GreeterDao {
	return &greeterDao{
		db:  db,
		log: log.NewHelper(logger),
	}
}

func (r *greeterDao) Save(ctx context.Context, g *dto.Greeter) (*dto.Greeter, error) {
	return g, nil
}

func (r *greeterDao) Update(ctx context.Context, g *dto.Greeter) (*dto.Greeter, error) {
	return g, nil
}

func (r *greeterDao) FindByID(context.Context, int64) (*dto.Greeter, error) {
	return nil, nil
}

func (r *greeterDao) ListByHello(context.Context, string) ([]*dto.Greeter, error) {
	return nil, nil
}

func (r *greeterDao) ListAll(context.Context) ([]*dto.Greeter, error) {
	return nil, nil
}
