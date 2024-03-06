-- Active: 1709698391238@@127.0.0.1@5432@filesync@public
ALTER TABLE record
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN deleted_at TIMESTAMP;

-- Create a trigger to update updated_at column automatically
CREATE OR REPLACE FUNCTION update_records_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER records_update_updated_at
BEFORE UPDATE ON records
FOR EACH ROW
EXECUTE FUNCTION update_records_updated_at();