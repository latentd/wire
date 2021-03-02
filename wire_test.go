package wire

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootRouting(t *testing.T) {

	tts := []struct {
		name   string
		path   string
		method string
		want   int
	}{
		{
			name:   "registered handler returns 200",
			path:   "/",
			method: http.MethodGet,
			want:   http.StatusOK,
		},
		{
			name:   "unregistered method returns 405",
			path:   "/",
			method: http.MethodPost,
			want:   http.StatusMethodNotAllowed,
		},
        {
            name: "unregistered path returns 404",
            path: "/test",
            method: http.MethodGet,
            want: http.StatusNotFound,
        },
	}

	r := NewRouter()

	r.GetF("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(r)
	defer srv.Close()

	for _, tt := range tts {

		t.Run(tt.name, func(t *testing.T) {

			cli := &http.Client{}

			req, err := http.NewRequestWithContext(
				context.TODO(),
				tt.method,
				srv.URL+tt.path,
				nil,
			)
			if err != nil {
				t.Fatal(err)
			}

			res, err := cli.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if res.StatusCode != tt.want {
				t.Errorf("want %d, got %d", tt.want, res.StatusCode)
			}
		})

	}

}
