package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Cookbook struct {
	Recipes map[string]Recipe
	client  *http.Client
	baseURL string
}

func NewCookbook(baseURL string, client *http.Client) (c *Cookbook) {
	if client == nil {
		client = &http.Client{}
	}
	recipes := make(map[string]Recipe)

	c = &Cookbook{
		Recipes: recipes,
		client:  client,
		baseURL: baseURL,
	}
	return c
}

type Recipe struct {
	// Id might prove difficult as its a uuid
	// also this should be private
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Steps       []string `json:"steps"`
	Ingredients []string `json:"ingredients"`
	Feedback    string   `json:"feedback"`
}

// Whew I don't know how i feel about just returning the id?
type recipeID struct {
	Id string `json:"id"`
}

type NoRecipeFound struct {
	Id string
}

func (e *NoRecipeFound) Error() string {
	return fmt.Sprintf("No recipe with id %s found", e.Id)
}

func (c *Cookbook) get(id string) (*Recipe, error) {
	res, err := c.client.Get(c.baseURL + "/" + id)

	if res.StatusCode != 200 {
		return nil, &NoRecipeFound{id}
	}

	if err != nil {
		return nil, err
	}

	var recipe Recipe
	err = json.NewDecoder(res.Body).Decode(&recipe)

	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

func (c *Cookbook) refresh() error {
	res, err := c.client.Get(c.baseURL)

	if err != nil {
		return errors.New("cookbook: Unable to refresh recipes")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var recipes []Recipe
	err = json.Unmarshal(body, &recipes)

	if err != nil {
		return err
	}

	for _, recipe := range recipes {
		c.Recipes[recipe.Id] = recipe
	}

	return nil
}

// This function is horrid?
func (c *Cookbook) save(recipe *Recipe) error {
	method := "POST"
	dest := c.baseURL
	if recipe.Id != "" {
		dest = dest + "/" + recipe.Id
		method = "PUT"
	}

	b, err := json.Marshal(recipe)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, dest, bytes.NewReader(b))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.client.Do(req)

	if err != nil {
		return err
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return err
	}

	var rid recipeID
	err = json.Unmarshal(body, &rid)

	if err != nil {
		return err
	}

	recipe.Id = rid.Id

	return nil
}

func (c *Cookbook) delete(recipe Recipe) (string, error) {
	return "", nil
}

// class RecipeSchema(Schema):
//     id = UUIDString(missing=str(uuid.uuid4()))
//     title = fields.Str(required=True)
//     steps = fields.List(fields.Str)
//     feedback = fields.Str()
//     ingredients = fields.List(fields.Str, required=True)
