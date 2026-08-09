package vps_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/paultibbetts/mythicbeasts-client-go"
	vpsapi "github.com/paultibbetts/mythicbeasts-client-go/vps"
)

func TestMakeDormant(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/vps/servers/my-id/dormant", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method=%s, want PUT", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Fatalf("Content-Type=%s, want application/json", ct)
		}

		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode req: %v", err)
		}
		if req["dormant"] != true {
			t.Fatalf("dormant=%v, want true", req["dormant"])
		}
		if _, ok := req["product"]; ok {
			t.Fatalf("product=%v, want omitted when making dormant", req["product"])
		}

		_, _ = w.Write([]byte(`{"message":"Operation successful"}`))
	})
	c, srv := newTestClient(t, mux)
	defer srv.Close()

	resp, err := c.VPS().MakeDormant(testContext(), "my-id")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Message != "Operation successful" {
		t.Fatalf("message=%q, want %q", resp.Message, "Operation successful")
	}
}

func TestReactivate(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/vps/servers/my-id/dormant", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method=%s, want PUT", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Fatalf("Content-Type=%s, want application/json", ct)
		}

		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode req: %v", err)
		}
		if req["dormant"] != false {
			t.Fatalf("dormant=%v, want false", req["dormant"])
		}
		if req["product"] != "VPSX1" {
			t.Fatalf("product=%v, want VPSX1", req["product"])
		}

		_, _ = w.Write([]byte(`{"message":"Operation successful"}`))
	})
	c, srv := newTestClient(t, mux)
	defer srv.Close()

	resp, err := c.VPS().Reactivate(testContext(), "my-id", "VPSX1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Message != "Operation successful" {
		t.Fatalf("message=%q, want %q", resp.Message, "Operation successful")
	}
}

func TestReactivate_MissingProduct(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/vps/servers/my-id/dormant", func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("no request should be made")
	})
	c, srv := newTestClient(t, mux)
	defer srv.Close()

	_, err := c.VPS().Reactivate(testContext(), "my-id", "  ")
	if err == nil || !strings.Contains(err.Error(), "product is required") {
		t.Fatalf("err=%v, want product required error", err)
	}
}

func TestMakeDormant_EmptyIdentifier(t *testing.T) {
	t.Parallel()
	c, _ := mythicbeasts.NewClient("", "")

	if _, err := c.VPS().MakeDormant(testContext(), "  "); !errors.Is(err, vpsapi.ErrEmptyIdentifier) {
		t.Fatalf("err=%v, want ErrEmptyIdentifier", err)
	}
}

func TestReactivate_EmptyIdentifier(t *testing.T) {
	t.Parallel()
	c, _ := mythicbeasts.NewClient("", "")

	if _, err := c.VPS().Reactivate(testContext(), "  ", "VPSX1"); !errors.Is(err, vpsapi.ErrEmptyIdentifier) {
		t.Fatalf("err=%v, want ErrEmptyIdentifier", err)
	}
}

func TestMakeDormant_UnexpectedStatus(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/vps/servers/my-id/dormant", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Server is already dormant"))
	})
	c, srv := newTestClient(t, mux)
	defer srv.Close()

	_, err := c.VPS().MakeDormant(testContext(), "my-id")
	if err == nil {
		t.Fatal("want error for 400 response")
	}

	want := "unexpected status 400: Server is already dormant"
	if err.Error() != want {
		t.Fatalf("err=%q want %q", err.Error(), want)
	}
}
