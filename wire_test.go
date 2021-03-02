package wire

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {

	r := NewRouter()

	r.GetF("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(r)
	defer srv.Close()

	cli := &http.Client{}
	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		srv.URL,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cli.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, res.StatusCode)
	}

}
