package main

import (
	"encoding/json"
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

// This test doesn't do much, just verifies json decode logic
func TestGet(t *testing.T) {
	cookbook, server, cleanup := setup()
	defer cleanup()

	id := "abcd"
	title := "Unspeakable things"
	ingredients := []string{"1 bag of garbage", "2 bags of garbage"}

	recip := new(Recipe)
	recip.Id = id
	recip.Title = title
	recip.Ingredients = ingredients

	server.HandleFunc("/"+id, func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		err := enc.Encode(recip)
		t.Logf("Rendering Response: %s", recip)

		if err != nil {
			t.Fatal("Unable to encode json")
		}
	})

	recipe, err := cookbook.get(id)
	if err != nil {
		t.Error(err)
	}

	if recipe.Id != id || recipe.Title != title {
		t.Error("Id or Title did not expected from server")
	}

	for i, v := range recipe.Ingredients {
		if v != ingredients[i] {
			t.Errorf("Element %d from server does match expected: %s actual: %s", i, ingredients[i], v)
		}
	}
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
