package api

import (
	"context"
	"github.com/stashapp/stash/internal/manager"
)

func (r *queryResolver) MetadataScan(ctx context.Context) (string, error) {
	manager.GetInstance().Scan()
	return "todo", nil
}

func (r *queryResolver) MetadataImport(ctx context.Context) (string, error) {
	manager.GetInstance().Import()
	return "todo", nil
}

func (r *queryResolver) MetadataExport(ctx context.Context) (string, error) {
	manager.GetInstance().Export()
	return "todo", nil
}

func (r *queryResolver) MetadataGenerate(ctx context.Context) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) MetadataClean(ctx context.Context) (string, error) {
	panic("not implemented")
}
