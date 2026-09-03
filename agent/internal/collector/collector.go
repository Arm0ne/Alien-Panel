package collector

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"xpanel-central/agent/internal/xpanel"
)

type Status struct {
	XrayRunning   bool    `json:"xray_running"`
	XrayVersion   string  `json:"xray_version,omitempty"`
	XPanelVersion string  `json:"xpanel_version,omitempty"`
	CPUUsage      float64 `json:"cpu_usage,omitempty"`
	MemoryUsed    int64   `json:"memory_used,omitempty"`
	MemoryTotal   int64   `json:"memory_total,omitempty"`
	DiskUsed      int64   `json:"disk_used,omitempty"`
	DiskTotal     int64   `json:"disk_total,omitempty"`
}

type Client struct {
	RemoteID   string `json:"remote_id"`
	Email      string `json:"email"`
	Enable     bool   `json:"enable"`
	ExpiryTime int64  `json:"expiry_time,omitempty"`
	Up         int64  `json:"up"`
	Down       int64  `json:"down"`
	AllTime    int64  `json:"all_time"`
	LastOnline int64  `json:"last_online,omitempty"`
}

type Inbound struct {
	RemoteID   int64    `json:"remote_id"`
	Tag        string   `json:"tag,omitempty"`
	Remark     string   `json:"remark,omitempty"`
	Protocol   string   `json:"protocol,omitempty"`
	Port       int64    `json:"port,omitempty"`
	Listen     string   `json:"listen,omitempty"`
	Enable     bool     `json:"enable"`
	ExpiryTime int64    `json:"expiry_time,omitempty"`
	Up         int64    `json:"up"`
	Down       int64    `json:"down"`
	AllTime    int64    `json:"all_time"`
	ConfigHash string   `json:"config_hash,omitempty"`
	Clients    []Client `json:"clients,omitempty"`
}

type Snapshot struct {
	NodeKey    string    `json:"node_key"`
	SyncID     string    `json:"sync_id"`
	ObservedAt string    `json:"observed_at"`
	Status     Status    `json:"status"`
	Inbounds   []Inbound `json:"inbounds"`
}

type Heartbeat struct {
	NodeKey    string `json:"node_key"`
	ObservedAt string `json:"observed_at"`
	Status     Status `json:"status"`
}

func (snapshot Snapshot) Heartbeat() Heartbeat {
	return Heartbeat{NodeKey: snapshot.NodeKey, ObservedAt: snapshot.ObservedAt, Status: snapshot.Status}
}

type Collector struct {
	client   *xpanel.Client
	nodeKey  string
	syncIDFn func(time.Time) string
}

func New(client *xpanel.Client, nodeKey string) (*Collector, error) {
	if client == nil {
		return nil, errors.New("xpanel client is required")
	}
	if strings.TrimSpace(nodeKey) == "" {
		return nil, errors.New("node key is required")
	}
	return &Collector{client: client, nodeKey: nodeKey, syncIDFn: newSyncID}, nil
}

func (collector *Collector) Collect(ctx context.Context) (Snapshot, error) {
	inboundResponse, err := collector.client.Get(ctx, "/inbounds/list")
	if err != nil {
		return Snapshot{}, fmt.Errorf("collect inbounds: %w", err)
	}
	inbounds, err := ParseInbounds(responsePayload(inboundResponse))
	if err != nil {
		return Snapshot{}, fmt.Errorf("parse inbounds: %w", err)
	}

	statusResponse, err := collector.client.Get(ctx, "/server/status")
	if err != nil {
		return Snapshot{}, fmt.Errorf("collect server status: %w", err)
	}
	status, err := ParseStatus(responsePayload(statusResponse))
	if err != nil {
		return Snapshot{}, fmt.Errorf("parse server status: %w", err)
	}

	now := time.Now().UTC()
	return Snapshot{
		NodeKey:    collector.nodeKey,
		SyncID:     collector.syncIDFn(now),
		ObservedAt: now.Format(time.RFC3339Nano),
		Status:     status,
		Inbounds:   inbounds,
	}, nil
}

