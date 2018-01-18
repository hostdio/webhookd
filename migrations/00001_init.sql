-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE webhooks (
  id              VARCHAR(255),
  url             VARCHAR(2048),
  header_override JSON,

  PRIMARY KEY(id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.


DROP TABLE webhooks;
