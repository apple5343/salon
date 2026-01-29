package models

import "salon/internal/models"

type Passport struct {
	Series   string `json:"series"`
	Number   string `json:"number"`
	IssuedBy string `json:"issued_by"`
}

func PassportToService(p Passport) models.Passport {
	return models.Passport{
		Series:   p.Series,
		Number:   p.Number,
		IssuedBy: p.IssuedBy,
	}
}

func PassportToHttp(p models.Passport) Passport {
	return Passport{
		Series:   p.Series,
		Number:   p.Number,
		IssuedBy: p.IssuedBy,
	}
}
