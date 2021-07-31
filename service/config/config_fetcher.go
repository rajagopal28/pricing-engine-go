package config

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"reflect"
  "log"
)

type ConfigFetcher struct{
  Path string
}


func (c *ConfigFetcher) ReadFileAndGetAsObject(filename string, class interface{}) (interface{}, error) {
  log.Println("Entering ReadFileAndGetAsObject")
	pwd, _ := os.Getwd()
	log.Println("Current Working Directory: ", pwd)
  jsonFile, err := os.Open(pwd+c.Path+filename)
	// txt, _ := ioutil.ReadFile(pwd+"/path/to/file.txt")
  // if we os.Open returns an error then handle it

  if err != nil {
			log.Printf("error opening file: %v", err)
      return nil, err
  }
  log.Println("Successfully Opened ", filename)
  // defer the closing of our jsonFile so that we can parse it later on
  defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	command := reflect.New(reflect.TypeOf(class))
  json.Unmarshal([]byte(byteValue), command.Interface())
	result := command.Elem().Interface()
  log.Println("Leaving ReadFileAndGetAsObject")
	return result, nil
}
