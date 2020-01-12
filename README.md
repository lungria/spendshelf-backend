# spendshelf-backend

## Overview

Spendshelf started as a desire to keep money spending in the personal or/and family budget under control.

![UI demo](https://github.com/lungria/spendshelf-ui/blob/master/demo.gif)

All incoming transactions are being sent to categories - and each category has a monthly budget (spent money limit). At the end of the budget period (usually at the end of the months or near to salary date) - you either close to monthly budget per each category (which is good) - or not (which is bad). 

## Quick start
### Build with Docker
* Change the mongo DB credentials 
    * `MONGO_INITDB_ROOT_USERNAME: some_user_name` 
    * `MONGO_INITDB_ROOT_PASSWORD: some_password` 
    * `MONGO_URI=mongodb://some_user_name:some_password@mongod:27017`
* Build the app `docker-compose up -d --build`
* Now you have access to the application on `localhost:80`

### API
#### POST /webhook 
Catch a transaction from the mono bank webhook and insert it into transactions collection.  

e.g request:
```
curl -X POST \
  http://localhost:8080/webhook \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "type": "StatementItem",
    "data": {
        "account": "fCfTvZgMihvz_URzZSqVEA",
        "statementItem": {
            "id": "kjhKjhkHkjh",
            "time": 1576008041,
            "description": "Vodafone\n+380935684777",
            "mcc": 5555,
            "amount": -4564,
            "operationAmount": -12456,
            "currencyCode": 456,
            "commissionRate": 0,
            "cashbackAmount": 0,
            "balance": 14630,
            "hold": true
        }
    }
}'
```  

e.g response
```
{
    "message": "Success"
}
```

#### GET /categories 
Returns all categories from categories collection.  

e.g request:
```
curl -X GET \
  http://localhost/categories \
  -H 'cache-control: no-cache'
```  

e.g response
```
{
    "categories": [
        {
            "id": "5e146165f22ead680880b21a",
            "name": "Shopping",
            "normalizedName": "SHOPPING"
        },
        {
            "id": "5e148c35332c1e86b27b2ce0",
            "name": "Entertainment",
            "normalizedName": "ENTERTAINMENT"
        }
    ]
}
```

#### POST /categories 
Inserts a new category to categories collection.  

e.g request
```
curl -X POST \
  http://localhost/categories \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "name": "test"
}'
```

e.g response:
```
{
    "id": "5e1080bb7d28a9208146f12d"
}
```
#### GET /transactions or /transactions?category= 
Returns all transactions from transactions collection.  

e.g request
```
curl -X GET \
  'http://localhost/transactions?category=' \
  -H 'cache-control: no-cache'
```

e.g response:
```
{
    "transactions": [
        {
            "_id": "5e19bf98b1bb94b446d7977a",
            "dateTime": "2019-12-10T20:00:41Z",
            "description": "Description",
            "category": {
                "id": "5e146165f22ead680880b21a",
                "name": "Shopping",
                "normalizedName": "SHOPPING"
            },
            "amount": -4564,
            "balance": 14630,
            "bank": "Mono Bank"
        },
        {
            "_id": "5e19c0f7b1bb94b446d7977c",
            "dateTime": "2019-12-10T20:00:41Z",
            "description": "Description",
            "amount": -4564,
            "balance": 14630,
            "bank": "Mono Bank"
        }
    ]
}
```

#### GET /transactions?category=with 
Returns all uncategorized transactions.  

e.g request
```
curl -X GET \
  'http://localhost/transactions?category=with' \
  -H 'Cache-Control: no-cache' \
```

e.g response:
```
{
    "transactions": [
        {
            "_id": "5e19bf98b1bb94b446d7977a",
            "dateTime": "2019-12-10T20:00:41Z",
            "description": "Description",
            "category": {
                "id": "5e146165f22ead680880b21a",
                "name": "Shopping",
                "normalizedName": "SHOPPING"
            },
            "amount": -4564,
            "balance": 14630,
            "bank": "Mono Bank"
        }
    ]
}
```

#### GET /transactions?category=without 
Returns all categorized transactions. 

e.g request
```
curl -X GET \
  'http://localhost/transactions?category=without' \
  -H 'Cache-Control: no-cache' \
```

e.g response:
```
{
    "transactions": [
        {
            "_id": "5e19c0f7b1bb94b446d7977c",
            "dateTime": "2019-12-10T20:00:41Z",
            "description": "Description",
            "amount": -4564,
            "balance": 14630,
            "bank": "Mono Bank"
        }
    ]
}
``` 

#### GET /transactions?category={category_id}
Returns all transactions related to the specified category.  

e.g request
```
curl -X GET \
  'http://localhost/transactions?category=5e146165f22ead680880b21a' \
  -H 'Cache-Control: no-cache' \
```

e.g response:
```
{
    "transactions": [
        {
            "_id": "5e19bf98b1bb94b446d7977a",
            "dateTime": "2019-12-10T20:00:41Z",
            "description": "Vodafone\n+380935684777",
            "category": {
                "id": "5e146165f22ead680880b21a",
                "name": "Shopping",
                "normalizedName": "SHOPPING"
            },
            "amount": -4564,
            "balance": 14630,
            "bank": "Mono Bank"
        }
    ]
}
```   

#### PATCH /transactions/{transactionID}
Sets the category for specify transaction.  

e.g request
```
curl -X PATCH \
  http://localhost/transactions/5e19bf98b1bb94b446d7977a \
  -H 'cache-control: no-cache' \
  -d '{
    "categoryId": "5e148c35332c1e86b27b2ce0"
}'
```

e.g response:
```
{
    "message": "Success"
}
```   
## Current status

Project is being actively developed by @suddengunter and @markelrep. At the moment only monobank webhooks are supported as incoming source, and web UI only has one page - we keep it simple and add only things we need.

## Contributors guide

All contributions are welcomed through pull requests


## Frontend

see frontend at https://github.com/lungria/spendshelf-ui
