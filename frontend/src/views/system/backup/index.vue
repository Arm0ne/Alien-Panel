<script setup lang="ts">
import { ref } from 'vue';
import { downloadSystemBackup, restoreSystemBackup } from '@/service/api';
import ModulePage from '@/components/project/module-page.vue';
import SvgIcon from '@/components/custom/svg-icon.vue';

defineOptions({ name: 'SystemBackup' });

const downloading = ref(false);
const restoring = ref(false);
const selectedFile = ref<File | null>(null);
const fileInput = ref<HTMLInputElement | null>(null);

function selectFile(event: Event) {
  const input = event.target as HTMLInputElement;
  selectedFile.value = input.files?.[0] || null;
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
    selectedFile.value = null;
    if (fileInput.value) fileInput.value.value = '';
    window.setTimeout(() => window.location.reload(), 900);
  } catch (error) {
    window.$message?.error(error instanceof Error ? error.message : '恢复数据失败，当前数据未改变');
  } finally {
    restoring.value = false;
  }
}
</script>

<template>
  <ModulePage title="数据备份" description="备份和迁移 Alien-Panel 的完整业务数据库" :empty="false" :loading="false">
    <template #actions>
      <NButton type="primary" :loading="downloading" @click="downloadBackup">
        <template #icon><SvgIcon icon="mdi:download" /></template>
        创建并下载备份
      </NButton>
    </template>

    <div class="p-16px">
      <NAlert type="info" :show-icon="true" class="mb-16px">
        备份文件包含全部历史数据和系统配置，请妥善保管。恢复只接受本系统生成且通过完整性校验的 SQLite 文件。
      </NAlert>
      <div class="border border-gray-200 rounded-4px p-16px dark:border-gray-700">
        <div class="mb-12px text-16px font-600">从备份恢复</div>
        <NSpace vertical :size="12">
          <input ref="fileInput" type="file" accept=".sqlite3,.db,application/vnd.sqlite3" @change="selectFile" />
          <div v-if="selectedFile" class="text-13px text-gray-500">已选择：{{ selectedFile.name }}</div>
          <NSpace>
            <NButton type="warning" :disabled="!selectedFile" :loading="restoring" @click="confirmRestore">
              <template #icon><SvgIcon icon="mdi:database-import" /></template>
              上传并恢复
            </NButton>
            <NButton v-if="selectedFile" secondary :disabled="restoring" @click="selectedFile = null">
              清除选择
            </NButton>
          </NSpace>
        </NSpace>
      </div>
    </div>
  </ModulePage>
</template>
