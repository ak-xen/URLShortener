CREATE TABLE urls (
                      id SERIAL PRIMARY KEY,
                      original_url TEXT NOT NULL,
                      short_code VARCHAR(20) UNIQUE NOT NULL,
                      created_at TIMESTAMP NOT NULL DEFAULT now()
);


CREATE INDEX idx_short_code ON urls(short_code);