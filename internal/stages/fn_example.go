package stages

import (
	"context"
	"errors"
	"go-pipeline/internal/model"
	"go-pipeline/internal/ports"
)

func ValidationFn() ports.StageFn[model.UserData] {
	return func(ctx context.Context, m model.UserData) (model.UserData, error) {
		if m.Email == "" || !containsAt2(m.Email) {
			return m, errors.New("validation: invalid email")
		}
		return m, nil
	}
}

func TransformFn() ports.StageFn[model.UserData] {
	return func(ctx context.Context, m model.UserData) (model.UserData, error) {
		if m.Name == "" {
			m.Name = "anonymous"
		}
		return m, nil
	}
}

func SinkFn(p ports.MessageQueueProducer) ports.StageFn[model.UserData] {
	return func(ctx context.Context, m model.UserData) (model.UserData, error) {
		return m, p.Produce(ctx, "users", m)
	}
}

func containsAt2(s string) bool {
	for i := range s {
		if s[i] == '@' {
			return true
		}
	}
	return false
}
