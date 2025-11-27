# MCP Novel Assistant 使用指南

MCP Novel Assistant 是一个专为小说创作设计的智能管理工具，通过 MCP（Model Context Protocol）以 JSON-RPC 2.0 over stdio 方式提供丰富的工具（tools）。它帮助作者管理复杂的小说结构、人物关系、时间线、物品与能力流转、情节线索、冲突检测以及纲要生成等核心工作。

## 特性概览

- 小说结构：支持 小说→分卷→章节→事件 的层次化管理
- 时间线：支持 世界→时期→时间段→事件 的时间轴管理
- 人物关系：支持双向亲密度与社会关系管理
- 人物记忆：记录人物在事件中的记忆与触发条件
- 物品流转：物品归属与事件流转记录
- 人物能力：能力创建、升级与使用记录
- 情节线索：多线索并行，阶段含 开始/进行中/关键点/结束
- 冲突检测：提供 10 类冲突检测入口（可扩展）
- 纲要生成：章节细纲、分卷总纲、小说总纲自动生成
- 文风参考：支持为小说设置参考正文用于风格模仿
- 持久化：内置 SQLite（`novel.db`），启动时自动迁移模型

## 系统要求

- Go 1.21+
- macOS（当前工作环境），其他平台可自行构建

## 构建与启动

```bash
# 拉取依赖
go mod tidy

# 编译全部包
go build ./...

# 生成可执行文件
go build -o mcp-novel ./cmd/mcp-novel

# 启动（从标准输入读入 JSON-RPC 请求）
./mcp-novel
```

启动后程序会在当前目录创建并使用 `novel.db`（SQLite 数据库），并自动执行模型迁移。模型定义参见 `internal/models/models.go:1`。

## MCP 交互模型

MCP 交互遵循 JSON-RPC 2.0，通过标准输入/输出进行：

- 初始化：`initialize`
- 列出工具：`tools/list`
- 调用工具：`tools/call`

示例（使用 shell 管道发送请求）：

```bash
# 初始化
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./mcp-novel

# 列出所有工具
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | ./mcp-novel
```

工具调用统一格式：

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "<toolName>",
    "arguments": { /* tool-specific args */ }
  }
}
```

## 可用工具与参数

- `dbHelper` 数据库管理
  - `action`: `init|export`
  - `path`: `string`（导出路径占位，当前返回默认 `novel.db`）
- `sqlHelper` SQL 操作（占位入口）
  - `entity`: `string`，`action`: `create|update|delete|get`，`data`: `object`
- `novelHelper` 小说管理
  - `action`: `create|get|export|outline`
  - `title`: `string`，`description`: `string`，`id`: `number`
- `volumeHelper` 分卷管理
  - `action`: `create`，`novelID`: `number`，`title`: `string`，`index`: `number`
- `chapterHelper` 章节管理
  - `action`: `create|update|export|outline`
  - `volumeID`: `number`，`title`: `string`，`index`: `number`，`status`: `string`，`id`: `number`，`content`: `string`
- `eventHelper` 事件管理
  - `action`: `create`
  - `chapterID|worldID|locationID|timeSegmentID`: `number`
  - `description`: `string`，`characters`: `number[]`，`items`: `number[]`
- `worldHelper` 世界管理
  - `action`: `create`，`name`: `string`，`description`: `string`
- `periodHelper` 时期管理
  - `action`: `create`，`worldID`: `number`，`name`: `string`，`index`: `number`
- `timeSegmentHelper` 时间段管理
  - `action`: `create`，`periodID`: `number`，`name`: `string`，`start|end`: `RFC3339 字符串`
- `characterHelper` 人物管理
  - `action`: `create`，`name`: `string`，`bio`: `string`
- `characterRelationshipHelper` 人物关系管理（双向）
  - `action`: `set`，`aid|bid`: `number`，`type`: `string`，`intimacy`: `number`
- `locationHelper` 地点管理
  - `action`: `create`，`worldID`: `number`，`name`: `string`，`description`: `string`
- `itemHelper` 物品管理
  - `action`: `create|transfer`，`name|ownerID|locationID|status` 或 `itemID|fromID|toID|eventID`
- `characterAbilityHelper` 人物能力管理
  - `action`: `create|upgrade|use`，`characterID|name|level` 或 `abilityID|level|eventID|note`
- `plotThreadHelper` 情节线索管理
  - `action`: `create|update`，`novelID|name|stage` 或 `plotID|stage`
- `characterMemoryHelper` 人物记忆管理
  - `action`: `create`，`characterID|eventID|content|trigger`
- `conflictDetectionHelper` 冲突检测
  - `action`: `run`（返回冲突列表 `[{type,detail}]`）
- `outlineGeneratorHelper` 纲要生成
  - `action`: `chapter|volume|novel`，`id`: `number`（返回 `outline` 字符串）
- `articleExportHelper` 文章导出
  - `action`: `chapter|volume|novel`，`id`: `number`（返回导出文本）
- `styleHelper` 文笔风格参考
  - `action`: `set|get`，`novelID`: `number`，`content`: `string`

## 快速上手示例

以下示例演示一个最小的世界与小说结构搭建流程：

```bash
# 1) 创建世界
echo '{"jsonrpc":"2.0","id":10,"method":"tools/call","params":{"name":"worldHelper","arguments":{"action":"create","name":"奇幻大陆","description":"魔法与王国"}}}' | ./mcp-novel

