CREATE TABLE account (
    "ID" varchar(50) PRIMARY KEY,
    "createdAt" timestamp NOT NULL,
    "description" varchar(50) NOT NULL,
    "balance" bigint NOT NULL,
    "currency" varchar(10) NOT NULL,
    "lastUpdatedAt" timestamp NOT NULL
);

CREATE TABLE category (
    "ID" integer PRIMARY KEY,
    "createdAt" timestamp NOT NULL,
    "name" varchar(25) NOT NULL,
    "logo" varchar(30) NOT NULL
);

CREATE TABLE transaction (
    "ID" varchar(50) PRIMARY KEY,
    "time" timestamp NOT NULL,
    "description" varchar(50) NOT NULL,
    "mcc" integer NOT NULL,
    "hold" boolean NOT NULL,
    "amount" bigint NOT NULL,
    "accountID" varchar(50) NOT NULL references account("ID") ON DELETE RESTRICT NOT NULL,
    "categoryID" integer references category("ID") ON DELETE RESTRICT NOT NULL,
    "lastUpdatedAt" timestamp NOT NULL,
    "comment" varchar(50)
);

CREATE TABLE budget (
    "ID" integer PRIMARY KEY,
    "startsAt" timestamp NOT NULL,
    "endsAt" timestamp NOT NULL,
	"createdAt" timestamp NOT NULL,
	unique ("startsAt", "endsAt")
);

CREATE TABLE "limit" (
    "budgetID" integer NOT NULL references budget("ID"),
    "categoryID" integer NOT NULL references category("ID"),
    "amount" integer NOT NULL,
    PRIMARY KEY ("budgetID", "categoryID")
);


/* Insert indexes */

BEGIN;

CREATE INDEX accountID_time_idx on transaction("accountID", "time")

CREATE INDEX time_categoryID_idx ON transaction ("time", "categoryID");

CREATE INDEX startsAt_idx ON budget ("startsAt");

CREATE INDEX createdAt_idx ON account ("createdAt");

CREATE INDEX budgetID_amount_idx ON "limit" ("budgetID_amount_idx");

COMMIT;