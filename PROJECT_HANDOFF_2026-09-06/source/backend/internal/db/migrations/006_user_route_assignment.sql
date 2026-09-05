-- Store the concrete route-exit binding selected for a user.  NULL means the
-- user follows the route's enabled weighted exit-IP pool.
ALTER TABLE user_routes ADD COLUMN route_exit_ip_id TEXT REFERENCES route_exit_ips(id);

CREATE INDEX IF NOT EXISTS idx_user_routes_active
  ON user_routes(user_id, is_primary, active_to);

CREATE INDEX IF NOT EXISTS idx_user_routes_route_exit_ip_id
  ON user_routes(route_exit_ip_id);
