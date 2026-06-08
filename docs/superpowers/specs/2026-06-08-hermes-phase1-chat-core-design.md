# Hermes USB — Phase 1: 基础聊天壳 设计文档

> 日期: 2026-06-08 | 状态: 设计完成，待审核

## 产品概述

Hermes 是一个即插即用的 U 盘 AI 助手。Phase 1 交付基础聊天能力：插上 U 盘、双击 exe、浏览器自动打开、扫码配置模型 Key、开始聊天。

## 里程碑路线图

| Phase | 内容 | 依赖 |
|-------|------|------|
| Phase 1 | 基础聊天壳（本 spec） | — |
| Phase 2 | 多层记忆系统 | Phase 1 |
| Phase 3 | Skill 自进化 | Phase 1, 2 |
| Phase 4 | 多平台聊天网关（QQ/飞书等） | Phase 1 |
| Phase 5 | 知识库 RAG | Phase 1, 2 |

## 技术栈

| 层 | 选择 | 理由 |
|---|---|---|
| 后端 | Go 1.22+ / Gin | 单二进制、跨平台编译、无运行时依赖 |
| 前端 | React + Tailwind CSS | 内嵌 SPA、生态丰富、流式渲染成熟 |
| 打包 | go:embed 内嵌 dist/ | 一个 exe 走天下 |
| 数据库 | SQLite (go-sqlite3) | 单文件、零配置、FTS5 全文索引 |
| AI 协议 | OpenAI 兼容 | 一键对接 DeepSeek/Kimi/Qwen/GLM 等 |
| 更新 | GitHub Releases | 免费 CDN、版本管理标准 |

## U 盘文件结构

```
Hermes/
├── hermes.exe          # 唯一可执行文件 (~35MB)
├── config.json          # 用户配置（首次运行自动生成）
└── data/
    └── hermes.db        # SQLite（聊天记录、记忆、技能数据）
```

---

## 架构

```
hermes.exe (Go 单二进制)
├── HTTP Server (Gin)
│   ├── /              → React SPA (go:embed)
│   ├── /api/chat      → 聊天 + SSE 流式输出
│   ├── /api/models     → 模型 CRUD
│   ├── /api/conversations → 会话管理
│   ├── /api/settings   → 配置读写
│   └── /api/health     → 健康检查
├── 核心服务层
│   ├── ChatService     → OpenAI 兼容协议调用 + SSE 转发
│   ├── ConfigService   → config.json 读写 + 默认值
│   ├── ConversationStore → SQLite 会话 & 消息 CRUD
│   └── UpdateService   → GitHub Releases 检查 + 下载 + 校验
└── 存储层
    ├── config.json     → JSON（可手改）
    └── data/hermes.db  → SQLite（WAL 模式）
```

### 关键设计约束

- HTTP 只监听 `127.0.0.1`，不暴露网络
- 固定端口 `21475`，冲突时自动 fallback 到随机端口
- 静默打开浏览器（Go `open.Start()`）
- 单实例运行：二次双击 exe 检测已有进程 → 直接打开浏览器

---

## 生命周期 & 系统托盘

### 状态机

```
[首次启动] → config.json 不存在 → 打开 /onboarding
                  ↓
[正常运行] → 托盘常驻 + SQLite 活跃
                  ↓
  双击 exe → 检测已有实例 → 打开浏览器到 /chat
  关闭浏览器 → 托盘仍在 → 右键/双击托盘可重开
  托盘右键 → "退出 Hermes" → 关闭 DB → 允许拔盘
[异常拔出] → 下次启动 WAL 恢复 → 提示"数据已自动恢复"
```

### 系统托盘菜单

```
┌─────────────────────┐
│ 🟢 Hermes 运行中     │
├─────────────────────┤
│ 🌐 打开聊天界面      │
│ 📊 查看状态          │
│ ─────────────────── │
│ 🔄 检查更新          │
│ ─────────────────── │
│ ⏹️ 退出 Hermes       │
└─────────────────────┘
```

