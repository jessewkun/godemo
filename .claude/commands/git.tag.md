---
description: Generate formatted tag message with changelog from last tag
---

## Smart Git Tag Message Generator

**使用 `generating-git-tags` skill 生成精炼的中文 tag 信息**

参数: $ARGUMENTS

### What This Command Does

1. **Runs the tag generation script** to collect commit data
2. **Script outputs**:
   - Raw commit data (JSON format with scope, description, original commit)
   - Suggested version number based on semantic versioning
   - Commit statistics
3. **AI (Claude) analyzes** the raw commit data and:
   - Groups commits by scope/module
   - Generates refined, concise Chinese descriptions (不超过30字)
   - Merges similar changes into single items
   - Creates a meaningful summary line
4. **Generates formatted changelog** with emoji categorization:
   - 🌟 新增功能 (feat) - New features
   - 🐛 问题修复 (fix) - Bug fixes
   - 🎨 UI优化 (style) - UI/UX improvements
   - 🔧 技术改进 (refactor, perf) - Technical improvements
   - 📚 文档更新 (docs) - Documentation updates
   - 🧪 测试相关 (test) - Test-related changes
   - 📦 其他更新 (chore, build) - Other updates
5. **Outputs ready-to-use git tag commands**

### Tag Message Format

The generated message follows this structure:

```
v{version} - {功能简述}

🌟 新增功能:
- {feature description}
...

🐛 问题修复:
- {fix description}
...

🎨 UI优化:
- {style improvement}
...

🔧 技术改进:
- {refactor/perf improvement}
...

📚 文档更新:
- {docs update}
...

📦 提交统计: {total} 个commit (feat: {count}, fix: {count}, ...)
```

### Command Usage Examples

```bash
# Generate tag message for the next version
/tag

# Generate tag message with specific version number
/tag v0.1.11

# Generate tag message from a specific tag
/tag --from v0.1.9

# Generate tag message and suggest version bump
/tag --suggest
```

### Version Number Suggestion

The command analyzes commits and suggests version bump based on:

- **Major (X.0.0)**: Breaking changes (BREAKING CHANGE in commit body)
- **Minor (0.X.0)**: New features (feat commits)
- **Patch (0.0.X)**: Bug fixes and minor improvements (fix, style, refactor)

### Implementation Guidelines

#### 1. Get Last Version Tag and Commits

```bash
# Get the latest version tag (cross-platform compatible)
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")

# If no version tag exists, use all commits
if [ -z "$LAST_TAG" ]; then
  COMMIT_RANGE="--all"
  echo "ℹ️ 没有找到版本tag,将分析所有commit"
else
  COMMIT_RANGE="$LAST_TAG..HEAD"
  echo "ℹ️ 从版本 $LAST_TAG 到 HEAD 分析commit"
fi

# Get commits since last version tag (use %B to get full commit body for BREAKING CHANGE detection)
git log $COMMIT_RANGE --pretty=format:"%B"
```

#### 2. Parse Conventional Commits

Extract commit type and description using regex patterns:

- `^feat(\(.*?\))?:\s*(.+)` → 🌟 新增功能
- `^fix(\(.*?\))?:\s*(.+)` → 🐛 问题修复
- `^style(\(.*?\))?:\s*(.+)` → 🎨 UI优化
- `^refactor(\(.*?\))?:\s*(.+)` → 🔧 技术改进
- `^perf(\(.*?\))?:\s*(.+)` → 🔧 技术改进
- `^docs(\(.*?\))?:\s*(.+)` → 📚 文档更新
- `^test(\(.*?\))?:\s*(.+)` → 🧪 测试相关
- `^chore(\(.*?\))?:\s*(.+)` → 📦 其他更新
- `^build(\(.*?\))?:\s*(.+)` → 📦 其他更新

#### 3. Filter Out Noise

