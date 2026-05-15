#!/bin/bash
# Install git hooks from .githooks/ to .git/hooks/
#
# 使用方式：
#   ./scripts/install-hooks.sh
#
# 這個腳本會：
# 1. 將 .githooks/ 目錄下的所有 hook 複製到 .git/hooks/
# 2. 設定執行權限
# 3. 設定 core.hooksPath（選擇性，視需求）

set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || (cd "$(dirname "$0")/.." && pwd))
cd "$REPO_ROOT"

GITHOOKS_DIR="$REPO_ROOT/.githooks"
GIT_HOOKS_DIR="$REPO_ROOT/.git/hooks"

if [ ! -d "$GITHOOKS_DIR" ]; then
  echo "Error: .githooks directory not found at $GITHOOKS_DIR"
  exit 1
fi

if [ ! -d "$GIT_HOOKS_DIR" ]; then
  echo "Error: .git/hooks directory not found. Are you in a git repository?"
  exit 1
fi

echo "Installing git hooks from .githooks/ to .git/hooks/..."
echo ""

INSTALLED=0
for hook in "$GITHOOKS_DIR"/*; do
  if [ -f "$hook" ]; then
    HOOK_NAME=$(basename "$hook")
    TARGET="$GIT_HOOKS_DIR/$HOOK_NAME"
    
    cp "$hook" "$TARGET"
    chmod +x "$TARGET"
    echo "  ✓ Installed: $HOOK_NAME"
    INSTALLED=$((INSTALLED + 1))
  fi
done

echo ""
echo "✅ $INSTALLED hook(s) installed successfully."
echo ""
echo "Hooks will run automatically on applicable git actions."
echo ""
echo "To skip checks (emergency only):"
echo "  SKIP_COMPILER_CHECK=1 git commit"
echo "  SKIP_GOVERNANCE_CHECK=1 git commit"
echo "  SKIP_COMPILER_CHECK=1 SKIP_GOVERNANCE_CHECK=1 git commit"
