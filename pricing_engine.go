package pricingengine

// GeneratePricingRequest is used for generate pricing requests, it holds the
// inputs that are used to provide pricing for a given user.
type GeneratePricingRequest struct {
  DateOfBirth string `json:"date_of_birth"`
  InsuranceGroup int `json:"insurance_group"`
  LicenseHeldSince string `json:"license_held_since"`
}

// GeneratePricingResponse - contains the list of all pricing generated for the request passed
// it typically has the input based on which the decision is taken
// IsEligible to indicate whether the user is eligible
// Message to state the reason for thich the Decline has happened
type GeneratePricingResponse struct {
	input GeneratePricingRequest
  IsEligible string `json:"is_eligible"`
  Message string `json:"message"`
  pricing_list []PricingItem `json:"pricing"`
}

// PricingItem - contains the pricing data generated for partucular group based on the request passed
type PricingItem struct {
	Premium float64 `json:"premium"`
  Currency string  `json:"currency"`
  FareGroup string `json:"fare_group"`
}
