<script setup lang="ts">
import { computed, ref } from 'vue';
import { downloadSystemBackup, restoreSystemBackup } from '@/service/api';

defineOptions({ name: 'SystemBackup' });

const downloading = ref(false);
const restoring = ref(false);
const selectedFile = ref<File | null>(null);
const fileInput = ref<HTMLInputElement | null>(null);
const copiedKey = ref('');

const backupContents = [
  '业务用户与 Client',
  '节点与 Inbound',
  '出口 IP 与路径',
  '账单与收费记录',
  '流量快照与汇总',
  '事件与审计日志',
  '系统配置'
];

const cliCommands = [
  {
    key: 'verify',
    title: '完整性校验',
    desc: '校验数据库或备份文件是否损坏，恢复前建议先跑一次。',
    command: 'xpanel-central-db verify --path ./data/panel.db'
  },
  {
    key: 'rollup',
    title: '重建小时流量汇总',
    desc: '汇总表与原始快照不一致时，用该命令重算（不影响原始数据）。',
    command: 'xpanel-central-db traffic-rollup --database ./data/panel.db'
  },
  {
    key: 'backup',
    title: '手动备份（保留 14 份）',
    desc: '在服务外生成一份带保留策略的备份，适合脚本化定时任务。',
    command: 'xpanel-central-db backup --database ./data/panel.db --backup-dir ./data/backups --retention 14'
  }
];

