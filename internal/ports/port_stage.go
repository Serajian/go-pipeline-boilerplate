package ports

import (
	"context"

	"go-pipeline/internal/model"
)

type RegistryStage interface {
	Name() string
	Execute(ctx context.Context, in chan model.UserData) error
}
