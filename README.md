## Run app

```
./start.sh
```

It deploys 2 docker containers: app at 8080 and mysql at 5432.
Mysql may take quite long to start on first installation, so the app may log error and keep retrying.

## Test

```
./test.sh
```

It runs both unit test and integration test.
Mysql took really long to init in my machine somehow, so i add retry inside the code.
The test report is located inside report folder.

## Project structure

### internal/services
- Handle business logic
- Handle validation
- Define interface for storage layer to implement, it accepts the fact that underlying storage is a sql dbms that supports transaction

### internal/transport
- Handle requests from outside world (we only have HTTP for now)
- Convert serialized obj into model and call service layer

### internal/storage
- Implements the interface exposed inside service layer 

### internal/common
- Acts as a wrapper for third party lib, for now we only have Mysql

### sqlc
- Let us write raw query and use tool to generate golang code for those queries

## Manual test

### PlaceWager

```
curl -i --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
    "total_wager_value": 50,
    "odds": 7,
    "selling_percentage": 20,
    "selling_price": 70
}'
```

Response:
```
{
   "id":1,
   "total_wager_value":50,
   "odds":7,
   "selling_percentage":20,
   "selling_price":70,
   "current_selling_price":70,
   "percentage_sold":null,
   "amount_sold":null,
   "placed_at":"2022-07-15T00:10:24Z"
}
```

## BuyWager

This requests depends on previous cmd in placing wager

#### Valid request
```
curl -i --location --request POST 'http://localhost:8080/buy/1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "buying_price": 30.0
}'
```

Response

```
{
   "id":1,
   "wager_id":1,
   "buying_price":30,
   "bought_at":"2022-07-15T00:15:29Z"
}
```

#### Invalid requests

```
curl -i --location --request POST 'http://localhost:8080/buy/1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "buying_price": 80.0
}'
```

```
curl -i --location --request POST 'http://localhost:8080/buy/1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "buying_price": -1.0
}'
```
Response

```
{"error":"INVALID BUYING PRICE"}
```


### WagerList


Create some wagers first then make this cmd

The items are order by placed_at desc by default

```
curl -i --location --request GET 'http://localhost:8080/wagers?page=1&limit=3'
```

Response

```
[
   {
      "id":3,
      "total_wager_value":50,
      "odds":7,
      "selling_percentage":20,
      "selling_price":70,
      "current_selling_price":70,
      "percentage_sold":null,
      "amount_sold":null,
      "placed_at":"2022-07-15T00:28:34Z"
   },
   {
      "id":2,
      "total_wager_value":50,
      "odds":7,
      "selling_percentage":20,
      "selling_price":70,
      "current_selling_price":70,
      "percentage_sold":null,
      "amount_sold":null,
      "placed_at":"2022-07-15T00:13:49Z"
   },
   {
      "id":1,
      "total_wager_value":50,
      "odds":7,
      "selling_percentage":20,
      "selling_price":70,
      "current_selling_price":40,
      "percentage_sold":43,
      "amount_sold":30,
      "placed_at":"2022-07-15T00:10:24Z"
   }
]
```

Page 2 limit 2

```
curl -i --location --request GET 'http://localhost:8080/wagers?page=2&limit=2'
```

Response

```
[
   {
      "id":1,
      "total_wager_value":50,
      "odds":7,
      "selling_percentage":20,
      "selling_price":70,
      "current_selling_price":40,
      "percentage_sold":43,
      "amount_sold":30,
      "placed_at":"2022-07-15T00:10:24Z"
   }
]
```



## TODO

- [ ] Proper authentication
- [ ] Add interceptor to recover from panic.
- [ ] The app needs more logging, and a place to collect and view logs for debug (GCP log console or maybe deploy Kibana stack).
- [ ] I'm skipping WagerList integration tests because it is hard to make the result item deterministic without mocking.
- [ ] More integration test for storage component.