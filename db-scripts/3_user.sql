create role merch_app with login password '426643';

grant all privileges on all tables    in schema public to merch_app;
grant all privileges on all sequences in schema public to merch_app;
grant all privileges on all functions in schema public to merch_app;
