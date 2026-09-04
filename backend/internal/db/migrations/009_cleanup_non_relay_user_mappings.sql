-- Business users are sourced only from relay-node Inbounds.  Older versions
-- could leave a user_inbounds row (and inbounds.user_id) on a landing or
-- otherwise non-relay node.  Remove only the invalid mapping/classification;
-- keep the users row because it is central business data and may be reused.
DELETE FROM user_inbounds
WHERE inbound_id IN (
  SELECT i.id
  FROM inbounds i
  JOIN nodes n ON n.id = i.node_id
  WHERE n.type <> 'relay' OR n.deleted_at IS NOT NULL
);

UPDATE inbounds
SET user_id = NULL,
    kind = 'infrastructure'
WHERE node_id IN (
  SELECT n.id
  FROM nodes n
  WHERE n.type <> 'relay' OR n.deleted_at IS NOT NULL
);
