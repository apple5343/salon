package models

type Passport struct {
	Series   string `validate:"required,numeric,len=4"`
	Number   string `validate:"required,numeric,len=6"`
	IssuedBy string `validate:"required"`
}