func ParseInbounds(payload json.RawMessage) ([]Inbound, error) {
	items, err := listItems(payload, "inbounds", "items", "list")
	if err != nil {
		return nil, err
	}
	inbounds := make([]Inbound, 0, len(items))
	for _, item := range items {
		inbound, err := parseInbound(item)
		if err != nil {
			return nil, err
		}
		inbounds = append(inbounds, inbound)
	}
	return inbounds, nil
}

func ParseStatus(payload json.RawMessage) (Status, error) {
	object, err := objectFromPayload(payload)
	if err != nil {
		return Status{}, err
	}
	status := Status{
		XrayRunning:   boolField(object, "xray_running", "xrayRunning", "xray_status", "xrayStatus"),
		XrayVersion:   stringField(object, "xray_version", "xrayVersion"),
		XPanelVersion: stringField(object, "xpanel_version", "xPanelVersion", "panel_version", "panelVersion"),
		CPUUsage:      floatField(object, "cpu_usage", "cpuUsage", "cpu_percent", "cpuPercent", "cpu"),
		MemoryUsed:    intField(object, "memory_used", "memoryUsed", "memUsed"),
		MemoryTotal:   intField(object, "memory_total", "memoryTotal", "memTotal"),
		DiskUsed:      intField(object, "disk_used", "diskUsed"),
		DiskTotal:     intField(object, "disk_total", "diskTotal"),
	}
	if nested, ok := objectMap(object, "xray"); ok {
		if _, present := objectValue(object, "xray_running", "xrayRunning", "xray_status", "xrayStatus"); !present {
			status.XrayRunning = boolField(nested, "running", "state", "xray_running")
		}
		if status.XrayVersion == "" {
			status.XrayVersion = stringField(nested, "version", "xray_version")
		}
	}
	if status.MemoryUsed == 0 || status.MemoryTotal == 0 {
		if nested, ok := objectMap(object, "mem", "memory", "memory_usage", "memoryUsage"); ok {
			if status.MemoryUsed == 0 {
				status.MemoryUsed = intField(nested, "current", "used", "memory_used", "memoryUsed")
			}
			if status.MemoryTotal == 0 {
				status.MemoryTotal = intField(nested, "total", "memory_total", "memoryTotal")
			}
		}
	}
	if status.DiskUsed == 0 || status.DiskTotal == 0 {
		if nested, ok := objectMap(object, "disk", "disk_usage", "diskUsage"); ok {
			if status.DiskUsed == 0 {
				status.DiskUsed = intField(nested, "current", "used", "disk_used", "diskUsed")
			}
			if status.DiskTotal == 0 {
				status.DiskTotal = intField(nested, "total", "disk_total", "diskTotal")
			}
		}
	}
	return status, nil
}

func responsePayload(response xpanel.Response) json.RawMessage {
	if len(response.Obj) > 0 && string(response.Obj) != "null" {
		return response.Obj
	}
	return response.Data
}

