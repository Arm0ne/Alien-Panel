package httpapi

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func exitIPAllocatedUserPredicate(exitIPReference string) string {
	return fmt.Sprintf(`(p.exit_ip_id = %s OR EXISTS (
  SELECT 1 FROM user_path_exit_ips upi WHERE upi.user_path_id = p.id AND upi.exit_ip_id = %s
)) AND p.active_to IS NULL AND u.deleted_at IS NULL AND u.status <> 'disabled'
AND datetime(u.created_at) <= datetime('now')
AND (u.expiry_time IS NULL OR datetime(u.expiry_time) >= datetime('now'))`, exitIPReference, exitIPReference)
}

func (s *Server) exitIPUsers(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusNotFound, notFoundCode, "exit IP not found")
		return
	}
	var exists int
	if err := s.db.QueryRow(`SELECT 1 FROM exit_ips WHERE id = ?`, id).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeFailure(w, http.StatusNotFound, notFoundCode, "exit IP not found")
			return
		}
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read exit IP")
		return
	}

	query := parseListQuery(r)
	where := []string{exitIPAllocatedUserPredicate("?")}
	args := []any{id, id}
	if query.keyword != "" {
		where = append(where, `(u.display_name LIKE ? OR n.name LIKE ? OR COALESCE(i.tag, '') LIKE ?)`)
		like := "%" + query.keyword + "%"
		args = append(args, like, like, like)
	}
	base := `FROM user_paths p
JOIN users u ON u.id = p.user_id
LEFT JOIN nodes n ON n.id = p.relay_node_id
LEFT JOIN user_inbounds ui ON ui.user_id = u.id AND ui.is_primary = 1 AND ui.active_to IS NULL
LEFT JOIN inbounds i ON i.id = ui.inbound_id AND i.deleted_at IS NULL
WHERE ` + strings.Join(where, " AND ")
	var total int
	if err := s.db.QueryRow(`SELECT COUNT(DISTINCT u.id) `+base, args...).Scan(&total); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count exit IP users")
		return
	}
	rows, err := s.db.Query(`SELECT u.id, COALESCE(NULLIF(TRIM(u.display_name), ''), u.id),
COALESCE(p.relay_node_id, ''), COALESCE(n.name, ''), COALESCE(i.tag, ''), COALESCE(p.mode, ''),
COALESCE(u.status, 'unknown'), COALESCE(u.billing_type, 'paid'), COALESCE(u.expiry_time, '')
`+base+` ORDER BY CASE WHEN u.expiry_time IS NULL THEN 1 ELSE 0 END, u.expiry_time ASC, u.display_name ASC, u.id ASC LIMIT ? OFFSET ?`,
		append(args, query.pageSize, query.offset)...)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read exit IP users")
		return
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var userID, name, nodeID, nodeName, inboundTag, pathMode, status, billingType, expiresAt string
		if err := rows.Scan(&userID, &name, &nodeID, &nodeName, &inboundTag, &pathMode, &status, &billingType, &expiresAt); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not decode exit IP users")
			return
		}
		items = append(items, map[string]any{
			"id": userID, "name": name, "nodeId": nodeID, "nodeName": nodeName,
			"inboundTag": nullableString(inboundTag), "pathMode": nullableString(pathMode),
			"status": status, "billingType": billingType, "expiresAt": nullableString(expiresAt),
		})
	}
	if err := rows.Err(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not finish reading exit IP users")
		return
	}
	writeSuccess(w, s.pageResponse(items, total, query))
}
