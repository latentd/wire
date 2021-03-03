package wire

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutingBasic(t *testing.T) {

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
		{
			name:   "unregistered root path on subrouer returns 404",
			path:   "/sub/",
			method: http.MethodGet,
			want:   http.StatusNotFound,
		},
		{
			name:   "registered handler on subrouter returns 200",
			path:   "/sub/test",
			method: http.MethodGet,
			want:   http.StatusOK,
		},
		{
			name:   "unregistered method on subrouter returns 405",
			path:   "/sub/test",
			method: http.MethodPost,
			want:   http.StatusMethodNotAllowed,
		},
		{
			name:   "unregistered path on subrouter returns 404",
			path:   "/sub/_test",
			method: http.MethodGet,
			want:   http.StatusNotFound,
		},
	}

	testF := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	r := NewRouter()
	r.GetF("/", testF)

	sr := r.SubRouter("/sub")
	sr.GetF("/test", testF)

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

func TestRoutingRegex(t *testing.T) {

	tts := []struct {
		name     string
		path     string
		method   string
		wantCode int
		wantVars string
	}{
		{
			name:     "single regex on path works",
			path:     "/test/1",
			method:   http.MethodGet,
			wantCode: http.StatusOK,
			wantVars: "1 ",
		},
		{
			name:     "multiple regex on path works",
			path:     "/test/1/test/2",
			method:   http.MethodGet,
			wantCode: http.StatusOK,
			wantVars: "1 2",
		},
	}

	testF := func(w http.ResponseWriter, r *http.Request) {

		vars := Vars(r)

		id := vars["id"]
		id2 := vars["id2"]

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(id + " " + id2))
	}

	r := NewRouter()
	r.GetF("/test/(id:[0-9]+)", testF)
	r.GetF("/test/(id:[0-9]+)/test/(id2:[0-9]+)", testF)

	srv := httptest.NewServer(r)
	defer srv.Close()

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			res, err := doRequest(t, tt.method, srv.URL+tt.path)
			if err != nil {
				t.Fatal(err)
			}
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			assertStatusCode(t, res.StatusCode, tt.wantCode)
			assertPathVars(t, string(b), tt.wantVars)
		})
	}
}

func TestRoutingMiddleware(t *testing.T) {

	tts := []struct {
		name     string
		path     string
		method   string
		wantCode int
		wantVars string
	}{
		{
			name:     "test middleware",
			path:     "/test",
			method:   http.MethodGet,
			wantCode: http.StatusOK,
			wantVars: "middleware\thandler",
		},
	}

	testF := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("handler"))
	}

	m := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("middleware\t"))
			h.ServeHTTP(w, r)
		})
	}

	r := NewRouter()
	r.Chain(m)
	r.GetF("/test", testF)

	srv := httptest.NewServer(r)
	defer srv.Close()

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			res, err := doRequest(t, tt.method, srv.URL+tt.path)
			if err != nil {
				t.Fatal(err)
			}
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			assertStatusCode(t, res.StatusCode, tt.wantCode)
			assertPathVars(t, string(b), tt.wantVars)
		})
	}
}

func assertStatusCode(t testing.TB, got int, want int) {
	t.Helper()

	if want != got {
		t.Errorf("got %d, want %d", got, want)
	}
}

func assertPathVars(t testing.TB, got, want string) {
	t.Helper()

	if want != got {
		t.Errorf("got %s, want %s", got, want)
	}
}

func doRequest(t testing.TB, method string, url string) (*http.Response, error) {
	t.Helper()

	cli := &http.Client{}

	req, err := http.NewRequestWithContext(context.TODO(), method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
