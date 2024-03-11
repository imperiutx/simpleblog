package db

import (
	"context"
)

type QuerierDynamic interface {
	ListPostsDynamic(ctx context.Context, title string, tags []string, filters Filters) ([]*Post, Metadata, error)
}

var _ QuerierDynamic = (*QueriesDynamic)(nil)