-- 保留策略按时间窗口删除历史数据，需要 started_at 上的索引避免全表扫描。
CREATE INDEX IF NOT EXISTS idx_sync_runs_started_at ON sync_runs(started_at);