- 双击托盘图标 → 打开聊天界面
- 退出时：关闭 SQLite 连接 → 清理临时文件 → 托盘消失 → 提示可安全拔盘

### 浏览器找回路径

1. 再次双击 `hermes.exe` → 检测已有实例 → 直接打开浏览器
2. 系统托盘右键 → "打开聊天界面"
3. 浏览器手动输入 `http://localhost:21475`

---

## UI 设计

### 视觉风格

暗色主题 + 渐变标签 + 气泡对话。

| 元素 | 设计 |
|------|------|
| 背景 | `#0a0a0b` 微渐变 |
| 用户气泡 | 紫色渐变 (`#6366f1 → #7c3aed`)，右对齐，带光晕 |
| AI 气泡 | 深灰卡片 (`#111113`)，左对齐，圆角 |
| 发送者标签 | 用户：紫色圆点 + "你"；AI：绿色方块 + "Hermes" |
| 品牌色 | 紫蓝渐变 (primary) / 翠绿 (success/ai) |
| 字体 | 系统无衬线（内容）+ SF Mono（代码） |

### 布局

自适应单栏：默认极简聊天窗，左侧汉堡菜单（☰）滑出功能面板。

**顶部栏**：Hermes logo + 绿色状态点 + 模型名 + 设置齿轮 + 汉堡菜单

**侧边面板**（渐进式）：
- Phase 1 只显示：+ 新对话、最近对话列表、设置
- Phase 2-5 的新功能随升级自然出现，不显示"即将上线"占位
- 底部叙事："Hermes 会越用越聪明"

### 设置页

卡片分组式，一页展示所有配置：
- 🤖 AI 模型卡片：模型选择下拉 + API Key + API 地址 + "测试连接"按钮
- 💬 聊天设置卡片：系统提示词
- ⚠️ 恢复出厂设置卡片：红色警告区，确认后删除所有数据

### 聊天空状态

首次打开聊天时显示：
- Hermes logo + "有什么可以帮你的？"
- 6 个建议问题 chips（根据时间段变化：早上→今日计划，下午→日报总结，晚上→复盘反思。Phase 2+ 加入历史对话个性化）
- 输入框常驻底部

---

## Onboarding 流程

3 步向导，首次启动自动触发：

| 步骤 | 页面 | 内容 |
|------|------|------|
| 1/3 | 欢迎页 | 品牌介绍 + 三个卖点（即插即用/数据本地/越用越聪明）→ "开始配置" |
| 2/3 | 选择模型 | 模型卡片列表（DeepSeek/Kimi/Qwen，预选 DeepSeek）+ API Key 输入 + 获取 Key 教程链接 → "测试连接" |
| 3/3 | 完成 | ✓ 动画 + "配置完成" + 提示"以后双击直接聊天" → "开始聊天" |

每步有进度条指示器（3 个点）。

---

## 流式输出

- SSE (Server-Sent Events) 逐 token 推送
- 前端逐字渲染，紫色闪烁光标表示生成中
- Markdown 渐进解析（代码块骨架占位）
- 用户可随时"停止生成"；停止后可"继续生成"
- 自动滚动到底部；用户手动上滚时暂停，回到底部恢复
- 等待响应时显示三点思考动画

---

## 错误处理

每个错误态包含：图标 + 错误描述 + 用户可理解的原因 + 操作按钮。

| 错误场景 | 文案 | 操作 |
|---------|------|------|
| API Key 未配置 | "需要配置 API Key" | 前往配置 |
| 网络不通 | "无法连接 AI 服务" | 检查网络 / 修改配置 |
| 模型返回错误 | 显示原始错误码 + 翻译 | 重试 / 切换模型 |
| 数据库异常 | "上次未正常退出，数据已自动恢复" | 自动修复，仅提示 |
| 端口冲突 | "端口已被占用，已切换到 XXXX" | 自动 fallback，仅提示 |

---

## 更新机制

