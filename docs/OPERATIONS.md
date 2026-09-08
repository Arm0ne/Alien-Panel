# 日常运维与故障排查

## 健康检查

```bash
curl -fsS http://127.0.0.1:18080/health/live
curl -fsS http://127.0.0.1:18080/health/ready
sudo docker compose -p alien-panel -f /opt/alien-panel/deploy/docker-compose.yml ps
```

`live` 表示进程存活，`ready` 表示数据库已打开并完成迁移。前端白屏时先确认 Nginx 容器和 `frontend-dist/index.html` 存在，再查看浏览器请求是否指向 `/api`。

## Agent 离线

在节点服务器执行：

```bash
sudo systemctl is-active xpanel-agent
sudo systemctl status xpanel-agent --no-pager
sudo journalctl -u xpanel-agent --since '30 minutes ago' --no-pager -o cat
sudo grep -E '^(central_url|node_key|node_type|xpanel_url|xpanel_base_path):' /etc/xpanel-agent/agent.yaml
```

排查顺序：

1. 节点能否访问中央 HTTPS 域名和 `/health/ready`；
2. Agent Token 是否被重新注册或节点是否被停用；
3. X-Panel 本机 URL、Base Path、账号和密码是否正确；
4. Agent 日志是否只有 `collect xpanel snapshot` 警告；这表示中央心跳仍可能正常，节点应保持在线但同步状态异常；
5. 修改配置后执行 `sudo systemctl restart xpanel-agent`。

正常日志会包含 `agent started`、周期 heartbeat 和 `xpanel snapshot synchronized`。Agent 版本应显示 `v1.0.3`，括号中的提交号只用于诊断。

## X-Panel 采集失败

采集器兼容常见 v1/v2 响应外壳、`obj`/`data.list`、字符串或对象 `settings`、`clientStats`/`clientTraffic` 和嵌套资源字段。HTTP 401 会重新登录一次；网络错误、超时、408/425/429/5xx 使用有界重试和退避。

确认本机 X-Panel 的管理地址包含正确的面板路径，例如：

```yaml
xpanel_url: 'https://panel.example.com:18086'
xpanel_base_path: '/Alien'
```

路径不能重复拼接。中央服务不会因为 X-Panel 采集失败而删除 Inbound、Client 或流量历史。

## 同步、缺失和流量

- 只有成功的完整同步才增加缺失计数；连续三次成功同步仍缺失才归档 Inbound。
- Agent 离线期间不增加缺失计数，最后一次成功快照保留。
- 累计流量下降写入 `traffic_reset`，当前值作为新基线，报表不会出现负流量。
- `traffic-check` 是只读核对工具，不访问 X-Panel，也不修改数据库。

## 用户和出口路径

用户按线路机 Inbound 聚合；Client/Email 只是设备。用户详情的固定出口 IP 支持多个同模式 IP，用户列表会完整显示所有地址。若某个 IP 被第二项或后续项使用，出口 IP 页面仍会显示配置归属用户数，删除和迁移会被阻止。

## 财务和续费

Agent 延长 Client 到期时间只生成待确认候选。管理员确认才产生收费记录；“非收费变更”不会进入实收。年费按订单服务区间分摊到月份。节点删除自动清理无关联业务用户，但不会删除收费记录和财务审计历史。

## 事件和日志

事件中心记录同步异常、节点离线/恢复、Inbound 缺失/归档、流量回退和待确认续费。查看中央日志：

```bash
sudo docker compose -p alien-panel -f /opt/alien-panel/deploy/docker-compose.yml logs -f --tail=200 central
```

不要在日志、截图或工单中粘贴管理员密码、Agent Token、X-Panel 密码、Session 或完整配置文件。
