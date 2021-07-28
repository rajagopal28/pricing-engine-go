package models
// import "encoding/json"
// import our encoding/json package


type BaseRate struct {
	Label string `json:"label"`
  Time int `json:"time"`
  Rate float64 `json:"rate"`
}

type DriverAgeFactor struct {
	Label string `json:"label"`
  Age int `json:"age"`
  IsEligible bool `json:"is-eligible"`
  Factor float64 `json:"factor"`
}

type InsuranceGroupFactor struct {
  Label string `json:"label"`
  Group int `json:"group"`
  IsEligible bool `json:"is-eligible"`
  factor float64 `json:"factor"`
}

type LicenceValidityFactor struct {
  Length int `json:"length"`
  Factor float64 `json:"factor"`
}

type RangeConfig struct {
  Start int
  End int
	IsEligible bool
	Value float64
	Label string
}
