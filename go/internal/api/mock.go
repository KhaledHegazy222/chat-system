package api

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockRedisClientWrapper is a mock of the RedisClientWrapper for testing
type MockRedisClientWrapper struct {
	mock.Mock
}

func (m *MockRedisClientWrapper) KeyExists(ctx context.Context, hashKey, field string) (bool, error) {
	args := m.Called(ctx, hashKey, field)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisClientWrapper) IncrementField(ctx context.Context, hashKey, field string) (int64, error) {
	args := m.Called(ctx, hashKey, field)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClientWrapper) PushToQueue(ctx context.Context, queueName string, data interface{}) error {
	args := m.Called(ctx, queueName, data)
	return args.Error(0)
}
