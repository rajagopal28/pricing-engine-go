package app

import (
	"context"
	"log"

	"pricingengine"
	"pricingengine/service/strategy"
	"pricingengine/service/config"
)

type App struct{
	Cache config.ConfigCache
}


// GeneratePricing will calculate how much a 'risk' be priced or if they should
// be denied.
// GeneratePricing method simply takes the input and peforms input validation/sanitization
// The pricing value is generated for all the base fare ranges available considering the user input
// Applies chain of command pattern to strategies that are to be executed based on the configs that are available
// Inputs ==> ctx context.Context, request *pricingengine.GeneratePricingRequest
// returns ==> *pricingengine.GeneratePricingResponse, error
func (a *App) GeneratePricing(ctx context.Context, request *pricingengine.GeneratePricingRequest) (*pricingengine.GeneratePricingResponse, error) {
	log.Println("Entering GeneratePricing")
	// Initialise with actual path if not present
	if a.Cache.TimeToLive == 0 {
		a.Cache.Fetcher = config.ConfigFetcher{Path: "/config/"}
	}
	a.Cache.InitialiseWithRefresh(false, 100000) // time to re-Initialise and see if the cache works

	result := pricingengine.GeneratePricingResponse{}
	result.Input = *request

	if len(request.DateOfBirth) == 0 {
		result.Message = "DateOfBirth cannot be empty"
		result.IsEligible = false
		return &result, nil
	}

	if request.InsuranceGroup <= 0 {
		result.Message = "InsuranceGroup should be a Positive number"
		result.IsEligible = false
		return &result, nil
	}

	if len(request.LicenseHeldSince) == 0 {
		result.Message = "LicenseHeldSince Date cannot be empty"
		result.IsEligible = false
		return &result, nil
	}

	var strategies = strategy.Strategy{}
	driver_factor_range, err := strategies.FindMatchingDriverAgeFactor(request, a.Cache.DriverAgeFactorList)
	if(err != nil) {
		log.Printf("error finding driver_factor_range: %v", err)
		result.Message = err.Error()
		result.IsEligible = false
		return &result, nil
	}

	insurance_factor_range, err := strategies.FindMatchingInsuranceGroupFactor(request, a.Cache.InsuranceGroupFactorList)
	if(err != nil) {
		log.Printf("error finding insurance_factor_range: %v", err)
		result.Message = err.Error()
		result.IsEligible = false
		return &result, nil
	}

	licence_factor_range, err := strategies.FindMatchingLicenceValidityFactor(request, a.Cache.LicenceValidityFactorList)
	if(err != nil) {
		log.Printf("error finding licence_factor_range: %v", err)
		result.Message = err.Error()
		result.IsEligible = false
		return &result, nil
	}

	var thirdStrategy = func( resp *pricingengine.PricingItem) (*pricingengine.PricingItem, error) {
		log.Println("Processing LicenceValidityFactor Stragegy")
		return strategies.ApplySubsecuentFactorsToPricing(request, resp, licence_factor_range, nil) // final call so no next strategy
	}

	var secondStrategy = func( resp *pricingengine.PricingItem) (*pricingengine.PricingItem, error) {
		log.Println("Got Previous Strategy in InsuranceGroupFactor Stragegy data here::", resp)
		return strategies.ApplySubsecuentFactorsToPricing(request, resp, insurance_factor_range, thirdStrategy)
	}

	var firstStrategy = func( resp *pricingengine.PricingItem) (*pricingengine.PricingItem, error) {
		log.Println("Got Previous Strategy in DriverAgeFactor Stragegy data here::", resp)
		return strategies.ApplySubsecuentFactorsToPricing(request, resp, driver_factor_range, secondStrategy)
	}

	price_items := []pricingengine.PricingItem{}
	for i:= 0; i < len(a.Cache.BaseRateList); i++ {
			item, err := strategies.ApplyBasePricing(request, &a.Cache.BaseRateList[i], firstStrategy)
			if(err != nil) {
				log.Printf("error finding ApplyBasePricing: %v", err)
				return &result, err
			}
			price_items = append(price_items, *item)
	}
	result.Message = "Success"
	result.IsEligible = true
	result.PricingList = price_items
	log.Println("Leaving GeneratePricing")
	return &result, nil
}

// GeneratePricingConfig fetch and cache the configs related to pricing computations
// Just forms a map[]{} based on the config in the cache
func (a *App) GeneratePricingConfig(ctx context.Context) (interface{}, error) {
	log.Println("Entering GeneratePricingConfig")
	// Initialise with actual path if not present
	if a.Cache.TimeToLive == 0 {
		a.Cache.Fetcher = config.ConfigFetcher{Path: "/config/"}
	}
	a.Cache.InitialiseWithRefresh(false, 100000) // time to live 100000s
	var result map[string]interface{} = make(map[string]interface{})

	result["base-rate"] = a.Cache.BaseRateList
	result["driver-age-factor"] = a.Cache.DriverAgeFactorList
	result["insurance-group-factor"] = a.Cache.InsuranceGroupFactorList
	result["licence-validity-factor"] = a.Cache.LicenceValidityFactorList
	log.Println("Leaving GeneratePricingConfig")
	return result, nil
}
