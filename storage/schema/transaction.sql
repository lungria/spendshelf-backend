CREATE TABLE transaction (
    "ID" varchar(50) PRIMARY KEY,
    "time" timestamp NOT NULL,
    "description" varchar(50) NOT NULL,
    "mcc" integer NOT NULL,
    "hold" boolean NOT NULL,
    "amount" bigint NOT NULL,
    "accountID" varchar(50) NOT NULL,
    "categoryID" integer references category("ID") ON DELETE RESTRICT NOT NULL,
    "lastUpdatedAt" timestamp NOT NULL
);
