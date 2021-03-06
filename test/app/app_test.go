package app

import (
  "testing"
  "context"
  "time"

  "pricingengine"
	"pricingengine/service/app"
	"pricingengine/service/config"
  "pricingengine/test/util"
)


func TestPriceGenerationAppWithActualConfigs(tp *testing.T){

  testApp := app.App{
    Cache: config.ConfigCache{
      TimeToLive : 1,
      Fetcher: config.ConfigFetcher{
        Path: "/../test_configs/",
      },
    },
  }
  tp.Run("TestPriceGenerationAppToGetValidGeneratedPriceList-FailureScenario-1", func(t *testing.T) {
    resp,_ := testApp.GeneratePricing(context.Background(),&pricingengine.GeneratePricingRequest{})

    util.AssertFalse(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 0, t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.Message, "DateOfBirth cannot be empty", t)
  })
  tp.Run("TestPriceGenerationAppToGetValidGeneratedPriceList-FailureScenario-2", func(t *testing.T) {

    request := pricingengine.GeneratePricingRequest{
      DateOfBirth: "2006-01-02",
    }
    resp,_ := testApp.GeneratePricing(context.Background(),&request)

    util.AssertFalse(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 0, t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.Message, "InsuranceGroup should be a Positive number", t)
  })

  tp.Run("TestPriceGenerationAppToGetValidGeneratedPriceList-FailureScenario-3", func(t *testing.T) {

    request := pricingengine.GeneratePricingRequest{
      DateOfBirth: "2006-01-02",
      InsuranceGroup: 20,
    }
    resp,_ := testApp.GeneratePricing(context.Background(),&request)

    util.AssertFalse(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 0, t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.Message, "LicenseHeldSince Date cannot be empty", t)
  })

  tp.Run("TestPriceGenerationAppToGetValidGeneratedPriceList-FailureScenario-4", func(t *testing.T) {

    request := pricingengine.GeneratePricingRequest{
      DateOfBirth: "2001-01-02",
      InsuranceGroup: 20,
      LicenseHeldSince: "2006-01-02",
    }
    resp,_ := testApp.GeneratePricing(context.Background(),&request)

    util.AssertFalse(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 0, t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.Message, "Declined due to :Insurance Group:8", t)
  })

  tp.Run("TestPriceGenerationAppToGetValidGeneratedPriceList-FailureScenario-5", func(t *testing.T) {

    now := time.Now()
    _t := now.AddDate(-15, 0, 0) //15 years before now
    dob := _t.Format("2006-01-02")
    request := pricingengine.GeneratePricingRequest{
      DateOfBirth: dob,
      InsuranceGroup: 20,
      LicenseHeldSince: "2006-01-02",
    }
    resp,_ := testApp.GeneratePricing(context.Background(),&request)

    util.AssertFalse(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 0, t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.Message, "Declined due to :Driver Age:0-16", t)
  })


  tp.Run("TestPriceGenerationAppToGetValidGeneratedPriceList-SuccessScenario-1", func(t *testing.T) {

    now := time.Now()
    _t := now.AddDate(-20, 0, 0) //20 years before now
    dob := _t.Format("2006-01-02")
    _t = now.AddDate(-7, 0, 0) //7 years before now
    licence := _t.Format("2006-01-02")
    request := pricingengine.GeneratePricingRequest{
      DateOfBirth: dob,
      InsuranceGroup: 7,
      LicenseHeldSince: licence,
    }
    resp,_ := testApp.GeneratePricing(context.Background(),&request)

    util.AssertTrue(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 2, t)
    util.AssertEqual(resp.Message, "Success", t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.PricingList[0], pricingengine.PricingItem{
      Premium: 259.349,
      Currency: "??",
      FareGroup: "0.5 hours, Driver Age:16-26, Insurance Group:1-8, Licence Validity:6",
      }, t)
      util.AssertEqual(resp.PricingList[1], pricingengine.PricingItem{
        Premium: 4943.8,
        Currency: "??",
        FareGroup: "96 hours / 4 days, Driver Age:16-26, Insurance Group:1-8, Licence Validity:6",
        }, t)
  })

  tp.Run("TestPriceGenerationAppToGetValidGeneratedPriceList-SuccessScenario-2", func(t *testing.T) {
    now := time.Now()
    _t := now.AddDate(-20, 0, 0) //20 years before now
    dob := _t.Format("2006-01-02")
    _t = now.AddDate(-5, 0, 0) //5 years before now
    licence := _t.Format("2006-01-02")
    request := pricingengine.GeneratePricingRequest{
      DateOfBirth: dob,
      InsuranceGroup: 7,
      LicenseHeldSince: licence,
    }
    resp,_ := testApp.GeneratePricing(context.Background(),&request)

    util.AssertTrue(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 2, t)
    util.AssertEqual(resp.Message, "Success", t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.PricingList[0], pricingengine.PricingItem{
      Premium: 300.3,
      Currency: "??",
      FareGroup: "0.5 hours, Driver Age:16-26, Insurance Group:1-8, Licence Validity:0-6",
      }, t)
      util.AssertEqual(resp.PricingList[1], pricingengine.PricingItem{
        Premium: 5724.4,
        Currency: "??",
        FareGroup: "96 hours / 4 days, Driver Age:16-26, Insurance Group:1-8, Licence Validity:0-6",
        }, t)
  })
}
