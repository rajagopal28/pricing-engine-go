package config

import (
  "log"
  "testing"
  "time"

  "pricingengine/service/config"
  "pricingengine/service/model"
  "pricingengine/test/util"
  )


func TestConfigCacheScenarios(tp *testing.T){
  cache := config.ConfigCache{
    Fetcher: config.ConfigFetcher {Path: "/../test_configs/"},
  }

  tp.Run("TestConfigCacheInitOnExistingCache", func(t *testing.T) {
    now := time.Now().Unix()
    cache.TimeToLive = now + 1000 // some future timestamp
    expected := []models.RangeConfig{
      models.RangeConfig {
        Start: 0,
        End: 10,
        Value: 1,
        Label: "Some Range",
      },
    }

    cache.BaseRateList = expected
    cache.DriverAgeFactorList = expected
    cache.InsuranceGroupFactorList = expected
    cache.LicenceValidityFactorList = expected

    cache.InitialiseWithRefresh(false, 1000)

    util.AssertEqual(now + 1000, cache.TimeToLive, t)
    util.AssertEqual(expected, cache.BaseRateList, t)
    util.AssertEqual(expected, cache.DriverAgeFactorList, t)
    util.AssertEqual(expected, cache.InsuranceGroupFactorList, t)
    util.AssertEqual(expected, cache.LicenceValidityFactorList, t)
    log.Printf("List : %+v", cache)
  })
  tp.Run("TestConfigCacheInitWithNoCacheRefreshButNoTTL", func(t *testing.T) {
    cache.TimeToLive = 0
    cache.InitialiseWithRefresh(false, 500)
    now := time.Now().Unix()
    util.AssertEqual(now + 500, cache.TimeToLive, t)
    util.AssertEqual(len(cache.BaseRateList), 2, t)
    util.AssertEqual(len(cache.DriverAgeFactorList), 3, t)
    util.AssertEqual(len(cache.InsuranceGroupFactorList),2 , t)
    util.AssertEqual(len(cache.LicenceValidityFactorList),2 , t)
    log.Printf("List : %+v", cache)
  })
  tp.Run("TestConfigCacheInitWithCacheRefreshFlag", func(t *testing.T) {
    now := time.Now().Unix()
    cache.TimeToLive = now + 500 // valid TTL still cache should force refresh with true
    actual := []models.RangeConfig{}
    cache.BaseRateList = actual
    cache.DriverAgeFactorList = actual
    cache.InsuranceGroupFactorList = actual
    cache.LicenceValidityFactorList = actual

    cache.InitialiseWithRefresh(true, 800)

    util.AssertEqual(now + 800, cache.TimeToLive, t)
    util.AssertEqual(len(cache.BaseRateList), 2, t)
    util.AssertEqual(len(cache.DriverAgeFactorList), 3, t)
    util.AssertEqual(len(cache.InsuranceGroupFactorList),2 , t)
    util.AssertEqual(len(cache.LicenceValidityFactorList),2 , t)
    log.Printf("List : %+v", cache)
  })
  tp.Run("TestConfigCacheInitWithCacheRefreshAfterTimeout", func(t *testing.T) {
    now := time.Now().Unix()
    cache.TimeToLive = now - 1000 // some past timestamp
    actual := []models.RangeConfig{
      models.RangeConfig {
        Start: 0,
        End: 10,
        Value: 1,
        Label: "Some Range",
      },
    }

    cache.BaseRateList = actual
    cache.DriverAgeFactorList = actual
    cache.InsuranceGroupFactorList = actual
    cache.LicenceValidityFactorList = actual

    cache.InitialiseWithRefresh(false, 1200)

    util.AssertEqual(now + 1200, cache.TimeToLive, t)
    util.AssertEqual(len(cache.BaseRateList), 2, t)
    util.AssertEqual(len(cache.DriverAgeFactorList), 3, t)
    util.AssertEqual(len(cache.InsuranceGroupFactorList),2 , t)
    util.AssertEqual(len(cache.LicenceValidityFactorList),2 , t)
    log.Printf("List : %+v", cache)
  })
}
