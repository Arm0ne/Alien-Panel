<script setup lang="ts">
import { computed, reactive, ref } from 'vue';
import type { VNode } from 'vue';
import type { FormInst } from 'naive-ui';
import { useAuthStore } from '@/store/modules/auth';
import { useRouterPush } from '@/hooks/common/router';
import { useSvgIcon } from '@/hooks/common/icon';
import { useFormRules } from '@/hooks/common/form';
import { updateAuthAccount } from '@/service/api';
import { $t } from '@/locales';

defineOptions({
  name: 'UserAvatar'
});

const authStore = useAuthStore();
const { toLogin } = useRouterPush();
const { SvgIconVNode } = useSvgIcon();
const accountVisible = ref(false);
const accountSaving = ref(false);
const accountFormRef = ref<FormInst | null>(null);
const accountForm = reactive({
  currentPassword: '',
  userName: '',
  newPassword: '',
  confirmPassword: ''
});

const accountRules = computed(() => {
  const { formRules, createRequiredRule } = useFormRules();
  return {
    currentPassword: [createRequiredRule($t('form.pwd.required'))],
    userName: formRules.userName,
    newPassword: [
      {
        validator: (_rule: unknown, value: string) => {
          if (!value || (value.length >= 6 && value.length <= 128)) return Promise.resolve();
          return Promise.reject($t('form.pwd.invalid'));
        },
        trigger: 'input'
      }
    ],
    confirmPassword: [
      {
        validator: (_rule: unknown, value: string) => {
          if (!accountForm.newPassword && !value) return Promise.resolve();
          if (accountForm.newPassword && value === accountForm.newPassword) return Promise.resolve();
          return Promise.reject($t('form.confirmPwd.invalid'));
        },
        trigger: 'input'
      }
    ]
  };
});

function loginOrRegister() {
  toLogin();
}

type DropdownKey = 'account' | 'logout';

type DropdownOption =
  | {
      key: DropdownKey;
      label: string;
      icon?: () => VNode;
    }
  | {
      type: 'divider';
      key: string;
    };

const options = computed(() => {
  const opts: DropdownOption[] = [
    {
      label: $t('common.accountSettings'),
      key: 'account',
      icon: SvgIconVNode({ icon: 'ph:user-gear', fontSize: 18 })
    },
    {
      label: $t('common.logout'),
      key: 'logout',
      icon: SvgIconVNode({ icon: 'ph:sign-out', fontSize: 18 })
    }
  ];

  return opts;
});

function logout() {
  window.$dialog?.info({
    title: $t('common.tip'),
    content: $t('common.logoutConfirm'),
    positiveText: $t('common.confirm'),
    negativeText: $t('common.cancel'),
    onPositiveClick: () => {
      authStore.resetStore();
    }
  });
}

function openAccountSettings() {
  accountForm.currentPassword = '';
  accountForm.userName = authStore.userInfo.userName;
  accountForm.newPassword = '';
  accountForm.confirmPassword = '';
  accountFormRef.value?.restoreValidation();
  accountVisible.value = true;
}

function closeAccountSettings() {
  if (!accountSaving.value) accountVisible.value = false;
}

async function saveAccountSettings() {
  try {
    await accountFormRef.value?.validate();
  } catch {
    return;
  }
  accountSaving.value = true;
  const { error } = await updateAuthAccount({
    currentPassword: accountForm.currentPassword,
    userName: accountForm.userName,
    newPassword: accountForm.newPassword || undefined
  });
  accountSaving.value = false;
  if (error) return;

  accountVisible.value = false;
  window.$notification?.success({
    title: $t('common.modifySuccess'),
    content: $t('common.accountUpdateSuccess'),
    duration: 4500
  });
  await authStore.resetStore();
}

function handleDropdown(key: DropdownKey) {
  if (key === 'account') {
    openAccountSettings();
  } else if (key === 'logout') {
    logout();
  }
}
</script>

<template>
  <NButton v-if="!authStore.isLogin" quaternary @click="loginOrRegister">
    {{ $t('page.login.common.loginOrRegister') }}
  </NButton>
  <NDropdown v-else placement="bottom" trigger="click" :options="options" @select="handleDropdown">
    <div>
      <ButtonIcon>
        <SvgIcon icon="ph:user-circle" class="text-icon-large" />
        <span class="text-16px font-medium">{{ authStore.userInfo.userName }}</span>
      </ButtonIcon>
    </div>
  </NDropdown>
  <NModal
    v-model:show="accountVisible"
    preset="card"
    :title="$t('common.accountSettingsTitle')"
    class="w-520px max-w-95vw"
    :mask-closable="false"
  >
    <NAlert type="info" :show-icon="true" class="mb-16px">{{ $t('common.accountUpdateHint') }}</NAlert>
    <NForm ref="accountFormRef" :model="accountForm" :rules="accountRules" label-placement="top">
      <NFormItem path="userName" :label="$t('common.userName')">
        <NInput v-model:value="accountForm.userName" maxlength="64" />
      </NFormItem>
      <NFormItem path="currentPassword" :label="$t('common.currentPassword')">
        <NInput
          v-model:value="accountForm.currentPassword"
          type="password"
          show-password-on="click"
          autocomplete="current-password"
        />
      </NFormItem>
      <NFormItem path="newPassword" :label="$t('common.newPassword')">
        <NInput
          v-model:value="accountForm.newPassword"
          type="password"
          show-password-on="click"
          autocomplete="new-password"
        />
      </NFormItem>
      <NFormItem path="confirmPassword" :label="$t('common.confirmNewPassword')">
        <NInput
          v-model:value="accountForm.confirmPassword"
          type="password"
          show-password-on="click"
          autocomplete="new-password"
        />
      </NFormItem>
    </NForm>
    <template #footer>
      <NSpace justify="end">
        <NButton :disabled="accountSaving" @click="closeAccountSettings">{{ $t('common.cancel') }}</NButton>
        <NButton type="primary" :loading="accountSaving" @click="saveAccountSettings">
          {{ $t('common.confirm') }}
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<style scoped></style>
