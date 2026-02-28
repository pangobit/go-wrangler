package main

type User struct {
	Name          string   `bind:"header,required"`
	Email         string   `bind:"query"`
	Password      string   `bind:"form,required"`
	Age           int      `validate:"min=18,max=120"`
	ID            string   `bind:"path,required" validate:"max=10"`
	Roles         []string `bind:"form=role"`
	PermissionIDs []int    `bind:"form=permission"`
}
