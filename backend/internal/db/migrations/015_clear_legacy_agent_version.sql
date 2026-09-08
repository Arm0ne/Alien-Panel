-- "online" was a one-time installer label, not an Agent binary version.
-- Clear it so the first heartbeat from an updated Agent writes its embedded
-- build identity without presenting the old label as an actual version.
UPDATE nodes SET agent_version = NULL WHERE agent_version = 'online';
