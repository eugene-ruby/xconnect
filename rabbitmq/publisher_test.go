package rabbitmq

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublisher_Publish_Success(t *testing.T) {
	mock := &mockChannel{}
	pub := NewPublisher(mock)

	err := pub.Publish("exchange", "key", []byte("test"))
	require.NoError(t, err)
	require.True(t, mock.published)
}

func TestPublisher_Publish_Error(t *testing.T) {
	mock := &mockChannel{publishErr: errors.New("publish failed")}
	pub := NewPublisher(mock)

	err := pub.Publish("exchange", "key", []byte("test"))
	require.Error(t, err)
	require.True(t, mock.published)
}

func TestPublisher_Close_Success(t *testing.T) {
	mock := &mockChannel{}
	pub := NewPublisher(mock)

	err := pub.Close()
	require.NoError(t, err)
	require.True(t, mock.closed)
}

func TestPublisher_Close_Error(t *testing.T) {
	mock := &mockChannel{closeErr: errors.New("close failed")}
	pub := NewPublisher(mock)

	err := pub.Close()
	require.Error(t, err)
	require.True(t, mock.closed)
}
