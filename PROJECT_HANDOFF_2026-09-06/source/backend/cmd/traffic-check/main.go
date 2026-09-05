// traffic-check independently recomputes traffic usage from stored X-Panel
// cumulative snapshots. It is intended for staging and incident verification,
// not for writing or repairing production data.
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	_ "modernc.org/sqlite"
)

type inboundUsage struct {
	NodeKey        string `json:"nodeKey"`
	NodeName       string `json:"nodeName"`
	InboundID      string `json:"inboundId"`
	InboundTag     string `json:"inboundTag,omitempty"`
	SnapshotCount  int    `json:"snapshotCount"`
	ResetCount     int    `json:"resetCount"`
	CalculatedUsed int64  `json:"calculatedUsedBytes"`
}

type report struct {
	Database      string         `json:"database"`
	From          string         `json:"from"`
	To            string         `json:"to"`
	NodeKey       string         `json:"nodeKey,omitempty"`
	TotalUsed     int64          `json:"totalUsedBytes"`
	TotalResets   int            `json:"totalResets"`
	InboundUsages []inboundUsage `json:"inbounds"`
	GeneratedAt   string         `json:"generatedAt"`
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("traffic-check", flag.ContinueOnError)
	flags.SetOutput(stderr)
	databasePath := flags.String("database", "", "path to the central SQLite database (required)")
	fromRaw := flags.String("from", "", "inclusive RFC3339 timestamp or YYYY-MM-DD (defaults to 30 days ago)")
	toRaw := flags.String("to", "", "exclusive RFC3339 timestamp or YYYY-MM-DD (defaults to now)")
	nodeKey := flags.String("node-key", "", "only verify one node_key")
	format := flags.String("format", "text", "output format: text or json")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if strings.TrimSpace(*databasePath) == "" {
		fmt.Fprintln(stderr, "--database is required")
		return 2
	}
	if *format != "text" && *format != "json" {
		fmt.Fprintln(stderr, "--format must be text or json")
		return 2
	}
	if _, err := os.Stat(*databasePath); err != nil {
		fmt.Fprintf(stderr, "read database: %v\n", err)
		return 1
	}

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -30)
	to := now
	var err error
	if strings.TrimSpace(*fromRaw) != "" {
		from, err = parseTime(*fromRaw, false)
		if err != nil {
			fmt.Fprintf(stderr, "parse --from: %v\n", err)
			return 2
		}
	}
	if strings.TrimSpace(*toRaw) != "" {
		to, err = parseTime(*toRaw, true)
		if err != nil {
			fmt.Fprintf(stderr, "parse --to: %v\n", err)
			return 2
		}
	}
	if !from.Before(to) {
		fmt.Fprintln(stderr, "--from must be earlier than --to")
		return 2
	}

	database, err := sql.Open("sqlite", readOnlySQLiteDSN(*databasePath))
	if err != nil {
		fmt.Fprintf(stderr, "open database: %v\n", err)
		return 1
	}
	defer database.Close()
	if err := database.Ping(); err != nil {
		fmt.Fprintf(stderr, "connect database: %v\n", err)
		return 1
	}

	result, err := calculateUsage(database, *databasePath, from, to, strings.TrimSpace(*nodeKey))
	if err != nil {
		fmt.Fprintf(stderr, "calculate traffic: %v\n", err)
		return 1
	}
	if *format == "json" {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			fmt.Fprintf(stderr, "write report: %v\n", err)
			return 1
		}
		return 0
	}
	writeTextReport(stdout, result)
	return 0
}

func parseTime(value string, dateIsEnd bool) (time.Time, error) {
	value = strings.TrimSpace(value)
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed.UTC(), nil
	}
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("use RFC3339 or YYYY-MM-DD")
	}
	if dateIsEnd {
		return date.UTC().AddDate(0, 0, 1), nil
	}
	return date.UTC(), nil
}