const selectedName = computed(() => selectedFile.value?.name || '');
const selectedSize = computed(() => {
  const size = selectedFile.value?.size || 0;
  if (!size) return '';
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / 1024 / 1024).toFixed(1)} MB`;
});

function selectFile(event: Event) {
  const input = event.target as HTMLInputElement;
  selectedFile.value = input.files?.[0] || null;
}

function clearFile() {
  selectedFile.value = null;
  if (fileInput.value) fileInput.value.value = '';
}

async function downloadBackup() {
  downloading.value = true;
  try {
    await downloadSystemBackup();
    window.$message?.success('备份已生成并开始下载');
  } catch (error) {
    window.$message?.error(error instanceof Error ? error.message : '生成备份失败');
  } finally {
    downloading.value = false;
  }
}

function confirmRestore() {
  if (!selectedFile.value || restoring.value) return;
  window.$dialog?.warning({
    title: '确认恢复数据',
    content:
      '恢复会覆盖当前所有用户、节点、入站、出口 IP、账单、流量和审计数据。系统会先保留恢复前安全快照，完成后当前登录会话将失效，需要重新登录。确定继续吗？',
    positiveText: '确认恢复',
    negativeText: '取消',
    maskClosable: false,
    onPositiveClick: restoreBackup
  });
}

async function restoreBackup() {
  if (!selectedFile.value) return;
  restoring.value = true;
  try {
    await restoreSystemBackup(selectedFile.value);
    window.$message?.success('数据恢复成功，请重新登录');
    clearFile();
    window.setTimeout(() => window.location.reload(), 900);
  } catch (error) {
    window.$message?.error(error instanceof Error ? error.message : '恢复数据失败，当前数据未改变');
  } finally {
    restoring.value = false;
  }
}

async function copyCommand(key: string, command: string) {
  try {
    await navigator.clipboard.writeText(command);
    copiedKey.value = key;
    window.$message?.success('命令已复制');
    window.setTimeout(() => {
      if (copiedKey.value === key) copiedKey.value = '';
    }, 2000);
  } catch {
    window.$message?.error('复制失败，请手动选择命令文本');
  }
}
</script>

<template>
  <div class="b-page">
    <header class="b-head">
      <div>
        <div class="b-eyebrow">系统 · 数据备份 · SQLite 全量快照</div>
        <h1 class="b-title">数据备份</h1>
        <p class="b-desc">
          备份和迁移 Alien-Panel 的完整业务数据库。备份文件包含全部历史数据与系统配置，恢复只接受本系统生成且通过完整性校验的 SQLite 文件。
        </p>
      </div>
      <div class="b-actions">
        <button class="b-btn is-primary" type="button" :disabled="downloading" @click="downloadBackup">
          <SvgIcon :icon="downloading ? 'mdi:loading' : 'mdi:download'" class="b-btn-ico" :class="{ 'is-spin': downloading }" />
          {{ downloading ? '生成中…' : '创建并下载备份' }}
        </button>
      </div>
    </header>

    <div class="b-callout is-danger">
      <SvgIcon icon="mdi:alert-outline" class="b-callout-ico" />
      <div>
        <b>恢复是覆盖性操作，不可撤销</b>
        恢复会替换当前数据库中的全部业务数据，并保留一份恢复前的安全快照；完成后所有登录会话失效，需要重新登录。
      </div>
    </div>

    <div class="b-grid">
      <section class="b-card">
        <div class="b-card-head">
          <h3>创建备份</h3>
          <span class="b-badge is-ok">推荐在变更前执行</span>
        </div>
        <div class="b-card-body">
          <p class="b-note">
            生成当前数据库的一致性快照并立即下载。下载文件可直接用于恢复，也可用 CLI 在服务器本地生成带保留策略的备份。
          </p>
          <div class="b-sub-title">备份包含</div>
          <div class="b-chips">
            <span v-for="item in backupContents" :key="item" class="b-chip">{{ item }}</span>
          </div>
          <div class="b-foot">
            <button class="b-btn is-primary" type="button" :disabled="downloading" @click="downloadBackup">
              <SvgIcon icon="mdi:database-export" class="b-btn-ico" />
              立即生成并下载
            </button>
            <span class="b-dim">生成过程在服务端完成，不会中断线上请求</span>
          </div>
        </div>
      </section>

      <section class="b-card">
        <div class="b-card-head">
          <h3>从备份恢复</h3>
          <span class="b-badge is-danger">覆盖当前数据</span>
        </div>
        <div class="b-card-body">
          <p class="b-note">
            选择一个本系统生成的 <code>.sqlite3</code> / <code>.db</code> 备份文件上传恢复。上传后先校验完整性，校验通过才会执行替换。
          </p>

          <label class="b-drop" :class="{ 'has-file': selectedFile }">
            <input
              ref="fileInput"
              type="file"
              accept=".sqlite3,.db,application/vnd.sqlite3"
              class="b-file"
              @change="selectFile"
            />
            <SvgIcon :icon="selectedFile ? 'mdi:file-check-outline' : 'mdi:tray-arrow-up'" class="b-drop-ico" />
            <div class="b-drop-main">
              <b>{{ selectedFile ? selectedName : '点击选择备份文件' }}</b>
              <span v-if="selectedFile">{{ selectedSize }} · 已选择，等待上传</span>
              <span v-else>仅接受 .sqlite3 / .db，且必须通过完整性校验</span>
            </div>
            <span class="b-drop-action">{{ selectedFile ? '重新选择' : '浏览文件' }}</span>
          </label>

          <div class="b-foot">
            <button class="b-btn is-danger" type="button" :disabled="!selectedFile || restoring" @click="confirmRestore">
              <SvgIcon icon="mdi:database-import" class="b-btn-ico" />
              {{ restoring ? '恢复中…' : '上传并恢复' }}
            </button>
            <button v-if="selectedFile" class="b-btn" type="button" :disabled="restoring" @click="clearFile">清除选择</button>
            <span class="b-dim">恢复完成后需要重新登录</span>
          </div>
        </div>
      </section>
    </div>

    <section class="b-card">
      <div class="b-card-head">
        <h3>运维命令（服务器本地执行）</h3>
        <span class="b-dim">以下操作没有 HTTP 接口，需在服务器上通过 CLI 执行</span>
      </div>
      <div class="b-card-body">
        <div v-for="item in cliCommands" :key="item.key" class="b-cli-row">
          <div class="b-cli-main">
            <b>{{ item.title }}</b>
            <span>{{ item.desc }}</span>
          </div>
          <code class="b-code">{{ item.command }}</code>
          <button class="b-btn b-btn-xs" type="button" @click="copyCommand(item.key, item.command)">
            <SvgIcon :icon="copiedKey === item.key ? 'mdi:check' : 'mdi:content-copy'" class="b-btn-ico" />
            {{ copiedKey === item.key ? '已复制' : '复制' }}
          </button>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.b-page {
  --b-card: #141519;
  --b-line: #23262c;
  --b-line-soft: #1c1f24;
  --b-text-1: #f2f3f5;
  --b-text-2: #bcc2cc;
  --b-text-3: #9aa1ad;
  --b-text-4: #868d99;
  --b-brand: #4c6fff;
  --b-ok: #2ed47a;
  --b-warn: #f5b93d;
  --b-danger: #f3656f;
  --b-radius: 14px;

  display: flex;
  flex-direction: column;
  gap: 15px;
}

/* ---------- 页头 ---------- */
.b-head {
  display: flex;
  align-items: flex-end;
  gap: 18px;
  flex-wrap: wrap;
}

.b-eyebrow {
  margin-bottom: 6px;
  font-size: 11.5px;
  letter-spacing: 0.12em;
  color: #7f8792;
}

.b-title {
  margin: 0;
  font-size: 23px;
  font-weight: 600;
  color: var(--b-text-1);
}

.b-desc {
  margin: 7px 0 0;
  max-width: 820px;
  font-size: 13px;
  line-height: 1.72;
  color: var(--b-text-3);
}

.b-actions {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
  padding-bottom: 2px;
}

/* ---------- 按钮 ---------- */
.b-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 32px;
  padding: 0 13px;
  border: 1px solid #2a2e35;
  border-radius: 9px;
  background: #1b1e24;
  color: #e6e8ec;
  font-size: 12.5px;
  cursor: pointer;
  transition: background 0.15s ease, border-color 0.15s ease;
}

.b-btn:hover:not(:disabled) {
  background: #232730;
  border-color: #363b44;
}

.b-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.b-btn.is-primary {
  border-color: rgba(76, 111, 255, 0.55);
  background: linear-gradient(180deg, #5b7dff, #3f5ee8);
  color: #fff;
}

.b-btn.is-danger {
  border-color: rgba(243, 101, 111, 0.5);
  background: rgba(243, 101, 111, 0.16);
  color: #ffb3b9;
}

.b-btn-xs {
  height: 27px;
  padding: 0 9px;
  font-size: 12px;
}

.b-btn-ico {
  width: 14px;
  height: 14px;
}

.is-spin {
  animation: b-spin 1s linear infinite;
}

@keyframes b-spin {
  to {
    transform: rotate(360deg);
  }
}

/* ---------- 提示条 ---------- */
.b-callout {
  display: flex;
  gap: 10px;
  padding: 12px 14px;
  border: 1px solid var(--b-line);
  border-radius: 11px;
  background: rgba(255, 255, 255, 0.022);
  font-size: 12px;
  line-height: 1.7;
  color: var(--b-text-2);
}

.b-callout-ico {
  width: 16px;
  height: 16px;
  flex: none;
  margin-top: 2px;
  color: var(--b-text-3);
}

.b-callout b {
  display: block;
  margin-bottom: 1px;
  font-size: 12.5px;
  font-weight: 600;
  color: var(--b-text-1);
}

.b-callout.is-danger {
  border-color: rgba(243, 101, 111, 0.28);
  background: rgba(243, 101, 111, 0.08);
}

.b-callout.is-danger .b-callout-ico {
  color: #ff9ba2;
}

/* ---------- 卡片 ---------- */
.b-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 15px;
}

.b-card {
  border: 1px solid var(--b-line);
  border-radius: var(--b-radius);
  background: var(--b-card);
  min-width: 0;
}

.b-card-head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 18px;
  border-bottom: 1px solid var(--b-line-soft);
}

.b-card-head h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--b-text-1);
  white-space: nowrap;
}

.b-card-head > .b-dim {
  margin-left: auto;
}

.b-card-body {
  padding: 16px 18px 18px;
}

.b-note {
  margin: 0 0 14px;
  font-size: 12.5px;
  line-height: 1.75;
  color: var(--b-text-3);
}

.b-note code {
  padding: 1px 5px;
  border-radius: 5px;
  background: #1d2026;
  color: #b8bec8;
  font-size: 11.5px;
}

.b-sub-title {
  margin-bottom: 8px;
  font-size: 12px;
  color: var(--b-text-3);
}

.b-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.b-chip {
  display: inline-flex;
  align-items: center;
  height: 24px;
  padding: 0 9px;
  border: 1px solid #2a2e35;
  border-radius: 7px;
  background: #1a1d22;
  font-size: 11.5px;
  color: var(--b-text-2);
}

.b-badge {
  display: inline-flex;
  align-items: center;
  height: 21px;
  padding: 0 8px;
  border: 1px solid #2a2e35;
  border-radius: 6px;
  background: #1d2026;
  font-size: 11.5px;
  color: var(--b-text-2);
  white-space: nowrap;
}

.b-badge.is-ok {
  border-color: rgba(46, 212, 122, 0.32);
  background: rgba(46, 212, 122, 0.14);
  color: #5ce49b;
}

.b-badge.is-danger {
  border-color: rgba(243, 101, 111, 0.32);
  background: rgba(243, 101, 111, 0.14);
  color: #ff9ba2;
}

.b-foot {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  margin-top: 16px;
  padding-top: 14px;
  border-top: 1px solid var(--b-line-soft);
}

.b-dim {
  font-size: 11.5px;
  color: var(--b-text-4);
}

/* ---------- 上传区 ---------- */
.b-drop {
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 16px;
  border: 1px dashed #333841;
  border-radius: 11px;
  background: #101114;
  cursor: pointer;
  transition: border-color 0.16s ease, background 0.16s ease;
}

.b-drop:hover {
  border-color: rgba(76, 111, 255, 0.45);
  background: rgba(76, 111, 255, 0.06);
}

.b-drop.has-file {
  border-style: solid;
  border-color: rgba(46, 212, 122, 0.4);
  background: rgba(46, 212, 122, 0.06);
}

.b-file {
  position: absolute;
  inset: 0;
  opacity: 0;
  cursor: pointer;
}

.b-drop-ico {
  width: 22px;
  height: 22px;
  flex: none;
  color: var(--b-text-3);
}

.b-drop.has-file .b-drop-ico {
  color: #5ce49b;
}

.b-drop-main {
  min-width: 0;
  flex: 1;
}

.b-drop-main b {
  display: block;
  font-size: 12.5px;
  font-weight: 550;
  color: var(--b-text-1);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.b-drop-main span {
  display: block;
  margin-top: 2px;
  font-size: 11.5px;
  color: var(--b-text-4);
}

.b-drop-action {
  flex: none;
  font-size: 12px;
  color: #9db2ff;
}

/* ---------- 运维命令 ---------- */
.b-cli-row {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 0;
  border-bottom: 1px solid var(--b-line-soft);
}

.b-cli-row:last-child {
  border-bottom: 0;
  padding-bottom: 0;
}

.b-cli-main {
  width: 190px;
  flex: none;
}

.b-cli-main b {
  display: block;
  font-size: 12.5px;
  font-weight: 550;
  color: var(--b-text-1);
}

.b-cli-main span {
  display: block;
  margin-top: 2px;
  font-size: 11.5px;
  line-height: 1.6;
  color: var(--b-text-4);
}

.b-code {
  flex: 1;
  min-width: 0;
  padding: 8px 11px;
  border: 1px solid #23262c;
  border-radius: 9px;
  background: #0d0e11;
  color: #b8bec8;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 11.5px;
  white-space: nowrap;
  overflow-x: auto;
}

@media (max-width: 1100px) {
  .b-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .b-cli-row {
    flex-wrap: wrap;
  }

  .b-cli-main {
    width: 100%;
  }
}
</style>
