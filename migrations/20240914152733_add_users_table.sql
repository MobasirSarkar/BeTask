-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
   id UUID  PRIMARY KEY,
   google_id VARCHAR(255) UNIQUE NOT NULL,
   profile_pic_url TEXT,
   name VARCHAR(255),
   email VARCHAR (255) UNIQUE NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