func parseInbound(raw json.RawMessage) (Inbound, error) {
	object, err := objectFromPayload(raw)
	if err != nil {
		return Inbound{}, err
	}
	inbound := Inbound{
		RemoteID:   intField(object, "id", "remote_id", "remoteId"),
		Tag:        stringField(object, "tag"),
		Remark:     stringField(object, "remark", "remarkName", "name"),
		Protocol:   stringField(object, "protocol", "scheme"),
		Port:       intField(object, "port"),
		Listen:     stringField(object, "listen"),
		Enable:     boolField(object, "enable", "enabled"),
		ExpiryTime: intField(object, "expiry_time", "expiryTime"),
		Up:         intField(object, "up", "upload"),
		Down:       intField(object, "down", "download"),
		AllTime:    intField(object, "all_time", "allTime", "total"),
	}
	if inbound.AllTime == 0 {
		inbound.AllTime = inbound.Up + inbound.Down
	}

	settings, _ := objectFromOptional(objectValueRaw(object, "settings"))
	clientItems, err := findClientItems(object, settings)
	if err != nil {
		return Inbound{}, fmt.Errorf("inbound %d clients: %w", inbound.RemoteID, err)
	}
	clients := make([]Client, 0, len(clientItems))
	for _, item := range clientItems {
		client, err := parseClient(item)
		if err != nil {
			return Inbound{}, fmt.Errorf("inbound %d client: %w", inbound.RemoteID, err)
		}
		clients = append(clients, client)
	}
	trafficItems, err := trafficItemsFrom(object)
	if err != nil {
		return Inbound{}, fmt.Errorf("inbound %d client traffic: %w", inbound.RemoteID, err)
	}
	mergeClientTraffic(clients, trafficItems)
	inbound.Clients = clients
	inbound.ConfigHash = configHash(object)
	return inbound, nil
}

func parseClient(raw json.RawMessage) (Client, error) {
	object, err := objectFromPayload(raw)
	if err != nil {
		return Client{}, err
	}
	email := stringField(object, "email", "id")
	remoteID := stringField(object, "id", "client_id", "clientId", "email")
	if remoteID == "" {
		remoteID = email
	}
	client := Client{
		RemoteID:   remoteID,
		Email:      email,
		Enable:     boolField(object, "enable", "enabled"),
		ExpiryTime: intField(object, "expiry_time", "expiryTime"),
		Up:         intField(object, "up", "upload"),
		Down:       intField(object, "down", "download"),
		AllTime:    intField(object, "all_time", "allTime", "total"),
		LastOnline: intField(object, "last_online", "lastOnline"),
	}
	if client.AllTime == 0 {
		client.AllTime = client.Up + client.Down
	}
	return client, nil
}

func findClientItems(object, settings map[string]json.RawMessage) ([]json.RawMessage, error) {
	for _, source := range []map[string]json.RawMessage{object, settings} {
		if source == nil {
			continue
		}
		for _, key := range []string{"clients", "clientList", "clientStats"} {
			if raw, ok := source[key]; ok {
				return rawList(raw)
			}
		}
	}
	return nil, nil
}

func trafficItemsFrom(object map[string]json.RawMessage) ([]json.RawMessage, error) {
	for _, key := range []string{"clientStats", "clientTraffics", "client_traffics", "clientTraffic"} {
		if raw, ok := object[key]; ok {
			return rawList(raw)
		}
	}
	return nil, nil
}

func mergeClientTraffic(clients []Client, trafficItems []json.RawMessage) {
	for _, raw := range trafficItems {
		traffic, err := objectFromPayload(raw)
		if err != nil {
			continue
		}
		email := stringField(traffic, "email")
		remoteID := stringField(traffic, "id", "client_id", "clientId")
		for index := range clients {
			if (email != "" && clients[index].Email == email) || (remoteID != "" && clients[index].RemoteID == remoteID) {
				clients[index].Up = intField(traffic, "up", "upload")
				clients[index].Down = intField(traffic, "down", "download")
				clients[index].AllTime = intField(traffic, "all_time", "allTime", "total")
				clients[index].LastOnline = intField(traffic, "last_online", "lastOnline")
				if clients[index].AllTime == 0 {
					clients[index].AllTime = clients[index].Up + clients[index].Down
				}
				break
			}
		}
	}
}

func configHash(object map[string]json.RawMessage) string {
	selected := make(map[string]json.RawMessage)
	for _, key := range []string{"tag", "remark", "enable", "expiry_time", "expiryTime", "port", "listen", "protocol", "settings", "stream_settings", "streamSettings", "sniffing"} {
		if raw, ok := object[key]; ok {
			selected[key] = raw
		}
	}
	encoded, err := json.Marshal(selected)
	if err != nil {
		return ""
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:])
}

