package config
import (
 "log"
 "time"

 "pricingengine/service/model"
 "pricingengine/service/util"
)

type ConfigCache struct{
  Fetcher ConfigFetcher
  TimeToLive int64
  BaseRateList []models.RangeConfig // all converted range list
  DriverAgeFactorList []models.RangeConfig
  InsuranceGroupFactorList []models.RangeConfig
  LicenceValidityFactorList []models.RangeConfig
}


func (c *ConfigCache) Initialise(TTL int64) (error) {
  log.Println("Initialising ConfigCache with new TTL: +v", TTL)
  err := c.FetchAndConvertBaseFareList()
  if err != nil {
    return err
  }
  err = c.FetchAndConvertDriverAgeFactorList()
  if err != nil {
    return err
  }
  err = c.FetchAndConvertInsuranceGroupFactorList()
  if err != nil {
    return err
  }
  err = c.FetchAndConvertLicenceValidityFactorList()
  if err != nil {
    return err
  }
  c.TimeToLive = time.Now().Unix() + TTL // time to live in epoch seconds
  return nil
}


func (c *ConfigCache) InitialiseWithRefresh(refresh_cache bool, TTL int64) (error) {
  now := time.Now().Unix()
  log.Println("Initialising ConfigCache with Refresh:", refresh_cache, "TimeToLive: ", c.TimeToLive, " Now: ", now)
  if(refresh_cache || c.TimeToLive == 0 || now > c.TimeToLive) {
    return c.Initialise(TTL) // reload all the file if it is fresh or TTL is expired
  }
  return nil
}


func (c *ConfigCache) FetchAndConvertBaseFareList() (error) {
  log.Println("In FetchAndConvertBaseFareList ")
  var temp []models.BaseRate
	res, err := c.Fetcher.ReadFileAndGetAsObject("base-rate.json" , temp)
	if err != nil {
		log.Println("error reading the config file:", err)
		return err
	}
	// println("result length =", len(result))
	log.Printf("List : %+v", res)
	temp = res.([]models.BaseRate)
  factorMapper := util.FactorMapper{}
  c.BaseRateList = factorMapper.BaseRateToRangeConfig(temp)
  log.Println("Mapped range config from file: +v", c.BaseRateList)
  return nil
}

func (c *ConfigCache) FetchAndConvertDriverAgeFactorList() (error) {
  log.Println("In FetchAndConvertDriverAgeFactorList ")
  var temp []models.DriverAgeFactor
	res, err := c.Fetcher.ReadFileAndGetAsObject("driver-age-factor.json" , temp)
	if err != nil {
		log.Println("error reading the config file:", err)
		return err
	}
	// println("result length =", len(result))
	log.Printf("List : %+v", res)
	temp = res.([]models.DriverAgeFactor)
  factorMapper := util.FactorMapper{}
  c.DriverAgeFactorList = factorMapper.DriverAgeFactorToRangeConfig(temp)
  log.Println("Mapped range config from file: +v", c.DriverAgeFactorList)
  return nil
}

func (c *ConfigCache) FetchAndConvertInsuranceGroupFactorList() (error) {
  log.Println("In FetchAndConvertInsuranceGroupFactorList ")
  var temp []models.InsuranceGroupFactor
	res, err := c.Fetcher.ReadFileAndGetAsObject("insurance-group-factor.json" , temp)
	if err != nil {
		log.Println("error reading the config file: ", err)
		return err
	}
	// println("result length =", len(result))
	log.Printf("List : %+v", res)
	temp = res.([]models.InsuranceGroupFactor)
  factorMapper := util.FactorMapper{}
 	c.InsuranceGroupFactorList = factorMapper.InsuranceGroupFactorToRangeConfig(temp)
  log.Println("Mapped range config from file: +v", c.InsuranceGroupFactorList)
  return nil
}

func (c *ConfigCache) FetchAndConvertLicenceValidityFactorList() (error) {
  log.Println("In FetchAndConvertLicenceValidityFactorList ")
  var temp []models.LicenceValidityFactor
  res, err := c.Fetcher.ReadFileAndGetAsObject("licence-validity-factor.json" , temp)
  if err != nil {
    log.Println("error reading the config file: ", err)
    return err
  }
  // println("result length =", len(result))
  log.Printf("List : %+v", res)
  temp = res.([]models.LicenceValidityFactor)
  factorMapper := util.FactorMapper{}
  c.LicenceValidityFactorList = factorMapper.LicenceValidityFactorToRangeConfig(temp)
  log.Println("Mapped range config from file: +v", c.LicenceValidityFactorList)
  return nil
}
