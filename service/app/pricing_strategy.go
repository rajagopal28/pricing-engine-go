package app

import (
	"log"

	"pricingengine"
	"pricingengine/service/model"
)


type Strategy struct{}

// Chained functional response that keeps the ball rolling with the
type StrartegyChain func(*pricingengine.GeneratePricingResponse) (*pricingengine.GeneratePricingResponse)

// GeneratePricing will calculate how much a 'risk' be priced or if they should
// be denied.
func (s *Strategy) ApplyBasePricing(input *pricingengine.GeneratePricingRequest, config *models.BaseRate, fn StrartegyChain) (*pricingengine.GeneratePricingResponse, error) {
  var result pricingengine.GeneratePricingResponse =  pricingengine.GeneratePricingResponse{}
  // for the current BaseRate and GeneratePricingRequest calculate outcome rate
  if fn != nil {
    log.Println("Found a chain function, Passing on the result for further computation")
    return fn(&result), nil
  }
  return &result, nil // return nil, errors.New("not implemented")
}
