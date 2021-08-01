package util
import (
  "testing"

	"pricingengine/service/util"
	"pricingengine/service/model"
)


func TestFactorMapperFeaturesWithMockConfigs(tp *testing.T){
  factorMapper := util.FactorMapper{}
  int_max := int((^uint(0))>> 1)
  tp.Run("TestFactorMapperBaseRateToRangeConfig-SuccessScenario", func(t *testing.T) {
    BaseRateList := []models.BaseRate{
      models.BaseRate{
        Label: "1 hour",
        Time: 3600,
        Rate: 100,
      },
      models.BaseRate{
        Label: "2 hour",
        Time: 7200,
        Rate: 200,
      },
    }
    rateConfigs := factorMapper.BaseRateToRangeConfig(BaseRateList)
    AssertEqual(len(rateConfigs), 2, t)
    AssertEqual(rateConfigs[0], models.RangeConfig{
      Start: 0,
      End: BaseRateList[0].Time,
      IsEligible: true,
      Value: BaseRateList[0].Rate,
      Label: BaseRateList[0].Label,
    }, t)
    AssertEqual(rateConfigs[1], models.RangeConfig{
      Start: BaseRateList[0].Time,
      End: BaseRateList[1].Time,
      IsEligible: true,
      Value: BaseRateList[1].Rate,
      Label: BaseRateList[1].Label,
    }, t)
  })
  tp.Run("TestFactorMapperDriverAgeFactorToRangeConfig-SuccessScenario", func(t *testing.T) {
    DriverAgeFactorList := []models.DriverAgeFactor{
      models.DriverAgeFactor{
        Label: "<10 years",
        Factor: 1.1,
        IsEligible: true,
        Age: 10,
      },
    }
    rateConfigs := factorMapper.DriverAgeFactorToRangeConfig(DriverAgeFactorList)
    AssertEqual(len(rateConfigs), 2, t)
    AssertEqual(rateConfigs[0], models.RangeConfig{
      Start: 0,
      End: DriverAgeFactorList[0].Age,
      IsEligible: true,
      Value: DriverAgeFactorList[0].Factor,
      Label: "Driver Age:0-10",
    }, t)
    AssertEqual(rateConfigs[1], models.RangeConfig{
      Start: DriverAgeFactorList[0].Age,
      End: int_max,
      IsEligible: true,
      Value: DriverAgeFactorList[0].Factor,
      Label: "Driver Age >10",
    }, t)
  })
  tp.Run("TestFactorMapperInsuranceGroupFactorToRangeConfig-SuccessScenario", func(t *testing.T) {
    InsuranceGroupFactorList := []models.InsuranceGroupFactor{
      models.InsuranceGroupFactor{
        Label: "<10 years",
        Factor: 1.1,
        IsEligible: true,
        Group: "1-10",
      },models.InsuranceGroupFactor{
        Label: ">10 years",
        Factor: 1,
        IsEligible: false,
        Group: "10",
      },
    }
    rateConfigs := factorMapper.InsuranceGroupFactorToRangeConfig(InsuranceGroupFactorList)
    AssertEqual(len(rateConfigs), 2, t)
    AssertEqual(rateConfigs[0], models.RangeConfig{
      Start: 1,
      End: 10,
      IsEligible: InsuranceGroupFactorList[0].IsEligible,
      Value: InsuranceGroupFactorList[0].Factor,
      Label: "Insurance Group:1-10",
    }, t)
    AssertEqual(rateConfigs[1], models.RangeConfig{
      Start: 10,
      End: int_max,
      IsEligible: InsuranceGroupFactorList[1].IsEligible,
      Value: InsuranceGroupFactorList[1].Factor,
      Label: "Insurance Group:10",
    }, t)
  })
  tp.Run("TestFactorMapperLicenceValidityFactorToRangeConfig-SuccessScenario", func(t *testing.T) {
    LicenceValidityFactorList := []models.LicenceValidityFactor{
      models.LicenceValidityFactor{
        Factor: 1.1,
        Length: "1-10",
      },models.LicenceValidityFactor{
        Factor: 1,
        Length: "10",
      },
    }
    rateConfigs := factorMapper.LicenceValidityFactorToRangeConfig(LicenceValidityFactorList)
    AssertEqual(len(rateConfigs), 2, t)
    AssertEqual(rateConfigs[0], models.RangeConfig{
      Start: 1,
      End: 10,
      IsEligible: true,
      Value: LicenceValidityFactorList[0].Factor,
      Label: "Licence Validity:1-10",
    }, t)
    AssertEqual(rateConfigs[1], models.RangeConfig{
      Start: 10,
      End: int_max,
      IsEligible: true,
      Value: LicenceValidityFactorList[1].Factor,
      Label: "Licence Validity:10",
    }, t)
  })
}
