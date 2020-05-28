--
-- Simple init sql file. migrate should be used as future improvement
--

CREATE TYPE reservation_status AS ENUM ('RESERVED', 'BACKORDER ', 'PENDING');

--
-- STOCK TABLE FOR PHYSICAL INVENTORY
--

CREATE TABLE product_inventory(
  product_id CHARACTER varying(16) PRIMARY KEY, --length based on back end product example. To define larger is required
  quantity INTEGER DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

--
-- RESERVATION TABLES FOR CUSTOMER RESERVATIONS 
--
CREATE TABLE reservation(
  id SERIAL PRIMARY KEY,
  status reservation_status,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reservation_line(
  reservation_id INTEGER REFERENCES reservation(id) DEFERRABLE INITIALLY IMMEDIATE,
  product_id CHARACTER varying(16) REFERENCES product_inventory(product_id),
  quantity INTEGER DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (reservation_id, product_id)
);

CREATE INDEX reservation_line_reservation_id_idx ON reservation_line(reservation_id);

--
--  AVAILABLE STOCK TABLE
--
CREATE TABLE product_inventory_availability(
  product_id CHARACTER varying(16) PRIMARY KEY REFERENCES product_inventory(product_id),
  quantity INTEGER DEFAULT 0,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

--
--  TEMP NEXT LINES FOR TEST PURPOSE - SHOULD BE REMOVED FOR PRODUCT
--
-- Create Products inventory stock with 10 of each products.
insert into product_inventory (product_id, quantity) VALUES('PIPR-JACKET-SIZM', 10);
insert into product_inventory (product_id, quantity) VALUES('PIPR-JACKET-SIZL', 10);
insert into product_inventory (product_id, quantity) VALUES('PIPR-JACKET-SIXL', 10);
insert into product_inventory (product_id, quantity) VALUES('PIPR-SVLMUG-GREN', 10);
insert into product_inventory (product_id, quantity) VALUES('PIPR-SVLMUG-YLOW', 10);
insert into product_inventory (product_id, quantity) VALUES('PIPR-MOSPAD-0000', 10);
insert into product_inventory (product_id, quantity) VALUES('PIPR-SMFRDG-0000', 10);
-- As a result, and for consistency, Available products should match product initial stocks
insert into product_inventory_availability (product_id, quantity) VALUES('PIPR-JACKET-SIZM', 10);
insert into product_inventory_availability (product_id, quantity) VALUES('PIPR-JACKET-SIZL', 10);
insert into product_inventory_availability (product_id, quantity) VALUES('PIPR-JACKET-SIXL', 10);
insert into product_inventory_availability (product_id, quantity) VALUES('PIPR-SVLMUG-GREN', 10);
insert into product_inventory_availability (product_id, quantity) VALUES('PIPR-SVLMUG-YLOW', 10);
insert into product_inventory_availability (product_id, quantity) VALUES('PIPR-MOSPAD-0000', 10);
insert into product_inventory_availability (product_id, quantity) VALUES('PIPR-SMFRDG-0000', 10);