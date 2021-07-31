package config

import (
  "log"
  "testing"

  "pricingengine/service/config"
  "pricingengine/test/util"
  )


func TestConfigFetchSuccess(t *testing.T){
  fetcher := config.ConfigFetcher {Path: "/../test_configs/"}
  var temp []map[string]interface{}
  res, err := fetcher.ReadFileAndGetAsObject("test-base-rate.json" , temp)
  if err != nil {
    log.Println("error reading the config file:", err)
    t.Errorf("got error: %q", err)
  }
  result := res.([]map[string]interface{})
  println("result length =", len(result))
  if len(result) != 2 {
    t.Errorf("Expected 2 records, found %q", len(result))
  }
  util.AssertEqual(result[0]["time"], float64(1800), t)
  util.AssertEqual(result[0]["label"], "0.5 hours", t)
  util.AssertEqual(result[0]["rate"], float64(273), t)


  util.AssertEqual(result[1]["time"], float64(345600), t)
  util.AssertEqual(result[1]["label"], "96 hours / 4 days", t)
  util.AssertEqual(result[1]["rate"], float64(5204), t)
  log.Printf("List : %+v", result)
}
