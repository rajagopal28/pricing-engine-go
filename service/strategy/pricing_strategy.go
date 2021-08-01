package strategy

import (
	"log"
  "time"
  "errors"
	"math"

	"pricingengine"
	"pricingengine/service/model"
)


type Strategy struct{}

// Chained functional response that keeps the ball rolling with the
type StrartegyChain func(*pricingengine.PricingItem) (*pricingengine.PricingItem, error)

// GeneratePricing will calculate how much a 'risk' be priced or if they should
// be denied.
func (s *Strategy) ApplyBasePricing(input *pricingengine.GeneratePricingRequest, config *models.RangeConfig, fn StrartegyChain) (*pricingengine.PricingItem, error) {
  var result pricingengine.PricingItem =  pricingengine.PricingItem{}
  // for the current BaseRate and GeneratePricingRequest calculate outcome rate
  // just check the base price and
  result.Premium = config.Value
  result.Currency = "£"
  result.FareGroup = config.Label
  if fn != nil {
    log.Println("Found a chain function, Passing on the result for further computation")
    return fn(&result)
  }
  return &result, nil // return nil, errors.New("not implemented")
}

func (s *Strategy) ApplySubsecuentFactorsToPricing(input *pricingengine.GeneratePricingRequest, previousPricingItem *pricingengine.PricingItem, config *models.RangeConfig, fn StrartegyChain) (*pricingengine.PricingItem, error) {
  var result pricingengine.PricingItem =  pricingengine.PricingItem{}
  // for the current BaseRate and GeneratePricingRequest calculate outcome rate
  // just check the base price and
  result.Premium = math.Floor(previousPricingItem.Premium * config.Value * 1000)/1000
  result.Currency = previousPricingItem.Currency
  result.FareGroup = previousPricingItem.FareGroup + ", " + config.Label
  if fn != nil {
    log.Println("Found a chain function, Passing on the result for further computation")
    return fn(&result)
  }
  return &result, nil // return nil, errors.New("not implemented")
}

func (s *Strategy) FindMatchingDriverAgeFactor(input *pricingengine.GeneratePricingRequest, allDriverAgeFactords []models.RangeConfig) (*models.RangeConfig, error) {
  date_of_birth := input.DateOfBirth
  parse_dob_t,_ := time.Parse("2006-01-02", date_of_birth)
  now := time.Now()
  age := int(now.Sub(parse_dob_t).Hours()/(24*30*12))
	log.Println("Checking the driver factor for date_of_birth=", date_of_birth, " parse_dob_t=", parse_dob_t, " age=", age)
  for i:= 0; i < len(allDriverAgeFactords); i++ {
		current := allDriverAgeFactords[i]
    if (current.Start < age && current.End >= age) {
      if (current.IsEligible) {
        return &current, nil
      } else {
        return &current, errors.New("Declined due to :"+current.Label)
      }
    }
  }
  return nil, errors.New("MatchingDriverAgeFactor not found!")
}

func (s *Strategy) FindMatchingInsuranceGroupFactor(input *pricingengine.GeneratePricingRequest, allInsuranceGroupFactors []models.RangeConfig) (*models.RangeConfig, error) {
  insurance_group := input.InsuranceGroup
  for i:= 0; i < len(allInsuranceGroupFactors); i++ {
    current := allInsuranceGroupFactors[i]
    if (current.Start < insurance_group && current.End >= insurance_group) {
      if (current.IsEligible) {
        return &current, nil
      } else {
        return &current, errors.New("Declined due to :"+current.Label)
      }
    }
  }
  return nil, errors.New("MatchingInsuranceGroupFactor not found!")
}


func (s *Strategy) FindMatchingLicenceValidityFactor(input *pricingengine.GeneratePricingRequest, allLicenceValidtyFactors []models.RangeConfig) (*models.RangeConfig, error) {
	date_of_birth := input.LicenseHeldSince
  parse_date_t,_ := time.Parse("2006-01-02", date_of_birth)
  now := time.Now()
  licence_length := int(now.Sub(parse_date_t).Hours()/(24*30*12))
  for i:= 0; i < len(allLicenceValidtyFactors); i++ {
    current := allLicenceValidtyFactors[i]
    if (current.Start < licence_length && current.End >= licence_length) {
      if (current.IsEligible) {
        return &current, nil
      } else {
        return &current, errors.New("Declined due to :"+current.Label)
      }
    }
  }
  return nil, errors.New("MatchingLicenceValidityFactor not found!")
}
