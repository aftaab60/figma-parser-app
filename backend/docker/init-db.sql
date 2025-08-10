-- The parser_db database is already created by POSTGRES_DB env var
-- Set basic configuration
SET client_encoding = 'UTF8';

SET timezone = 'UTC';

-- Table for storing information about parsed Figma files
CREATE TABLE IF NOT EXISTS figma_files (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) UNIQUE NOT NULL,
    file_key VARCHAR(255) UNIQUE NOT NULL,
    image_url VARCHAR(255),
    thumbnails TEXT, -- Storing as TEXT as it's a simple JSON string
    canvas_width DOUBLE PRECISION,
    canvas_height DOUBLE PRECISION,
    parsed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    active BOOLEAN DEFAULT TRUE NOT NULL
);

-- Table for storing extracted Figma components
CREATE TABLE IF NOT EXISTS components (
    id SERIAL PRIMARY KEY,
    figma_file_id INTEGER NOT NULL,
    node_id VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    description TEXT,
    x DOUBLE PRECISION,
    y DOUBLE PRECISION,
    width DOUBLE PRECISION,
    height DOUBLE PRECISION,
    z_index INTEGER, -- to control stacking order
    properties JSONB,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    active BOOLEAN DEFAULT TRUE NOT NULL,
    CONSTRAINT fk_components_figma_file FOREIGN KEY (figma_file_id) REFERENCES figma_files (id) ON DELETE CASCADE
);

-- Table for storing instances of components
CREATE TABLE IF NOT EXISTS instances (
    id SERIAL PRIMARY KEY,
    component_id INTEGER NOT NULL,
    node_id VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    x DOUBLE PRECISION,
    y DOUBLE PRECISION,
    width DOUBLE PRECISION,
    height DOUBLE PRECISION,
    z_index INTEGER, -- to control stacking order
    properties JSONB,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    active BOOLEAN DEFAULT TRUE NOT NULL,
    CONSTRAINT fk_instances_component FOREIGN KEY (component_id) REFERENCES components (id) ON DELETE CASCADE
);

-- Add indexes for faster lookups on foreign keys and soft delete columns
CREATE INDEX IF NOT EXISTS idx_components_figma_file_id ON components (figma_file_id);

CREATE INDEX IF NOT EXISTS idx_instances_component_id ON instances (component_id);