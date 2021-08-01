package service

import (
  "testing"
  "fmt"
  "net/http"
  "net/http/httptest"
	"encoding/json"
	"io/ioutil"
  "strings"

  "pricingengine"
	// "pricingengine/service"
	"pricingengine/service/rpc"
	"pricingengine/service/app"
	"pricingengine/service/config"


  "pricingengine/test/util"
)


func TestServiceEndpointsIntegrationTestWithActualConfigs(tp *testing.T){
  /*
  service := service.Service{}
	service.Start("2020")
  service.Stop()
  */
  tp.Run("TestRESTAPIEndpointToGetValidGeneratedPriceList-FailureScenario-1", func(t *testing.T) {

    resp,_ := MakeHttpRequestAndGetResponse(&pricingengine.GeneratePricingRequest{})
    println(resp)
    util.AssertFalse(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 0, t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.Message, "DateOfBirth cannot be empty", t)
  })
  tp.Run("TestRESTAPIEndpointToGetValidGeneratedPriceList-FailureScenario-2", func(t *testing.T) {

    request := pricingengine.GeneratePricingRequest{
      DateOfBirth: "2006-01-02",
    }
    resp,_ := MakeHttpRequestAndGetResponse(&request)
    println(resp)
    util.AssertFalse(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 0, t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.Message, "InsuranceGroup should be a Positive number", t)
  })

  tp.Run("TestRESTAPIEndpointToGetValidGeneratedPriceList-FailureScenario-3", func(t *testing.T) {

    request := pricingengine.GeneratePricingRequest{
      DateOfBirth: "2006-01-02",
      InsuranceGroup: 20,
    }
    resp,_ := MakeHttpRequestAndGetResponse(&request)
    println(resp)
    util.AssertFalse(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 0, t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.Message, "LicenseHeldSince Date cannot be empty", t)
  })

  tp.Run("TestRESTAPIEndpointToGetValidGeneratedPriceList-FailureScenario-4", func(t *testing.T) {

    request := pricingengine.GeneratePricingRequest{
      DateOfBirth: "2001-01-02",
      InsuranceGroup: 20,
      LicenseHeldSince: "2006-01-02",
    }
    resp,_ := MakeHttpRequestAndGetResponse(&request)
    println(resp)
    util.AssertFalse(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 0, t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.Message, "Declined due to :Insurance Group:8", t)
  })

  tp.Run("TestRESTAPIEndpointToGetValidGeneratedPriceList-FailureScenario-5", func(t *testing.T) {

    request := pricingengine.GeneratePricingRequest{
      DateOfBirth: "2006-01-02",
      InsuranceGroup: 20,
      LicenseHeldSince: "2006-01-02",
    }
    resp,_ := MakeHttpRequestAndGetResponse(&request)
    println(resp)
    util.AssertFalse(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 0, t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.Message, "Declined due to :Driver Age:0-16", t)
  })


  tp.Run("TestRESTAPIEndpointToGetValidGeneratedPriceList-SuccessScenario-1", func(t *testing.T) {

    request := pricingengine.GeneratePricingRequest{
      DateOfBirth: "2001-01-02",
      InsuranceGroup: 7,
      LicenseHeldSince: "2006-01-02",
    }
    resp,_ := MakeHttpRequestAndGetResponse(&request)
    println(resp)
    util.AssertTrue(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 2, t)
    util.AssertEqual(resp.Message, "Success", t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.PricingList[0], pricingengine.PricingItem{
      Premium: 259.349,
      Currency: "£",
      FareGroup: "0.5 hours, Driver Age:16-26, Insurance Group:1-8, Licence Validity:6",
      }, t)
      util.AssertEqual(resp.PricingList[1], pricingengine.PricingItem{
        Premium: 4943.8,
        Currency: "£",
        FareGroup: "96 hours / 4 days, Driver Age:16-26, Insurance Group:1-8, Licence Validity:6",
        }, t)
  })

  tp.Run("TestRESTAPIEndpointToGetValidGeneratedPriceList-SuccessScenario-2", func(t *testing.T) {
    request := pricingengine.GeneratePricingRequest{
      DateOfBirth: "2001-01-02",
      InsuranceGroup: 7,
      LicenseHeldSince: "2020-01-02",
    }
    resp,_ := MakeHttpRequestAndGetResponse(&request)
    println(resp)
    util.AssertTrue(resp.IsEligible, t)
    util.AssertEqual(len(resp.PricingList), 2, t)
    util.AssertEqual(resp.Message, "Success", t)
    util.AssertTrue(len(resp.Message) > 0, t)
    util.AssertEqual(resp.PricingList[0], pricingengine.PricingItem{
      Premium: 300.3,
      Currency: "£",
      FareGroup: "0.5 hours, Driver Age:16-26, Insurance Group:1-8, Licence Validity:0-6",
      }, t)
      util.AssertEqual(resp.PricingList[1], pricingengine.PricingItem{
        Premium: 5724.4,
        Currency: "£",
        FareGroup: "96 hours / 4 days, Driver Age:16-26, Insurance Group:1-8, Licence Validity:0-6",
        }, t)
  })
}

func MakeHttpRequestAndGetResponse( requestTo  *pricingengine.GeneratePricingRequest) (*pricingengine.GeneratePricingResponse, error) {

  jsonValue,err := json.Marshal(requestTo)
  println(jsonValue)
  if err != nil {
		fmt.Printf("failed to marshall data %v: %v", requestTo, err)
		return nil, err
	}
  rpc := rpc.RPC{
    App: &app.App{
      Cache: config.ConfigCache{
        TimeToLive : 1,
        Fetcher: config.ConfigFetcher{
          Path: "/../test_configs/",
        },
      },
    },
  }

  request := httptest.NewRequest(http.MethodPost, "/generate_pricing", strings.NewReader(string(jsonValue[:])))
  responseRecorder := httptest.NewRecorder()

  handler := http.HandlerFunc(rpc.GeneratePricing)
  handler.ServeHTTP(responseRecorder, request)


	if responseRecorder.Code != 200 {
		println("Want status '%d', got '%d'", 200, responseRecorder.Code)
	}

  println("Want '%s', got '%s'", "some", responseRecorder.Body)
	// if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
		// t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
	// }
  body, err := ioutil.ReadAll(responseRecorder.Body)
	if err != nil {
		// response(w, err)
		return nil, err
	}
  result := pricingengine.GeneratePricingResponse{}
  json.Unmarshal(body,&result)
  return &result, nil
}
