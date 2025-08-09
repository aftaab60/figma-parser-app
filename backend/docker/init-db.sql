-- Simple database initialization for health checks
-- The parser_db database is already created by POSTGRES_DB env var

-- Set basic configuration
SET client_encoding = 'UTF8';

SET timezone = 'UTC';

-- Simple confirmation that database is ready
SELECT 'Database parser_db is ready for connections' as status;