- 更新源：GitHub Releases
- 检查频率：启动时 + 每 24 小时后台轮询
- 有更新时托盘图标变更（绿点 → 蓝点 + 气泡提示）
- 更新流程：下载 `hermes.new` → SHA256 校验 → 旧文件重命名 `.bak` → 新文件替换 → 自动重启
- 更新失败：回滚到 `.bak`，提示错误

---

## 恢复出厂设置

- 入口：设置页底部红色警告卡片
- 确认：弹窗二次确认 "此操作不可撤销，将删除所有配置和对话历史"
- 行为：停止服务 → 删除 `config.json` + `data/` 目录 → 重启 exe → 进入 onboarding 流程

---

## API 设计

### `/api/chat` (POST, SSE)

```json
// Request
{
  "conversation_id": "uuid",
  "message": "帮我写周报",
  "model": "deepseek-v3"
}

// SSE Response Stream
event: token
data: {"token": "好"}

event: token
data: {"token": "的"}

event: done
data: {"total_tokens": 156, "conversation_id": "uuid"}
```

### `/api/conversations` (REST)

- `GET /api/conversations` — 会话列表
- `POST /api/conversations` — 新建会话
- `GET /api/conversations/:id` — 会话详情 + 消息列表
- `DELETE /api/conversations/:id` — 删除会话

### `/api/models` (REST)

- `GET /api/models` — 模型列表
- `PUT /api/models/:id` — 更新模型配置
- `POST /api/models/:id/test` — 测试连接

### `/api/settings` (REST)

- `GET /api/settings` — 获取所有设置
- `PUT /api/settings` — 更新设置
- `POST /api/settings/reset` — 恢复出厂

### `/api/health` (GET)

```json
{
  "status": "ok",
  "version": "0.1.0",
  "model": "deepseek-v3",
  "model_connected": true,
  "uptime": 3600
}
```

---

## 数据模型

### config.json

```json
{
  "version": 1,
  "models": [
    {
      "id": "default",
      "name": "DeepSeek V3",
      "api_base": "https://api.deepseek.com",
      "api_key": "sk-xxx",
      "is_default": true
    }
  ],
  "chat": {
    "system_prompt": "你是一个有用的AI助手"
  }
}
```

### SQLite 表结构

```sql
CREATE TABLE conversations (
  id TEXT PRIMARY KEY,
  title TEXT,
  created_at INTEGER,
  updated_at INTEGER
);

CREATE TABLE messages (
  id TEXT PRIMARY KEY,
  conversation_id TEXT REFERENCES conversations(id),
  role TEXT,  -- 'user' | 'assistant'
  content TEXT,
  tokens_used INTEGER,
  created_at INTEGER
);

-- Phase 1 模型配置存储在 config.json 中，SQLite 仅存会话数据
-- Phase 2+ 模型配置迁移到 SQLite
```

---

## CLI 设计

Phase 1 CLI 最小集：

```bash
hermes.exe              # 启动服务 + 打开浏览器
hermes.exe --version     # 输出版本号
hermes.exe --reset       # 恢复出厂设置（确认交互）
hermes.exe --port 21475  # 指定端口启动
hermes.exe --no-browser  # 只启动服务，不打开浏览器
```

---

## 测试策略

| 层级 | 范围 | 工具 |
|------|------|------|
| 单元测试 | Go 核心服务（ChatService, ConfigService, ConversationStore） | Go testing |
| 集成测试 | API 端点 + SQLite 读写 + SSE 流 | Go testing + httptest |
| 前端测试 | React 组件渲染 + 流式输出 | Vitest + Testing Library |
| E2E | 完整 onboarding → 聊天流程 | Playwright |
| 平台测试 | Windows 10/11 实际 U 盘插拔 | 手动 |

---

## 不在 Phase 1 范围内

- 记忆系统（Phase 2）
- Skill 自动生成（Phase 3）
- 聊天平台网关 QQ/飞书（Phase 4）
- 知识库 RAG（Phase 5）
- 多语言国际化
- 移动端适配
- macOS/Linux 原生支持（Phase 1 仅 Windows；Go 跨平台编译保留扩展能力）
