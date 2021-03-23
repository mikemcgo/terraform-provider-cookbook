package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// https://www.thegreatcodeadventure.com/mocking-http-requests-in-golang/
// Actually stole this from https://github.com/google/go-github/blob/760872962c7c542bb8aa07fdeb4aa5e3c0276512/github/github_test.go#L34
func setup() (c *Cookbook, mux *http.ServeMux, cleanup func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	c = NewCookbook(server.URL, server.Client())

	return c, mux, server.Close
}

func TestSave(t *testing.T) {
	cookbook, server, cleanup := setup()
	defer cleanup()

	expected_id := "abcd"

	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Not sure this test is good since the id isn't a uuid?
		fmt.Fprintf(w, `{"id": "%s"}`, expected_id)
	})

	r := Recipe{
		Ingredients: []string{"1/3 cup dirt"},
		Title:       "Test Recipe",
	}
	err := cookbook.save(&r)

	if err != nil {
		t.Errorf(`%v`, err)
	} else if r.Id != expected_id {
		t.Errorf(`id is "%v", not "%s"`, r.Id, expected_id)
	}

	server.HandleFunc("/"+expected_id, func(w http.ResponseWriter, r *http.Request) {
		// Not sure this test is good since the id isn't a uuid?
		fmt.Fprintf(w, `{"id": "%s"}`, expected_id)
	})

	r.Ingredients = []string{"2/3 cups of dirt", "garbage"}

	err = cookbook.save(&r)

	if err != nil {
		t.Errorf(`%v`, err)
	} else if r.Id != expected_id {
		t.Errorf(`id is "%v", not "%s"`, r.Id, expected_id)
	}
}

func TestGet(t *testing.T) {
	cookbook, server, cleanup := setup()
	defer cleanup()

	id := "abcd"
	server.HandleFunc("/"+id, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"":""}`)
	})

	cookbook.get(id)

	t.Error("Not a real test")
}

func TestRefresh(t *testing.T) {
	cookbook, _, _ := setup()
	cookbook.refresh()
	t.Error("Not a real test")
}

func TestDelete(t *testing.T) {
	cookbook, _, _ := setup()
	cookbook.delete(Recipe{})
	t.Error("Not a real test")
}
