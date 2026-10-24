package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

type testingStruct struct {
	ID uuid.UUID
}

var testUser = testingStruct{
	ID: uuid.New(),
}
var tokenSecret = "ThisSecretHere"

func TestCompare(t *testing.T) {
	password := "Plunge2"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatal("There was an error ", err)
	}
	itWorks, err := CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Fatal("There was an error: ", err)
	}
	if !itWorks {
		t.Errorf("expected True, got %v", itWorks)
	}
}

func TestCompareFalse(t *testing.T) {
	password := "rando23"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Got err: %v", err)
	}
	itDoesNotWork, err := CheckPasswordHash("Notright25", hashedPassword)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if itDoesNotWork == true {
		t.Errorf("It's not working right, password should not match hash")
	}
}

func TestExpiration(t *testing.T) {
	tokenString, err := MakeJWT(testUser.ID, tokenSecret, time.Duration(-time.Hour))
	if err != nil {
		print(err)
		t.Fatal("Wasn't expecting an error from makeJWT")
	}
	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Fatal("Was expecting an error because of an expired token")
	}
}

func TestUserIsUser(t *testing.T) {
	tokenString, err := MakeJWT(testUser.ID, tokenSecret, time.Duration(time.Hour))
	if err != nil {
		print(err)
		t.Fatal("Wasn't expecting an error from makeJWT")
	}
	userID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("Validate err: %v", err)
	}
	if userID != testUser.ID {
		t.Fatalf("The Id's don't match")
	}
}

func TestWrongSecret(t *testing.T) {
	tokenString, err := MakeJWT(testUser.ID, tokenSecret, time.Duration(time.Hour))
	if err != nil {
		print(err)
		t.Fatal("Wasn't expecting an error from makeJWT")
	}
	_, err = ValidateJWT(tokenString, "wrong-secret")
	if err == nil {
		t.Fatal("expecting err")
	}
}
