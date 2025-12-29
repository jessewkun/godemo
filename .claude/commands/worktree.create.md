---
description: Create a new git worktree with specified name and base branch
---

## Git Worktree 创建工具

创建新的git worktree：$ARGUMENTS

### 当前仓库状态

- 当前分支：!`git branch --show-current`
- 当前worktree列表：!`git worktree list`

### 功能描述

这个命令用于创建新的git worktree，让你可以在多个分支上并行开发而不用频繁切换分支。

### 参数说明

- `$1` - worktree名称（必传）- 将创建为 `../feature-{worktree名称}` 目录
- `$2` - base分支（选传）- 默认为当前分支

### 命令用法

```bash
# 从当前分支创建新worktree
/worktree reading-homework

# 从指定分支创建新worktree
/worktree reading-homework feature/codebase-251002

# 从main分支创建新worktree
/worktree new-feature main
```

### 创建逻辑

#### 1. 参数解析
```bash
WORKTREE_NAME="$1"
BASE_BRANCH="$2"

# 如果没有指定base分支，使用当前分支
if [ -z "$BASE_BRANCH" ]; then
  BASE_BRANCH=$(git branch --show-current)
fi
```

#### 2. 目录命名规则
```bash
# worktree目录命名规则
WORKTREE_DIR="../feature-${WORKTREE_NAME}"
BRANCH_NAME="feature/${WORKTREE_NAME}"
```

#### 3. 创建worktree
```bash
# 创建新分支和worktree
git worktree add -b ${BRANCH_NAME} ${WORKTREE_DIR} ${BASE_BRANCH}
```

#### 4. 验证创建结果
```bash
# 显示所有worktree
echo "✅ Worktree创建成功！"
echo ""
git worktree list
```

### 安全检查

在创建worktree之前，会进行以下检查：

1. **检查worktree名称是否为空**
   ```bash
   if [ -z "$WORKTREE_NAME" ]; then
     echo "❌ 错误：必须提供worktree名称"
     echo "用法: /worktree <worktree名称> [base分支]"
     return 1
   fi
   ```

2. **检查目标目录是否已存在**
   ```bash
   if [ -d "$WORKTREE_DIR" ]; then
     echo "❌ 错误：目录 $WORKTREE_DIR 已存在"
     return 1
   fi
   ```

3. **检查base分支是否存在**
   ```bash
   if ! git rev-parse --verify "$BASE_BRANCH" >/dev/null 2>&1; then
     echo "❌ 错误：分支 $BASE_BRANCH 不存在"
     return 1
   fi
   ```

### 创建后操作

创建成功后，你可以：

```bash
# 进入新worktree目录
cd ../feature-${WORKTREE_NAME}

# 查看当前分支
git branch

# 开始在新分支上开发
# ... 你的代码修改 ...

# 完成后删除worktree（可选）
cd ../$(basename $(git rev-parse --show-toplevel))
git worktree remove ../feature-${WORKTREE_NAME}
```

### Worktree管理命令

```bash
# 列出所有worktree
git worktree list

# 删除worktree
git worktree remove ../feature-${WORKTREE_NAME}

# 清理已删除的worktree分支
git branch -d feature/${WORKTREE_NAME}

# 切换到特定worktree目录
cd ../feature-${WORKTREE_NAME}
```

### 实际例子

假设你要从 `feature/codebase-251002` 分支创建一个阅读作业功能的worktree：

```bash
# 执行命令
/worktree reading-homework feature/codebase-251002

# 创建结果：
# ✅ Worktree创建成功！
```

### 注意事项

1. **目录结构**：worktree会创建在上级目录下，命名为 `feature-{worktree名称}`
2. **分支命名**：新分支会命名为 `feature/{worktree名称}`
3. **独立工作空间**：每个worktree都有独立的工作目录，可以同时在不同分支上开发
4. **共享仓库**：所有worktree共享同一个git仓库的历史记录
5. **清理建议**：完成功能开发后，记得清理不需要的worktree

### 错误处理

常见错误及解决方案：

- **目录已存在**：换个worktree名称或手动删除现有目录
- **分支不存在**：检查base分支名称是否正确，使用 `git branch -a` 查看所有分支
- **权限问题**：确保对上级目录有写权限

---

### Implementation Instructions for Claude

当执行此命令时，应该：

#### Step 1: 解析参数
```typescript
const args = "$ARGUMENTS".trim().split(" ")
const worktreeName = args[0]
const baseBranch = args[1] || $(git branch --show-current).trim()
```

#### Step 2: 验证参数
```bash
# 检查worktree名称
if [ -z "$worktreeName" ]; then
  echo "❌ 错误：必须提供worktree名称"
  echo "用法: /worktree <worktree名称> [base分支]"
  exit 1
fi

# 检查base分支是否存在
if ! git rev-parse --verify "$baseBranch" >/dev/null 2>&1; then
  echo "❌ 错误：分支 $baseBranch 不存在"
  exit 1
fi

# 检查目录是否已存在
WORKTREE_DIR="../feature-${worktreeName}"
if [ -d "$WORKTREE_DIR" ]; then
  echo "❌ 错误：目录 $WORKTREE_DIR 已存在"
  exit 1
fi
```

#### Step 3: 创建worktree
```bash
BRANCH_NAME="feature/${worktreeName}"

echo "🔨 正在创建worktree..."
echo "   Worktree名称: $worktreeName"
echo "   目标分支: $BRANCH_NAME"
echo "   基础分支: $baseBranch"
echo "   工作目录: $WORKTREE_DIR"
echo ""

# 创建worktree和分支
git worktree add -b "$BRANCH_NAME" "$WORKTREE_DIR" "$baseBranch"

if [ $? -eq 0 ]; then
  echo ""
  echo "✅ Worktree创建成功！"
  echo ""
  echo "📋 当前所有worktree："
  git worktree list
  echo ""
  echo "💡 使用提示："
  echo "   进入新worktree: cd $WORKTREE_DIR"
  echo "   查看当前分支: git branch"
  echo "   删除worktree: git worktree remove $WORKTREE_DIR"
else
  echo "❌ Worktree创建失败！"
  exit 1
fi
```

#### Step 4: 显示使用提示
```bash
echo ""
echo "🎯 下一步操作："
echo "1. cd $WORKTREE_DIR"
echo "2. 开始你的功能开发..."
echo "3. 完成后可以删除worktree: git worktree remove $WORKTREE_DIR"
```
