make migrate-create NAME=add_sessions

atlas migrate hash --dir "file://migrations"

make migrate-up