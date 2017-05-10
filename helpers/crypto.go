// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT BY HAND BUT IDEALLY YOU SHOULDN'T

package helpers

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"strings"
)

// HashPasswordPBKDF2 returns a string (hash) or an error if
// something went wrong.
//
// Uses pbkdf2 with sha512 in compliance with NIST standards
func HashPasswordPBKDF2(password string, iterations uint32, saltLen uint16) (string, error) {

	// Generate salt
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// hash salt
	saltHash := base32.StdEncoding.EncodeToString(salt)

	// Use sha512 as the HMAC hash function.
	// sha256 uses 32b operations in it's hashing process, and can be accellerated
	// massively by GPU based attacks. sha512 uses 64b operations, which are less
	// well suited for GPU attacks, and thus more resistant to bruteforce.
	hashByes := pbkdf2.Key([]byte(password), salt, int(iterations), 32, sha512.New)

	// Base32 encode the resulting hash
	passHash := base32.StdEncoding.EncodeToString(hashByes)

	// Append salt to hash
	fullHash := fmt.Sprintf("%s:%s", passHash, saltHash)

	return fullHash, nil
}

// CheckPasswordPBKDF2 compares the hash to the attempted password.
//
// Returns true if the passwords match, and false otherwise and err if there was an error
func CheckPasswordPBKDF2(userHash, attemptedPassword string, iterations uint32) (bool, error) {

	// Get hash:salt
	encoded := strings.Split(userHash, ":")
	// sanity check
	if len(encoded) != 2 {
		return false, errors.New("userHash invalid! This password is unretrievable")
	}

	// Get password hash and salt
	passHash, err := base32.StdEncoding.DecodeString(encoded[0])
	if err != nil {
		return false, err
	}
	salt, err := base32.StdEncoding.DecodeString(encoded[1])
	if err != nil {
		return false, err
	}

	// hash attmpted password with same salt
	attemptedHash := pbkdf2.Key([]byte(attemptedPassword), salt, int(iterations), 32, sha512.New)

	// check
	if bytes.Compare(passHash, attemptedHash) != 0 {
		return false, nil
	}

	return true, nil
}

// HashPasswordBcrypt returns a string (hash) or an error if
// something went wrong.
//
// Uses bcrypt so that secrets (passwords) may later be strengthened
func HashPasswordBcrypt(password string, cost int) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CheckPasswordBcrypt compares the hash to the attempted password
func CheckPasswordBcrypt(userHash, attemptedPassword string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(userHash), []byte(attemptedPassword)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// UpgradeHashBcrypt checks to see if the users password hash meets current complexity
// standards and upgrades it if not.
//
// Returns true if the hash has been upgraded
func UpgradeHashBcrypt(currentHash, password string, cost uint16, allowHashDowngrades bool) (upgrade bool, new string) {
	curCost, err := bcrypt.Cost([]byte(currentHash))
	if err != nil {
		return
	}

	// if the current cost is too low
	// or if the current host is too high and downgrades are allowed
	// regenerate the hash
	if curCost < int(cost) || (curCost > int(cost) && allowHashDowngrades) {
		newHash, err := bcrypt.GenerateFromPassword([]byte(password), int(cost))
		if err != nil {
			return
		}
		return true, string(newHash)
	}

	return false, ""
}
