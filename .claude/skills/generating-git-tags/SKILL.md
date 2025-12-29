---
name: generating-git-tags
description: Collects commit data since the last tag and outputs structured information for AI to generate refined Chinese tag messages. The script gathers raw commit data (scope, description) while the AI model analyzes and summarizes it into concise Chinese descriptions. Use when creating release tags or generating changelogs.
---

# Git Tag Generator

Collects commit data and provides structured information for AI analysis. The script gathers raw commit data (with scope and description), while the AI model (Claude) analyzes the data and generates refined, concise Chinese tag messages instead of using raw commit text.

## Quick start

1. **Script**: Get the last version tag and commit range
2. **Script**: Collect and parse commits using conventional commit format (extract type, scope, description)
3. **Script**: Output raw commit data as JSON for AI analysis
4. **AI (Claude)**: Analyze commit data and group by scope/module
5. **AI (Claude)**: Generate refined Chinese summaries instead of using raw commit messages
6. **Script**: Suggest version number based on semantic versioning
7. **AI (Claude)**: Generate final tag message with refined descriptions and git commands

## Workflow

### Step 1: Determine commit range

```bash
# Get last version tag
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")

# Determine range
if [ -z "$LAST_TAG" ]; then
  RANGE="--all"
  echo "ℹ️ No version tag found, analyzing all commits"
else
  RANGE="$LAST_TAG..HEAD"
  echo "ℹ️ Analyzing commits from $LAST_TAG to HEAD"
fi
```

### Step 2: Collect commits

```bash
# Get commit subjects (for categorization)
git log $RANGE --pretty=format:"%s" --no-merges

# Get full commit bodies (for BREAKING CHANGE detection)
git log $RANGE --pretty=format:"%B" --no-merges
```

### Step 3: Parse and collect commits (Script)

Parse each commit using conventional commit format and collect raw data:

- Pattern: `^(feat|fix|style|refactor|perf|docs|test|chore|build)(\(([^)]+)\))?:\s*(.+)`
- Extract: type, scope (optional), description
- Skip merge commits: `^Merge (pull request|branch)`
- Remove PR numbers: `\s*\(#\d+\)$`
- Store as structured data: `{ scope: string, desc: string, original: string }`

**Category mapping**:
- `feat` → 🌟 新增功能
- `fix` → 🐛 问题修复
- `style` → 🎨 UI优化
- `refactor`, `perf` → 🔧 技术改进
- `docs` → 📚 文档更新
- `test` → 🧪 测试相关
- `chore`, `build` → 📦 其他更新

### Step 3.5: Output raw commit data (Script)

Output structured commit data as JSON for AI analysis:

```javascript
{
  "feat": [
    { "scope": "homework", "desc": "add analysis component", "original": "..." },
    { "scope": "homework", "desc": "update ui", "original": "..." }
  ],
  "fix": [
    { "scope": "router", "desc": "rename reading-homework routes", "original": "..." }
  ],
  // ... other categories
}
```

Along with AI instructions:
- 分析每个类别的 commits,按 scope 分组
- 为每个分组生成简洁的中文描述(不超过30字)
- 相似的功能应该合并成一条描述
- 使用用户易懂的语言,避免技术术语
- 保持专业性和准确性

### Step 3.6: Analyze and summarize commits (AI/Claude)

**AI analyzes the raw commit data and generates refined Chinese descriptions:**

1. **Group by scope**: Group commits with the same scope together
2. **Identify patterns**: Look for similar functionality or related changes
3. **Generate concise descriptions**:
   - Single commit: Translate and refine to clear Chinese (e.g., "路由重命名优化")
   - Multiple commits in same scope: Summarize as grouped item (e.g., "作业模块功能完善 (3项)")
4. **Apply domain knowledge**: Use context about the codebase to generate meaningful descriptions
5. **Avoid raw commit text**: Don't just translate; summarize and refine

**Example transformation:**
- Input: `feat(homework): add analysis component`, `feat(homework): update ui`
- AI Output: `作业分析功能和界面优化`

- Input: `fix(router): rename reading-homework routes to english-exercise`
- AI Output: `路由命名规范化(阅读作业→英语练习)`

### Step 4: Suggest version

```typescript
// Determine bump type
let bumpType = "patch"
if (counts.feat > 0) bumpType = "minor"
if (hasBreakingChange) bumpType = "major"

// For 0.x versions, increment patch for minor bumps
const [major, minor, patch] = currentVersion.replace(/^v/, "").split(".").map(Number)
const newVersion = bumpType === "minor" && major === 0
  ? `v${major}.${minor}.${patch + 1}`
  : bumpType === "minor"
    ? `v${major}.${minor + 1}.0`
    : `v${major}.${minor}.${patch + 1}`
```

