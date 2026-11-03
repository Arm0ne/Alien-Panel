-- 用户在面板里手工改过名称后，不再被 Agent 同步覆盖。
-- 0 = 名称仍与 Agent 同步值一致，允许跟随 X-Panel 的 Inbound 名称更新
-- 1 = 已人工编辑，保留运营自定义的名称
ALTER TABLE users ADD COLUMN display_name_locked INTEGER NOT NULL DEFAULT 0;

-- Agent 上一次写入的名称基准。只有当 display_name 仍等于该基准时，
-- 才认为是「未人工编辑」，可以跟随同步更新；否则视为人工修改并加锁。
-- 这样通过面板接口或直接改库改名，都能被正确保护。
ALTER TABLE users ADD COLUMN display_name_auto TEXT NOT NULL DEFAULT '';
