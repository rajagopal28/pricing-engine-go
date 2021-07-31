package config

import (
  "log"
  "testing"
  "reflect"
  "pricingengine/service/config"
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
  expected := [2]map[string]interface{}{
    map[string]interface{}{
      "time":1800,
      "label":"0.5 hours",
      "rate":273,
    },
    map[string]interface{}{
      "time":345600,
      "label":"96 hours / 4 days",
      "rate":5204,
    },
  }
  AssertEqual(result[0]["time"], float64(1800), t)
  AssertEqual(result[0]["label"], "0.5 hours", t)
  AssertEqual(result[0]["rate"], float64(273), t)


  AssertEqual(result[1]["time"], float64(345600), t)
  AssertEqual(result[1]["label"], "96 hours / 4 days", t)
  AssertEqual(result[1]["rate"], float64(5204), t)
  log.Printf("expected : %+v", expected)
  log.Printf("List : %+v", result)
}



func AssertEqual(a interface{}, b interface{}, t *testing.T) {
  if a == b || reflect.DeepEqual(a, b) {
    return
  }
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}
