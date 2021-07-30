package app

import (
	"context"
	"log"

	"pricingengine"
	"pricingengine/service/config"
)

type App struct{
	Cache config.ConfigCache
}


// GeneratePricing will calculate how much a 'risk' be priced or if they should
// be denied.
func (a *App) GeneratePricing(ctx context.Context, request *pricingengine.GeneratePricingRequest) (*pricingengine.GeneratePricingResponse, error) {
	log.Println("Entering GeneratePricing")
	a.Cache.Fetcher = config.ConfigFetcher{Path: "/config/"} // Initialise with actual path
	a.Cache.InitialiseWithRefresh(false, 100000) // time to re-Initialise and see if the cache works

	result := pricingengine.GeneratePricingResponse{}
	result.Input = *request

	var strategies = Strategy{}
	driver_factor_range, err := strategies.FindMatchingDriverAgeFactor(request, a.Cache.DriverAgeFactorList)
	if(err != nil) {
		log.Fatalf("error finding driver_factor_range: %v", err)
		result.Message = err.Error()
		result.IsEligible = false
		return &result, nil
	}

	insurance_factor_range, err := strategies.FindMatchingInsuranceGroupFactor(request, a.Cache.InsuranceGroupFactorList)
	if(err != nil) {
		log.Fatalf("error finding insurance_factor_range: %v", err)
		result.Message = err.Error()
		result.IsEligible = false
		return &result, nil
	}

	licence_factor_range, err := strategies.FindMatchingLicenceValidityFactor(request, a.Cache.LicenceValidityFactorList)
	if(err != nil) {
		log.Fatalf("error finding licence_factor_range: %v", err)
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
				log.Fatalf("error finding ApplyBasePricing: %v", err)
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
