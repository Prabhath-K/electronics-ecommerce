package role

import "context"

type Repository interface {
	FindByName(ctx context.Context, name string) (*Role, error)
}
