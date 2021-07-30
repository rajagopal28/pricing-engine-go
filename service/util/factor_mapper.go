package util

import (
  "sort"
  "strings"
  "strconv"
	"pricingengine/service/model"
)


type FactorMapper struct{}


func (f *FactorMapper) BaseRateToRangeConfig(baseRates []models.BaseRate) (rangeRates []models.RangeConfig) {
  sort.Slice(baseRates, func(i, j int) bool {
    return baseRates[i].Time < baseRates[j].Time
  })
  prev := 0
  result := []models.RangeConfig{}
  for i:= 0; i < len(baseRates); i++ {
    curr := baseRates[i]
    c_range := models.RangeConfig{Start: prev, End: curr.Time, Label: curr.Label, Value: curr.Rate, IsEligible: true}
    result = append(result, c_range)
    prev = curr.Time
  }
  return result
}


func (f *FactorMapper) DriverAgeFactorToRangeConfig(ageFactors []models.DriverAgeFactor) (rangeRates []models.RangeConfig) {
  sort.Slice(ageFactors, func(i, j int) bool {
    return ageFactors[i].Age < ageFactors[j].Age
  })
  prev := 0
  result := []models.RangeConfig{}
  for i:= 0; i < len(ageFactors); i++ {
    curr := ageFactors[i]
    c_range := models.RangeConfig{Start: prev, End: curr.Age, Label: "Driver Age:"+strconv.Itoa(prev) +"-"+strconv.Itoa(curr.Age), Value: curr.Factor, IsEligible: curr.IsEligible}
    result = append(result, c_range)
    if(i == len(ageFactors)-1) {
      // for last item add boundary range
      end := int((^uint(0))>> 1) // max int range
      c_range = models.RangeConfig{Start: curr.Age, End: end, Label: "Driver Age >"+strconv.Itoa(curr.Age), Value: curr.Factor, IsEligible: true}
      result = append(result, c_range)
    }
    prev = curr.Age
  }
  return result
}

func (f *FactorMapper) InsuranceGroupFactorToRangeConfig(insuranceGroups []models.InsuranceGroupFactor) (rangeRates []models.RangeConfig) {
  result := []models.RangeConfig{}
  for i:= 0; i < len(insuranceGroups); i++ {
    curr := insuranceGroups[i]
    s := strings.Split(curr.Group, "-")
    var start, end int
    start,_ = strconv.Atoi(s[0])
  	if(len(s) > 1) {
  		end,_ = strconv.Atoi(s[1])
  	} else {
      end = int((^uint(0))>> 1) // max int range
    }
    c_range := models.RangeConfig{Start: start, End: end, Label: "Insurance Group:"+curr.Group, Value: curr.Factor, IsEligible: curr.IsEligible}
    result = append(result, c_range)
  }
  return result
}

func (f *FactorMapper) LicenceValidityFactorToRangeConfig(licenceValidities []models.LicenceValidityFactor) (rangeRates []models.RangeConfig) {
  result := []models.RangeConfig{}
  for i:= 0; i < len(licenceValidities); i++ {
    curr := licenceValidities[i]
    s := strings.Split(curr.Length, "-")
    var start, end int
    start,_ = strconv.Atoi(s[0])
  	if(len(s) > 1) {
  		end,_ = strconv.Atoi(s[1])
  	} else {
      end = int((^uint(0))>> 1) // max int range
    }
    c_range := models.RangeConfig{Start: start, End: end, Label: "Licence Validity:"+curr.Length, Value: curr.Factor, IsEligible: true}
    result = append(result, c_range)
  }
  return result
}
