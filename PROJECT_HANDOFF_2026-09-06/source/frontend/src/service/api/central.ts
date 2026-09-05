import { request } from '../request';

/** Fetch central operations dashboard data. */
export function fetchDashboard() {
  return request<Api.Central.DashboardSummary>({ url: '/dashboard' });
}

export function fetchUsers(params: Api.Central.PageParams = {}) {
  return request<Api.Central.PageResult<Api.Central.UserSummary>>({ url: '/users', params });
}

/** Read a single Inbound-based business user with its X-Panel snapshots. */
export function fetchUserDetail(id: string) {
  return request<Api.Central.UserDetail>({ url: `/users/${encodeURIComponent(id)}` });
}

export function fetchUserRenewals(id: string) {
  return request<Api.Central.UserRenewalCandidate[]>({ url: `/users/${encodeURIComponent(id)}/renewals` });
}

export function confirmUserRenewal(
  id: string,
  candidateId: string,
  data: { billingCycle?: Api.Central.BillingCycle; amount?: number; paidAt?: string; notes?: string }
) {
  return request<Api.Central.UserDetail>({
    url: `/users/${encodeURIComponent(id)}/renewals/${encodeURIComponent(candidateId)}/confirm`,
    method: 'post',
    data
  });
}

export function rejectUserRenewal(id: string, candidateId: string, data: { notes?: string } = {}) {
  return request<Api.Central.UserDetail>({
    url: `/users/${encodeURIComponent(id)}/renewals/${encodeURIComponent(candidateId)}/reject`,
    method: 'post',
    data
  });
}

/** Read cumulative traffic growth trend data for a relay business Inbound. */
export function fetchUserTraffic(id: string, range: Api.Central.TrafficRange = '1h') {
  return request<Api.Central.UserTrafficTrend>({
    url: `/users/${encodeURIComponent(id)}/traffic`,
    params: { range }
  });
}

/** Read the authoritative resources available for one user's network path. */
export function fetchUserPathAssets(id: string) {
  return request<Api.Central.UserPathAssets>({ url: `/users/${encodeURIComponent(id)}/path-assets` });
}

/** Update central-owned business metadata; never writes back to X-Panel. */
export function updateUser(id: string, data: Api.Central.UpdateUserPayload) {
  return request<Api.Central.UserDetail>({
    url: `/users/${encodeURIComponent(id)}`,
    method: 'patch',
    data
  });
}

/** Assign a direct per-user path: relay is derived from the primary inbound. */
export function assignUserPath(id: string, data: Api.Central.AssignUserPathPayload) {
  return request<Api.Central.UserDetail>({
    url: `/users/${encodeURIComponent(id)}/path`,
    method: 'put',
    data
  });
}

/** Clear the user's current direct path while retaining its history. */
export function clearUserPath(id: string) {
  return request<Api.Central.UserDetail>({
    url: `/users/${encodeURIComponent(id)}/path`,
    method: 'delete'
  });
}

/** Legacy route assignment kept for old clients during migration. */
export function assignUserRoute(id: string, data: Api.Central.AssignUserRoutePayload) {
  return request<Api.Central.UserDetail>({
    url: `/users/${encodeURIComponent(id)}/route`,
    method: 'put',
    data
  });
}

/** Remove the user's active route assignment while retaining its history. */
export function clearUserRoute(id: string) {
  return request<Api.Central.UserDetail>({
    url: `/users/${encodeURIComponent(id)}/route`,
    method: 'delete'
  });
}

export function fetchNodes(params: Api.Central.PageParams = {}) {
  return request<Api.Central.PageResult<Api.Central.NodeSummary>>({ url: '/nodes', params });
}

/** Create a central node record and issue a one-time Agent token. */
export function createNode(data: Api.Central.NodeCreatePayload) {
  return request<Api.Central.NodeRegistrationResult>({ url: '/nodes', method: 'post', data });
}

/** Issue a short-lived one-time token for the online Agent installer. */
export function issueNodeInstallToken(id: string) {
  return request<Api.Central.NodeInstallTokenResult>({
    url: `/nodes/${encodeURIComponent(id)}/install-token`,
    method: 'post'
  });
}

/** Update central node metadata or enable/disable a node. */
export function updateNode(id: string, data: Api.Central.NodeUpdatePayload) {
  return request<Api.Central.NodeSummary>({
    url: `/nodes/${encodeURIComponent(id)}`,
    method: 'patch',
    data
  });
}

/** Permanently delete a node and all node-owned operational data. */
export function deleteNode(id: string) {
  return request<{ id: string; deleted: boolean }>({
    url: `/nodes/${encodeURIComponent(id)}`,
    method: 'delete'
  });
}

export function fetchNodeDetail(id: string) {
  return request<Api.Central.NodeDetail>({ url: `/nodes/${encodeURIComponent(id)}` });
}

/** Queue an on-demand snapshot; the node Agent executes it on its next cycle. */
export function requestNodeSync(id: string) {
  return request<Api.Central.NodeSyncRequest>({
    url: `/nodes/${encodeURIComponent(id)}/sync`,
    method: 'post'
  });
}

