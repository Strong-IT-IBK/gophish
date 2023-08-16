
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE targets ADD COLUMN "department" varchar(255);
ALTER TABLE targets ADD COLUMN "dep_number" varchar(255);
ALTER TABLE targets ADD COLUMN "age" varchar(255);
ALTER TABLE targets ADD COLUMN "gender" varchar(255);
ALTER TABLE targets ADD COLUMN "site" varchar(255);
ALTER TABLE targets ADD COLUMN "phone" varchar(255);
ALTER TABLE targets ADD COLUMN "degree" varchar(255);
ALTER TABLE targets ADD COLUMN "desc" varchar(255);

ALTER TABLE results ADD COLUMN "department" VARCHAR(255);
ALTER TABLE results ADD COLUMN "dep_number" VARCHAR(255);
ALTER TABLE results ADD COLUMN "age" VARCHAR(255);
ALTER TABLE results ADD COLUMN "gender" VARCHAR(255);
ALTER TABLE results ADD COLUMN "site" VARCHAR(255);
ALTER TABLE results ADD COLUMN "phone" VARCHAR(255);
ALTER TABLE results ADD COLUMN "degree" VARCHAR(255);
ALTER TABLE results ADD COLUMN "desc" VARCHAR(255);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

