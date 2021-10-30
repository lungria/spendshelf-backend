BEGIN;

CREATE INDEX accountID_time_idx on transaction("accountID", "time")

CREATE INDEX time_categoryID_idx ON transaction ("time", "categoryID");

CREATE INDEX startsAt_idx ON budget ("startsAt");

CREATE INDEX createdAt_idx ON account ("createdAt");

CREATE INDEX budgetID_amount_idx ON "limit" ("budgetID_amount_idx");

COMMIT;