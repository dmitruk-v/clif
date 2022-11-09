package auth

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

type authenticator struct {
	cfgpath string
}

func NewAuthenticator() *authenticator {
	exepath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	cfgpath := path.Join(path.Dir(exepath), "auth.ini")
	return &authenticator{
		cfgpath: cfgpath,
	}
}

func (au *authenticator) signUp() error {
	fmt.Println("---------------------------------")
	fmt.Println("This is the first authentication.")
	fmt.Println("---------------------------------")
	fmt.Print("1. Please create your password: ")
	var password, passwordRep string
	sc := bufio.NewScanner(os.Stdin)
	if sc.Scan() {
		password = sc.Text()
	}
	fmt.Print("2. Repeat password: ")
	if sc.Scan() {
		passwordRep = sc.Text()
	}
	if sc.Err() != nil {
		return sc.Err()
	}
	if password != passwordRep {
		return errors.New("passwords are not equal")
	}
	fmt.Println("3. Generating hash...")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	fmt.Println("4. Creating config file...")
	if err := os.WriteFile(au.cfgpath, hash, 0400); err != nil {
		return err
	}
	fmt.Println("5. Done!")
	fmt.Println("---------------------------------")
	return nil
}

func (au *authenticator) SignIn() ([]byte, error) {
	formatError := func(err error) error {
		return fmt.Errorf("signing in: %v", err)
	}
	if _, err := os.Stat(au.cfgpath); err != nil {
		if err := au.signUp(); err != nil {
			return nil, formatError(err)
		}
	}
	hash, err := os.ReadFile(au.cfgpath)
	if err != nil {
		return nil, formatError(err)
	}
	fmt.Println("Please enter a password:")
	password, err := term.ReadPassword(0)
	if err != nil {
		return nil, formatError(err)
	}
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		return nil, formatError(errors.New("wrong password"))
	}
	return password, nil
}
