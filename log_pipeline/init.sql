\echo 'Creating database logpipeline...'
CREATE DATABASE logpipeline;
\echo 'Creating user postgres...'
CREATE USER postgres WITH PASSWORD 'postgres';
\echo 'Granting privileges...'
GRANT ALL PRIVILEGES ON DATABASE logpipeline TO postgres;
\echo 'Initialization complete.'
