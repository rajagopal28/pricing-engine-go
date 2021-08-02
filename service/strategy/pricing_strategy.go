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

// ApplyBasePricing method will compute the base fare based on the RangeConfig config
// If there is a subsecuent decision to be made, the corresponding chain will be invoked with the generated PricingItem data
// returns the computed PricingItem or error if any happened during the computation
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

// ApplySubsecuentFactorsToPricing method will generically compute the next fare based on the Consecutive RangeConfig
// If there is a subsecuent decision to be made, the corresponding chain will be invoked with the generated PricingItem data
// returns the computed PricingItem or error if any happened during the computation
func (s *Strategy) ApplySubsecuentFactorsToPricing(input *pricingengine.GeneratePricingRequest, previousPricingItem *pricingengine.PricingItem, config *models.RangeConfig, fn StrartegyChain) (*pricingengine.PricingItem, error) {
  var result pricingengine.PricingItem =  pricingengine.PricingItem{}
  result.Premium = math.Floor(previousPricingItem.Premium * config.Value * 1000)/1000
  result.Currency = previousPricingItem.Currency
  result.FareGroup = previousPricingItem.FareGroup + ", " + config.Label
  if fn != nil {
    log.Println("Found a chain function, Passing on the result for further computation")
    return fn(&result)
  }
  return &result, nil
}

// FindMatchingDriverAgeFactor method will find the appropriate DriverAgeFactor RangeConfig
// based on the DateOfBirth data passed in the input GeneratePricingRequest
// returns the found DriverAgeFactor
//  error will be thrown if the field level validation fails or a matching config is not found
func (s *Strategy) FindMatchingDriverAgeFactor(input *pricingengine.GeneratePricingRequest, allDriverAgeFactors []models.RangeConfig) (*models.RangeConfig, error) {
  date_of_birth := input.DateOfBirth
  parse_dob_t, err := time.Parse("2006-01-02", date_of_birth)
	if err != nil {
		return nil, errors.New("Error wile Parsing DateOfBirth date. Error: "+ err.Error())
	}
  now := time.Now()
  age := int(now.Sub(parse_dob_t).Hours()/(24*30*12))
	log.Println("Checking the driver factor for date_of_birth=", date_of_birth, " parse_dob_t=", parse_dob_t, " age=", age)
  for i:= 0; i < len(allDriverAgeFactors); i++ {
		current := allDriverAgeFactors[i]
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

// FindMatchingInsuranceGroupFactor method will find the appropriate InsuranceGroupFactor RangeConfig
// based on the InsuranceGroup data passed in the input GeneratePricingRequest
// returns the found InsuranceGroupFactor
//  error will be thrown if the field level validation fails or a matching config is not found
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

// FindMatchingLicenceValidityFactor method will find the appropriate LicenceValidityFactor RangeConfig
// based on the LicenseHeldSince data passed in the input GeneratePricingRequest
// returns the found LicenceValidityFactor
//  error will be thrown if the field level validation fails or a matching config is not found
func (s *Strategy) FindMatchingLicenceValidityFactor(input *pricingengine.GeneratePricingRequest, allLicenceValidtyFactors []models.RangeConfig) (*models.RangeConfig, error) {
	licence_date := input.LicenseHeldSince
  parse_date_t, err := time.Parse("2006-01-02", licence_date)
	if err != nil {
		return nil, errors.New("Error wile Parsing LicenseHeldSince date. Error: "+ err.Error())
	}
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
