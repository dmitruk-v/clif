package controllers

import "fmt"

type cliBlaController struct{}

func NewCliBlaController() *cliBlaController {
	return &cliBlaController{}
}

func (ctrl *cliBlaController) Handle(req map[string]string) error {
	fmt.Println("handling bla command:", req)
	return nil
}
