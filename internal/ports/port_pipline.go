package ports

import (
	"context"

	"go-pipeline/internal/model"
)

type RegistryPipeline interface {
	Run(ctx context.Context, in chan model.UserData) error
}