# 假设返回 World.ID = 1

# 2) 创建时期
echo '{"jsonrpc":"2.0","id":11,"method":"tools/call","params":{"name":"periodHelper","arguments":{"action":"create","worldID":1,"name":"第一纪元","index":1}}}' | ./mcp-novel

# 假设返回 Period.ID = 1

# 3) 创建时间段（RFC3339 格式）
echo '{"jsonrpc":"2.0","id":12,"method":"tools/call","params":{"name":"timeSegmentHelper","arguments":{"action":"create","periodID":1,"name":"王朝建立","start":"2025-01-01T00:00:00Z","end":"2025-12-31T23:59:59Z"}}}' | ./mcp-novel

# 假设返回 TimeSegment.ID = 1

# 4) 创建地点
echo '{"jsonrpc":"2.0","id":13,"method":"tools/call","params":{"name":"locationHelper","arguments":{"action":"create","worldID":1,"name":"都城","description":"王室所在"}}}' | ./mcp-novel

# 假设返回 Location.ID = 1

# 5) 创建人物
echo '{"jsonrpc":"2.0","id":14,"method":"tools/call","params":{"name":"characterHelper","arguments":{"action":"create","name":"艾伦","bio":"年轻的骑士"}}}' | ./mcp-novel
echo '{"jsonrpc":"2.0","id":15,"method":"tools/call","params":{"name":"characterHelper","arguments":{"action":"create","name":"莉亚","bio":"王国公主"}}}' | ./mcp-novel

# 假设返回 Character.ID = 1,2

# 6) 建立人物关系（双向）
echo '{"jsonrpc":"2.0","id":16,"method":"tools/call","params":{"name":"characterRelationshipHelper","arguments":{"action":"set","aid":1,"bid":2,"type":"同盟","intimacy":0.8}}}' | ./mcp-novel

# 7) 创建小说、分卷与章节
echo '{"jsonrpc":"2.0","id":17,"method":"tools/call","params":{"name":"novelHelper","arguments":{"action":"create","title":"王冠之路","description":"成长与选择"}}}' | ./mcp-novel
# 假设 Novel.ID = 1
echo '{"jsonrpc":"2.0","id":18,"method":"tools/call","params":{"name":"volumeHelper","arguments":{"action":"create","novelID":1,"title":"起始之章","index":1}}}' | ./mcp-novel
# 假设 Volume.ID = 1
echo '{"jsonrpc":"2.0","id":19,"method":"tools/call","params":{"name":"chapterHelper","arguments":{"action":"create","volumeID":1,"title":"骑士的旅途","index":1,"status":"草稿"}}}' | ./mcp-novel
# 假设 Chapter.ID = 1

# 8) 创建事件（绑定世界/地点/时间段与人物）
echo '{"jsonrpc":"2.0","id":20,"method":"tools/call","params":{"name":"eventHelper","arguments":{"action":"create","chapterID":1,"worldID":1,"locationID":1,"timeSegmentID":1,"description":"艾伦在都城宣誓","characters":[1,2],"items":[]}}}' | ./mcp-novel

# 9) 填写章节内容
echo '{"jsonrpc":"2.0","id":21,"method":"tools/call","params":{"name":"chapterHelper","arguments":{"action":"update","id":1,"content":"艾伦步入殿堂，光辉自穹顶洒落..."}}}' | ./mcp-novel