func newSyncID(now time.Time) string {
	random := make([]byte, 4)
	if _, err := rand.Read(random); err != nil {
		return now.Format("20060102T150405.000000000Z07:00")
	}
	return now.Format("20060102T150405.000000000Z07:00") + "-" + hex.EncodeToString(random)
}

func listItems(raw json.RawMessage, keys ...string) ([]json.RawMessage, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return []json.RawMessage{}, nil
	}
	if items, err := rawList(raw); err == nil {
		return items, nil
	}
	object, err := objectFromPayload(raw)
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		if nested, ok := object[key]; ok {
			return listItems(nested, keys...)
		}
	}
	return nil, errors.New("expected JSON array or list object")
}

func rawList(raw json.RawMessage) ([]json.RawMessage, error) {
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func objectFromPayload(raw json.RawMessage) (map[string]json.RawMessage, error) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(raw, &object); err != nil {
		return nil, err
	}
	if object == nil {
		return nil, errors.New("expected JSON object")
	}
	return object, nil
}

func objectFromOptional(raw json.RawMessage) (map[string]json.RawMessage, bool) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, false
	}
	var object map[string]json.RawMessage
	if json.Unmarshal(raw, &object) == nil && object != nil {
		return object, true
	}
	var encoded string
	if json.Unmarshal(raw, &encoded) == nil && json.Unmarshal([]byte(encoded), &object) == nil {
		return object, object != nil
	}
	return nil, false
}

func objectMap(object map[string]json.RawMessage, keys ...string) (map[string]json.RawMessage, bool) {
	for _, key := range keys {
		if raw, ok := object[key]; ok {
			if nested, ok := objectFromOptional(raw); ok {
				return nested, true
			}
		}
	}
	return nil, false
}

func objectValue(object map[string]json.RawMessage, keys ...string) (json.RawMessage, bool) {
	for _, key := range keys {
		if value, ok := object[key]; ok {
			return value, true
		}
	}
	return nil, false
}

func objectValueRaw(object map[string]json.RawMessage, keys ...string) json.RawMessage {
	value, _ := objectValue(object, keys...)
	return value
}

func stringField(object map[string]json.RawMessage, keys ...string) string {
	raw, ok := objectValue(object, keys...)
	if !ok {
		return ""
	}
	var value string
	if json.Unmarshal(raw, &value) == nil {
		return value
	}
	return strings.TrimSpace(string(raw))
}

func intField(object map[string]json.RawMessage, keys ...string) int64 {
	raw, ok := objectValue(object, keys...)
	if !ok {
		return 0
	}
	var number json.Number
	if json.Unmarshal(raw, &number) == nil {
		value, _ := strconv.ParseInt(number.String(), 10, 64)
		return value
	}
	var text string
	if json.Unmarshal(raw, &text) == nil {
		value, _ := strconv.ParseInt(strings.TrimSpace(text), 10, 64)
		return value
	}
	return 0
}

func floatField(object map[string]json.RawMessage, keys ...string) float64 {
	raw, ok := objectValue(object, keys...)
	if !ok {
		return 0
	}
	var number float64
	if json.Unmarshal(raw, &number) == nil {
		return number
	}
	var text string
	if json.Unmarshal(raw, &text) == nil {
		number, _ = strconv.ParseFloat(strings.TrimSpace(text), 64)
		return number
	}
	return 0
}

func boolField(object map[string]json.RawMessage, keys ...string) bool {
	raw, ok := objectValue(object, keys...)
	if !ok {
		return false
	}
	var value bool
	if json.Unmarshal(raw, &value) == nil {
		return value
	}
	text := strings.ToLower(strings.TrimSpace(stringField(object, keys...)))
	return text == "true" || text == "1" || text == "running" || text == "on"
}
