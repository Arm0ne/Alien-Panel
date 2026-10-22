<script setup lang="ts">
import { computed, reactive } from 'vue';
import { useAuthStore } from '@/store/modules/auth';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';

defineOptions({
  name: 'PwdLogin'
});

const authStore = useAuthStore();
const { formRef, validate } = useNaiveForm();

interface FormModel {
  userName: string;
  password: string;
}

const model: FormModel = reactive({
  userName: '',
  password: ''
});

const rules = computed<Record<keyof FormModel, App.Global.FormRule[]>>(() => {
  const { formRules } = useFormRules();

  return {
    userName: [formRules.userName[0]],
    // 密码由中央后端验证；不要沿用模板的字符集规则，否则包含连字符等
    // 合法字符的部署密码会在浏览器中被拦截而无法提交。
    password: [formRules.pwd[0]]
  };
});

async function handleSubmit() {
  try {
    await validate();
  } catch {
    return;
  }
  await authStore.login(model.userName, model.password);
}
</script>

<template>
  <NForm
    ref="formRef"
    class="pwd-login"
    :model="model"
    :rules="rules"
    size="large"
    :show-label="true"
    label-placement="top"
    @keyup.enter="handleSubmit"
  >
    <NFormItem label="用户名" path="userName" :show-require-mark="false">
      <NInput v-model:value="model.userName" placeholder="请输入用户名" :input-props="{ autocomplete: 'username' }" />
    </NFormItem>

    <NFormItem label="密码" path="password" :show-require-mark="false">
      <NInput
        v-model:value="model.password"
        type="password"
        show-password-on="click"
        placeholder="请输入密码"
        :input-props="{ autocomplete: 'current-password' }"
      />
    </NFormItem>

    <NButton
      class="submit"
      type="primary"
      size="large"
      block
      :loading="authStore.loginLoading"
      @click="handleSubmit"
    >
      登 录
    </NButton>
  </NForm>
</template>

<style scoped>
.pwd-login {
  --n-border-radius: 8px;
}

/* 字段间距与标签样式，对齐视觉稿 */
.pwd-login :deep(.n-form-item) {
  margin-bottom: 20px;
}

.pwd-login :deep(.n-form-item:last-of-type) {
  margin-bottom: 28px;
}

.pwd-login :deep(.n-form-item-label) {
  margin-bottom: 9px;
  font-size: 13px;
  font-weight: 500;
  letter-spacing: 0.01em;
  color: #d8d8e0;
  padding: 0;
  line-height: 1.4;
}

.pwd-login :deep(.n-form-item-blank) {
  min-height: auto;
}

/* 错误提示：让它正常占据空间显示出来，用负 margin 抵消默认占位，
   这样没错误时不会多出空隙，有错误时提示能正常可见 */
.pwd-login :deep(.n-form-item-feedback-wrapper) {
  min-height: 0;
}

.pwd-login :deep(.n-form-item-feedback) {
  font-size: 12px;
  line-height: 1.5;
  padding-top: 5px;
}

/* 输入框：浅色实心，与深色面板形成对比 */
.pwd-login :deep(.n-input) {
  --n-height: 48px;
  --n-color: #f4f4f6;
  --n-color-focus: #ffffff;
  --n-text-color: #101015;
  --n-placeholder-color: #9d9da8;
  --n-border: 1px solid transparent;
  --n-border-hover: 1px solid transparent;
  --n-border-focus: 1px solid #7683f0;
  --n-box-shadow-focus: 0 0 0 3px rgba(102, 112, 230, 0.28);
  --n-border-radius: 8px;
  font-size: 14.5px;
}

.pwd-login :deep(.n-input:hover) {
  --n-color: #fafafb;
}

.pwd-login :deep(.n-input__input-el) {
  height: 46px;
}

/* 登录按钮 */
.submit {
  width: 100%;
  height: 48px;
  --n-border-radius: 8px;
  --n-color: #5961c8;
  --n-color-hover: #636bd8;
  --n-color-pressed: #4c54b8;
  --n-color-focus: #5961c8;
  --n-text-color: #ffffff;
  --n-text-color-hover: #ffffff;
  --n-text-color-pressed: #ffffff;
  --n-text-color-focus: #ffffff;
  font-size: 14.5px;
  font-weight: 600;
  letter-spacing: 0.02em;
  box-shadow:
    0 12px 28px -12px rgba(78, 89, 207, 0.78),
    0 0 0 1px rgba(255, 255, 255, 0.05) inset;
  transition: filter 0.16s ease, box-shadow 0.16s ease;
}

.submit:hover {
  filter: brightness(1.1);
  box-shadow:
    0 16px 34px -12px rgba(78, 89, 207, 0.92),
    0 0 0 1px rgba(255, 255, 255, 0.08) inset;
}

@media (max-width: 900px) {
  .pwd-login :deep(.n-input) {
    --n-height: 46px;
  }

  .submit {
    height: 46px;
  }
}
</style>