# 10) 设置文风参考
echo '{"jsonrpc":"2.0","id":22,"method":"tools/call","params":{"name":"styleHelper","arguments":{"action":"set","novelID":1,"content":"参考正文片段..."}}}' | ./mcp-novel

# 11) 生成纲要（章节/分卷/小说）
echo '{"jsonrpc":"2.0","id":23,"method":"tools/call","params":{"name":"outlineGeneratorHelper","arguments":{"action":"novel","id":1}}}' | ./mcp-novel

# 12) 冲突检测
echo '{"jsonrpc":"2.0","id":24,"method":"tools/call","params":{"name":"conflictDetectionHelper","arguments":{"action":"run"}}}' | ./mcp-novel

# 13) 导出文本（章节/分卷/小说）
echo '{"jsonrpc":"2.0","id":25,"method":"tools/call","params":{"name":"articleExportHelper","arguments":{"action":"novel","id":1}}}' | ./mcp-novel
```

## 数据持久化与导出

- 数据库文件：`novel.db`（SQLite）
- 初始化与迁移：应用启动时自动执行，逻辑见 `internal/mcp/server.go:332`
- 导出接口：
  - 小说：`novelHelper` `action=export` 返回整本文本
  - 分卷：`articleExportHelper` `action=volume`
  - 章节：`articleExportHelper` `action=chapter`

## 冲突检测

冲突检测入口为 `conflictDetectionHelper`，当前实现了以下 10 类检测（规则可扩展，实现位于 `internal/conflict/conflict.go:1`）：

- 时间冲突：时间段重叠、无效时间段
- 事件冲突：必需引用缺失（世界/地点）
- 人物冲突：基础字段有效性（如姓名缺失）
- 地点冲突：世界引用缺失
- 引用完整性：事件引用章节不存在
- 关系逻辑：自我关系非法等
- 状态一致性：章节状态枚举校验
- 物品能力冲突：物品或能力不存在时的使用/流转
- 线索冲突：线索阶段缺失
- 人物地点关系冲突：事件的地点/世界引用缺失

## 纲要生成

纲要生成由 `outlineGeneratorHelper` 提供：

- `action=chapter`：生成章节细纲（基于事件）
- `action=volume`：生成分卷总纲（汇总章节细纲）
- `action=novel`：生成小说总纲（汇总分卷纲要）

实现见 `internal/outline/outline.go:7`。

## 文笔风格参考

为小说设置参考正文用于后续创作风格对齐：

- 设置：`styleHelper` `action=set` `novelID` `content`
- 获取：`styleHelper` `action=get` `novelID`

数据模型见 `internal/models/models.go:StyleRef`，服务实现见 `internal/helpers/helpers.go:Services.SetStyleRef`、`internal/helpers/helpers.go:Services.GetStyleRef`。

## 错误与返回格式

- 成功：`result` 字段返回对象/文本
- 失败：`error` 字段包含 `code` 与 `message`
- JSON-RPC 2.0 规范：`internal/mcp/server.go:24`（请求）与 `internal/mcp/server.go:31`（响应）结构

## 最佳实践

- 先搭建世界观与时间轴，再铺设地点与人物关系，最后编排事件与章节
- 事件中维持引用完整（章节/世界/地点/时间段），人物与物品/能力的演进绑定事件记录
- 持续运行冲突检测，避免时间/引用/状态类问题
- 通过分卷与章节的纲要生成，迭代调整线索阶段与事件顺序

## 代码导航

- 入口：`cmd/mcp-novel/main.go:1`
- MCP 服务与工具路由：`internal/mcp/server.go:84`（循环）、`internal/mcp/server.go:118`（工具列表）、`internal/mcp/server.go:143`（调用路由）
- 模型：`internal/models/models.go:1`
- 服务层：`internal/helpers/helpers.go:16`
- 冲突检测：`internal/conflict/conflict.go:1`
- 纲要生成：`internal/outline/outline.go:1`

## 常见问题

- 时间格式错误：`timeSegmentHelper` 的 `start|end` 需为 RFC3339 格式，如 `2025-01-01T00:00:00Z`
- 关系未双向：使用 `characterRelationshipHelper` 会自动建立反向关系，无需手动再建一遍
- 导出为空：请确认章节内容已通过 `chapterHelper` `action=update` 写入

## 集成说明

该服务通过 stdio 提供 MCP 能力。外部 MCP 客户端可在初始化后调用 `tools/list` 获取能力清单，再通过 `tools/call` 调用各工具。若需要以进程形式嵌入，请保持标准输入输出为该协议的传输通道。