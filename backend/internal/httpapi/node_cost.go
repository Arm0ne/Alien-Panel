package httpapi

import (
	"database/sql"
	"errors"
	"math"
	"net/http"
	"strings"
	"time"
)

// nodeCostRequest describes a central-owned temporal cost record. Dates are
// deliberately separate from Agent/X-Panel data and are stored as YYYY-MM-DD.
type nodeCostRequest struct {
	Category      *string  `json:"category"`
	MonthlyAmount *float64 `json:"monthlyAmount"`
	Currency      *string  `json:"currency"`
	EffectiveFrom *string  `json:"effectiveFrom"`
	EffectiveTo   *string  `json:"effectiveTo"`
	Notes         *string  `json:"notes"`
}

type nodeCostRecord struct {
	ID            string
	NodeID        string
	NodeName      string
	Category      string
	MonthlyAmount float64
	Currency      string
	EffectiveFrom string
	EffectiveTo   string
	Notes         string
	CreatedAt     string
}

type nodeCostValues struct {
	Category      string
	MonthlyAmount float64
	Currency      string
	EffectiveFrom string
	EffectiveTo   string
	Notes         string
}

func normalizeCostDate(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", nil
	}
	if parsed, err := time.Parse("2006-01-02", value); err == nil {
		return parsed.Format("2006-01-02"), nil
	}
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed.UTC().Format("2006-01-02"), nil
	}
	return "", errors.New("invalid date")
}

func validateNodeCostValues(values nodeCostValues, requireFrom bool) error {
	if strings.TrimSpace(values.Category) == "" {
		return errors.New("category is required")
	}
	if len(values.Category) > 120 {
		return errors.New("category is too long")
	}
	if math.IsNaN(values.MonthlyAmount) || math.IsInf(values.MonthlyAmount, 0) || values.MonthlyAmount < 0 || values.MonthlyAmount > 100000000 {
		return errors.New("monthlyAmount must be a non-negative number")
	}
	if values.Currency != "CNY" {
		return errors.New("currency currently supports CNY only")
	}
	if requireFrom && values.EffectiveFrom == "" {
		return errors.New("effectiveFrom is required")
	}
	from, err := normalizeCostDate(values.EffectiveFrom)
	if err != nil || from == "" {
		return errors.New("effectiveFrom must be YYYY-MM-DD or RFC3339")
	}
	to, err := normalizeCostDate(values.EffectiveTo)
	if err != nil {
		return errors.New("effectiveTo must be YYYY-MM-DD or RFC3339")
	}
	if to != "" && to < from {
		return errors.New("effectiveTo must not be earlier than effectiveFrom")
	}
	if len(values.Notes) > 2000 {
		return errors.New("node cost notes are too long")
	}
	return nil
}

func nodeCostValuesFromRequest(request nodeCostRequest, existing *nodeCostRecord) (nodeCostValues, error) {
	values := nodeCostValues{Currency: "CNY"}
	if existing != nil {
		values = nodeCostValues{
			Category: existing.Category, MonthlyAmount: existing.MonthlyAmount, Currency: existing.Currency,
			EffectiveFrom: existing.EffectiveFrom, EffectiveTo: existing.EffectiveTo, Notes: existing.Notes,
		}
	}
	if request.Category != nil {
		values.Category = strings.TrimSpace(*request.Category)
	}
	if request.MonthlyAmount != nil {
		values.MonthlyAmount = *request.MonthlyAmount
	}
	if request.Currency != nil {
		values.Currency = strings.TrimSpace(*request.Currency)
	}
	if request.EffectiveFrom != nil {
		values.EffectiveFrom = strings.TrimSpace(*request.EffectiveFrom)
	}
	if request.EffectiveTo != nil {
		values.EffectiveTo = strings.TrimSpace(*request.EffectiveTo)
	}
	if request.Notes != nil {
		values.Notes = strings.TrimSpace(*request.Notes)
	}
	var err error
	if values.EffectiveFrom, err = normalizeCostDate(values.EffectiveFrom); err != nil {
		return nodeCostValues{}, errors.New("effectiveFrom must be YYYY-MM-DD or RFC3339")
	}
	if values.EffectiveTo, err = normalizeCostDate(values.EffectiveTo); err != nil {
		return nodeCostValues{}, errors.New("effectiveTo must be YYYY-MM-DD or RFC3339")
	}
	return values, validateNodeCostValues(values, true)
}

func nodeCostRecordData(record nodeCostRecord) map[string]any {
	return map[string]any{
		"id": record.ID, "nodeId": record.NodeID, "nodeName": record.NodeName,
		"category": record.Category, "monthlyAmount": record.MonthlyAmount, "currency": record.Currency,
		"effectiveFrom": record.EffectiveFrom, "effectiveTo": nullableString(record.EffectiveTo),
		"notes": nullableString(record.Notes), "createdAt": record.CreatedAt,
	}
}

