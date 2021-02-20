alter table category
add column visible bool;

update category
set visible = true;

alter table category
alter column visible set not null;

create table budget (
    "ID" integer primary key,
    "startsAt" timestamp not null,
    "endsAt" timestamp not null,
	"createdAt" timestamp not null,
	unique ("startsAt", "endsAt")
);

create table limit (
    "budgetID" integer not null references budget("ID"),
    "categoryID" integer not null references category("ID"),
    "amount" integer not null,
    primary key ("budgetID", "categoryID")
);
