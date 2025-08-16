-- Drop devices table and related objects
DROP TRIGGER IF EXISTS update_devices_updated_at_trigger ON devices;
DROP FUNCTION IF EXISTS update_devices_updated_at();
DROP TABLE IF EXISTS devices;