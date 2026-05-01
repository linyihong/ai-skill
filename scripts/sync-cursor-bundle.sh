#!/usr/bin/env bash
# 將本 repository 收斂到 ~/.cursor/bundles/ 下兩塊：
#   bundles/shared-rules  -> 共用規則
#   bundles/ai-skill/     -> 各 skill（僅含 skills/*，不含 shared-rules）
# 再讓 ~/.cursor/shared-rules 與 ~/.cursor/skills/<name> 指向上述 bundle（與其他 `.cursor` 內容分流）。
#
# 禁止在 repo 內出現「指回父目錄」的同名 symlink（例 shared-rules/shared-rules），否則 IDE 會無限巢狀。
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUNDLE_SKILLS="${HOME}/.cursor/bundles/ai-skill"
BUNDLE_RULES="${HOME}/.cursor/bundles/shared-rules"
CURSOR_SHARED="${HOME}/.cursor/shared-rules"

# 移除 repo 內會造成無限展開的 symlink（開頭與結尾各跑一次，避免同步過程或外部工具中途建立）。
_strip_repo_loop_symlinks() {
  local sr="${REPO_ROOT}/shared-rules"
  local sk="${REPO_ROOT}/skills"

  # shared-rules 底下任意深度的「名為 shared-rules 的 symlink」在此專案皆不合法；-P 不跟隨目錄 symlink 遞迴
  if [[ -d "${sr}" ]]; then
    while IFS= read -r -d '' bad; do
      echo "Note: removing forbidden symlink ${bad}"
      rm -f "${bad}"
    done < <(find -P "${sr}" -mindepth 1 -name 'shared-rules' -type l -print0 2>/dev/null || true)
  fi

  if [[ -d "${sk}" ]]; then
    while IFS= read -r -d '' dir; do
      local name="${dir%/}"
      name="$(basename "${name}")"
      [[ "${name}" == "_template" ]] && continue
      local nested="${dir%/}/${name}"
      if [[ -L "${nested}" ]]; then
        echo "Note: removing forbidden symlink ${nested}"
        rm -f "${nested}"
      fi
    done < <(find -P "${sk}" -mindepth 1 -maxdepth 1 -type d -print0 2>/dev/null || true)
  fi
}

# bundle 路徑絕不可落在 repo 內，否則 ln 可能把連結建進 shared-rules 形成套娃（字首比對，不依賴路徑已存在）
_repo_contains_path() {
  local t="${1%/}"
  case "${t}" in
    "${REPO_ROOT}" | "${REPO_ROOT}/"*) return 0 ;;
    *) return 1 ;;
  esac
}

if _repo_contains_path "${BUNDLE_RULES}" || _repo_contains_path "${BUNDLE_SKILLS}" || _repo_contains_path "${CURSOR_SHARED}"; then
  echo "ERROR: bundle paths must live outside REPO_ROOT (got REPO_ROOT=${REPO_ROOT}). Refusing to run." >&2
  exit 1
fi

_strip_repo_loop_symlinks

mkdir -p "${BUNDLE_SKILLS}" "${HOME}/.cursor/skills" "${HOME}/.cursor/bundles"

# 舊版曾把 shared-rules 放在 ai-skill 底下；改為 bundles/shared-rules 後移除殘留連結
rm -f "${BUNDLE_SKILLS}/shared-rules"

if [[ -e "${CURSOR_SHARED}" && ! -L "${CURSOR_SHARED}" ]]; then
  backup="${CURSOR_SHARED}.bak.${RANDOM}"
  echo "Note: ${CURSOR_SHARED} is not a symlink; moving aside -> ${backup}"
  mv "${CURSOR_SHARED}" "${backup}"
fi

# 僅寫入 ~/.cursor/bundles，永不寫 REPO_ROOT/shared-rules 內
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
done < <(find -P "${REPO_ROOT}/skills" -mindepth 1 -maxdepth 1 -type d -print0)

_strip_repo_loop_symlinks

echo "Synced bundles:"
echo "  ${BUNDLE_RULES} -> ${REPO_ROOT}/shared-rules"
echo "  ${CURSOR_SHARED} -> ${BUNDLE_RULES}"
echo "  ${BUNDLE_SKILLS}/<name>/ -> ${REPO_ROOT}/skills/<name>/"
echo "  ~/.cursor/skills/<name> -> ${BUNDLE_SKILLS}/<name>"
echo "Reload Cursor (Developer: Reload Window) if skills do not refresh."
