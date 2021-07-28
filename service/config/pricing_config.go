package config

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"reflect"
  "log"
)

type Config struct{}


func (a *Config) ReadFileAndGetAsObject(path string, class interface{}) (interface{}, error) {
	pwd, _ := os.Getwd()
	log.Print("Current Working Directory: ", pwd)
  jsonFile, err := os.Open(pwd+path)
	// txt, _ := ioutil.ReadFile(pwd+"/path/to/file.txt")
  // if we os.Open returns an error then handle it

  if err != nil {
			log.Fatalf("error opening file: %v", err)
      return nil, err
  }
  log.Print("Successfully Opened ", path)
  // defer the closing of our jsonFile so that we can parse it later on
  defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)


	command := reflect.New(reflect.TypeOf(class))
  json.Unmarshal([]byte(byteValue), command.Interface())
	result := command.Elem().Interface()
	println(result)
	return result, nil
}
