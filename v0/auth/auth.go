package auth

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

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

func (au *authenticator) IsRegistered() bool {
	if _, err := os.Stat(au.cfgpath); err != nil {
		return false
	}
	return true
}

func (au *authenticator) SignIn() ([]byte, error) {
	formatError := func(err error) error {
		return fmt.Errorf("[ERROR]: sign in: %v", err)
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

func (au *authenticator) Register() ([]byte, error) {
	formatError := func(err error) error {
		return fmt.Errorf("[ERROR]: register: %v", err)
	}
	fmt.Println("---------------------------------")
	fmt.Println("This is the first authentication.")
	fmt.Println("---------------------------------")
	fmt.Println("Please create your password...")
	password, err := au.createPassword()
	if err != nil {
		return nil, formatError(err)
	}
	fmt.Println("Generating hash...")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, formatError(err)
	}
	fmt.Println("Creating config file...")
	if err := os.WriteFile(au.cfgpath, hash, 0400); err != nil {
		return nil, formatError(err)
	}
	fmt.Println("Done!")
	fmt.Println("---------------------------------")
	return []byte(password), nil
}

func (au *authenticator) Unregister() error {
	return os.Remove(au.cfgpath)
}

func (au *authenticator) createPassword() (string, error) {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("Password: ")
	password, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	fmt.Print("Repeat: ")
	passwordRep, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	if password != passwordRep {
		return "", errors.New("passwords are not equal")
	}
	return strings.TrimSpace(password), nil
}
