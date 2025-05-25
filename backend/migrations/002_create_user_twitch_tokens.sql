-- Migration: 002_create_user_twitch_tokens.sql
-- Description: Create table for storing encrypted Twitch OAuth tokens

-- User Twitch Tokens table - stores encrypted OAuth tokens
CREATE TABLE user_twitch_tokens (
    clerk_user_id VARCHAR(255) PRIMARY KEY, -- From Clerk JWT, specific to the Clerk environment
    twitch_user_id VARCHAR(255) NOT NULL UNIQUE,
    encrypted_access_token TEXT NOT NULL,
    encrypted_refresh_token TEXT NOT NULL,
    scopes TEXT, -- Comma-separated list of granted scopes
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_user_twitch_tokens_twitch_user_id ON user_twitch_tokens(twitch_user_id);
CREATE INDEX idx_user_twitch_tokens_expires_at ON user_twitch_tokens(expires_at);

-- Create trigger to update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_twitch_tokens_updated_at
BEFORE UPDATE ON user_twitch_tokens
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column(); 