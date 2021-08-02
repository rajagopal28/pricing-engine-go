package strategy

import (
  "testing"
  "time"

  "pricingengine"
	"pricingengine/service/model"
	"pricingengine/service/strategy"
  "pricingengine/test/util"
)


func TestPricingStrategyVariedConfigs(tp *testing.T){
  var strategies = strategy.Strategy{}
  request := pricingengine.GeneratePricingRequest{
    DateOfBirth: "2001-01-02",
    InsuranceGroup: 7,
    LicenseHeldSince: "2006-01-02",
  }
  tp.Run("TestPricingStrategyVariedConfigsToApplyBaseStrategyWithChainStrategy", func(t *testing.T) {
    strategyExecuted := false
  	var firstStrategy = func(resp *pricingengine.PricingItem) (*pricingengine.PricingItem, error) {
  		strategyExecuted = true
      util.AssertTrue(resp != nil, t)
      util.AssertTrue(resp.Premium != 0, t)
      util.AssertEqual(resp.Premium, 249.99, t)
      util.AssertEqual(resp.Currency, "£", t)
      util.AssertEqual(resp.FareGroup, "Base Fare Range", t)
      util.AssertTrue(len(resp.FareGroup) > 0, t)
  		return resp, nil
  	}
    baseRate := models.RangeConfig{
      Start: 0,
      End: 1200,
    	IsEligible: true,
    	Value: 249.99,
      Label: "Base Fare Range",
    }
    resp, err := strategies.ApplyBasePricing(&request, &baseRate, firstStrategy)

    util.AssertTrue(strategyExecuted, t)
    util.AssertTrue(err == nil, t)
    util.AssertTrue(resp.Premium != 0, t)
    util.AssertEqual(resp.Currency, "£", t)
    util.AssertEqual(resp.FareGroup, "Base Fare Range", t)
    util.AssertTrue(len(resp.FareGroup) > 0, t)
    // util.AssertEqual(resp.Message, "DateOfBirth cannot be empty", t)
  })
  tp.Run("TestPricingStrategyVariedConfigsToApplyBaseStrategyWithoutChainStrategy", func(t *testing.T) {
    baseRate := models.RangeConfig{
      Start: 0,
      End: 1200,
    	IsEligible: true,
    	Value: 249.99,
      Label: "Base Fare Range",
    }
    resp, err := strategies.ApplyBasePricing(&request, &baseRate, nil)

    util.AssertTrue(err == nil, t)
    util.AssertEqual(resp.Premium, 249.99, t)
    util.AssertEqual(resp.Currency, "£", t)
    util.AssertEqual(resp.FareGroup, "Base Fare Range", t)
    util.AssertTrue(len(resp.FareGroup) > 0, t)
  })
  tp.Run("TestPricingStrategyVariedConfigsToApplySubsecuentFactorsToPricingWithChainStrategy", func(t *testing.T) {
      baseRate := models.RangeConfig{
        Start: 0,
        End: 1200,
      	IsEligible: true,
      	Value: 149.99,
        Label: "Base Fare Range",
      }

      secondRate := models.RangeConfig{
        Start: 1201,
        End: 1000,
      	IsEligible: true,
      	Value: 1.1,
        Label: "Secondary Fare Range",
      }

      firstStrategyExecuted := false
      secondStrategyExecuted := false
      var secondStrategy = func( resp *pricingengine.PricingItem) (*pricingengine.PricingItem, error) {
        secondStrategyExecuted = true
        return resp, nil
      }
    	var firstStrategy = func(resp *pricingengine.PricingItem) (*pricingengine.PricingItem, error) {
    		firstStrategyExecuted = true
        util.AssertTrue(resp != nil, t)
        util.AssertTrue(resp.Premium != 0, t)
        util.AssertEqual(resp.Premium, 149.99, t)
        util.AssertEqual(resp.Currency, "£", t)
        util.AssertEqual(resp.FareGroup, "Base Fare Range", t)
        util.AssertTrue(len(resp.FareGroup) > 0, t)
    		return strategies.ApplySubsecuentFactorsToPricing(&request, resp, &secondRate, secondStrategy)
    	}

      resp, err := strategies.ApplyBasePricing(&request, &baseRate, firstStrategy)

      util.AssertTrue(firstStrategyExecuted, t)
      util.AssertTrue(secondStrategyExecuted, t)
      util.AssertTrue(err == nil, t)
      util.AssertTrue(resp.Premium != 0, t)
      util.AssertEqual(resp.Premium, 164.989, t)
      util.AssertEqual(resp.Currency, "£", t)
      util.AssertEqual(resp.FareGroup, "Base Fare Range, Secondary Fare Range", t)
      util.AssertTrue(len(resp.FareGroup) > 0, t)
  })
  tp.Run("TestPricingStrategyVariedConfigsToApplySubsecuentFactorsToPricingWithoutChainStrategy", func(t *testing.T) {
    baseRate := models.RangeConfig{
      Start: 0,
      End: 1200,
      IsEligible: true,
      Value: 149.99,
      Label: "Base Fare Range",
    }

    secondRate := models.RangeConfig{
      Start: 1201,
      End: 1300,
      IsEligible: true,
      Value: 1.1,
      Label: "Secondary Fare Range",
    }

    firstStrategyExecuted := false
    var firstStrategy = func(resp *pricingengine.PricingItem) (*pricingengine.PricingItem, error) {
      firstStrategyExecuted = true
      util.AssertTrue(resp != nil, t)
      util.AssertTrue(resp.Premium != 0, t)
      util.AssertEqual(resp.Premium, 149.99, t)
      util.AssertEqual(resp.Currency, "£", t)
      util.AssertEqual(resp.FareGroup, "Base Fare Range", t)
      util.AssertTrue(len(resp.FareGroup) > 0, t)
      return strategies.ApplySubsecuentFactorsToPricing(&request, resp, &secondRate, nil)
    }

    resp, err := strategies.ApplyBasePricing(&request, &baseRate, firstStrategy)

    util.AssertTrue(firstStrategyExecuted, t)
    util.AssertTrue(err == nil, t)
    util.AssertTrue(resp.Premium != 0, t)
    util.AssertEqual(resp.Premium, 164.989, t)
    util.AssertEqual(resp.Currency, "£", t)
    util.AssertEqual(resp.FareGroup, "Base Fare Range, Secondary Fare Range", t)
    util.AssertTrue(len(resp.FareGroup) > 0, t)
  })
  tp.Run("TestPricingStrategyVariedConfigsToFindMatchingLicenceValidityFactor-Success-Senario", func(t *testing.T) {
    now := time.Now()
    _t := now.AddDate(-3, 0, 0) //3 years before now
    licence := _t.Format("2006-01-02")
    request.LicenseHeldSince = licence
    LicenceValidityFactorList := []models.RangeConfig{
      models.RangeConfig{
        Start: 0,
        End: 2,
        IsEligible: true,
        Value: 3.0,
        Label: "<2 years",
      },
      models.RangeConfig{
        Start: 2,
        End: 4,
        IsEligible: true,
        Value: 2.0,
        Label: "2-4",
      },
    }

    licence_factor_range, err := strategies.FindMatchingLicenceValidityFactor(&request, LicenceValidityFactorList)
  	util.AssertTrue(err == nil, t)
    util.AssertTrue(licence_factor_range != nil, t)
    util.AssertEqual(*licence_factor_range, LicenceValidityFactorList[1], t)
  })
  tp.Run("TestPricingStrategyVariedConfigsToFindMatchingLicenceValidityFactor-NotFound-Senario", func(t *testing.T) {
    now := time.Now()
    _t := now.AddDate(-3, 0, 0) //3 years before now
    licence := _t.Format("2006-01-02")
    request.LicenseHeldSince = licence
    LicenceValidityFactorList := []models.RangeConfig{
      models.RangeConfig{
        Start: 0,
        End: 2,
        IsEligible: true,
        Value: 3.0,
        Label: "<2 years",
      },
    }

    licence_factor_range, err := strategies.FindMatchingLicenceValidityFactor(&request, LicenceValidityFactorList)
  	util.AssertTrue(err != nil, t)
    util.AssertEqual(err.Error(), "MatchingLicenceValidityFactor not found!", t)
    util.AssertTrue(licence_factor_range == nil, t)
  })

  tp.Run("TestPricingStrategyVariedConfigsToFindMatchingLicenceValidityFactor-InvalidDate-Senario", func(t *testing.T) {
    licence := "22-Jan-2020"
    request.LicenseHeldSince = licence
    LicenceValidityFactorList := []models.RangeConfig{
      models.RangeConfig{
        Start: 0,
        End: 2,
        IsEligible: true,
        Value: 3.0,
        Label: "<2 years",
      },
    }

    licence_factor_range, err := strategies.FindMatchingLicenceValidityFactor(&request, LicenceValidityFactorList)
  	util.AssertTrue(err != nil, t)
    util.AssertEqual(err.Error(), "Error wile Parsing LicenseHeldSince date. Error: parsing time \"22-Jan-2020\" as \"2006-01-02\": cannot parse \"an-2020\" as \"2006\"", t)
    util.AssertTrue(licence_factor_range == nil, t)
  })
  tp.Run("TestPricingStrategyVariedConfigsToFindMatchingInsuranceGroupFactor-Success-Senario", func(t *testing.T) {
    request.InsuranceGroup = 3
    InsuranceGroupFactorList := []models.RangeConfig{
      models.RangeConfig{
        Start: 0,
        End: 2,
        IsEligible: true,
        Value: 3.0,
        Label: "<2 years",
      },
      models.RangeConfig{
        Start: 2,
        End: 4,
        IsEligible: true,
        Value: 2.0,
        Label: "2-4",
      },
    }

    insurance_factor_range, err := strategies.FindMatchingInsuranceGroupFactor(&request, InsuranceGroupFactorList)
    util.AssertTrue(err == nil, t)
    util.AssertTrue(insurance_factor_range != nil, t)
    util.AssertEqual(*insurance_factor_range, InsuranceGroupFactorList[1], t)
  })
  tp.Run("TestPricingStrategyVariedConfigsToFindMatchingInsuranceGroupFactor-NotFound-Senario", func(t *testing.T) {
    request.InsuranceGroup = 3
    InsuranceGroupFactorList := []models.RangeConfig{
      models.RangeConfig{
        Start: 0,
        End: 2,
        IsEligible: true,
        Value: 3.0,
        Label: "<2 years",
      },
    }

    insurance_factor_range, err := strategies.FindMatchingInsuranceGroupFactor(&request, InsuranceGroupFactorList)
    util.AssertTrue(err != nil, t)
    util.AssertTrue(insurance_factor_range == nil, t)
    util.AssertEqual(err.Error(), "MatchingInsuranceGroupFactor not found!", t)
  })
  tp.Run("TestPricingStrategyVariedConfigsToFindMatchingDriverAgeFactor-Success-Senario", func(t *testing.T) {
    now := time.Now()
    _t := now.AddDate(-3, 0, 0) //3 years before now
    request.DateOfBirth = _t.Format("2006-01-02")
    DriverAgeFactorList := []models.RangeConfig{
      models.RangeConfig{
        Start: 0,
        End: 2,
        IsEligible: true,
        Value: 3.0,
        Label: "<2 years",
      },
      models.RangeConfig{
        Start: 2,
        End: 4,
        IsEligible: true,
        Value: 2.0,
        Label: "2-4",
      },
    }

    driver_factor_range, err := strategies.FindMatchingDriverAgeFactor(&request, DriverAgeFactorList)
    util.AssertTrue(err == nil, t)
    util.AssertTrue(driver_factor_range != nil, t)
    util.AssertEqual(*driver_factor_range, DriverAgeFactorList[1], t)
  })
  tp.Run("TestPricingStrategyVariedConfigsToFindMatchingDriverAgeFactor-NotFound-Senario", func(t *testing.T) {
    now := time.Now()
    _t := now.AddDate(-1, 0, 0) //1 year before now
    request.DateOfBirth = _t.Format("2006-01-02")
    DriverAgeFactorList := []models.RangeConfig{
      models.RangeConfig{
        Start: 2,
        End: 4,
        IsEligible: true,
        Value: 2.0,
        Label: "2-4",
      },
    }

    driver_factor_range, err := strategies.FindMatchingDriverAgeFactor(&request, DriverAgeFactorList)
    util.AssertTrue(err != nil, t)
    util.AssertEqual(err.Error(), "MatchingDriverAgeFactor not found!", t)
    util.AssertTrue(driver_factor_range == nil, t)
  })

  tp.Run("TestPricingStrategyVariedConfigsToFindMatchingDriverAgeFactor-InvalidDate-Senario", func(t *testing.T) {
    request.DateOfBirth = "22-Jan-2020"
    DriverAgeFactorList := []models.RangeConfig{
      models.RangeConfig{
        Start: 2,
        End: 4,
        IsEligible: true,
        Value: 2.0,
        Label: "2-4",
      },
    }

    driver_factor_range, err := strategies.FindMatchingDriverAgeFactor(&request, DriverAgeFactorList)
    util.AssertTrue(err != nil, t)
    util.AssertEqual(err.Error(), "Error wile Parsing DateOfBirth date. Error: parsing time \"22-Jan-2020\" as \"2006-01-02\": cannot parse \"an-2020\" as \"2006\"", t)
    util.AssertTrue(driver_factor_range == nil, t)
  })
}
