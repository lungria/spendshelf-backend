# spendshelf-backend

![main](https://github.com/lungria/spendshelf-backend/workflows/main/badge.svg)

## Current status

Project is paused indefinitely. I might get back to it later. Reasons:

- my current bank does not expose any API's for implementing automated transactions export
- budget is configured via manual DB calls each month
- front-end that uses current spendshelf API is written in old deprecated React version and it's too high maintanance of a framework for such a small project. Next FE version will be done in either HTMX or Elixir Phoenix

## Overview

Spendshelf started as a desire to keep money spending in the personal or/and family budget under control.
 
All incoming transactions are being sent to categories - and each category has a monthly budget (spent money limit). At the end of the budget period (usually at the end of the months or near to salary date) - you either close to monthly budget per each category (which is good) - or not (which is bad). 

## Current status

### Working
- Import transactions from bank
- Readonly API for categories
- API for transactions that allows to list transactions / set transaction category
- Basic monthly report
- Budget keeping (PoC) - we have read-only API that uses budget from DB. Budget is created manualy in the DB.
### Not working / Not Implemented
- Budget keeping: CRUD API for budget
- Automated transaction category matching (based on description/mcc/amount)
- Savings accounts tracking: monobank doesn't have any API for it, but we could keep track of those accounts using some internal logic.

## Contributors guide

Any contribution is welcome through pull requests

## Frontend (iOS)

Native iOS app was being actively developed in a private repository, but at the moment it's abandoned. 

There are lots of hardcoded stuff in there. Contact [me](https://github.com/suddengunter) if you want to get access to the repository.
Preview screenshots of native iOS client:

![main](https://raw.githubusercontent.com/lungria/spendshelf-backend/main/.github/img/1.png)

## Frontend (Web)

React app is being slowly developed in a private repository, but at the moment there are lots of hardcoded stuff in there. Contact [me](https://github.com/suddengunter) if you want to get access to the repository.
Preview screenshots of React app:

![main](https://raw.githubusercontent.com/lungria/spendshelf-backend/main/.github/img/2.png)
