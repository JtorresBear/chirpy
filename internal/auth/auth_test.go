package auth

import "testing"

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
