package main

import (
	"net/http"
	"snippetbox/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {
	t.Parallel()

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, "pong", body)
}
