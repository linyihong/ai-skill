#!/usr/bin/env bash
# 將本 repository 收斂到 ~/.cursor/bundles/ 下兩塊：
#   bundles/shared-rules  -> 共用規則
#   bundles/ai-skill/     -> 各 skill（僅含 skills/*，不含 shared-rules）
# 再讓 ~/.cursor/shared-rules 與 ~/.cursor/skills/<name> 指向上述 bundle（與其他 ~/.cursor 內容分流）。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUNDLE_SKILLS="${HOME}/.cursor/bundles/ai-skill"
BUNDLE_RULES="${HOME}/.cursor/bundles/shared-rules"
CURSOR_SHARED="${HOME}/.cursor/shared-rules"

mkdir -p "${BUNDLE_SKILLS}" "${HOME}/.cursor/skills" "${HOME}/.cursor/bundles"

# 舊版曾把 shared-rules 放在 ai-skill 底下；改為 bundles/shared-rules 後移除殘留連結
rm -f "${BUNDLE_SKILLS}/shared-rules"

if [[ -e "${CURSOR_SHARED}" && ! -L "${CURSOR_SHARED}" ]]; then
  backup="${CURSOR_SHARED}.bak.${RANDOM}"
  echo "Note: ${CURSOR_SHARED} is not a symlink; moving aside -> ${backup}"
  mv "${CURSOR_SHARED}" "${backup}"
fi

ln -sf "${REPO_ROOT}/shared-rules" "${BUNDLE_RULES}"

rm -f "${CURSOR_SHARED}"
ln -sf "${BUNDLE_RULES}" "${CURSOR_SHARED}"

while IFS= read -r -d '' dir; do
  name="$(basename "${dir}")"
  case "${name}" in _template) continue ;; esac
  if [[ ! -f "${dir}/SKILL.md" ]]; then
    continue
  fi
  ln -sf "${dir}" "${BUNDLE_SKILLS}/${name}"
  rm -f "${HOME}/.cursor/skills/${name}"
  ln -sf "${BUNDLE_SKILLS}/${name}" "${HOME}/.cursor/skills/${name}"
done < <(find "${REPO_ROOT}/skills" -mindepth 1 -maxdepth 1 -type d -print0)

echo "Synced bundles:"
echo "  ${BUNDLE_RULES} -> ${REPO_ROOT}/shared-rules"
echo "  ${CURSOR_SHARED} -> ${BUNDLE_RULES}"
echo "  ${BUNDLE_SKILLS}/<name>/ -> ${REPO_ROOT}/skills/<name>/"
echo "  ~/.cursor/skills/<name> -> ${BUNDLE_SKILLS}/<name>"
echo "Reload Cursor (Developer: Reload Window) if skills do not refresh."
