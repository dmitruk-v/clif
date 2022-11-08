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

func (au *authenticator) init() error {
	fmt.Println("---------------------------------")
	fmt.Println("This is the first authentication.")
	fmt.Println("---------------------------------")
	fmt.Print("1. Please create your password: ")
	sc := bufio.NewScanner(os.Stdin)
	var password, passwordRep string
	if sc.Scan() {
		password = sc.Text()
	}
	fmt.Print("2. Repeat password: ")
	if sc.Scan() {
		passwordRep = sc.Text()
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

func (au *authenticator) CheckPassword() ([]byte, error) {
	if _, err := os.Stat(au.cfgpath); err != nil {
		if err := au.init(); err != nil {
			return nil, err
		}
	}
	hash, err := os.ReadFile(au.cfgpath)
	if err != nil {
		return nil, err
	}
	fmt.Println("Please enter a password:")
	password, err := term.ReadPassword(0)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		return nil, errors.New("wrong password")
	}
	return password, nil
}