- Skip merge commits: `^Merge (pull request|branch)`
- Skip PR merge messages: `^.*\(#\d+\)$` (only if it's the full message)
- Remove duplicate descriptions
- Trim commit descriptions to reasonable length (< 80 chars)

#### 4. Generate Categorized Changelog

Group commits by category and format with bullet points:

```
🌟 新增功能:
- 添加token刷新钩子
- 阅读材料管理功能

🐛 问题修复:
- 修复比赛SSE连接问题
- 修复登录页面用户检测
```

#### 5. Calculate Statistics

Count total commits and commits per category:

```
📦 提交统计: 38个commit (feat: 3, fix: 5, style: 4, refactor: 3, chore: 2)
```

### Best Practices

- **Review before creating tag**: Always review the generated message before creating the actual tag
- **Edit as needed**: The generated message is a starting point - feel free to edit it
- **Semantic versioning**: Follow semantic versioning principles when choosing version numbers
- **Meaningful summaries**: Add a concise summary after the version number
- **Keep it clean**: Remove redundant or noise commits from the changelog
- **Group related changes**: Combine similar commits into single bullet points when appropriate

### Important Notes

- This command only **generates the tag message** - it does NOT create the actual git tag
- You need to manually create the tag using the provided HEREDOC command (safer than inline strings)
- Always verify the generated message matches the actual changes
- For production releases, review the changelog carefully before tagging
- If commits don't follow conventional commit format, they may be categorized as "其他更新"

### Example Output

```
建议版本号: v0.1.11 (minor版本升级,因为有新功能)

=== Tag Message ===

v0.1.11 - SSE连接优化和Token刷新功能

🌟 新增功能:
- 添加token刷新钩子到HTTP客户端
- 阅读材料管理功能

🐛 问题修复:
- 修复比赛SSE连接采用自定义方式
- 修复登录页面用户检测

🎨 UI优化:
- 更新电子书阅读UI
- 更新全局样式

🔧 技术改进:
- 重构顶部栏样式
- 清理未使用代码并优化组件
- 简化文件上传逻辑和错误处理

📦 提交统计: 38个commit (feat: 3, fix: 5, style: 4, refactor: 3, chore: 23)

=== 创建tag命令 ===

复制以下命令创建tag:

```bash
git tag -a v0.1.11 -m "$(cat <<'EOF'
v0.1.11 - SSE连接优化和Token刷新功能

🌟 新增功能:
- 添加token刷新钩子到HTTP客户端
- 阅读材料管理功能

🐛 问题修复:
- 修复比赛SSE连接采用自定义方式
- 修复登录页面用户检测

🎨 UI优化:
- 更新电子书阅读UI
- 更新全局样式

🔧 技术改进:
- 重构顶部栏样式
- 清理未使用代码并优化组件
- 简化文件上传逻辑和错误处理

📦 提交统计: 38个commit (feat: 3, fix: 5, style: 4, refactor: 3, chore: 23)
EOF
)"
```

⚠️ **高风险操作确认**
操作类型：创建并推送 Git tag
影响范围：版本标记和远程仓库同步
风险评估：会影响版本管理和发布流程

⚠️ 艹！检测到危险操作！
你真要创建并推送这个tag吗？输入"确认"继续执行。

然后推送tag到远程:
```bash
git push origin v0.1.11
```

### Tips

- Use `--from` parameter to generate changelog between specific tags
- The generated message can be directly used with `git tag -a -m`
- For major releases, consider manually writing a more detailed summary
- Keep the first line (version + summary) under 72 characters
- Review PR merge commits - they might need better descriptions

---

### Implementation Instructions for Claude

**IMPORTANT: This command uses the `generating-git-tags` skill. Follow the skill's workflow.**

When this command is invoked:

**Step 1: Invoke the generating-git-tags skill**

The skill will:
1. Run the tag generation script: `node .claude/skills/generating-git-tags/scripts/generate-tag.cjs`
2. Parse the script output to get:
   - Raw commit data in JSON format
   - Suggested version number
   - Commit statistics

**Step 2: Analyze the raw commit data (AI's responsibility)**

Following the skill's guidelines, analyze each category of commits:

1. **Group by scope**: Group commits with the same scope together
2. **Generate refined Chinese descriptions**:
   - For single commits: Translate and refine to concise Chinese (e.g., "路由重命名优化")
   - For multiple commits in same scope: Create aggregated description (e.g., "作业模块功能完善(3项)")
3. **Merge similar changes**: Combine related commits into single items
4. **Use domain knowledge**: Apply understanding of the codebase to generate meaningful descriptions
5. **Keep descriptions concise**: Maximum 30 Chinese characters per item

**Step 3: Generate the final tag message**

Create a complete tag message with:
1. Version line: `{version} - {refined Chinese summary}`
2. Changelog sections with refined Chinese descriptions (NOT raw commit text)
3. Statistics from the script
4. Git commands

**Key principles (from generating-git-tags skill):**
- ✅ Generate refined, user-friendly Chinese descriptions
- ❌ Don't use raw commit text directly
- ✅ Group and merge similar changes
- ✅ Apply domain knowledge about the codebase

---

### Example AI Analysis Process

**Input (from script):**
```json
{
  "feat": [
    { "scope": "homework", "desc": "english writing homework analysis support view ai polished result", "original": "..." },
    { "scope": "", "desc": "add PrecisionLearning", "original": "..." }
  ],
  "fix": [
    { "scope": "homework", "desc": "homework publish page enable time select", "original": "..." },
    { "scope": "", "desc": "homework navigation and student detail style", "original": "..." }
  ],
  "style": [
    { "scope": "homework", "desc": "update analysis component ui", "original": "..." }
  ],
  "refactor": [
    { "scope": "router", "desc": "rename reading-homework routes to english-exercise", "original": "..." }
  ]
}
```

**AI Analysis Output:**
```
🌟 新增功能:
- 作业分析功能支持查看AI润色结果
- 精准学习模块

🐛 问题修复:
- 作业发布页面时间选择功能优化
- 作业导航和学生详情页面样式调整

🎨 UI优化:
- 作业分析组件界面优化

🔧 技术改进:
- 路由命名规范化(阅读作业→英语练习)
```

**Summary generation:**
- Analyze: Most commits in `homework` scope, with new features
- Generate: "作业模块功能完善和优化"
