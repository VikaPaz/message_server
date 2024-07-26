-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


create table if not exists messages
(
    id         uuid default uuid_generate_v4() primary key,
    message    text,
    status     varchar(100),
    created_at timestamp,
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table messages;
drop extension "uuid-ossp";
-- +goose StatementEnd
