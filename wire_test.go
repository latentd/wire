package wire

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func assertStatusCode(t testing.TB, want int, got int) {
	t.Helper()

	if want != got {
		t.Errorf("got %d, want %d", got, want)
	}
}

func doRequest(t testing.TB, method string, url string) (*http.Response, error) {
	t.Helper()

	cli := &http.Client{}

	req, err := http.NewRequestWithContext(
		context.TODO(),
		method,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

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
			name:   "unregistered path returns 404",
			path:   "/test",
			method: http.MethodGet,
			want:   http.StatusNotFound,
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
			res, err := doRequest(t, tt.method, srv.URL+tt.path)
			if err != nil {
				t.Fatal(err)
			}
			assertStatusCode(t, res.StatusCode, tt.want)

		})

	}

}
