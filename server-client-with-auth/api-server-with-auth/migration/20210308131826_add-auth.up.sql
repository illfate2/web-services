CREATE TYPE provider_type AS ENUM ('google', 'github', 'discord');

CREATE TABLE auth (
    id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    provider provider_type NOT NULL,
    provider_id text NOT NULL,
    user_id int NOT NULL REFERENCES users(id)
);


CREATE UNIQUE INDEX auth_provider_unique_idx ON auth(provider, provider_id) ;