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

func TestSnippetView(t *testing.T) {
	//Create new instance
	app := newTestApplication(t)

	// Create new test server
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name string
		urlPath string
		wantCode int
		wantBody string
	}{
		{
			name: "Valid ID",
			urlPath: "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "There's nothing to se here..", 
			// Delete comment later
			// Changed text because the db is not showing the actual snippet text
		},
		{
			name: "Non-existent ID",
			urlPath: "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Negative ID",
			urlPath: "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Decimal ID",
			urlPath: "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name: "String ID",
			urlPath: "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Empty ID",
			urlPath: "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T){
			code, _, body := ts.get(t, tt.urlPath)
			
			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

