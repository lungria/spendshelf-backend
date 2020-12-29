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
    "comment" varchar(50) /* todo: migration */
);