/** Add a central-owned temporal monthly cost record for a node. */
export function createNodeCost(nodeId: string, data: Api.Central.NodeCostCreatePayload) {
  return request<Api.Central.NodeCostRecord>({
    url: `/nodes/${encodeURIComponent(nodeId)}/costs`,
    method: 'post',
    data
  });
}

export function fetchNodeCosts(nodeId: string) {
  return request<Api.Central.NodeCostRecord[]>({ url: `/nodes/${encodeURIComponent(nodeId)}/costs` });
}

/** Update cost metadata without changing its effective date range. */
export function updateNodeCost(nodeId: string, costId: string, data: Api.Central.NodeCostUpdatePayload) {
  return request<Api.Central.NodeCostRecord>({
    url: `/nodes/${encodeURIComponent(nodeId)}/costs/${encodeURIComponent(costId)}`,
    method: 'patch',
    data
  });
}

export function fetchRoutes(params: Api.Central.PageParams = {}) {
  return request<Api.Central.PageResult<Api.Central.RouteSummary>>({ url: '/routes', params });
}

export function fetchRouteDetail(id: string) {
  return request<Api.Central.RouteDetail>({ url: `/routes/${encodeURIComponent(id)}` });
}

export function createRoute(data: Api.Central.RouteUpsertPayload) {
  return request<Api.Central.RouteDetail>({ url: '/routes', method: 'post', data });
}

export function updateRoute(id: string, data: Partial<Api.Central.RouteUpsertPayload>) {
  return request<Api.Central.RouteDetail>({
    url: `/routes/${encodeURIComponent(id)}`,
    method: 'patch',
    data
  });
}

export function deleteRoute(id: string) {
  return request<{ id: string; deleted: boolean }>({
    url: `/routes/${encodeURIComponent(id)}`,
    method: 'delete'
  });
}

export function fetchRouteExitIpBindings(routeId: string) {
  return request<Api.Central.RouteExitIpBinding[]>({ url: `/routes/${encodeURIComponent(routeId)}/exit-ips` });
}

export function bindRouteExitIp(routeId: string, data: Api.Central.RouteExitIpBindPayload) {
  return request<Api.Central.RouteExitIpBinding>({
    url: `/routes/${encodeURIComponent(routeId)}/exit-ips`,
    method: 'post',
    data
  });
}

export function updateRouteExitIp(routeId: string, exitIpId: string, data: Api.Central.RouteExitIpUpdatePayload) {
  return request<Api.Central.RouteExitIpBinding>({
    url: `/routes/${encodeURIComponent(routeId)}/exit-ips/${encodeURIComponent(exitIpId)}`,
    method: 'patch',
    data
  });
}

export function unbindRouteExitIp(routeId: string, exitIpId: string) {
  return request<{ id: string; routeId: string; exitIpId: string; deleted: boolean }>({
    url: `/routes/${encodeURIComponent(routeId)}/exit-ips/${encodeURIComponent(exitIpId)}`,
    method: 'delete'
  });
}

export function fetchExitIps(params: Api.Central.PageParams = {}) {
  return request<Api.Central.PageResult<Api.Central.ExitIpSummary>>({ url: '/exit-ips', params });
}

export function fetchExitIpDetail(id: string) {
  return request<Api.Central.ExitIpDetail>({ url: `/exit-ips/${encodeURIComponent(id)}` });
}

export function createExitIp(data: Api.Central.ExitIpUpsertPayload) {
  return request<Api.Central.ExitIpDetail>({ url: '/exit-ips', method: 'post', data });
}

export function updateExitIp(id: string, data: Partial<Api.Central.ExitIpUpsertPayload>) {
  return request<Api.Central.ExitIpDetail>({
    url: `/exit-ips/${encodeURIComponent(id)}`,
    method: 'patch',
    data
  });
}

export function deleteExitIp(id: string) {
  return request<{ id: string; deleted: boolean }>({
    url: `/exit-ips/${encodeURIComponent(id)}`,
    method: 'delete'
  });
}

export function fetchFinanceSummary(params: { period?: string } = {}) {
  return request<Api.Central.FinanceSummary>({ url: '/costs/summary', params });
}

export function fetchEvents(params: Api.Central.PageParams = {}) {
  return request<Api.Central.PageResult<Api.Central.EventSummary>>({ url: '/events', params });
}

export function fetchEventSummary() {
  return request<Api.Central.EventSummaryCounts>({ url: '/events/summary' });
}

export function markEventRead(id: string) {
  return request<{ id: string; read: boolean }>({
    url: `/events/${encodeURIComponent(id)}/read`,
    method: 'post'
  });
}

export function markAllEventsRead() {
  return request<{ read: boolean }>({ url: '/events/read-all', method: 'post' });
}

export function resolveEvent(id: string) {
  return request<{ id: string; resolved: boolean }>({
    url: `/events/${encodeURIComponent(id)}/resolve`,
    method: 'post'
  });
}
