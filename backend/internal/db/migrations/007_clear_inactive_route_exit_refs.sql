-- Fixed exit bindings are operational state.  Once a user route is no longer
-- active, release the foreign-key reference so the binding can be retired;
-- active assignments remain protected and retain the selected binding.
UPDATE user_routes SET route_exit_ip_id = NULL WHERE active_to IS NOT NULL;
