package main

import (
	"net/http"
	"testing"

	"github.com/Cerecero/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	//Create new instance
	app := newTestApplication(t)

	// Create new test server
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}
