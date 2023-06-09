package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/betelgeusexru/golang-hotel-reservation/db"
	"github.com/betelgeusexru/golang-hotel-reservation/testutils"
	"github.com/betelgeusexru/golang-hotel-reservation/types"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email: "james@foo.com",
		LastName: "foo",
		FirstName: "james",
		Password: "supersecurepassword",
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = userStore.InsertUser(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := testutils.SetupDatabase(t)
	defer tdb.Teardown(t)

	insertedUser := insertTestUser(t, tdb)

	app := testutils.SetupFiberApp()

	authHandler := NewAuthHandler(tdb.UserStore)

	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email: "james@foo.com",
		Password: "supersecurepassword",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status of 200 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}


	if authResp.Token == "" {
		t.Fatalf("expected the JWT token to be present in the auth response")
	}

	insertedUser.EncryptedPassword = ""

	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("expected the user to be the inserted user")
	} 
}

func TestAuthenticateWithWrongPasswordFailure(t *testing.T) {
	tdb := testutils.SetupDatabase(t)
	defer tdb.Teardown(t)

	insertTestUser(t, tdb)

	app := testutils.SetupFiberApp()

	authHandler := NewAuthHandler(tdb.UserStore)

	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email: "james@foo.com",
		Password: "notcorrect",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status of 400 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if authResp.Token != "" {
		t.Fatalf("expected the JWT token to be not present in the auth response")
	}
}