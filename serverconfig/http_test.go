package serverconfig

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartServer_ShutdownIn5ms(t *testing.T) {
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	assert.NotPanics(t, func() { StartServer(ctxShutDown) })
}
