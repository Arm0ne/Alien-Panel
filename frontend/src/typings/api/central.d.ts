declare namespace Api {
  namespace Central {
    type UserStatus = 'active' | 'expiring' | 'expired' | 'disabled' | 'unknown';
    type NodeType = 'relay' | 'landing' | 'unknown';
    type NodeStatus = 'online' | 'degraded' | 'offline' | 'disabled' | 'unknown';
    type EventSeverity = 'info' | 'warning' | 'error';
    type EventCategory = 'business' | 'node' | 'sync' | 'system';
    type EventStatus = 'open' | 'resolved' | 'dismissed';

    interface PageParams {
      page?: number;
      page_size?: number;
      sort?: string;
      order?: 'asc' | 'desc';
      keyword?: string;
      status?: string;
      severity?: EventSeverity;
      category?: EventCategory;
      node_id?: string;
      node_type?: NodeType;
      billing_type?: BillingType;
    }

    interface PageResult<T> {
      items: T[];
      total: number;
      page: number;
      pageSize: number;
      dataAt: string | null;
    }

    interface UserSummary {
      id: string;
      name: string;
      nodeId: string;
      nodeName: string;
      inboundTag: string;
      routeName?: string | null;
      landingNodeName?: string | null;
      exitIpAddress?: string | null;
      exitIpOwnerNodeName?: string | null;
      pathMode?: 'relay' | 'landing' | 'external' | null;
      pathConfigured?: boolean;
      status: UserStatus;
      billingType: BillingType;
      expiresAt?: string;
      clientCount: number;
      trafficBytes: number;
      lastActivityAt?: string;
      dataAt?: string | null;
    }

    interface UserListStats {
      active: number;
      paid: number;
      free: number;
    }

    interface UserListResult extends PageResult<UserSummary> {
      stats: UserListStats;
    }

    interface UserGroupStats {
      total: number;
      active: number;
      expiring: number;
      expired: number;
      disabled: number;
      paid: number;
      free: number;
      trafficBytes: number;
    }

    interface UserGroupSummary {
      nodeId: string;
      nodeName: string;
      nodeType: NodeType;
      status: NodeStatus;
      syncStatus: 'unknown' | 'success' | 'failed';
      enabled: boolean;
      lastSyncAt: string | null;
      stats: UserGroupStats;
    }

    interface UserGroupListResult extends PageResult<UserGroupSummary> {
      stats: UserListStats;
    }

    interface UserInboundDetail {
      id: string | null;
      remoteId: string | null;
      tag: string | null;
      remark: string | null;
      protocol: string | null;
      port: number;
      enabled: boolean;
      clientCount: number;
      up: number;
      down: number;
      allTime: number;
      lastSeenAt: string | null;
    }

    interface UserNodeDetail {
      id: string | null;
      name: string | null;
      type: NodeType | null;
    }

    interface UserClientDetail {
      remoteId: string;
      email: string | null;
      enabled: boolean;
      expiresAt: string | null;
      up: number;
      down: number;
      allTime: number;
      lastOnlineAt: string | null;
      lastSeenAt: string | null;
    }

    type TrafficRange = '1h' | '6h' | '1d' | '7d';
    type BillingCycle = 'monthly' | 'annual';
    type BillingType = 'paid' | 'free';

    interface UserRenewalCandidate {
      id: string;
      inboundId: string | null;
      oldExpiryAt: string;
      newExpiryAt: string;
      detectedAt: string;
      billingCycle: BillingCycle;
      suggestedAmount: number;
      currency: 'CNY';
      status: 'pending' | 'confirmed' | 'rejected';
      processedAt: string | null;
      notes: string | null;
    }

    interface UserBillingRecord {
      id: string;
      billingCycle: BillingCycle;
      amount: number;
      currency: 'CNY';
      serviceFrom: string;
      serviceTo: string;
      paidAt: string | null;
      status: 'pending' | 'confirmed' | 'rejected' | 'cancelled';
      source: 'manual' | 'agent_detected';
      orderType: 'initial' | 'renewal' | 'recovery';
      orderNo: string | null;
      notes: string | null;
      createdAt: string;
      origin: 'live' | 'historical_import';
      verificationStatus: 'verified' | 'unverified';
    }

    interface BillingImportResult {
      dryRun: boolean;
      imported: number;
      unverified: number;
      records: Array<{
        row: number;
        userId: string;
        billingCycle: BillingCycle;
        amount: number;
        serviceFrom: string;
        serviceTo: string;
        paidAt: string | null;
        orderType: 'initial' | 'renewal' | 'recovery';
        status: 'confirmed' | 'pending';
        verificationStatus: 'verified' | 'unverified';
      }>;
      errors: Array<{ row: number; message: string }>;
    }

    interface UserTrafficPoint {
      time: string;
      uploadBytes: number;
      downloadBytes: number;
      totalBytes: number;
      uploadRate: number;
      downloadRate: number;
      sampleCount: number;
      resetDetected: boolean;
      hasGap: boolean;
    }

    interface UserTrafficSummary {
      uploadBytes: number;
      downloadBytes: number;
      totalBytes: number;
      averageUploadRate: number;
      averageDownloadRate: number;
      peakUploadRate: number;
      peakDownloadRate: number;
      sampleCount: number;
      coverage: number;
    }

    interface UserTrafficTrend {
      range: TrafficRange;
      from: string;
      to: string;
      bucket: string;
      dataAt: string | null;
      points: UserTrafficPoint[];
      summary: UserTrafficSummary;
    }

    interface UserRouteDetail {
      id: string;
      name: string;
      relayNodeName: string | null;
      landingNodeName: string | null;
      landingInboundId: string | null;
      landingInboundTag: string | null;
      enabled: boolean;
      isPrimary: boolean;
      activeFrom: string | null;
      activeTo: string | null;
      routeExitIpId: string | null;
      exitIpId: string | null;
      exitIpAddress: string | null;
      exitIpScope: 'relay' | 'landing' | 'external' | null;
      exitIpSourceType: 'node' | 's5' | null;
      exitIpOwnerNodeName: string | null;
      assignmentMode: 'fixed' | 'pool';
    }

    interface UserDetail {
      id: string;
      displayName: string;
      status: UserStatus;
      billingType: BillingType;
      freeReason: string | null;
      monthlyFee: number;
      billingCycle: BillingCycle;
      billingAmount: number;
      currency: 'CNY';
      notes: string | null;
      expiresAt: string | null;
      inbound: UserInboundDetail;
      node: UserNodeDetail;
      clients: UserClientDetail[];
      routes: UserRouteDetail[];
      path: UserPathDetail | null;
      pathHistory: UserPathDetail[];
      traffic: Array<{
        collectedAt: string;
        up: number;
        down: number;
        allTime: number;
        resetDetected: boolean;
      }>;
      renewalCandidates: UserRenewalCandidate[];
      billingRecords: UserBillingRecord[];
    }

    interface UpdateUserPayload {
      displayName?: string;
      monthlyFee?: number;
      billingCycle?: BillingCycle;
      billingAmount?: number;
      billingType?: BillingType;
      freeReason?: string;
      currency?: 'CNY';
      notes?: string;
    }

    interface UserPathDetail {
      id: string;
      relayNodeId: string;
      relayNodeName: string | null;
      landingNodeId: string | null;
      landingNodeName: string | null;
      landingInboundId: string | null;
      landingInboundTag: string | null;
      exitIpId: string;
      exitIpAddress: string | null;
      exitIpIds?: string[];
      exitIps?: Array<{
        id: string;
        address: string;
        sourceType: 'node' | 's5';
        ownerNodeId: string | null;
        ownerNodeName: string | null;
        enabled?: boolean;
      }>;
      exitIpSourceType: 'node' | 's5';
      exitIpOwnerNodeId: string | null;
      exitIpOwnerNodeName: string | null;
      mode: 'relay' | 'landing' | 'external';
      notes: string | null;
      activeFrom: string | null;
      activeTo: string | null;
      valid: boolean;
    }

    interface AssignUserPathPayload {
      landingNodeId?: string | null;
      landingInboundId?: string | null;
      exitIpId?: string;
      exitIpIds: string[];
      notes?: string | null;
    }

    type UserPathAssetInboundState = 'ready' | 'pending' | 'empty';

    interface UserPathAssetNode extends NodeSummary {
      type: 'landing';
      inboundState: UserPathAssetInboundState;
      inbounds: Array<NodeInboundDetail & { purpose: 'infrastructure' }>;
      exitIps: ExitIpSummary[];
    }

    interface UserPathAssets {
      generatedAt: string;
      dataAt: string | null;
      relay: Pick<NodeSummary, 'id' | 'name' | 'type' | 'status' | 'enabled' | 'host' | 'managementUrl'> & {
        lastSeenAt: string | null;
        lastSyncAt: string | null;
      };
      landingNodes: UserPathAssetNode[];
      relayExitIps: ExitIpSummary[];
      externalExitIps: ExitIpSummary[];
    }

    interface AssignUserRoutePayload {
      routeId: string;
      /** Route-exit binding ID. Omit or set null to use the route's weighted pool. */
      routeExitIpId?: string | null;
    }

    interface NodeSummary {
      id: string;
      name: string;
      type: NodeType;
      status: NodeStatus;
      enabled: boolean;
      host: string;
      region?: string | null;
      managementUrl?: string | null;
      xrayVersion?: string;
      syncStatus?: 'unknown' | 'success' | 'failed';
      lastSyncError?: string | null;
      cpuUsage?: number;
      memoryUsed?: number;
      memoryTotal?: number;
      diskUsed?: number;
      diskTotal?: number;
      exitIpCount: number;
      lastSyncAt?: string;
      monthlyCost?: number;
      currency?: string;
      userCount: number;
      paidUserCount: number;
      freeUserCount: number;
      clientCount: number;
    }

    interface NodeListStats {
      total: number;
      online: number;
      relay: number;
      landing: number;
    }

    interface NodeListResult extends PageResult<NodeSummary> {
      stats: NodeListStats;
    }

    interface NodeCreatePayload {
      /** Optional advanced override; the central service generates one when omitted. */
      nodeKey?: string;
      name: string;
      type: 'relay' | 'landing';
      hostname?: string;
      managementUrl?: string;
      publicIp?: string;
      /** Optional node-owned exit addresses; one address per item. */
      exitIps?: string[];
      region?: string;
      provider?: string;
      panelBasePath?: string;
    }

    type NodeInboundStatus = 'active' | 'disabled' | 'archived';
    type NodeSyncRunStatus = 'queued' | 'running' | 'success' | 'failed';

    type InboundPurpose = 'business' | 'infrastructure';

    interface NodeInboundDetail {
      id: string;
      remoteId: string;
      tag: string | null;
      remark: string | null;
      kind: string;
      purpose?: InboundPurpose;
      protocol: string | null;
      port: number;
      listen: string | null;
      enabled: boolean;
      status: NodeInboundStatus;
      expiresAt: string | null;
      clientCount: number;
      up: number;
      down: number;
      allTime: number;
      lastSeenAt: string | null;
      deletedAt: string | null;
    }

    interface NodeExitIpDetail {
      id: string;
      address: string;
      sourceType: 'node' | 's5';
      ownerNodeId: string | null;
      ownerNodeName: string | null;
      ownerNodeType: 'relay' | 'landing' | null;
      region: string | null;
      provider: string | null;
      family: number;
      status: 'active' | 'disabled';
      monthlyCost: number;
      currency: string;
      validFrom: string | null;
      validTo: string | null;
      checkedAt: string | null;
    }

    interface NodeCostRecord {
      id: string;
      nodeId: string;
      nodeName: string;
      category: string;
      monthlyAmount: number;
      currency: 'CNY';
      effectiveFrom: string;
      effectiveTo: string | null;
      notes: string | null;
      createdAt: string;
    }

    interface NodeCostCreatePayload {
      category: string;
      monthlyAmount: number;
      currency?: 'CNY';
      effectiveFrom: string;
      effectiveTo?: string | null;
      notes?: string | null;
    }

    interface NodeCostUpdatePayload {
      category?: string;
      monthlyAmount?: number;
      currency?: 'CNY';
      effectiveFrom?: string;
      effectiveTo?: string | null;
      notes?: string | null;
    }

    interface NodeSyncRun {
      id: string;
      syncId: string;
      startedAt: string;
      finishedAt: string | null;
      status: NodeSyncRunStatus;
      inboundCount: number;
      clientCount: number;
      errorMessage: string | null;
    }

    interface NodeStatusHistory {
      id: string;
      type: string;
      severity: EventSeverity;
      title?: string;
      message: string;
      occurredAt: string;
      acknowledged: boolean;
      requiresAction?: boolean;
      status?: EventStatus | string;
    }

    interface NodeDetail extends NodeSummary {
      nodeKey: string;
      hostname: string | null;
      managementUrl: string | null;
      publicIp: string | null;
      region: string | null;
      provider: string | null;
      panelBasePath: string | null;
      agentVersion: string | null;
      syncStatus: 'unknown' | 'success' | 'failed';
      lastSyncError: string | null;
      enabled: boolean;
      cpuUsage: number | null;
      memoryUsed: number | null;
      memoryTotal: number | null;
      diskUsed: number | null;
      diskTotal: number | null;
      userCount: number;
      paidUserCount: number;
      freeUserCount: number;
      clientCount: number;
      trafficBytes: number;
      dataAt: string | null;
      inbounds: NodeInboundDetail[];
      exitIps: NodeExitIpDetail[];
      costs: NodeCostRecord[];
      syncRuns: NodeSyncRun[];
      statusHistory: NodeStatusHistory[];
    }

    interface NodeSyncRequest {
      requestId: string;
      nodeId: string;
      nodeName: string;
      status: 'queued';
      requestedAt: string;
    }

    interface NodeRegistrationResult {
      nodeId: string;
      nodeKey: string;
      name: string;
      type: 'relay' | 'landing';
      token: string;
      installerToken: string;
      installerTokenExpiresAt: string;
      enabled: boolean;
      exitIpCount: number;
    }

    interface NodeInstallTokenResult {
      nodeId: string;
      nodeName: string;
      installerToken: string;
      installerTokenExpiresAt: string;
    }

    interface NodeUpdatePayload {
      name?: string;
      type?: NodeType;
      hostname?: string;
      managementUrl?: string;
      publicIp?: string;
      region?: string;
      provider?: string;
      panelBasePath?: string;
      enabled?: boolean;
    }

    interface RouteSummary {
      id: string;
      name: string;
      relayNodeId: string;
      relayNodeName: string;
      landingNodeId: string;
      landingNodeName: string;
      landingInboundTag?: string;
      exitIpCount: number;
      allocatedUserCount: number;
      status: 'active' | 'disabled' | 'unknown';
      dataAt?: string | null;
    }

    interface RouteDetail {
      id: string;
      name: string;
      relayNodeId: string;
      relayNodeName: string | null;
      landingNodeId: string;
      landingNodeName: string | null;
      relayOutboundTag: string | null;
      landingInboundId: string | null;
      landingInboundTag: string | null;
      enabled: boolean;
      validFrom: string | null;
      validTo: string | null;
      notes: string | null;
      exitIpCount: number;
      allocatedUserCount: number;
      status: 'active' | 'disabled' | 'unknown';
    }

    interface RouteExitIpBinding {
      id: string;
      routeId: string;
      exitIpId: string;
      address: string;
      family: 4 | 6;
      sourceType: 'node' | 's5';
      ownerNodeId: string | null;
      ownerNodeName: string | null;
      ownerNodeType: 'relay' | 'landing' | null;
      region: string | null;
      scope: 'relay' | 'landing' | 'external';
      landingNodeId: string | null;
      landingNodeName: string | null;
      allocationWeight: number;
      enabled: boolean;
    }

    type RouteExitIpScope = 'relay' | 'landing' | 'external';

    interface RouteExitIpBindPayload {
      exitIpId: string;
      scope?: 'relay' | 'landing' | 'external';
      allocationWeight?: number;
      enabled?: boolean;
    }

    interface RouteExitIpUpdatePayload {
      scope?: 'relay' | 'landing' | 'external';
      allocationWeight?: number;
      enabled?: boolean;
    }

    interface RouteUpsertPayload {
      name: string;
      relayNodeId: string;
      landingNodeId: string;
      relayOutboundTag?: string | null;
      landingInboundId?: string | null;
      landingInboundTag?: string | null;
      enabled?: boolean;
      validFrom?: string | null;
      validTo?: string | null;
      notes?: string | null;
    }

    interface ExitIpSummary {
      id: string;
      address: string;
      sourceType: 'node' | 's5';
      ownerNodeId: string | null;
      ownerNodeName: string | null;
      ownerNodeType: 'relay' | 'landing' | null;
      region: string | null;
      landingNodeId: string | null;
      landingNodeName: string | null;
      family?: 4 | 6;
      provider?: string | null;
      status: 'active' | 'disabled' | 'unknown';
      monthlyCost?: number;
      currency?: string;
      allocatedUserCount: number;
      checkedAt?: string | null;
    }

    interface ExitIpCountryStat {
      name: string;
      count: number;
    }

    interface ExitIpListStats {
      total: number;
      active: number;
      assigned: number;
      unassigned: number;
      countries: ExitIpCountryStat[];
    }

    interface ExitIpListResult extends PageResult<ExitIpSummary> {
      stats: ExitIpListStats;
    }

    interface ExitIpDetail {
      id: string;
      address: string;
      sourceType: 'node' | 's5';
      ownerNodeId: string | null;
      ownerNodeName: string | null;
      ownerNodeType: 'relay' | 'landing' | null;
      region: string | null;
      landingNodeId: string | null;
      landingNodeName: string | null;
      family: 4 | 6;
      provider: string | null;
      monthlyCost: number;
      currency: 'CNY';
      enabled: boolean;
      validFrom: string | null;
      validTo: string | null;
      notes: string | null;
      allocatedUserCount: number;
      checkedAt: string | null;
      status: 'active' | 'disabled' | 'unknown';
    }

    interface ExitIpUpsertPayload {
      address: string;
      sourceType?: 'node' | 's5';
      ownerNodeId?: string | null;
      region?: string | null;
      /** @deprecated use sourceType + ownerNodeId */
      landingNodeId?: string;
      family?: 4 | 6;
      provider?: string | null;
      monthlyCost?: number;
      currency?: 'CNY';
      enabled?: boolean;
      validFrom?: string | null;
      validTo?: string | null;
      notes?: string | null;
    }

    interface FinanceSummary {
      period: string;
      currency: string;
      effectiveUserCount: number;
      monthIncome: number;
      cashIncome: number;
      cumulativeCashIncome: number;
      paidUserCount: number;
      freeUserCount: number;
      monthCost: number;
      cashBalance: number;
      grossProfit: number;
      breakdown?: Array<{
        label: string;
        amount: number;
      }>;
      dataAt?: string | null;
      cashBreakdown?: Array<{ orderType: 'initial' | 'renewal' | 'recovery' | string; count: number; amount: number }>;
      orders?: Array<{
        id: string;
        orderNo: string | null;
        userId: string;
        userName: string;
        orderType: 'initial' | 'renewal' | 'recovery' | string;
        billingCycle: BillingCycle;
        amount: number;
        currency: string;
        serviceFrom: string;
        serviceTo: string;
        paidAt: string | null;
        source: string;
        inboundTag: string | null;
        nodeName: string | null;
      }>;
      costDetails?: Array<{
        id: string;
        type: string;
        name: string;
        category: string;
        amount: number;
        currency: string;
        effectiveFrom: string;
        effectiveTo: string | null;
        notes: string | null;
      }>;
      renewal?: {
        dueUsers: number;
        renewedUsers: number;
        notRenewedUsers: number;
        renewalRate: number;
        recoveryUsers: number;
      };
    }

    interface EventSummary {
      id: string;
      type: string;
      category: EventCategory | string;
      severity: EventSeverity;
      title: string;
      nodeId?: string | null;
      nodeName?: string;
      message: string;
      occurredAt: string;
      acknowledged: boolean;
      requiresAction: boolean;
      status: EventStatus | string;
      resourceType?: string | null;
      resourceId?: string | null;
      actionType?: string | null;
      payload?: Record<string, unknown> | null;
      readAt?: string | null;
      resolvedAt?: string | null;
      resolvedBy?: string | null;
      source?: string | null;
      correlationId?: string | null;
    }

    interface EventSummaryCounts {
      pendingCount: number;
      unreadCount: number;
      generatedAt: string;
    }

    type DashboardRange = 'today' | '7d' | '30d' | 'custom';

    interface DashboardQuery {
      range?: DashboardRange;
      from?: string;
      to?: string;
    }

    interface DashboardTrafficPoint {
      time: string;
      uploadBytes: number;
      downloadBytes: number;
      totalBytes: number;
      sampleCount: number;
      resetDetected: boolean;
      hasGap: boolean;
    }

    interface DashboardTrafficSummary {
      uploadBytes: number;
      downloadBytes: number;
      totalBytes: number;
      sampleCount: number;
      coverage: number;
    }

    interface DashboardTrafficTrend {
      range: DashboardRange | 'month';
      from: string;
      to: string;
      bucket: '1h' | '1d' | string;
      dataAt: string | null;
      points: DashboardTrafficPoint[];
      summary: DashboardTrafficSummary;
    }

    interface DashboardNodeTrafficItem {
      nodeId: string;
      nodeName: string;
      nodeType: NodeType;
      status: NodeStatus;
      uploadBytes: number;
      downloadBytes: number;
      totalBytes: number;
      share: number;
      dataAt: string | null;
    }

    interface DashboardUserTrafficItem {
      userId: string;
      userName: string;
      nodeId: string;
      nodeName: string;
      inboundId: string;
      inboundTag: string;
      clientCount: number;
      status: UserStatus;
      expiresAt: string | null;
      uploadBytes: number;
      downloadBytes: number;
      totalBytes: number;
      lastActivityAt: string | null;
      dataAt: string | null;
    }

    interface DashboardEventItem {
      id: string;
      type: string;
      category: EventCategory | string;
      severity: EventSeverity;
      title: string;
      nodeId: string | null;
      nodeName: string;
      message: string;
      occurredAt: string;
      acknowledged: boolean;
      requiresAction: boolean;
      status: EventStatus | string;
    }

    interface DashboardSummary {
      generatedAt: string;
      dataAt: string | null;
      range: {
        name: DashboardRange;
        from: string;
        to: string;
        bucket: '1h' | '1d' | string;
      };
      nodes: {
        total: number;
        online: number;
        relay: number;
        landing: number;
      };
      users: {
        active: number;
        expiring: number;
        expired: number;
      };
      traffic: {
        todayBytes: number;
        monthBytes: number;
        uploadBytes: number;
        downloadBytes: number;
        totalBytes: number;
        businessScope: 'relay-user-inbound';
      };
      trafficTrend: DashboardTrafficTrend;
      nodeTrafficRanking: DashboardNodeTrafficItem[];
      userTrafficRanking: DashboardUserTrafficItem[];
      expiry: {
        active: number;
        normalAfter30d: number;
        expiring7d: number;
        expiring30d: number;
        expired: number;
      };
      events: {
        pendingCount: number;
        items: DashboardEventItem[];
      };
      finance: {
        monthIncome: number;
        cashIncome: number;
        monthCost: number;
        grossProfit: number;
        currency: string;
      };
    }
  }
}
