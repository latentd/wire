package wire

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutingMethods(t *testing.T) {

	tts := []struct {
		name   string
		path   string
		method string
		want   int
	}{
		{
			name:   "registered GET path with handler allows GET",
			path:   "/h",
			method: http.MethodGet,
			want:   http.StatusOK,
		},
		{
			name:   "registered POST path with handler allows POST",
			path:   "/h",
			method: http.MethodPost,
			want:   http.StatusOK,
		},
		{
			name:   "registered GET path with handlerFunc allows GET",
			path:   "/f",
			method: http.MethodGet,
			want:   http.StatusOK,
		},
		{
			name:   "registered POST path with handlerFunc allows POST",
			path:   "/f",
			method: http.MethodPost,
			want:   http.StatusOK,
		},
		{
			name:   "registered ALL path with handler allows GET",
			path:   "/ah",
			method: http.MethodGet,
			want:   http.StatusOK,
		},
		{
			name:   "registered ALL path with handlerFunc allows GET",
			path:   "/af",
			method: http.MethodGet,
			want:   http.StatusOK,
		},
	}

	r := NewRouter()

	testF := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	testH := http.HandlerFunc(testF)

	r.Get("/h", testH)
	r.Post("/h", testH)

	r.GetF("/f", testF)
	r.PostF("/f", testF)

	r.All("/ah", testH)
	r.AllF("/af", testF)

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
