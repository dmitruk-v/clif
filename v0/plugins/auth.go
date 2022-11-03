package plugins

import (
	"errors"
	"fmt"

	"github.com/dmitruk-v/clim/v0"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

type authPlugin struct {
	PwdHash []byte
}

func NewAuthPlugin(hash []byte) *authPlugin {
	return &authPlugin{
		PwdHash: hash,
	}
}

func (plug *authPlugin) Execute(app *clim.App) error {
	fmt.Println("Please enter a password:")
	fmt.Print("> ")
	password, err := term.ReadPassword(0)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword(plug.PwdHash, password); err != nil {
		return errors.New("wrong password")
	}
	return nil
}
