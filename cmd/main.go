package main

import (
	"pricingengine/service"
)

// Main method that invokes the service and starts it at default port
func main() {
	service := service.Service{}
	service.Start("")
}
