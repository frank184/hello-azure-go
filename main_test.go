package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	server := httptest.NewServer(newHandler())
	defer server.Close()

	response, err := server.Client().Get(server.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if got := response.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	var body map[string]string
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body) != 1 || body["status"] != "ok" {
		t.Errorf("body = %v, want only status:ok", body)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		t.Errorf("unexpected trailing response data: %v", err)
	}
}

func TestRouting(t *testing.T) {
	cases := []struct {
		name   string
		method string
		path   string
		status int
	}{
		{"head health", http.MethodHead, "/health", http.StatusOK},
		{"reject post", http.MethodPost, "/health", http.StatusMethodNotAllowed},
		{"unknown path", http.MethodGet, "/unknown", http.StatusNotFound},
		{"health subpath", http.MethodGet, "/health/extra", http.StatusNotFound},
	}
	server := httptest.NewServer(newHandler())
	defer server.Close()

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(test.method, server.URL+test.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			response, err := server.Client().Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != test.status {
				t.Errorf("status = %d, want %d", response.StatusCode, test.status)
			}
			if test.method == http.MethodHead {
				body, err := io.ReadAll(response.Body)
				if err != nil {
					t.Fatal(err)
				}
				if len(body) != 0 {
					t.Errorf("HEAD returned a body: %q", body)
				}
			}
		})
	}
}

func BenchmarkHealth(b *testing.B) {
	handler := newHandler()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	b.ReportAllocs()
	for b.Loop() {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			b.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	}
}
