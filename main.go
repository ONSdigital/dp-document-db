package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/ONSdigital/dp-recipe-api/recipe"
	log "github.com/daiLlew/funkylog"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type Identity struct {
	ID          string   `json:"id,omitempty"`
	Identifier  string   `json:"identifier,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

func main() {
	log.Customise(log.Configuration{
		Namespace: "doc-db-demo",
		TimeFmt:   time.RFC3339,
		InfoStyle: log.NewStyle(color.FgHiMagenta, ":unicorn_face:"),
		WarnStyle: log.NewStyle(color.FgHiYellow, ":warning:"),
		ErrStyle:  log.NewStyle(color.FgHiRed, ":rotating_light:"),
	})

	if err := run(); err != nil {
		log.Err("application error: %+v :rotating_light:\n", err)
		os.Exit(1)
	}
}

func run() error {
	root := &cobra.Command{
		Use:   "poc",
		Short: "TODO",
	}

	root.AddCommand(poc())

	return root.Execute()
}

func poc() *cobra.Command {
	return &cobra.Command{
		Use:   "post-recipe",
		Short: "post a the test recipe",
		RunE: func(cmd *cobra.Command, args []string) error {
			identity, err := getIdentity()
			if err != nil {
				return err
			}

			r, err := createRecipe(identity)
			if err != nil {
				return err
			}

			log.Info("create recipe completed successfully : ID: %s Alias %s\n", r.ID, r.Alias)

			r2, err := getRecipe(r.ID, identity)
			if err != nil {
				return err
			}

			recipeJson, err := json.MarshalIndent(r2, "", "  ")
			if err != nil {
				return err
			}

			log.Info("recipe response json:\n\n%s\n", string(recipeJson))

			log.Info("demo complete :rocket: :tada:\n")
			return nil
		},
	}
}

func createRecipe(identity *Identity) (*recipe.Response, error) {
	log.Info("posting new recipe")

	recipeBody, err := getRecipeBody()
	if err != nil {
		return nil, err
	}

	docDBEndpoint := os.Getenv("DOC_DB_POC_IP")
	if docDBEndpoint == "" {
		return nil, fmt.Errorf("env var %q expected but not found", "DOC_DB_POC_IP")
	}

	newRecipeEndpoint := fmt.Sprintf("http://%s:22300/recipes", docDBEndpoint)
	respBody, status, err := execRequest(http.MethodPost, newRecipeEndpoint, identity.ID, recipeBody)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("incorrect http status for post recipie expected %d but was %d", http.StatusOK, status)
	}

	log.Info("post recipe response status OK")

	var r recipe.Response
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func getRecipeBody() (*bytes.Buffer, error) {
	recipeBytes, err := ioutil.ReadFile("example-recipe.json")
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(recipeBytes), nil
}

func getIdentity() (*Identity, error) {
	identityBytes, err := ioutil.ReadFile("poc/bin/identity_data.json")
	if err != nil {
		return nil, err
	}

	var identities map[string]Identity
	err = json.Unmarshal(identityBytes, &identities)
	if err != nil {
		return nil, err
	}

	var identity Identity
	for _, item := range identities {
		identity = item
		break
	}

	return &identity, nil
}

func getRecipe(id string, identity *Identity) (*recipe.Response, error) {
	log.Info("retrieving recipe from API")

	docDBEndpoint := os.Getenv("DOC_DB_POC_IP")
	if docDBEndpoint == "" {
		return nil, fmt.Errorf("env var %q expected but not found", "DOC_DB_POC_IP")
	}

	getRecipeEndpoint := fmt.Sprintf("http://%s:22300/recipes/%s", docDBEndpoint, id)
	respBody, status, err := execRequest(http.MethodGet, getRecipeEndpoint, identity.ID, nil)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("incorrect http status for get recipie expected %d but was %d", http.StatusOK, status)
	}

	log.Info("get recipe response status OK")

	var r recipe.Response
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return nil, err
	}

	log.Info("get recipe completed successfully : ID: %s Alias %s\n", r.ID, r.Alias)
	return &r, nil
}

func execRequest(method, url, token string, reqBody io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	cli := http.Client{Timeout: time.Second * 5}

	log.Info("executing request to Recipe API")
	resp, err := cli.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
}
