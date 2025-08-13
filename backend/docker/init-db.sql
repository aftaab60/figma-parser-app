-- The parser_db database is already created by POSTGRES_DB env var
-- Set basic configuration
SET client_encoding = 'UTF8';

SET timezone = 'UTC';

-- Table for storing information about parsed Figma files
CREATE TABLE IF NOT EXISTS figma_files (
    id SERIAL PRIMARY KEY,
    name VARCHAR(500) NOT NULL, -- Increased for long Figma file names
    url TEXT, -- TEXT for long URLs
    file_key VARCHAR(255) NOT NULL, -- File keys are typically shorter (removed UNIQUE to allow multiple parses)
    image_url TEXT, -- TEXT for long thumbnail URLs
    thumbnails TEXT, -- Storing as TEXT as it's a simple JSON string
    canvas_width DOUBLE PRECISION,
    canvas_height DOUBLE PRECISION,
    parsed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN DEFAULT TRUE NOT NULL
);

-- Table for storing extracted Figma components
CREATE TABLE IF NOT EXISTS components (
    id SERIAL PRIMARY KEY,
    figma_file_id INTEGER NOT NULL,
    node_id VARCHAR(100) NOT NULL, -- Node IDs are typically shorter (removed UNIQUE)
    name VARCHAR(500) NOT NULL, -- Increased for long component names
    type VARCHAR(50) NOT NULL,
    description TEXT,
    x DOUBLE PRECISION,
    y DOUBLE PRECISION,
    width DOUBLE PRECISION,
    height DOUBLE PRECISION,
    properties JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN DEFAULT TRUE NOT NULL,
    CONSTRAINT fk_components_figma_file FOREIGN KEY (figma_file_id) REFERENCES figma_files (id) ON DELETE CASCADE,
    CONSTRAINT unique_component_per_file UNIQUE (figma_file_id, node_id) -- Composite unique constraint
);

-- Table for storing instances of components
CREATE TABLE IF NOT EXISTS instances (
    id SERIAL PRIMARY KEY,
    component_id INTEGER NOT NULL,
    node_id VARCHAR(100) NOT NULL, -- Node IDs are typically shorter (removed UNIQUE)
    name VARCHAR(500) NOT NULL, -- Increased for long instance names
    x DOUBLE PRECISION,
    y DOUBLE PRECISION,
    width DOUBLE PRECISION,
    height DOUBLE PRECISION,
    properties JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN DEFAULT TRUE NOT NULL,
    CONSTRAINT fk_instances_component FOREIGN KEY (component_id) REFERENCES components (id) ON DELETE CASCADE,
    CONSTRAINT unique_instance_per_component UNIQUE (component_id, node_id) -- Composite unique constraint
);

-- Add indexes for faster lookups on foreign keys and soft delete columns
CREATE INDEX IF NOT EXISTS idx_components_figma_file_id ON components (figma_file_id);

CREATE INDEX IF NOT EXISTS idx_instances_component_id ON instances (component_id);