func readOnlySQLiteDSN(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}
	return "file:" + filepath.ToSlash(absPath) + "?mode=ro&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)"
}

func calculateUsage(database *sql.DB, databasePath string, from, to time.Time, nodeKey string) (report, error) {
	where := []string{"t.collected_at >= ?", "t.collected_at < ?"}
	args := []any{from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano)}
	if nodeKey != "" {
		where = append(where, "n.node_key = ?")
		args = append(args, nodeKey)
	}
	rows, err := database.Query(`SELECT n.node_key, n.name, i.remote_inbound_id, COALESCE(i.tag, ''),
t.all_time, t.reset_detected,
COALESCE((SELECT p.all_time FROM traffic_snapshots p WHERE p.inbound_id = t.inbound_id AND p.collected_at < t.collected_at ORDER BY p.collected_at DESC LIMIT 1), -1)
FROM traffic_snapshots t
JOIN inbounds i ON i.id = t.inbound_id
JOIN nodes n ON n.id = t.node_id
WHERE `+strings.Join(where, " AND ")+`
ORDER BY n.node_key ASC, i.remote_inbound_id ASC, t.collected_at ASC`, args...)
	if err != nil {
		return report{}, err
	}
	defer rows.Close()

	byInbound := make(map[string]*inboundUsage)
	for rows.Next() {
		var nodeKeyValue, nodeName, inboundID, inboundTag string
		var current, previous int64
		var resetDetected int
		if err := rows.Scan(&nodeKeyValue, &nodeName, &inboundID, &inboundTag, &current, &resetDetected, &previous); err != nil {
			return report{}, err
		}
		key := nodeKeyValue + "\x00" + inboundID
		usage := byInbound[key]
		if usage == nil {
			usage = &inboundUsage{NodeKey: nodeKeyValue, NodeName: nodeName, InboundID: inboundID, InboundTag: inboundTag}
			byInbound[key] = usage
		}
		usage.SnapshotCount++
		if resetDetected == 1 || (previous >= 0 && current < previous) {
			usage.CalculatedUsed += current
			usage.ResetCount++
			continue
		}
		if previous >= 0 {
			usage.CalculatedUsed += current - previous
		}
	}
	if err := rows.Err(); err != nil {
		return report{}, err
	}

	items := make([]inboundUsage, 0, len(byInbound))
	var totalUsed int64
	var totalResets int
	for _, usage := range byInbound {
		items = append(items, *usage)
		totalUsed += usage.CalculatedUsed
		totalResets += usage.ResetCount
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].NodeKey == items[j].NodeKey {
			return items[i].InboundID < items[j].InboundID
		}
		return items[i].NodeKey < items[j].NodeKey
	})
	return report{
		Database: databasePath, From: from.Format(time.RFC3339Nano), To: to.Format(time.RFC3339Nano), NodeKey: nodeKey,
		TotalUsed: totalUsed, TotalResets: totalResets, InboundUsages: items, GeneratedAt: time.Now().UTC().Format(time.RFC3339Nano),
	}, nil
}

func writeTextReport(output io.Writer, result report) {
	fmt.Fprintf(output, "Traffic verification period: %s to %s (end exclusive)\n", result.From, result.To)
	if result.NodeKey != "" {
		fmt.Fprintf(output, "Node filter: %s\n", result.NodeKey)
	}
	writer := tabwriter.NewWriter(output, 0, 4, 2, ' ', 0)
	fmt.Fprintln(writer, "NODE\tINBOUND\tTAG\tSNAPSHOTS\tRESETS\tCALCULATED USED (BYTES)")
	for _, usage := range result.InboundUsages {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%d\t%d\t%d\n", usage.NodeKey, usage.InboundID, usage.InboundTag, usage.SnapshotCount, usage.ResetCount, usage.CalculatedUsed)
	}
	_ = writer.Flush()
	fmt.Fprintf(output, "Total calculated usage: %d bytes; reset baselines: %d\n", result.TotalUsed, result.TotalResets)
}