func (s *Server) readNodeCost(id string) (nodeCostRecord, error) {
	var record nodeCostRecord
	err := s.db.QueryRow(`SELECT c.id, c.node_id, n.name, c.category, c.monthly_amount, c.currency,
 c.effective_from, COALESCE(c.effective_to, ''), COALESCE(c.notes, ''), c.created_at
FROM node_costs c JOIN nodes n ON n.id = c.node_id AND n.deleted_at IS NULL WHERE c.id = ?`, id).Scan(
		&record.ID, &record.NodeID, &record.NodeName, &record.Category, &record.MonthlyAmount, &record.Currency,
		&record.EffectiveFrom, &record.EffectiveTo, &record.Notes, &record.CreatedAt,
	)
	return record, err
}

func (s *Server) readNodeCosts(nodeID string) ([]nodeCostRecord, error) {
	rows, err := s.db.Query(`SELECT c.id, c.node_id, n.name, c.category, c.monthly_amount, c.currency,
 c.effective_from, COALESCE(c.effective_to, ''), COALESCE(c.notes, ''), c.created_at
FROM node_costs c JOIN nodes n ON n.id = c.node_id AND n.deleted_at IS NULL WHERE c.node_id = ?
ORDER BY c.effective_from DESC, c.created_at DESC, c.category ASC`, nodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]nodeCostRecord, 0)
	for rows.Next() {
		var record nodeCostRecord
		if err := rows.Scan(&record.ID, &record.NodeID, &record.NodeName, &record.Category, &record.MonthlyAmount, &record.Currency,
			&record.EffectiveFrom, &record.EffectiveTo, &record.Notes, &record.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, record)
	}
	return items, rows.Err()
}

func (s *Server) createNodeCost(w http.ResponseWriter, r *http.Request) {
	nodeID := strings.TrimSpace(r.PathValue("id"))
	if nodeID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "node id is required")
		return
	}
	var nodeName string
	if err := s.db.QueryRow(`SELECT name FROM nodes WHERE id = ? AND deleted_at IS NULL`, nodeID).Scan(&nodeName); errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "node not found")
		return
	} else if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node")
		return
	}
	var request nodeCostRequest
	if err := decodeJSON(r, &request); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid node cost payload")
		return
	}
	values, err := nodeCostValuesFromRequest(request, nil)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	id := newID()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := s.db.Exec(`INSERT INTO node_costs (id, node_id, category, monthly_amount, currency, effective_from, effective_to, notes, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, id, nodeID, values.Category, values.MonthlyAmount, values.Currency, values.EffectiveFrom,
		nullableRouteValue(values.EffectiveTo), nullableRouteValue(values.Notes), now); err != nil {
		s.logger.Error("create node cost", "node_id", nodeID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not create node cost")
		return
	}
	record := nodeCostRecord{ID: id, NodeID: nodeID, NodeName: nodeName, Category: values.Category, MonthlyAmount: values.MonthlyAmount,
		Currency: values.Currency, EffectiveFrom: values.EffectiveFrom, EffectiveTo: values.EffectiveTo, Notes: values.Notes, CreatedAt: now}
	s.writeAuditLog(r, "node_cost.create", "node_cost", id, nil, nodeCostRecordData(record))
	writeSuccess(w, nodeCostRecordData(record))
}

func (s *Server) listNodeCosts(w http.ResponseWriter, r *http.Request) {
	nodeID := strings.TrimSpace(r.PathValue("id"))
	if nodeID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "node id is required")
		return
	}
	var exists int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM nodes WHERE id = ? AND deleted_at IS NULL`, nodeID).Scan(&exists); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node")
		return
	}
	if exists == 0 {
		writeFailure(w, http.StatusNotFound, notFoundCode, "node not found")
		return
	}
	costs, err := s.readNodeCosts(nodeID)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node costs")
		return
	}
	items := make([]map[string]any, 0, len(costs))
	for _, cost := range costs {
		items = append(items, nodeCostRecordData(cost))
	}
	writeSuccess(w, items)
}

func (s *Server) updateNodeCost(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("costId"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "node cost id is required")
		return
	}
	existing, err := s.readNodeCost(id)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "node cost not found")
		return
	}
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read node cost")
		return
	}
	if nodeID := strings.TrimSpace(r.PathValue("id")); nodeID == "" || nodeID != existing.NodeID {
		writeFailure(w, http.StatusNotFound, notFoundCode, "node cost not found")
		return
	}
	var request nodeCostRequest
	if err := decodeJSON(r, &request); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid node cost payload")
		return
	}
	values, err := nodeCostValuesFromRequest(request, &existing)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	// Effective dates define the historical version. Create a new record for a
	// date change so finance reports never rewrite an earlier period's meaning.
	if values.EffectiveFrom != existing.EffectiveFrom || values.EffectiveTo != existing.EffectiveTo {
		writeFailure(w, http.StatusConflict, validationCode, "effective dates are immutable; create a new node cost record")
		return
	}
	if _, err := s.db.Exec(`UPDATE node_costs SET category = ?, monthly_amount = ?, currency = ?, notes = ? WHERE id = ?`,
		values.Category, values.MonthlyAmount, values.Currency, nullableRouteValue(values.Notes), id); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not update node cost")
		return
	}
	updated := existing
	updated.Category, updated.MonthlyAmount, updated.Currency, updated.Notes = values.Category, values.MonthlyAmount, values.Currency, values.Notes
	s.writeAuditLog(r, "node_cost.update", "node_cost", id, nodeCostRecordData(existing), nodeCostRecordData(updated))
	writeSuccess(w, nodeCostRecordData(updated))
}
