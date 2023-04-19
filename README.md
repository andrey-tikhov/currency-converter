# Currency converter
compile guaranteed for go 1.20.3

## Quickstart
    git clone git@github.com:andrey-tikhov/currency-converter.git
    go run main.go
Accepts request on the following endpoints on http://localhost:8000 
 - [/convert](#/convert-endpoint) allows to convert amount of one currency to another currency based on country central bank rate provided
 - [/get_exchange_rates](#/get_exchange_rates-endpoint) allows to load all central bank rates for provided country

### /convert endpoint
Accepts the following requests.
If country is omitted the default central bank will be applied (defined in config/base.yaml defaults)
```
    curl -X "POST" "http://localhost:8000/convert" \
     -d $'{
        "country": "russia",
        "source_currency": "JPY",
        "target_currency": "VND",
        "amount":20
      }'
```
Expected response
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 19 Apr 2023 20:47:29 GMT
Content-Length: 20
Connection: close

{"amount":3523.6376}
```

### /get_exchange_rates endpoint
Accepts the following requests
```
POST /get_exchange_rates HTTP/1.1
Host: localhost:8000
Connection: close
Content-Length: 24

{
"country": "russia"
 }
```
Expected response
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 19 Apr 2023 20:47:56 GMT
Connection: close
Transfer-Encoding: chunked

{"rates":{"AED":{"nominal":1,"base_currency":"RUR","target_currency":"AED","rate_target_to_base":22.2335},"AMD":{"nominal":100,"base_currency":"RUR","target_currency":"AMD","rate_target_to_base":21.0635}}}
```

# Architecture

4 layers service
- handler: accepts the incoming http requests, transforms them to internal entities
- controller: orchestrates internal calls between layers cleaning up technical data from downstream systems (e.g. Exchange rates data from repository is cleaned from Timezone and DateLoaded data). Isolates implementation of the data repository from handler.
- [repository](#repository-implementation): provides the interface to operate with service data. Orchestrates gateways that loads data and internal cache. Currently in-memory cache is implemented but db can be plugged in to this layer if needed. 
- gateway: provide basic integration endpoints with Central Banks API. Thus allowing integration of new sources of data easily.

On top of that 
- entities represent key objects that service operates with
- mapper functions transform entities between each other, providing better testability

## Fx dependency ingestion
Service leverages open-sourced Uber dependency ingestion framework [fx](https://pkg.go.dev/go.uber.org/fx)
In short this framework allows you to register constructors for various Interfaces and then provide them as params to the functions called.
The key benefit of the framework for service implementation is that constructors are called in lazy manner.
Meaning that once created object will be reused.

## repository implementation
Currently in memory cache based on Map and repository with sync.RWMutex is used for repository interface implementation.
```
type cbr struct {
	sync.RWMutex

	TimeNow    func() time.Time
	Gateways   map[string]gateway.CBGateway    // maps country to respective gateway
	RatesCache map[string]entity.ExchangeRates // maps country to ExchangeRatesObject
}
```
As constructors are used in lazy manner once created the repository implementation is reused and cached Central Bank rates are used in order to prevent unlimited requests to the central bank's endpoints.
This is tested in the integration test.
Cache reloading is protected by .Lock() as well as loading data from cache is protected by .RLock() to ensure no data races are happening while http handlers are executed in independet coroutines.
Cache reloading happens with the following logic (simplified, detailed logic is available in the [repository implementation](https://github.com/andrey-tikhov/currency-converter/blob/main/repository/cbr_repository.go#L126)).
1. Incoming request is enriched with current date in the time zone of central bank where data is requested.
2. Repository loads the cached data for the rates. If date loaded in the cache doesn't equal the date in the incoming request we reload the cache.
Key assumption here is that central banks provide the rates that are valid DAILY in the timezone of the respective central bank.
Nothing prevents us from replacing the implementation with database or anything else if needed.

## adding new source
Requires implementing new [CBGateway](https://github.com/andrey-tikhov/currency-converter/blob/main/gateway/CBAPI.go) for the respective central bank and updating repository implementation with that gateway.
As Gateway implementation is totally independent we do not care if the data we receive is JSON, XNL, plain text or whatever.
Configuration for the gateway must include the timezone where central bank is located.