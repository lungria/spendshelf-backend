CREATE TABLE account (
    "ID" varchar(50) PRIMARY KEY,
    "createdAt" timestamp NOT NULL,
    "description" varchar(50) NOT NULL,
    "balance" bigint NOT NULL,
    "currency" varchar(10) NOT NULL,
    "lastUpdatedAt" timestamp NOT NULL
);

ALTER TABLE transaction ADD FOREIGN KEY ("accountID") REFERENCES account("ID") ON DELETE RESTRICT;
