
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

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

