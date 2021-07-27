# Pricing Engine - GoLang
Application in goLang to work in a simple pricing engine considering various factors from the customer who is trying to rent a vehicle from a rental service provider.


## Requirements
The requirement of this application is to bring up a go based REST API which will determine the rental cost of a car bases on the attributes of the customer who is purchasing it. You can go through the detailed requirements [Here](./REQUIREMENTS.md). The essential attributes that play key aspects in the pricing calculation are as follows:
- Rent Duration
- Age of the customer
- Insurance Group
- Validity of the licence held by the customer


## The Pricing formula
> Total Rental Cost = BaseFare based on the Duration * Factor based on Age * Factor based on the Insurance Group * Factor based on Licence Validity

## References:
- setting up go environment : https://pkg.go.dev/cmd/go
- go packages : https://medium.com/rungo/everything-you-need-to-know-about-packages-in-go-b8bac62b74cc
