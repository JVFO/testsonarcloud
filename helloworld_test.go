package testsonarcloud

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
	"testing"
)

func TestName(t *testing.T) {
	_, _ = Receive(context.Background(), event.New())
}
