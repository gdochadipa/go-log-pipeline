-- create enum type in postgre
CREATE TYPE level AS ENUM ('INFO', 'WARN', 'ERROR','DEBUG');

-- create table
CREATE TABLE logs (
    id BIGSERIAL,
    timestamp TIMESTAMPTZ NOT NULL,
    service VARCHAR(255),
    level level not null,
    message TEXT,
    metadata JSONB,
    PRIMARY KEY (id, timestamp)
) PARTITION BY RANGE (timestamp);


/* call extension pg_partman, there are several extension on postgree */
CREATE EXTENSION pg_partman;

/* check is the extention already exists */
SELECT * FROM pg_extension WHERE extname = 'pg_partman';

/* check is partman already exists in pg_namespace */
SELECT nspname FROM pg_namespace WHERE nspname = 'partman';

/* create table partion/sharding settings */
SELECT public.create_parent(
    p_parent_table := 'public.logs',
    p_control := 'timestamp',
    p_interval := '1 day',
    p_premake := 4
);

/* create indexing */
create index idx_logs_service on logs (service);
create index idx_logs_level on logs (level);

/* indexing using GIN, because JSONB and and text typedata have spesific method */
/* If you need to search for keywords within the message field, you can create a full-text search index. */
create index idx_logs_metadata_gin on logs USING GIN (metadata);

create index idx_logs_message_fts on logs using GIN (to_tsvector('english',message));

/*
 * data ingestion
 *
 * using copy for bulking data
 *  */
