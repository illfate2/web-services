ALTER TABLE museum_items
ADD COLUMN created_by INT REFERENCES users(id);

ALTER TABLE museum_funds
ADD COLUMN created_by INT REFERENCES users(id);

ALTER TABLE museum_item_sets
ADD COLUMN created_by INT REFERENCES users(id);

ALTER TABLE museum_item_movements
ADD COLUMN created_by INT REFERENCES users(id);