### Step 5: Generate summary (AI/Claude)

AI generates a concise Chinese summary (≤30 chars) by analyzing commit patterns:

1. **Analyze dominant scope**: Identify the most frequently changed module/area
2. **Determine main theme**:
   - Multiple features in same module → "{模块名}功能完善"
   - Major refactoring → "代码重构和优化"
   - Bug fixes → "问题修复和稳定性改进"
3. **Use domain knowledge**: Apply understanding of the codebase to create meaningful summary
4. **Keep concise**: Limit to 30 Chinese characters maximum

**Examples:**
- Commits mainly in `homework` scope with new features → "作业模块功能完善"
- Mixed fixes and UI updates → "问题修复和界面优化"
- Router refactoring → "路由架构优化"

### Step 6: Format changelog (AI/Claude)

AI generates refined changelog sections in this order:
1. 🌟 新增功能 (if any feat commits)
2. 🐛 问题修复 (if any fix commits)
3. 🎨 UI优化 (if any style commits)
4. 🔧 技术改进 (refactor + perf)
5. 📚 文档更新 (if any docs commits)
6. 🧪 测试相关 (if any test commits)
7. 📦 其他更新 (chore + build, only if ≤5 items)

**AI generates each section by:**
1. Grouping commits by scope within each category
2. Generating concise Chinese descriptions for each group
3. Merging similar changes into single items
4. Avoiding raw commit text - using refined, user-friendly language

**Each section format:**
```
🌟 新增功能:
- {AI生成的精炼中文描述 1}
- {AI生成的精炼中文描述 2}

🐛 问题修复:
- {AI生成的精炼中文描述}
```

**Not this (raw commit text):**
```
🌟 新增功能:
- add analysis component
- update ui
```

**But this (refined Chinese):**
```
🌟 新增功能:
- 作业分析功能和界面优化
```

### Step 7: Add statistics (Script)

Script generates commit statistics:

```
📦 提交统计: {total}个commit (feat: {count}, fix: {count}, ...)
```

### Step 8: Final output format (AI/Claude)

AI generates the final tag message by:
1. Taking the version number from the script
2. Creating a refined Chinese summary
3. Generating refined changelog sections from raw commit data
4. Including the statistics from the script
5. Providing ready-to-use git commands

```
建议版本号: {version} ({bumpType}版本升级,因为{reason})

=== Tag Message ===

{version} - {summary}

{changelog sections}

{statistics}

=== 创建tag命令 ===

复制以下命令创建tag:

```bash
git tag -a {version} -m "$(cat <<'EOF'
{full message}
EOF
)"
```

⚠️ **高风险操作确认**
操作类型：创建并推送 Git tag
影响范围：版本标记和远程仓库同步
风险评估：会影响版本管理和发布流程

⚠️ 检测到危险操作！
你真要创建并推送这个tag吗？输入"确认"继续执行。

然后推送tag到远程:
```bash
git push origin {version}
```
```

## Implementation

**Recommended workflow:**

1. **Run the script** to collect raw commit data:
   ```bash
   node .claude/skills/generating-git-tags/scripts/generate-tag.cjs [fromTag] [targetVersion]
   ```

2. **Script outputs:**
   - Raw commit data as JSON (to stderr)
   - AI instructions (to stderr)
   - Suggested version number
   - Basic statistics

3. **AI (Claude) analyzes** the raw commit data and:
   - Groups commits by scope/module
   - Generates refined Chinese descriptions
   - Creates concise summary line
   - Formats final tag message with git commands

**Key principle**: Script collects data, AI generates refined content. No hardcoded translations in the script.

## Error handling

- **No commits found**: "ℹ️ 从上一个tag到现在没有新的commit"
- **No tags exist**: Analyze all commits and suggest v0.1.0
- **Invalid commit format**: Categorize as "其他更新"

## Best practices

- Review generated message before creating tag
- Edit summary line to be meaningful (limit 30 chars)
- Remove redundant commits from changelog
- Group related changes when appropriate
- For 0.x versions, use patch increment for minor bumps (common practice)

## Notes

- **Script role**: Collects and structures commit data only
- **AI role**: Analyzes data and generates refined Chinese descriptions
- **No hardcoded translations**: All Chinese content is generated by AI based on context
- Only generates the tag message; does NOT create the actual git tag
- Always verify the AI-generated message matches actual changes
- For production releases, review changelog carefully
- Commits not following conventional format are categorized as "其他更新"
