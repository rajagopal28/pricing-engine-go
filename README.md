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


## Understanding the existing setup
The current application is built on top on an existing code as per the [Requirements](./REQUIREMENTS.md). The setup typically had a layered data flow among the components as follows:
- The main.go file which typically invoked the service.
- The service that serves various REST end point for the application and passing on the data to RPC layer.
- The RPC which typically takes care of doing the conversion from HTTP middleware understandable data to the Application understandable data.
- Then the RPC passes the application understandable data to the *App* component which has all the business logics of computing the pricing and serving the required data.


## The Approach in solving the problem


## The Flow


## The Endpoints
### Generate Pricing for a customer
#### Add a new transaction
##### Request
```http
POST /generate_pricing HTTP/1.1
Host: localhost:3000
Content-Type: application/json

{
    "date_of_birth": "1970-12-04",
    "insurance_group": 12,
    "license_held_since": "1988-08-01"
}
```
### Request format:
*date_of_birth* – Date of Birth of the customer who is trying to rent the vehicle
*insurance_group* – the insurance group to which the customer belongs to
*license_held_since* – The date of acquiring of the Driver's licence by the existing customer

##### Response
Returns: Empty body with one of the following:

200 – in case of success


```http
HTTP/1.1 200 OK
Content-Type: application/json
{
  "input": {
          "date_of_birth": "1970-12-04",
          "insurance_group": 12,
          "license_held_since": "1988-08-01"
      },
      "is-eligible": true,
      "message": "Success",
      "pricing": [
          {
              "premium": 278.28254999999996,
              "currency": "£",
              "fare_group": "0.5 hours, Driver Age >26, Insurance Group:9-16, Licence Validity:6"
          },
          {....},
          {....},
          .....
      ]
}

```

#### Get current pricing configuration ranges
##### Request
```http
GET /statistics HTTP/1.1
Host: localhost:3000


```

##### Response
Returns: The computed statistics for all the posted transaction or empty if there are no transactions present.

```http
HTTP/1.1 200
Content-Type: application/json;charset=UTF-8
{
  "base-rate": [
        {
            "Start": 0,
            "End": 1800,
            "IsEligible": true,
            "Value": 273,
            "Label": "0.5 hours"
        },
        {....}
        ....
        ],
    "driver-age-factor": [
        {
            "Start": 0,
            "End": 16,
            "IsEligible": false,
            "Value": 0,
            "Label": "Driver Age:0-16"
        },
        {...}
        ....
      ],
      "insurance-group-factor": [
          {
              "Start": 1,
              "End": 8,
              "IsEligible": true,
              "Value": 1,
              "Label": "Insurance Group:1-8"
          },
          {...},
          ....
      ],
      "licence-validity-factor": [
        {
            "Start": 0,
            "End": 1,
            "IsEligible": true,
            "Value": 1.1,
            "Label": "Licence Validity:0-1"
        },
        {...}
        .....
      ]
    }
```



## Setting up the environment
This application is purely build on goLang. To run this application in your local you need to follow the below steps

#### Build:
```
go build ./cmd/.
```

#### Run
```
go run ./cmd/.
```


#### Test
```
go test ./test/.
```

## References:
- setting up go environment : https://pkg.go.dev/cmd/go
- go packages : https://medium.com/rungo/everything-you-need-to-know-about-packages-in-go-b8bac62b74cc
- maps and struts : https://tour.golang.org/moretypes/19
- getting query params : https://golangbyexample.com/net-http-package-get-query-params-golang/
- current working directory: https://stackoverflow.com/a/31464648
- parsing json: https://www.sohamkamani.com/golang/parsing-json/
- reading json file : https://tutorialedge.net/golang/parsing-json-with-golang/
- generic struct type parsing unmarshall : https://stackoverflow.com/a/49002939
- fix for `implicit assignment of unexported field` issue: https://stackoverflow.com/a/48130761
- logging in go: https://www.loggly.com/use-cases/logging-in-golang-how-to-start/
- logging in detail: https://www.honeybadger.io/blog/golang-logging/
- passing functions as parameters: https://stackoverflow.com/a/12655719
- fixing cast exception from `interface{}` to other `structs array` : https://stackoverflow.com/a/42740448
- unbounded list creation, updation : https://stackoverflow.com/a/3387362
- sorting struct list : https://yourbasic.org/golang/how-to-sort-in-go/
- https://stackoverflow.com/a/42382736
- https://stackoverflow.com/a/29000001
- time/date formatting : https://stackoverflow.com/a/54314594
- days between dates: https://yourbasic.org/golang/days-between-dates/
- https://stackoverflow.com/a/46298277
- collections-functions : https://gobyexample.com/collection-functions
- String to int: https://stackoverflow.com/a/4279644
- int to String : https://stackoverflow.com/a/10105983
- max integer value : https://stackoverflow.com/a/6878625
- check empty string : https://stackoverflow.com/a/18594463
- splitting string: https://stackoverflow.com/a/16551613
- https://www.educative.io/edpresso/how-to-split-a-string-in-golang
- pointers and references: https://medium.com/wesionary-team/pointers-and-passby-value-reference-in-golang-a00c8c59b7f1
- string formatting : https://gobyexample.com/string-formatting
- https://pkg.go.dev/fmt
- handling error messages: https://stackoverflow.com/a/22171548
- no nil check for struct variables: https://stackoverflow.com/a/25350005
- https://stackoverflow.com/a/28447372
