-- Initialize PostgreSQL database
-- This script runs when the container is first created

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE portfolio_db TO portfolio;
