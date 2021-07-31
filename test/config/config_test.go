package config

import (
  "log"
  "testing"
  "strings"

  "pricingengine/service/config"
  "pricingengine/test/util"
  )


func TestConfigFetch(tp *testing.T){
  fetcher := config.ConfigFetcher {Path: "/../test_configs/"}
  var temp []map[string]interface{}
  tp.Run("TestConfigFetchSuccess", func(t *testing.T) {
    res, err := fetcher.ReadFileAndGetAsObject("base-rate.json" , temp)
    if err != nil {
      log.Println("error reading the config file:", err)
      t.Errorf("got error: %q", err)
    }
    result := res.([]map[string]interface{})
    log.Println("result length =", len(result))
    util.AssertEqual(2, len(result), t)
    util.AssertEqual(result[0]["time"], float64(1800), t)
    util.AssertEqual(result[0]["label"], "0.5 hours", t)
    util.AssertEqual(result[0]["rate"], float64(273), t)


    util.AssertEqual(result[1]["time"], float64(345600), t)
    util.AssertEqual(result[1]["label"], "96 hours / 4 days", t)
    util.AssertEqual(result[1]["rate"], float64(5204), t)
    log.Printf("List : %+v", result)
  })
  tp.Run("TestConfigFetchFileNotFoundError", func(t *testing.T) {
     _,err := fetcher.ReadFileAndGetAsObject("some-base-rate.json" , temp)
    log.Println("Getting result:", err)
    if err == nil || strings.Contains(err.Error(), "error opening file: open"){
      log.Println("Should not come here!")
      t.Errorf("should get error")
    }
    log.Println("Coming here!")
  })
}
