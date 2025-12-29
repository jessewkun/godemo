#!/usr/bin/env node

/**
 * Git Tag Message Generator
 *
 * Analyzes commits since last version tag and generates formatted changelog
 * with semantic version suggestion.
 *
 * Usage:
 *   node generate-tag.js [fromTag] [targetVersion]
 *
 * Examples:
 *   node generate-tag.js                    # From last tag to HEAD
 *   node generate-tag.js v0.1.9            # From v0.1.9 to HEAD
 *   node generate-tag.js v0.1.9 v0.1.10   # From v0.1.9, use v0.1.10
 */

const { execSync } = require('node:child_process')
const process = require('node:process')

// Parse arguments
const fromTag = process.argv[2] || ''
const targetVersion = process.argv[3] || ''

// Get last version tag if not provided
let lastTag = fromTag
if (!lastTag) {
  try {
    lastTag = execSync('git describe --tags --abbrev=0 2>/dev/null', { encoding: 'utf-8' }).trim()
  }
  catch {
    lastTag = ''
  }
}

// Determine commit range
const range = lastTag ? `${lastTag}..HEAD` : '--all'
if (!lastTag) {
  console.error('ℹ️ No version tag found, analyzing all commits')
}
else {
  console.error(`ℹ️ Analyzing commits from ${lastTag} to HEAD`)
}

// Get commit subjects
let commits = []
try {
  const output = execSync(`git log ${range} --pretty=format:"%s" --no-merges`, { encoding: 'utf-8' })
  commits = output.trim().split('\n').filter(c => c.trim())
}
catch {
  console.error('ℹ️ No commits found in range')
  process.exit(0)
}

// Categories - store raw commit info for AI analysis
const rawCommits = {
  feat: [],
  fix: [],
  style: [],
  refactor: [],
  perf: [],
  docs: [],
  test: [],
  chore: [],
  build: [],
}

const counts = {
  feat: 0,
  fix: 0,
  style: 0,
  refactor: 0,
  perf: 0,
  docs: 0,
  test: 0,
  chore: 0,
  build: 0,
}

// Parse commits - collect raw commit data (no translation, just collection)
commits.forEach((commit) => {
  // Skip merge commits
  if (/^Merge (?:pull request|branch)/i.test(commit)) {
    return
  }

  // Match conventional commit format
  const match = commit.match(/^(feat|fix|style|refactor|perf|docs|test|chore|build)(?:\(([^)]+)\))?:\s*(.+)/i)
  if (match) {
    const type = match[1].toLowerCase()
    const scope = match[2] || ''
    let desc = match[3].trim()

    // Remove PR numbers
    desc = desc.replace(/\s*\(#\d+\)$/, '').replace(/\s*#\d+$/, '')

    if (rawCommits[type]) {
      rawCommits[type].push({ scope, desc, original: commit })
      counts[type]++
    }
  }
  else if (/^Test\//i.test(commit)) {
    rawCommits.test.push({ scope: 'test', desc: commit, original: commit })
    counts.test++
  }
})

// Output raw commits data for AI to analyze
// This will be used by the AI model to generate refined summaries
const categories = {
  feat: rawCommits.feat,
  fix: rawCommits.fix,
  style: rawCommits.style,
  refactor: [...rawCommits.refactor, ...rawCommits.perf],
  perf: [],
  docs: rawCommits.docs,
  test: rawCommits.test,
  chore: [...rawCommits.chore, ...rawCommits.build],
  build: [],
}

// Output raw commit data for AI to analyze
// Format: JSON structure that AI can easily parse and summarize
const commitData = {
  feat: categories.feat,
  fix: categories.fix,
  style: categories.style,
  refactor: categories.refactor,
  docs: categories.docs,
  test: categories.test,
  chore: categories.chore,
}

// Output structured commit data for AI analysis
console.error('\n=== Raw Commit Data (for AI analysis) ===')
console.error(JSON.stringify(commitData, null, 2))
console.error('\n=== AI Instructions ===')
console.error('请根据上述 commit 数据生成精炼的中文 tag 信息:')
console.error('1. 分析每个类别的 commits,按 scope 分组')
console.error('2. 为每个分组生成简洁的中文描述(不超过30字)')
console.error('3. 相似的功能应该合并成一条描述')
console.error('4. 使用用户易懂的语言,避免技术术语')
console.error('5. 保持专业性和准确性')
console.error('')

// Placeholder sections - AI will generate the actual content
const sections = []

// Stats
const total = Object.values(counts).reduce((a, b) => a + b, 0)
const statsItems = Object.entries(counts)
  .filter(([, count]) => count > 0)
  .map(([type, count]) => `${type}: ${count}`)
const stats = `📦 提交统计: ${total}个commit (${statsItems.join(', ')})`

// Version suggestion
let version = targetVersion
let bumpType = 'patch'
let reason = '修复和优化'

if (!version) {
  const currentVersion = lastTag || 'v0.0.0'
  const versionMatch = currentVersion.match(/^v?(\d+)\.(\d+)\.(\d+)/)

  if (versionMatch) {
    const major = Number.parseInt(versionMatch[1], 10)
    const minor = Number.parseInt(versionMatch[2], 10)
    const patch = Number.parseInt(versionMatch[3], 10)

    // Check for breaking changes (would need commit bodies)
    // For now, assume no breaking changes

    if (counts.feat > 0) {
      bumpType = 'minor'
      reason = '有新功能'
      // For 0.x versions, increment patch for minor bumps
      if (major === 0) {
        version = `v${major}.${minor}.${patch + 1}`
      }
      else {
        version = `v${major}.${minor + 1}.0`
      }
    }
    else {
      version = `v${major}.${minor}.${patch + 1}`
    }
  }
  else {
    version = 'v0.1.0'
  }
}

// Summary - placeholder, AI will generate the actual summary
const summary = '[AI将根据commit数据生成精炼的中文描述]'

// Assemble message
const firstLine = `${version} - ${summary}`
const message = [firstLine, '', ...sections, stats].join('\n')

// Output
console.log(`建议版本号: ${version} (${bumpType}版本升级,因为${reason})\n\n=== Tag Message ===\n\n${message}\n\n=== 创建tag命令 ===\n\n复制以下命令创建tag:\n\n\`\`\`bash\ngit tag -a ${version} -m "$(cat <<'EOF'\n${message}\nEOF\n)"\n\`\`\`\n\n⚠️ **高风险操作确认**\n操作类型：创建并推送 Git tag\n影响范围：版本标记和远程仓库同步\n风险评估：会影响版本管理和发布流程\n\n⚠️ 检测到危险操作！\n你真要创建并推送这个tag吗？输入"确认"继续执行。\n\n然后推送tag到远程:\n\`\`\`bash\ngit push origin ${version}\n\`\`\``)
