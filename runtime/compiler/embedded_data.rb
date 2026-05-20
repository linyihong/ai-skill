# frozen_string_literal: true

# Embedded Runtime Configuration Data
# This file is auto-generated from runtime/**/*.yaml
# DO NOT EDIT MANUALLY — edit the YAML files and regenerate
#
# Generated: 2026-05-17 03:47:55 UTC
#
module EmbeddedRuntimeData

  # Source: runtime/budget/token-budget.yaml
  BUDGET_TOKEN_BUDGET =   {
    "token_budget_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/budget",
    "default_budget": {
      "max_tokens": 120000,
      "warning_threshold": 0.7,
      "hard_stop_threshold": 0.9
    },
    "per_model": [
      {
        "model": "claude-sonnet-4",
        "max_tokens": 200000,
        "warning_threshold": 0.7,
        "hard_stop_threshold": 0.9
      },
      {
        "model": "claude-opus-4",
        "max_tokens": 200000,
        "warning_threshold": 0.7,
        "hard_stop_threshold": 0.9
      },
      {
        "model": "gpt-4o",
        "max_tokens": 128000,
        "warning_threshold": 0.7,
        "hard_stop_threshold": 0.9
      },
      {
        "model": "gemini-2.5-pro",
        "max_tokens": 1000000,
        "warning_threshold": 0.8,
        "hard_stop_threshold": 0.95
      }
    ],
    "layer_budget": {
      "bootstrap": {
        "max_tokens": 2000,
        "priority": "P0"
      },
      "skill_index": {
        "max_tokens": 500,
        "priority": "P0"
      },
      "activation_rules": {
        "max_tokens": 1000,
        "priority": "P1"
      },
      "summaries": {
        "max_tokens": 3000,
        "priority": "P1"
      },
      "full_source": {
        "max_tokens": 30000,
        "priority": "P2"
      },
      "tool_output": {
        "max_tokens": 50000,
        "priority": "P2"
      },
      "conversation": {
        "max_tokens": 30000,
        "priority": "P3"
      }
    },
    "on_warning": [
      {
        "action": "log",
        "message": "⚠️ Token usage at {percentage}% ({used_tokens}/{max_tokens})"
      },
      {
        "action": "suggest_prune",
        "message": "建議 prune 過期 context：TTL-expired tasks、deferred rules、old feedback lessons"
      },
      {
        "action": "suggest_new_session",
        "message": "建議建立新 session 以重置 context"
      }
    ],
    "on_hard_stop": [
      {
        "action": "halt_agent",
        "message": "⛔ Token budget exceeded ({used_tokens}/{max_tokens}). Agent halted."
      },
      {
        "action": "force_prune",
        "strategy": "aggressive",
        "preserve": [
          "core-bootstrap",
          "agent-goals",
          "conversation-goal-ledger"
        ]
      },
      {
        "action": "suggest_new_session",
        "message": "請建立新 session 繼續工作。當前 context 已存檔至 memory/working/。"
      }
    ]
  }

  # Source: runtime/compiler/compiler-rules.yaml
  COMPILER_COMPILER_RULES =   {
    "version": 1.1,
    "status": "active",
    "owner_layer": "runtime/compiler",
    "last_updated": "2026-05-17",
    "compilation_scope": {
      "description": "只有以下類別的內容會被編譯為 SQLite（runtime.db）。 所有需要 human judgment、heuristics、troubleshooting 的內容永遠留在 prose。\n",
      "deterministic_only": [
        "current_phase",
        "allowed_actions",
        "forbidden_actions",
        "blocking_gates (id, severity, check.type, check.verification)",
        "required_artifacts (name, path, format)",
        "open_obligations (id, phase, verification_criteria)",
        "transaction_state (state, allowed_actions, transitions)",
        "phase_transition_rules (from, to, blocked, reason)",
        "language_policy (core_rules, anti_drift, tool_overrides, validation)",
        "output_rules (format_rules, content_boundary, sanitization, validation)",
        "governance_gates (id, phase, severity, check.type, check.verification)",
        "runtime_budget (model_name, per_model config, default_budget, per_layer)",
        "context_ttl_policy (ttl_type, ttl_rules)",
        "circuit_breaker (guard_name, guard config)",
        "context_pollution (signal_name, pollution signals)",
        "context_health_score (dimension, health dimensions)",
        "intelligence_routing (rule_id, routing rules)",
        "obligation_ledger (obligation_id, obligations)",
        "phase_machine (phase_id, phases)",
        "pipeline_context_flow (level, context flow levels)",
        "guard_chain (stage, guard stages)",
        "relevance_engine (component, scoring components)",
        "session_lifecycle (stage, session stages)",
        "prompt_artifact_templates (template_name, templates)",
        "prompt_composition_rules (rule_id, composition rules)",
        "recovery_strategies (strategy_id, recovery strategies)",
        "state_repair (procedure_id, repair procedures)",
        "obligation_rebuild (procedure_id, rebuild procedures)",
        "phase_reconciliation (procedure_id, reconciliation procedures)",
        "execution_queue (queue_name, queue config)",
        "priority_scheduler (priority_level, scheduler config)",
        "activation_rules_mirror (rule_id, activation rules mirror)",
        "transaction_templates_ext (template_name, extended templates)",
        "distributed_locks (lock_name, distributed locks)",
        "multi_agent_coordination (rule_id, coordination rules)",
        "async_job_lifecycle (state, job lifecycle states)",
        "capability_checkpoints (checkpoint_id, capability checkpoints)"
      ],
      "never_compile": [
        "heuristics",
        "debugging judgment",
        "troubleshooting strategy",
        "architectural tradeoffs",
        "domain intelligence",
        "anti-pattern reasoning",
        "failure analysis",
        "decision rationale",
        "contextual examples"
      ]
    },
    "source_target_mapping": [
      {
        "source": "workflow/*/execution-flow.md",
        "target": "sqlite:runtime.db → generated_surfaces (type='workflow_phases')",
        "extract": [
          "phase_definitions",
          "allowed_actions",
          "forbidden_actions",
          "blocking_gates"
        ],
        "compile_rule": "從 execution-flow 的章節標題提取 phase 定義，從 gate 描述提取 blocking gates，寫入 runtime.db generated_surfaces 表"
      },
      {
        "source": "workflow/*/artifact-gates.md",
        "target": "sqlite:runtime.db → generated_surfaces (type='workflow_artifacts')",
        "extract": [
          "required_artifacts",
          "artifact_verification_criteria"
        ],
        "compile_rule": "從 artifact gates 的檢查清單提取 required artifacts 與 verification criteria，寫入 runtime.db generated_surfaces 表"
      },
      {
        "source": "enforcement/dependency-reading.md",
        "target": "sqlite:runtime.db → transaction_states",
        "extract": [
          "transaction_states",
          "transaction_rules"
        ],
        "compile_rule": "從 writeback transaction 章節提取 state machine 定義，寫入 runtime.db transaction_states 表"
      },
      {
        "source": "enforcement/goal-action-validation.md",
        "target": "sqlite:runtime.db → gates (type='validation_gate')",
        "extract": [
          "validation_gates",
          "validation_criteria"
        ],
        "compile_rule": "從 validation gate 描述提取 gate 定義，寫入 runtime.db gates 表"
      },
      {
        "source": "enforcement/failure-learning-system.md",
        "target": "sqlite:runtime.db → generated_surfaces (type='failure_recovery')",
        "extract": [
          "failure_patterns",
          "recovery_strategies"
        ],
        "compile_rule": "從 failure taxonomy 與 recovery 描述提取 pattern 與 strategy，寫入 runtime.db generated_surfaces 表"
      },
      {
        "source": "enforcement/neutral-language.md",
        "target": "sqlite:runtime.db → language_policy",
        "extract": [
          "language_core_rules",
          "language_anti_drift",
          "language_validation"
        ],
        "compile_rule": "從 neutral-language.md 的核心原則與寫作規則提取 language policy 定義，寫入 runtime.db language_policy 表"
      },
      {
        "source": "enforcement/sanitization.md",
        "target": "sqlite:runtime.db → generated_surfaces (type='sanitization_rules')",
        "extract": [
          "sanitization_rules",
          "placeholder_rules"
        ],
        "compile_rule": "從 sanitization.md 的去敏規則提取 sanitization 定義，寫入 runtime.db generated_surfaces 表"
      },
      {
        "source": "enforcement/tool-neutral-documentation.md",
        "target": "sqlite:runtime.db → generated_surfaces (type='tool_neutrality_rules')",
        "extract": [
          "tool_neutrality_rules",
          "content_boundary_rules"
        ],
        "compile_rule": "從 tool-neutral-documentation.md 的核心規則與允許位置提取 tool neutrality 定義，寫入 runtime.db generated_surfaces 表"
      },
      {
        "source": "plans/active/*.md",
        "target": "sqlite:runtime.db → plans_index",
        "extract": [
          "plan_id",
          "plan_status",
          "plan_priority",
          "phases",
          "affected_files"
        ],
        "compile_rule": "從 plans/active/*.md 的 front matter、phase 標題、受影響檔案表格提取 plan index，寫入 runtime.db plans_index 表"
      },
      {
        "source": "governance/lifecycle/knowledge-update-flow.md",
        "target": "sqlite:runtime.db → generated_surfaces (type='knowledge_update_phases')",
        "extract": [
          "knowledge_update_steps",
          "step_entry_conditions",
          "step_exit_conditions",
          "step_reference_sources"
        ],
        "compile_rule": "從 knowledge-update-flow.md 的 11 個步驟標題與判斷表格提取 phase 定義，寫入 runtime.db generated_surfaces 表"
      },
      {
        "source": "governance/lifecycle/knowledge-update-flow.md",
        "target": "sqlite:runtime.db → generated_surfaces (type='classification_rules')",
        "extract": [
          "classification_dimensions",
          "dimension_decision_conditions",
          "known_languages",
          "known_frameworks",
          "known_platforms"
        ],
        "compile_rule": "從 knowledge-update-flow.md Step 2.4 的決策樹與 intelligence/engineering/ 的 README 提取分類維度定義與已知語言/框架/平台列表，寫入 runtime.db generated_surfaces 表"
      },
      {
        "source": "governance/lifecycle/system-upgrade-governance.md",
        "target": "sqlite:runtime.db → generated_surfaces (type='system_upgrade_governance')",
        "extract": [
          "upgrade_conditions",
          "checklist_categories",
          "checklist_items",
          "forced_rules"
        ],
        "compile_rule": "從系統升級治理要則的 §1 條件表格、§2 檢查清單分類與項目、§3 強制規則提取升級治理定義，寫入 runtime.db generated_surfaces 表"
      },
      {
        "source": "analysis/apk/workflows/*.md",
        "target": "analysis/apk/workflows/generated-phases.yaml",
        "extract": [
          "step_definitions",
          "prerequisites",
          "output_formats"
        ],
        "compile_rule": "從 analysis/apk/workflows/*.md 的「步驟 N：」標題提取 step 定義，從「前置準備」提取 prerequisites，從「成功產出格式」提取 output format"
      }
    ],
    "compilation_rules": [
      {
        "id": "compile.rule.deterministic_only",
        "name": "Deterministic Only",
        "description": "只編譯 deterministic 內容，不編譯需要 judgment 的內容",
        "rule": "如果內容需要 human judgment 才能理解 → 不編譯。如果內容是機械式的 state/action/gate → 編譯。",
        "severity": "critical"
      },
      {
        "id": "compile.rule.header_required",
        "name": "Generated Header Required",
        "description": "每個 generated YAML 必須包含檔頭標註來源",
        "rule": "檔頭包含 generated_from、generated_at、compiler_version、status（synced/stale/orphan）",
        "severity": "high"
      },
      {
        "id": "compile.rule.no_manual_edit",
        "name": "No Manual Edit",
        "description": "generated YAML 不應手動編輯",
        "rule": "所有修改應在 prose source 進行，再透過 compiler 更新 generated YAML",
        "severity": "high"
      },
      {
        "id": "compile.rule.sync_on_commit",
        "name": "Sync on Commit",
        "description": "commit 前必須確保 prose 與 generated YAML 一致",
        "rule": "pre-commit hook 檢查：如果 prose source 修改了，對應的 generated YAML 必須同步更新",
        "severity": "critical"
      },
      {
        "id": "compile.rule.idempotent",
        "name": "Idempotent Compilation",
        "description": "相同 prose source 應產生相同 generated YAML",
        "rule": "compiler 應為 deterministic function：相同輸入 → 相同輸出",
        "severity": "high"
      },
      {
        "id": "compile.rule.backward_compatible",
        "name": "Backward Compatible",
        "description": "compiler 升級不應破壞現有 generated YAML 格式",
        "rule": "新增欄位應為 optional，不應移除現有欄位",
        "severity": "medium"
      }
    ],
    "compiler_workflow": {
      "description": "Compiler 的標準執行流程",
      "steps": [
        {
          "step": 1,
          "action": "SCAN_SOURCES",
          "description": "掃描所有 prose source 檔案",
          "command": "依 source_target_mapping 找出所有需要編譯的 source"
        },
        {
          "step": 2,
          "action": "CHECK_MODIFIED",
          "description": "檢查哪些 source 有修改",
          "command": "比對 git diff 或檔案修改時間"
        },
        {
          "step": 3,
          "action": "EXTRACT_STATE",
          "description": "從修改的 source 中提取 deterministic state",
          "command": "依 extract 欄位定義提取內容"
        },
        {
          "step": 4,
          "action": "GENERATE_YAML",
          "description": "產生對應的 YAML",
          "command": "依 compile_rule 轉換為 YAML 格式"
        },
        {
          "step": 5,
          "action": "UPDATE_HEADER",
          "description": "更新 generated YAML 檔頭",
          "command": "設定 generated_at、status=synced"
        },
        {
          "step": 6,
          "action": "VERIFY",
          "description": "驗證 generated YAML 格式",
          "command": "檢查 YAML 語法、必填欄位、reference 可解析"
        }
      ]
    }
  }

  # Source: runtime/context/ttl-policy.yaml
  CONTEXT_TTL_POLICY =   {
    "ttl_policy_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/context",
    "description": "Context TTL（Time-To-Live）系統。每條 context 載入後依本政策決定存活時間。 過期 context 應被 prune，不永久留在 context graph。\n",
    "default_ttl": {
      "type": "session",
      "value": 1,
      "reason": "預設只活一個 task/session，防止 context 無限制累積"
    },
    "rules": [
      {
        "id": "core-bootstrap",
        "type": "rule-weight",
        "ttl": {
          "type": "session",
          "value": 1
        },
        "reason": "規則權重每個 session 重新載入即可",
        "cacheable": false
      },
      {
        "id": "core-bootstrap-dependency",
        "type": "dependency-reading",
        "ttl": {
          "type": "session",
          "value": 1
        },
        "reason": "依賴讀取規則每個 session 重新載入",
        "cacheable": false
      },
      {
        "id": "core-bootstrap-goal",
        "type": "conversation-goal-ledger",
        "ttl": {
          "type": "conversation",
          "value": true
        },
        "reason": "對話目標 ledger 需要在整個對話中保留",
        "cacheable": true
      },
      {
        "id": "lazy-linked-updates",
        "type": "linked-updates",
        "ttl": {
          "type": "task",
          "value": 1
        },
        "reason": "連動更新規則只在相關 task 中需要",
        "cacheable": false
      },
      {
        "id": "lazy-failure-learning",
        "type": "failure-learning-system",
        "ttl": {
          "type": "task",
          "value": 1
        },
        "reason": "失效學習規則只在發生 failure 時需要",
        "cacheable": false
      },
      {
        "id": "lazy-decision-efficiency",
        "type": "decision-efficiency",
        "ttl": {
          "type": "task",
          "value": 1
        },
        "reason": "決策效率規則只在路線選擇時需要",
        "cacheable": false
      },
      {
        "id": "lazy-tool-neutral",
        "type": "tool-neutral-documentation",
        "ttl": {
          "type": "task",
          "value": 1
        },
        "reason": "工具中立規則只在寫文件時需要",
        "cacheable": false
      },
      {
        "id": "lazy-sanitization",
        "type": "sanitization",
        "ttl": {
          "type": "task",
          "value": 1
        },
        "reason": "去敏規則只在寫 lesson 時需要",
        "cacheable": false
      },
      {
        "id": "lazy-authorization",
        "type": "authorization-scope",
        "ttl": {
          "type": "task",
          "value": 1
        },
        "reason": "授權範圍規則只在涉及外部系統時需要",
        "cacheable": false
      },
      {
        "id": "skill-workflow",
        "type": "skill-workflow",
        "ttl": {
          "type": "task",
          "value": 1
        },
        "reason": "Skill workflow 只在相關 task 中需要",
        "cacheable": true
      },
      {
        "id": "skill-checklist",
        "type": "skill-checklist",
        "ttl": {
          "type": "task",
          "value": 1
        },
        "reason": "Checklist 只在相關 task 中需要",
        "cacheable": true
      },
      {
        "id": "feedback-lesson",
        "type": "feedback-lesson",
        "ttl": {
          "type": "task",
          "value": 1
        },
        "reason": "Feedback lesson 只在相關操作時需要",
        "cacheable": false
      },
      {
        "id": "architecture-roadmap",
        "type": "architecture-roadmap",
        "ttl": {
          "type": "session",
          "value": 1
        },
        "reason": "Roadmap 每個 session 確認一次即可",
        "cacheable": false
      },
      {
        "id": "runtime-registry",
        "type": "runtime-routing-registry",
        "ttl": {
          "type": "conversation",
          "value": true
        },
        "reason": "Routing registry 可在整個對話中 cache，因為它輕量且不常變",
        "cacheable": true
      },
      {
        "id": "knowledge-summary",
        "type": "knowledge-summary",
        "ttl": {
          "type": "conversation",
          "value": true
        },
        "reason": "Summary 輕量（300-500 tokens）且可跨 task 使用",
        "cacheable": true
      },
      {
        "id": "knowledge-graph",
        "type": "knowledge-graph",
        "ttl": {
          "type": "conversation",
          "value": true
        },
        "reason": "Knowledge graph edges 輕量且可跨 task 使用",
        "cacheable": true
      },
      {
        "id": "skills-index",
        "type": "skills-index",
        "ttl": {
          "type": "conversation",
          "value": true
        },
        "reason": "Skills index 輕量且可跨 task 使用",
        "cacheable": true
      },
      {
        "id": "model-profile",
        "type": "model-profile",
        "ttl": {
          "type": "conversation",
          "value": true
        },
        "reason": "Model profile 不常變，可在整個對話中 cache",
        "cacheable": true
      },
      {
        "id": "tool-adapter",
        "type": "tool-adapter",
        "ttl": {
          "type": "session",
          "value": 1
        },
        "reason": "Tool adapter 每個 session 確認一次即可",
        "cacheable": false
      },
      {
        "id": "agent-goals",
        "type": "agent-goals",
        "ttl": {
          "type": "conversation",
          "value": true
        },
        "reason": "Active goals 需要在整個對話中保留",
        "cacheable": true
      }
    ],
    "prune": {
      "trigger": [
        {
          "type": "token_threshold",
          "threshold": "80%"
        },
        {
          "type": "task_boundary",
          "when": "task_completed"
        },
        {
          "type": "session_boundary",
          "when": "new_session"
        }
      ],
      "strategy": {
        "expired_first": true,
        "lowest_priority_first": true,
        "largest_first": true
      },
      "preserve": [
        "core-bootstrap",
        "agent-goals",
        "conversation-goal-ledger"
      ]
    }
  }

  # Source: runtime/discovery/capability-checkpoints.yaml
  DISCOVERY_CAPABILITY_CHECKPOINTS =   {
    "discovery_checkpoints_version": "v1",
    "status": "active",
    "owner_layer": "runtime/discovery",
    "description": "Phase-aware capability discovery checkpoints. 定義 agent 在每個 phase 中 應主動探索的能力類型與搜尋來源。Discovery 是 meta-cognitive 層的操作， 與 lazy-loading（activation engine）互補：activation 被動回應已知條件， discovery 主動搜尋未知能力。\n",
    "checkpoints": [
      {
        "phase": "phase.checkpoint",
        "trigger": "before_execution",
        "description": "在進入 execution 前，搜尋與本輪任務相關的能力。 此 checkpoint 確保 agent 不會因為不知道某個 workflow 或 intelligence atom 而採用次佳策略。\n",
        "discovery_targets": [
          {
            "type": "workflow",
            "description": "搜尋與任務意圖相關的工作流程",
            "search_sources": [
              "knowledge/indexes/README.md",
              "knowledge/runtime/routing-registry.yaml"
            ],
            "min_confidence": "medium",
            "failure_action": "warn + continue"
          },
          {
            "type": "intelligence",
            "description": "搜尋與任務領域相關的 intelligence atoms",
            "search_sources": [
              "knowledge/graphs/README.md",
              "intelligence/README.md"
            ],
            "min_confidence": "low",
            "failure_action": "warn + continue"
          },
          {
            "type": "validation_rules",
            "description": "搜尋與任務相關的驗證規則與 failure patterns",
            "search_sources": [
              "enforcement/failure-patterns/README.md",
              "validation/README.md"
            ],
            "min_confidence": "low",
            "failure_action": "ignore"
          },
          {
            "type": "governance",
            "description": "搜尋與任務相關的治理規則（系統升級、知識更新等）",
            "search_sources": [
              "governance/lifecycle/README.md",
              "governance/validation/README.md"
            ],
            "min_confidence": "low",
            "failure_action": "ignore"
          },
          {
            "type": "analysis",
            "description": "搜尋與任務領域相關的分析方法與觀察框架（analysis/）。當任務涉及分析外部系統、工具、程式庫或新領域時，應檢查 analysis/ 下是否有對應的分析方法。",
            "search_sources": [
              "analysis/README.md"
            ],
            "min_confidence": "medium",
            "failure_action": "warn + continue"
          }
        ]
      },
      {
        "phase": "phase.execution",
        "trigger": "during_execution",
        "description": "在執行過程中，當遇到未知領域或技術問題時，主動搜尋相關 intelligence。 此 checkpoint 是「on-demand discovery」，只在需要時觸發。\n",
        "discovery_targets": [
          {
            "type": "language_specific",
            "description": "搜尋與當前程式語言相關的 intelligence atoms",
            "search_sources": [
              "intelligence/engineering/language-specific/README.md"
            ],
            "min_confidence": "low",
            "failure_action": "ignore"
          },
          {
            "type": "analytical_reasoning",
            "description": "搜尋與當前分析任務相關的啟發式規則",
            "search_sources": [
              "intelligence/engineering/analytical-reasoning/README.md"
            ],
            "min_confidence": "low",
            "failure_action": "ignore"
          },
          {
            "type": "failure_patterns",
            "description": "搜尋與當前操作相關的已知失效模式",
            "search_sources": [
              "enforcement/failure-patterns/README.md"
            ],
            "min_confidence": "low",
            "failure_action": "ignore"
          }
        ]
      },
      {
        "phase": "phase.validation",
        "trigger": "before_validation",
        "description": "在執行驗證前，搜尋與本輪產出相關的驗證規則與 scenarios。 此 checkpoint 確保 agent 不會遺漏重要的驗證項目。\n",
        "discovery_targets": [
          {
            "type": "validation_scenarios",
            "description": "搜尋與本輪產出相關的 validation scenarios",
            "search_sources": [
              "validation/README.md"
            ],
            "min_confidence": "medium",
            "failure_action": "warn + continue"
          },
          {
            "type": "output_governance",
            "description": "搜尋 output governance 規則",
            "search_sources": [
              "runtime/output-governance/README.md"
            ],
            "min_confidence": "medium",
            "failure_action": "warn + continue"
          },
          {
            "type": "linked_updates",
            "description": "搜尋需要連動更新的相關檔案",
            "search_sources": [
              "enforcement/linked-updates.md"
            ],
            "min_confidence": "high",
            "failure_action": "warn + continue"
          }
        ]
      },
      {
        "phase": "phase.finalize",
        "trigger": "before_finalize",
        "description": "在完成前，搜尋 writeback 需求與 transaction 相關規則。 此 checkpoint 確保 agent 不會遺漏回寫步驟。\n",
        "discovery_targets": [
          {
            "type": "writeback_requirements",
            "description": "搜尋回寫需求與 transaction gates",
            "search_sources": [
              "runtime/transactions/transaction-machine.yaml",
              "enforcement/dependency-reading.md"
            ],
            "min_confidence": "high",
            "failure_action": "warn + continue"
          },
          {
            "type": "completion_criteria",
            "description": "搜尋完成條件與 close-loop 規則",
            "search_sources": [
              "enforcement/conversation-goal-ledger.md",
              "scripts/ai-skill-close-loop.sh"
            ],
            "min_confidence": "high",
            "failure_action": "warn + continue"
          }
        ]
      }
    ],
    "search_strategy": {
      "priority_order": [
        "knowledge/indexes/README.md",
        "knowledge/graphs/README.md",
        "knowledge/runtime/routing-registry.yaml"
      ],
      "fallback": [
        {
          "type": "directory_scan",
          "description": "掃描相關目錄結構，尋找未在索引中註冊的檔案",
          "max_depth": 2
        },
        {
          "type": "grep_search",
          "description": "使用 regex 搜尋關鍵字",
          "max_results": 10
        }
      ],
      "min_confidence_threshold": "low"
    }
  }

  # Source: runtime/distributed/async-job-lifecycle.yaml
  DISTRIBUTED_ASYNC_JOB_LIFECYCLE =   {
    "async_job_lifecycle_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/distributed",
    "description": "Background job / async task lifecycle management. 定義 job 從提交到完成的 完整生命週期，以及 job 之間的依賴關係。\n",
    "job_states": [
      {
        "id": "job.state.pending",
        "name": "Pending",
        "description": "Job 已提交，等待排程",
        "allowed_transitions": [
          "job.state.running",
          "job.state.cancelled"
        ]
      },
      {
        "id": "job.state.running",
        "name": "Running",
        "description": "Job 正在執行中",
        "allowed_transitions": [
          "job.state.completed",
          "job.state.failed",
          "job.state.cancelled"
        ]
      },
      {
        "id": "job.state.completed",
        "name": "Completed",
        "description": "Job 執行成功",
        "allowed_transitions": [
  
        ]
      },
      {
        "id": "job.state.failed",
        "name": "Failed",
        "description": "Job 執行失敗",
        "allowed_transitions": [
          "job.state.pending",
          "job.state.cancelled"
        ]
      },
      {
        "id": "job.state.cancelled",
        "name": "Cancelled",
        "description": "Job 已被取消",
        "allowed_transitions": [
  
        ]
      }
    ],
    "job_lifecycle_rules": [
      {
        "id": "job.rule.lock_before_execution",
        "name": "Lock Before Execution",
        "description": "Job 執行前必須取得對應的 lock",
        "rule": "Job 在進入 running state 前，必須先取得 lock.transaction",
        "severity": "critical"
      },
      {
        "id": "job.rule.timeout",
        "name": "Job Timeout",
        "description": "Job 有最大執行時間限制",
        "rule": "Job 在 running state 超過 max_execution_seconds 後自動標記為 failed",
        "severity": "high",
        "default_max_execution_seconds": 600
      },
      {
        "id": "job.rule.retry_policy",
        "name": "Retry Policy",
        "description": "Failed job 可依 retry policy 重新執行",
        "rule": "Failed job 可重新進入 pending state，最多重試 max_retries 次",
        "severity": "medium",
        "default_max_retries": 3,
        "default_retry_backoff_seconds": 30
      },
      {
        "id": "job.rule.dependency_check",
        "name": "Dependency Check",
        "description": "Job 執行前檢查 dependencies 是否已完成",
        "rule": "如果 job 有 depends_on，所有 dependency job 必須為 completed state 才能執行",
        "severity": "critical"
      },
      {
        "id": "job.rule.result_persistence",
        "name": "Result Persistence",
        "description": "Job 完成後必須持久化結果",
        "rule": "Job 在 completed 或 failed state 時，必須將結果寫入指定的 output 位置",
        "severity": "high"
      }
    ],
    "job_templates": [
      {
        "id": "job.template.compilation",
        "name": "Compilation Job",
        "description": "非同步執行 compiler 編譯至 runtime.db",
        "default_max_execution_seconds": 120,
        "default_max_retries": 3,
        "lock_required": "lock.generated_surface",
        "output": "runtime/runtime.db"
      },
      {
        "id": "job.template.validation",
        "name": "Validation Job",
        "description": "非同步執行 validation checks",
        "default_max_execution_seconds": 60,
        "default_max_retries": 2,
        "lock_required": nil,
        "output": "validation results"
      },
      {
        "id": "job.template.recovery",
        "name": "Recovery Job",
        "description": "非同步執行 state recovery",
        "default_max_execution_seconds": 300,
        "default_max_retries": 1,
        "lock_required": "lock.phase_transition",
        "output": "recovery report"
      }
    ]
  }

  # Source: runtime/distributed/distributed-locks.yaml
  DISTRIBUTED_DISTRIBUTED_LOCKS =   {
    "distributed_locks_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/distributed",
    "description": "Distributed lock / lease 機制。當多個 agent 可能同時操作同一 runtime state 時， 需取得對應的 lock 才能執行操作。Lock 有 lease time，逾時自動釋放。\n",
    "locks": [
      {
        "id": "lock.phase_transition",
        "name": "Phase Transition Lock",
        "description": "防止多 agent 同時進行 phase transition",
        "scope": "global",
        "lease_time_seconds": 30,
        "acquire_retry": {
          "max_retries": 3,
          "backoff_seconds": 5
        },
        "release_on": [
          "phase_transition_complete",
          "lease_expired"
        ]
      },
      {
        "id": "lock.obligation_ledger",
        "name": "Obligation Ledger Lock",
        "description": "防止多 agent 同時修改 obligation ledger",
        "scope": "global",
        "lease_time_seconds": 60,
        "acquire_retry": {
          "max_retries": 5,
          "backoff_seconds": 2
        },
        "release_on": [
          "ledger_update_complete",
          "lease_expired"
        ]
      },
      {
        "id": "lock.generated_surface",
        "name": "Generated Surface Lock",
        "description": "防止多 agent 同時編譯 generated YAML",
        "scope": "global",
        "lease_time_seconds": 120,
        "acquire_retry": {
          "max_retries": 3,
          "backoff_seconds": 10
        },
        "release_on": [
          "compilation_complete",
          "lease_expired"
        ]
      },
      {
        "id": "lock.transaction",
        "name": "Transaction Lock",
        "description": "防止多 agent 同時操作同一 transaction",
        "scope": "per_transaction",
        "lease_time_seconds": 300,
        "acquire_retry": {
          "max_retries": 10,
          "backoff_seconds": 5
        },
        "release_on": [
          "transaction_closed",
          "lease_expired"
        ]
      }
    ],
    "lock_rules": [
      {
        "id": "lock.rule.acquire_before_write",
        "rule": "任何 agent 在修改 runtime state 前必須先取得對應的 lock",
        "severity": "critical"
      },
      {
        "id": "lock.rule.release_after_complete",
        "rule": "操作完成後必須釋放 lock，避免 deadlock",
        "severity": "critical"
      },
      {
        "id": "lock.rule.lease_timeout_recovery",
        "rule": "Lease 逾時後，agent 必須執行 recovery 程序才能重新取得 lock",
        "severity": "high",
        "recovery_procedure": "runtime/recovery/phase-reconciliation.yaml"
      },
      {
        "id": "lock.rule.no_nested_locks",
        "rule": "避免 nested lock acquisition，防止 deadlock",
        "severity": "medium"
      }
    ]
  }

  # Source: runtime/distributed/multi-agent-coordination.yaml
  DISTRIBUTED_MULTI_AGENT_COORDINATION =   {
    "multi_agent_coordination_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/distributed",
    "description": "Multi-agent coordination rules. 當多個 agent 同時操作同一 repository 時， 需遵守以下協調規則確保 state consistency。\n",
    "coordination_rules": [
      {
        "id": "coord.rule.phase_check_before_action",
        "name": "Phase Check Before Action",
        "description": "執行任何 action 前先讀取 phase-machine.yaml 確認目前 phase",
        "rule": "Agent 在執行任何 action 前，必須先讀取 runtime/phases/phase-machine.yaml， 確認目前 phase 與 allowed_actions。如果 phase 與預期不符，表示有其他 agent 已變更 state，需執行 reconciliation。\n",
        "severity": "critical",
        "check": {
          "type": "file_read",
          "target": "runtime/phases/phase-machine.yaml",
          "verification": "current_phase_matches_expectation == true"
        }
      },
      {
        "id": "coord.rule.generated_surface_version_check",
        "name": "Generated Surface Version Check",
        "description": "使用 runtime.db 前檢查版本是否最新",
        "rule": "Agent 在查詢 runtime.db 前，必須確認 generated_surfaces 表的記錄 是否為最新版本。如果 prose source 的修改時間晚於最後一次 compiler 執行時間， 表示 runtime.db 已 stale，需先執行 compiler 重新編譯。\n",
        "severity": "high",
        "check": {
          "type": "sqlite_query",
          "target": "runtime/runtime.db",
          "query": "SELECT MAX(updated_at) FROM generated_surfaces",
          "verification": "db_updated_at >= prose_source_mtime"
        }
      },
      {
        "id": "coord.rule.obligation_ledger_sync",
        "name": "Obligation Ledger Sync",
        "description": "修改 obligation ledger 前先取得 lock",
        "rule": "Agent 在修改 runtime/obligations/obligation-ledger.yaml 前，必須先取得 lock.obligation_ledger。修改完成後釋放 lock。如果無法取得 lock， 表示有其他 agent 正在操作，需等待或進入 recovery。\n",
        "severity": "critical",
        "lock_required": "lock.obligation_ledger"
      },
      {
        "id": "coord.rule.transaction_isolation",
        "name": "Transaction Isolation",
        "description": "同一時間只有一個 agent 可以操作同一 transaction",
        "rule": "Agent 在操作 transaction 前，必須先取得 lock.transaction。 Transaction 的 scope 為 per_transaction，不同 transaction 可並行操作。\n",
        "severity": "critical",
        "lock_required": "lock.transaction"
      },
      {
        "id": "coord.rule.phase_transition_exclusive",
        "name": "Phase Transition Exclusive",
        "description": "Phase transition 為獨占操作，同一時間只有一個 agent 可以執行",
        "rule": "Agent 在執行 phase transition 前，必須先取得 lock.phase_transition。 Phase transition 完成後釋放 lock。如果無法取得 lock，表示有其他 agent 正在進行 phase transition，需等待。\n",
        "severity": "critical",
        "lock_required": "lock.phase_transition"
      },
      {
        "id": "coord.rule.stale_state_recovery",
        "name": "Stale State Recovery",
        "description": "偵測到 stale state 時執行 recovery",
        "rule": "如果 agent 偵測到 phase inconsistency、obligation drift、或 generated surface stale，必須先執行 runtime/recovery/phase-reconciliation.yaml 的 reconciliation 程序，才能繼續執行。\n",
        "severity": "critical",
        "recovery_procedure": "runtime/recovery/phase-reconciliation.yaml"
      }
    ],
    "consistency_checks": [
      {
        "id": "check.phase_consistency",
        "name": "Phase Consistency Check",
        "description": "檢查目前 phase 是否與預期一致",
        "check_interval": "每次 action 前",
        "failure_action": "run_reconciliation",
        "failure_message": "Phase inconsistency detected: expected {expected_phase}, actual {actual_phase}"
      },
      {
        "id": "check.obligation_drift",
        "name": "Obligation Drift Check",
        "description": "檢查 obligation ledger 是否有未預期的變更",
        "check_interval": "每次 phase entry",
        "failure_action": "run_reconciliation",
        "failure_message": "Obligation drift detected: obligation ledger has unexpected changes"
      },
      {
        "id": "check.generated_surface_stale",
        "name": "Generated Surface Stale Check",
        "description": "檢查 generated YAML 是否為最新版本",
        "check_interval": "每次讀取 generated YAML 前",
        "failure_action": "run_compiler",
        "failure_message": "Generated surface is stale: {target} generated_at {generated_at}, prose modified at {prose_mtime}"
      }
    ]
  }

  # Source: runtime/gates/blocking-gates.yaml
  GATES_BLOCKING_GATES =   {
    "blocking_gates_version": "v1",
    "status": "active",
    "owner_layer": "runtime/gates",
    "description": "Per-phase blocking gates. 定義每個 phase 中必須通過的檢查點。 Gate 是 phase transition 的必要條件。未通過 critical/high gate 時， agent 不能進入下一個 phase。\n",
    "gates": [
      {
        "id": "gate.bootstrap.core_rules_loaded",
        "phase": "phase.bootstrap",
        "name": "Core Rules Loaded",
        "description": "CORE_BOOTSTRAP.md 的 3 條核心規則已載入",
        "severity": "critical",
        "check": {
          "type": "file_read",
          "target": "CORE_BOOTSTRAP.md",
          "verification": "core_bootstrap_loaded == true"
        },
        "failure_action": "halt_bootstrap",
        "failure_message": "無法繼續：CORE_BOOTSTRAP.md 未載入"
      },
      {
        "id": "gate.bootstrap.layout_loaded",
        "phase": "phase.bootstrap",
        "name": "Layout Loaded",
        "description": "README.md（OS layout）已載入",
        "severity": "high",
        "check": {
          "type": "file_read",
          "target": "README.md",
          "verification": "layout_loaded == true"
        },
        "failure_action": "warn + continue",
        "failure_message": "警告：README.md 未載入，可能缺少 OS layout 資訊"
      },
      {
        "id": "gate.checkpoint.obligation_ledger_loaded",
        "phase": "phase.checkpoint",
        "name": "Obligation Ledger Loaded",
        "description": "Obligation ledger 已讀取，pending obligations 已確認",
        "severity": "critical",
        "check": {
          "type": "file_read",
          "target": "runtime/obligations/obligation-ledger.yaml",
          "verification": "obligation_ledger_loaded == true AND pending_obligations_reviewed == true"
        },
        "failure_action": "block_transition",
        "failure_message": "無法進入 execution：obligation ledger 未讀取或有未完成的 obligation"
      },
      {
        "id": "gate.checkpoint.blocking_gates_checked",
        "phase": "phase.checkpoint",
        "name": "Blocking Gates Checked",
        "description": "本輪的 blocking gates 已檢查，無 blocking condition",
        "severity": "critical",
        "check": {
          "type": "file_read",
          "target": "runtime/gates/blocking-gates.yaml",
          "verification": "blocking_gates_checked == true AND no_blocking_condition == true"
        },
        "failure_action": "block_transition",
        "failure_message": "無法進入 execution：有 blocking gate 未通過"
      },
      {
        "id": "gate.checkpoint.knowledge_update_checked",
        "phase": "phase.checkpoint",
        "name": "Knowledge Update Checked",
        "description": "知識更新檢查已完成（快速路徑：runtime.db generated_surfaces 表）",
        "severity": "high",
        "check": {
          "type": "sqlite_query",
          "target": "runtime/runtime.db",
          "query": "SELECT content FROM generated_surfaces WHERE type = 'knowledge_update_phases'",
          "verification": "knowledge_update_checked == true"
        },
        "failure_action": "warn + continue",
        "failure_message": "警告：知識更新檢查未完成，建議在 checkpoint 階段執行"
      },
      {
        "id": "gate.checkpoint.system_upgrade_governance_checked",
        "phase": "phase.checkpoint",
        "name": "System Upgrade Governance Checked",
        "description": "系統升級治理檢查已完成（快速路徑：runtime.db generated_surfaces 表）。 如有進行中的大型升級計畫（plans/active/），必須確認治理檢查清單已執行。\n",
        "severity": "high",
        "check": {
          "type": "sqlite_query",
          "target": "runtime/runtime.db",
          "query": "SELECT content FROM generated_surfaces WHERE type = 'system_upgrade_governance'",
          "verification": "system_upgrade_governance_checked == true AND (no_active_upgrade_plan == true OR checklist_items_reviewed == true)"
        },
        "failure_action": "warn + continue",
        "failure_message": "警告：系統升級治理檢查未完成，如有進行中的大型升級請確認檢查清單"
      },
      {
        "id": "gate.checkpoint.capability_discovery_completed",
        "phase": "phase.checkpoint",
        "name": "Capability Discovery Completed",
        "description": "能力探索檢查已完成。Agent 已根據 runtime/discovery/capability-checkpoints.yaml 搜尋與本輪任務相關的工作流程、intelligence atoms、驗證規則與治理規則。 此 gate 為 medium severity，未通過時記錄警告但不 block transition。\n",
        "severity": "medium",
        "check": {
          "type": "file_read",
          "target": "runtime/discovery/capability-checkpoints.yaml",
          "verification": "capability_discovery_completed == true AND discovery_targets_searched == true"
        },
        "failure_action": "warn + continue",
        "failure_message": "警告：能力探索檢查未完成，建議在 checkpoint 階段執行 discovery 以發現可能不知道的能力"
      },
      {
        "id": "gate.execution.goal_defined",
        "phase": "phase.execution",
        "name": "Goal Defined",
        "description": "本輪的執行目標已明確定義",
        "severity": "critical",
        "check": {
          "type": "verification",
          "target": ".agent-goals/ 或 conversation context",
          "verification": "execution_goal_defined == true"
        },
        "failure_action": "block_execution",
        "failure_message": "無法執行：本輪目標未定義"
      },
      {
        "id": "gate.execution.scope_clear",
        "phase": "phase.execution",
        "name": "Scope Clear",
        "description": "執行範圍已明確，無 scope creep 風險",
        "severity": "high",
        "check": {
          "type": "verification",
          "target": "goal definition",
          "verification": "execution_scope_clear == true"
        },
        "failure_action": "warn + confirm",
        "failure_message": "警告：執行範圍不明確，建議先釐清 scope"
      },
      {
        "id": "gate.validation.all_obligations_met",
        "phase": "phase.validation",
        "name": "All Obligations Met",
        "description": "本輪所有 obligation 已完成",
        "severity": "critical",
        "check": {
          "type": "verification",
          "target": "runtime/obligations/obligation-ledger.yaml",
          "verification": "all_phase_obligations_fulfilled == true"
        },
        "failure_action": "block_transition",
        "failure_message": "無法進入 commit：有未完成的 obligation"
      },
      {
        "id": "gate.validation.linked_updates_complete",
        "phase": "phase.validation",
        "name": "Linked Updates Complete",
        "description": "所有連動更新已執行（依 enforcement/linked-updates.md）",
        "severity": "critical",
        "check": {
          "type": "verification",
          "target": "enforcement/linked-updates.md",
          "verification": "all_linked_updates_applied == true OR no_linked_updates_needed == true"
        },
        "failure_action": "block_transition",
        "failure_message": "無法進入 commit：有連動更新未完成"
      },
      {
        "id": "gate.validation.artifacts_complete",
        "phase": "phase.validation",
        "name": "Artifacts Complete",
        "description": "本輪產出的 artifacts 完整且符合規範",
        "severity": "high",
        "check": {
          "type": "verification",
          "target": "artifact list",
          "verification": "all_required_artifacts_produced == true"
        },
        "failure_action": "warn + suggest_fix",
        "failure_message": "警告：部分 artifacts 可能不完整"
      },
      {
        "id": "gate.validation.no_forbidden_actions_used",
        "phase": "phase.validation",
        "name": "No Forbidden Actions Used",
        "description": "本輪未使用任何 forbidden actions",
        "severity": "critical",
        "check": {
          "type": "audit",
          "target": "action log",
          "verification": "no_forbidden_actions_in_log == true"
        },
        "failure_action": "block_transition + escalate",
        "failure_message": "嚴重：偵測到 forbidden action 使用，需進入 recovery"
      },
      {
        "id": "gate.commit.staging_complete",
        "phase": "phase.commit",
        "name": "Staging Complete",
        "description": "所有變更已 git add",
        "severity": "critical",
        "check": {
          "type": "command",
          "target": "git diff --cached",
          "verification": "staging_complete == true"
        },
        "failure_action": "block_commit",
        "failure_message": "無法 commit：有變更未 staging"
      },
      {
        "id": "gate.commit.commit_message_valid",
        "phase": "phase.commit",
        "name": "Commit Message Valid",
        "description": "Commit message 符合 conventional commit 格式",
        "severity": "high",
        "check": {
          "type": "verification",
          "target": "commit message",
          "verification": "commit_message_format_valid == true"
        },
        "failure_action": "warn + suggest_format",
        "failure_message": "警告：commit message 格式不標準"
      },
      {
        "id": "gate.commit.pre_commit_hooks_passed",
        "phase": "phase.commit",
        "name": "Pre-commit Hooks Passed",
        "description": "Pre-commit hook 執行成功",
        "severity": "critical",
        "check": {
          "type": "command",
          "target": "git commit 輸出",
          "verification": "pre_commit_hooks_exit_code == 0"
        },
        "failure_action": "block_commit",
        "failure_message": "無法 commit：pre-commit hook 失敗"
      },
      {
        "id": "gate.push.remote_reachable",
        "phase": "phase.push",
        "name": "Remote Reachable",
        "description": "遠端 repository 可連線",
        "severity": "critical",
        "check": {
          "type": "command",
          "target": "git remote -v",
          "verification": "remote_reachable == true"
        },
        "failure_action": "block_push",
        "failure_message": "無法 push：遠端無法連線"
      },
      {
        "id": "gate.push.push_successful",
        "phase": "phase.push",
        "name": "Push Successful",
        "description": "git push 成功",
        "severity": "critical",
        "check": {
          "type": "command",
          "target": "git push 輸出",
          "verification": "push_exit_code == 0"
        },
        "failure_action": "block_transition",
        "failure_message": "無法進入 readback：push 失敗"
      },
      {
        "id": "gate.readback.working_tree_clean",
        "phase": "phase.readback",
        "name": "Working Tree Clean",
        "description": "git status 顯示工作樹乾淨",
        "severity": "critical",
        "check": {
          "type": "command",
          "target": "git status --short",
          "verification": "working_tree_clean == true"
        },
        "failure_action": "block_transition",
        "failure_message": "無法進入 finalize：工作樹不乾淨"
      },
      {
        "id": "gate.readback.remote_synced",
        "phase": "phase.readback",
        "name": "Remote Synced",
        "description": "本地與遠端同步",
        "severity": "high",
        "check": {
          "type": "command",
          "target": "git log --oneline origin/main..main",
          "verification": "local_remote_synced == true"
        },
        "failure_action": "warn + suggest_push",
        "failure_message": "警告：本地與遠端不同步，建議重新 push"
      },
      {
        "id": "gate.finalize.all_obligations_closed",
        "phase": "phase.finalize",
        "name": "All Obligations Closed",
        "description": "本輪所有 obligation 已關閉",
        "severity": "critical",
        "check": {
          "type": "verification",
          "target": "runtime/obligations/obligation-ledger.yaml",
          "verification": "all_obligations_closed == true"
        },
        "failure_action": "block_finalize",
        "failure_message": "無法 finalize：有 obligation 未關閉"
      },
      {
        "id": "gate.finalize.knowledge_updated",
        "phase": "phase.finalize",
        "name": "Knowledge Updated",
        "description": "知識更新 checkpoint 已執行",
        "severity": "high",
        "check": {
          "type": "verification",
          "target": "governance/lifecycle/knowledge-update-flow.md",
          "verification": "knowledge_checkpoint_completed == true"
        },
        "failure_action": "warn + suggest_update",
        "failure_message": "警告：知識更新 checkpoint 未執行"
      },
      {
        "id": "gate.recovery.repair_strategy_defined",
        "phase": "phase.recovery",
        "name": "Repair Strategy Defined",
        "description": "修復策略已定義，有明確的修復步驟",
        "severity": "critical",
        "check": {
          "type": "verification",
          "target": "recovery plan",
          "verification": "repair_strategy_defined == true"
        },
        "failure_action": "halt_recovery",
        "failure_message": "無法修復：修復策略未定義"
      }
    ],
    "metadata": {
      "priority": "P0",
      "confidence": "high",
      "compatibility_state": "new",
      "last_checked": "2026-05-15"
    }
  }

  # Source: runtime/guards/circuit-breaker.yaml
  GUARDS_CIRCUIT_BREAKER =   {
    "circuit_breaker_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/guards",
    "recursive_depth": {
      "description": "防止 agent 在同一 task 中無限遞迴",
      "max_depth": 4,
      "reset_condition": "task_completed",
      "on_exceed": [
        {
          "action": "halt_agent",
          "message": "⛔ Recursive depth exceeded ({depth}/{max_depth}). Agent halted."
        },
        {
          "action": "suggest_decomposition",
          "message": "建議將目前 task 拆解為更小的子任務，使用 .agent-goals/ 管理。"
        }
      ]
    },
    "tool_calls": {
      "description": "防止 tool 呼叫爆炸（search → search → search loop）",
      "max_calls_per_task": 20,
      "max_calls_per_session": 100,
      "window_minutes": 5,
      "max_calls_per_window": 15,
      "on_exceed": [
        {
          "action": "warn",
          "message": "⚠️ Tool call rate high ({calls}/{max_calls_per_window} in {window_minutes}m). 請確認是否陷入 tool loop。"
        },
        {
          "action": "suggest_consolidate",
          "message": "建議合併多個 tool call 為一次綜合查詢。"
        }
      ]
    },
    "context_growth": {
      "description": "防止 context 在單一 task 中快速膨脹",
      "max_growth_per_task": 0.3,
      "max_growth_per_session": 0.8,
      "check_interval": "task_boundary",
      "on_exceed": [
        {
          "action": "warn",
          "message": "⚠️ Context growth exceeded ({growth}/{max_growth}). 建議 prune 或建立新 session。"
        },
        {
          "action": "force_prune",
          "strategy": "moderate"
        }
      ]
    },
    "hallucination_risk": {
      "description": "偵測高 hallucination 風險情境",
      "risk_factors": [
        {
          "factor": "stale_summary",
          "weight": 0.3,
          "signal": "summary 與 source 不一致"
        },
        {
          "factor": "conflicting_rules",
          "weight": 0.3,
          "signal": "多條規則對同一情境有不同指示"
        },
        {
          "factor": "insufficient_source",
          "weight": 0.2,
          "signal": "未讀 canonical source 就做出結論"
        },
        {
          "factor": "tool_output_truncated",
          "weight": 0.2,
          "signal": "tool output 被截斷但 agent 仍做判斷"
        }
      ],
      "risk_threshold": 0.7,
      "on_exceed": [
        {
          "action": "halt_agent",
          "message": "⛔ Hallucination risk score: {score}. Agent halted. 請先讀取 canonical source 再繼續。"
        },
        {
          "action": "suggest_source_reading",
          "message": "建議先讀取以下 canonical source：{sources}"
        }
      ]
    },
    "conflict_rules": {
      "description": "偵測目前 context 中的規則衝突",
      "check_on": "task_start",
      "sources": [
        "enforcement/rule-weight.md",
        "enforcement/README.md"
      ],
      "on_conflict": [
        {
          "action": "warn",
          "message": "⚠️ 偵測到規則衝突。依 rule-weight.md 判斷優先序。"
        },
        {
          "action": "suggest_resolution",
          "message": "P0（安全/source）> P1（validation）> P2（user goal）> P3（tool adapter）> P4（efficiency）"
        }
      ]
    }
  }

  # Source: runtime/guards/context-pollution.yaml
  GUARDS_CONTEXT_POLLUTION =   {
    "context_pollution_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/guards",
    "signals": [
      {
        "name": "conversation_length",
        "description": "對話輪數過多",
        "threshold": 50,
        "severity": "warning",
        "action": [
          "⚠️ 對話已達 {turns} 輪。Context 可能已膨脹。建議建立新 session。",
          "使用 memory/summary/ 保存目前進度後建立新 session。"
        ]
      },
      {
        "name": "repetitive_edits",
        "description": "同一文件反覆修改",
        "threshold": 5,
        "severity": "warning",
        "action": [
          "⚠️ 檔案 {file} 已修改 {edits} 次。建議確認修改方向是否穩定。",
          "若持續修改同一檔案，考慮拆分子任務。"
        ]
      },
      {
        "name": "module_count",
        "description": "同時載入的 modules/rules 過多",
        "threshold": 20,
        "severity": "critical",
        "action": [
          "⛔ Context 中有 {count} 個 modules/rules。超過 {threshold} 門檻。",
          "建議 prune 非必要 modules，或建立新 session。"
        ]
      },
      {
        "name": "cross_reference_depth",
        "description": "跨文件引用深度過深",
        "threshold": 5,
        "severity": "warning",
        "action": [
          "⚠️ 跨文件引用深度達 {depth} 層。可能導致 context 快速膨脹。",
          "建議先讀 summary 再決定是否需要展開完整引用鏈。"
        ]
      },
      {
        "name": "token_utilization",
        "description": "Token 使用率過高",
        "threshold": 0.85,
        "severity": "critical",
        "action": [
          "⛔ Token 使用率達 {percentage}%。Context 污染風險高。",
          "建議立即 prune 或建立新 session。"
        ]
      }
    ],
    "pollution_score": {
      "formula": "weighted_average(signals)",
      "thresholds": {
        "healthy": "< 0.3",
        "warning": "0.3 - 0.6",
        "critical": "> 0.6"
      },
      "on_critical": [
        {
          "action": "halt_agent",
          "message": "⛔ Context pollution score: {score}. 建議建立新 session 並存檔目前進度。"
        },
        {
          "action": "auto_archive",
          "target": "memory/working/",
          "format": "session-archive-{timestamp}.md"
        }
      ]
    }
  }

  # Source: runtime/health/context-health-score.yaml
  HEALTH_CONTEXT_HEALTH_SCORE =   {
    "context_health_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/health",
    "dimensions": [
      {
        "name": "relevance",
        "description": "目前 context 與當前 task 的相關性",
        "weight": 0.35,
        "signals": [
          {
            "type": "task_intent_match",
            "source": "routing-registry.yaml"
          },
          {
            "type": "activation_rule_match",
            "source": "runtime/router/activation-rules.yaml"
          }
        ],
        "scoring": {
          "0.0-0.3": "critical（多數 context 與 task 無關）",
          "0.3-0.6": "warning（部分 context 過期）",
          "0.6-1.0": "healthy"
        }
      },
      {
        "name": "duplication",
        "description": "Context 中重複知識的比例",
        "weight": 0.2,
        "signals": [
          {
            "type": "cross_reference_overlap"
          },
          {
            "type": "summary_source_duplicate"
          }
        ],
        "scoring": {
          "0.0-0.2": "healthy",
          "0.2-0.5": "warning（有可合併的 context）",
          "0.5-1.0": "critical（大量重複）"
        }
      },
      {
        "name": "staleness",
        "description": "Context 中過期知識的比例",
        "weight": 0.25,
        "signals": [
          {
            "type": "ttl_expired",
            "source": "runtime/context/ttl-policy.yaml"
          },
          {
            "type": "refresh_needed",
            "source": "knowledge/runtime/refresh-policy.yaml"
          }
        ],
        "scoring": {
          "0.0-0.2": "healthy",
          "0.2-0.4": "warning（部分 context 需 refresh）",
          "0.4-1.0": "critical（大量過期 context）"
        }
      },
      {
        "name": "conflict",
        "description": "Context 中規則衝突的比例",
        "weight": 0.2,
        "signals": [
          {
            "type": "rule_conflict",
            "source": "enforcement/rule-weight.md"
          },
          {
            "type": "source_mirror_drift",
            "source": "enforcement/failure-patterns/source-mirror-write-drift.md"
          }
        ],
        "scoring": {
          "0.0-0.1": "healthy",
          "0.1-0.3": "warning（有可解決的衝突）",
          "0.3-1.0": "critical（需人工介入）"
        }
      }
    ],
    "composite": {
      "formula": "relevance * 0.35 + (1 - duplication) * 0.20 + (1 - staleness) * 0.25 + (1 - conflict) * 0.20",
      "thresholds": {
        "healthy": ">= 0.75",
        "warning": ">= 0.50",
        "critical": "< 0.50"
      }
    },
    "actions": {
      "healthy": [
        {
          "action": "continue",
          "message": "✅ Context health score: {score}. No action needed."
        }
      ],
      "warning": [
        {
          "action": "suggest_prune",
          "message": "⚠️ Context health score: {score}. 建議 prune 低 relevance / 過期 context。"
        },
        {
          "action": "suggest_refresh",
          "message": "建議 refresh stale summaries 與 routing registry。"
        }
      ],
      "critical": [
        {
          "action": "halt_agent",
          "message": "⛔ Context health score: {score}. Agent halted. 建議建立新 session。"
        },
        {
          "action": "force_prune",
          "strategy": "aggressive",
          "preserve": [
            "core-bootstrap",
            "agent-goals"
          ]
        }
      ]
    }
  }

  # Source: runtime/intelligence/intelligence-routing.yaml
  INTELLIGENCE_INTELLIGENCE_ROUTING =   {
    "version": 1.0,
    "status": "candidate",
    "owner_layer": "runtime/intelligence",
    "last_updated": "2026-05-15",
    "routing_rules": [
      {
        "id": "route.intelligence.analytical_reasoning",
        "domain": "analytical_reasoning",
        "path": "intelligence/engineering/analytical-reasoning/",
        "applicable_phases": [
          "phase.execution",
          "phase.validation"
        ],
        "task_intent_patterns": [
          "分析",
          "逆向",
          "偵測",
          "triage",
          "capture",
          "證據"
        ],
        "load_strategy": "lazy",
        "priority": "P2",
        "description": "分析推理 heuristic，適用於 APK 分析、流量分析、偵測等任務"
      },
      {
        "id": "route.intelligence.agent_architecture",
        "domain": "agent_architecture",
        "path": "intelligence/engineering/agent-architecture/",
        "applicable_phases": [
          "phase.checkpoint",
          "phase.finalize"
        ],
        "task_intent_patterns": [
          "架構",
          "規則",
          "設計",
          "migration",
          "重構"
        ],
        "load_strategy": "lazy",
        "priority": "P2",
        "description": "Agent 架構知識，適用於規則設計、架構遷移等任務"
      },
      {
        "id": "route.intelligence.development",
        "domain": "development",
        "path": "intelligence/engineering/development/",
        "applicable_phases": [
          "phase.execution",
          "phase.validation"
        ],
        "task_intent_patterns": [
          "開發",
          "BDD",
          "測試",
          "contract",
          "程式碼"
        ],
        "load_strategy": "lazy",
        "priority": "P2",
        "description": "開發流程 heuristic，適用於 BDD closure、contract governance 等任務"
      },
      {
        "id": "route.intelligence.heuristics",
        "domain": "heuristics",
        "path": "intelligence/engineering/heuristics/",
        "applicable_phases": [
          "phase.execution"
        ],
        "task_intent_patterns": [
          "heuristic",
          "經驗",
          "判斷",
          "通用"
        ],
        "load_strategy": "lazy",
        "priority": "P3",
        "description": "通用 engineering heuristics"
      },
      {
        "id": "route.intelligence.anti_patterns",
        "domain": "anti_patterns",
        "path": "intelligence/engineering/anti-patterns/",
        "applicable_phases": [
          "phase.validation",
          "phase.finalize"
        ],
        "task_intent_patterns": [
          "反模式",
          "anti-pattern",
          "可疑",
          "重複問題"
        ],
        "load_strategy": "lazy",
        "priority": "P3",
        "description": "Engineering anti-patterns，適用於 code review 與設計審查"
      },
      {
        "id": "route.intelligence.tradeoffs",
        "domain": "tradeoffs",
        "path": "intelligence/engineering/tradeoffs/",
        "applicable_phases": [
          "phase.execution",
          "phase.checkpoint"
        ],
        "task_intent_patterns": [
          "技術決策",
          "tradeoff",
          "比較",
          "選擇"
        ],
        "load_strategy": "lazy",
        "priority": "P3",
        "description": "技術 tradeoff 分析"
      },
      {
        "id": "route.intelligence.distributed_systems",
        "domain": "distributed_systems",
        "path": "intelligence/engineering/distributed-systems/",
        "applicable_phases": [
          "phase.execution"
        ],
        "task_intent_patterns": [
          "分散式",
          "distributed",
          "一致性",
          "共識"
        ],
        "load_strategy": "lazy",
        "priority": "P3",
        "description": "分散式系統模式與 heuristic"
      },
      {
        "id": "route.intelligence.domain",
        "domain": "domain",
        "path": "intelligence/engineering/domain/",
        "applicable_phases": [
          "phase.execution"
        ],
        "task_intent_patterns": [
          "領域建模",
          "domain",
          "aggregate",
          "bounded context"
        ],
        "load_strategy": "lazy",
        "priority": "P3",
        "description": "領域驅動設計 heuristic"
      },
      {
        "id": "route.intelligence.failure",
        "domain": "failure",
        "path": "intelligence/engineering/failure/",
        "applicable_phases": [
          "phase.validation",
          "phase.recovery"
        ],
        "task_intent_patterns": [
          "失效",
          "failure",
          "錯誤",
          "異常",
          "blocking gate"
        ],
        "load_strategy": "lazy",
        "priority": "P2",
        "description": "工程 failure pattern，適用於 blocking gate 阻斷或驗證失敗時"
      },
      {
        "id": "route.intelligence.travel",
        "domain": "travel",
        "path": "intelligence/travel/",
        "applicable_phases": [
          "phase.execution"
        ],
        "task_intent_patterns": [
          "旅行",
          "旅遊",
          "travel",
          "itinerary"
        ],
        "load_strategy": "lazy",
        "priority": "P2",
        "description": "旅行規劃 heuristic"
      },
      {
        "id": "route.intelligence.business",
        "domain": "business",
        "path": "intelligence/business/",
        "applicable_phases": [
          "phase.execution"
        ],
        "task_intent_patterns": [
          "商業",
          "business",
          "定價",
          "pricing",
          "SaaS"
        ],
        "load_strategy": "lazy",
        "priority": "P3",
        "description": "商業分析 heuristic"
      },
      {
        "id": "route.intelligence.ide",
        "domain": "ide",
        "path": "intelligence/ide/",
        "applicable_phases": [
          "phase.execution"
        ],
        "task_intent_patterns": [
          "IDE",
          "VSCode",
          "設定",
          "extension"
        ],
        "load_strategy": "lazy",
        "priority": "P3",
        "description": "IDE 設定與使用 heuristic"
      }
    ],
    "fallback": {
      "strategy": "full_search",
      "target": "intelligence/README.md",
      "description": "若無 routing rule 匹配，fallback 到 intelligence/README.md 全文搜尋"
    }
  }

  # Source: runtime/obligations/obligation-ledger.yaml
  OBLIGATIONS_OBLIGATION_LEDGER =   {
    "obligation_ledger_version": "v1",
    "status": "active",
    "owner_layer": "runtime/obligations",
    "description": "Per-phase obligation ledger. 定義每個 phase 中 agent 必須完成的 obligation。 Obligation 是 atomic duty，完成後標記 fulfilled。所有 obligation 都 fulfilled 後才能進行 phase transition（除非進入 recovery）。\n",
    "obligations": [
      {
        "id": "obligation.bootstrap.init_ledger",
        "phase": "phase.bootstrap",
        "name": "Initialize Obligation Ledger",
        "description": "初始化 obligation ledger，建立空的 obligation 追蹤表",
        "verification": [
          "obligation_ledger_initialized == true"
        ],
        "severity": "critical",
        "depends_on": [
  
        ],
        "linked_gates": [
          "gate.bootstrap.core_rules_loaded"
        ]
      },
      {
        "id": "obligation.checkpoint.review_pending_obligations",
        "phase": "phase.checkpoint",
        "name": "Review Pending Obligations",
        "description": "讀取 obligation ledger，檢查上一輪是否有未完成的 obligation",
        "verification": [
          "pending_obligations_reviewed == true",
          "no_stale_obligations == true OR stale_obligations_documented == true"
        ],
        "severity": "critical",
        "depends_on": [
          "obligation.bootstrap.init_ledger"
        ],
        "linked_gates": [
          "gate.checkpoint.obligation_ledger_loaded"
        ]
      },
      {
        "id": "obligation.checkpoint.check_knowledge_update_needed",
        "phase": "phase.checkpoint",
        "name": "Check Knowledge Update Needed",
        "description": "快速路徑（預設）：查詢 runtime.db（generated_surfaces 表 WHERE type='knowledge_update_phases'） 了解 11 個步驟的結構與參考文件，自問本輪是否新增可重用知識。 完整路徑（首次或需要細節時）：讀 governance/lifecycle/knowledge-update-flow.md（~522 lines）\n",
        "verification": [
          "knowledge_update_check_completed == true"
        ],
        "severity": "high",
        "depends_on": [
  
        ],
        "linked_gates": [
          "gate.checkpoint.blocking_gates_checked"
        ]
      },
      {
        "id": "obligation.checkpoint.check_system_upgrade_governance",
        "phase": "phase.checkpoint",
        "name": "Check System Upgrade Governance",
        "description": "快速路徑（預設）：查詢 runtime.db（generated_surfaces 表 WHERE type='system_upgrade_governance'） 了解大型系統升級的條件、檢查清單分類與強制規則。 如果 plans/active/ 中有大型升級計畫，必須逐項確認檢查清單。 完整路徑：讀 governance/lifecycle/system-upgrade-governance.md（~244 lines）\n",
        "verification": [
          "system_upgrade_governance_checked == true",
          "if_active_upgrade_plan: checklist_items_reviewed == true"
        ],
        "severity": "high",
        "depends_on": [
          "obligation.checkpoint.check_knowledge_update_needed"
        ],
        "linked_gates": [
          "gate.checkpoint.system_upgrade_governance_checked"
        ]
      },
      {
        "id": "obligation.checkpoint.run_capability_discovery",
        "phase": "phase.checkpoint",
        "name": "Run Capability Discovery",
        "description": "執行能力探索檢查點。讀取 runtime/discovery/capability-checkpoints.yaml， 根據本輪任務意圖搜尋可能不知道的工作流程、intelligence atoms、 驗證規則與治理規則。使用 knowledge indexes、graphs 與 routing registry 作為搜尋來源。此 obligation 解決 Capability Discovery Problem： lazy loading 本身不會產生探索意識，agent 需要主動搜尋未知能力。\n",
        "verification": [
          "capability_discovery_completed == true",
          "discovery_targets_searched == true",
          "if_relevant_workflow_found: workflow_loaded == true",
          "if_relevant_intelligence_found: intelligence_loaded == true"
        ],
        "severity": "medium",
        "depends_on": [
          "obligation.checkpoint.check_system_upgrade_governance"
        ],
        "linked_gates": [
          "gate.checkpoint.capability_discovery_completed"
        ]
      },
      {
        "id": "obligation.execution.log_actions",
        "phase": "phase.execution",
        "name": "Log Actions",
        "description": "記錄本輪執行的所有 action，供 validation 階段檢查",
        "verification": [
          "action_log_maintained == true"
        ],
        "severity": "medium",
        "depends_on": [
  
        ],
        "linked_gates": [
  
        ]
      },
      {
        "id": "obligation.execution.track_artifacts",
        "phase": "phase.execution",
        "name": "Track Artifacts",
        "description": "追蹤本輪產出的 artifacts（檔案修改、命令輸出等）",
        "verification": [
          "artifact_list_maintained == true"
        ],
        "severity": "medium",
        "depends_on": [
  
        ],
        "linked_gates": [
  
        ]
      },
      {
        "id": "obligation.validation.run_knowledge_runtime_validator",
        "phase": "phase.validation",
        "name": "Run Knowledge Runtime Validator",
        "description": "執行 ruby scripts/validate-knowledge-runtime.rb 檢查 knowledge runtime 一致性",
        "verification": [
          "validator_exit_code == 0",
          "validator_output_no_errors == true"
        ],
        "severity": "critical",
        "depends_on": [
          "obligation.execution.log_actions"
        ],
        "linked_gates": [
          "gate.validation.all_obligations_met"
        ]
      },
      {
        "id": "obligation.validation.check_linked_updates",
        "phase": "phase.validation",
        "name": "Check Linked Updates",
        "description": "依 enforcement/linked-updates.md 檢查所有連動更新是否完整",
        "verification": [
          "linked_updates_checked == true",
          "all_linked_updates_applied == true OR no_linked_updates_needed == true"
        ],
        "severity": "critical",
        "depends_on": [
          "obligation.execution.track_artifacts"
        ],
        "linked_gates": [
          "gate.validation.linked_updates_complete"
        ]
      },
      {
        "id": "obligation.validation.check_output_governance",
        "phase": "phase.validation",
        "name": "Check Output Governance",
        "description": "檢查 output governance gates：語言一致性、去敏狀態、工具中立性、中性語言",
        "verification": [
          "output_governance_checked == true",
          "gate.output.language_consistency == passed",
          "gate.output.sanitization_check == passed",
          "gate.output.tool_neutrality == passed OR warn_only == true",
          "gate.output.neutral_language == passed OR warn_only == true"
        ],
        "severity": "critical",
        "depends_on": [
          "obligation.validation.run_knowledge_runtime_validator"
        ],
        "linked_gates": [
          "gate.output.language_consistency",
          "gate.output.sanitization_check",
          "gate.output.tool_neutrality",
          "gate.output.neutral_language"
        ]
      },
      {
        "id": "obligation.commit.write_commit_message",
        "phase": "phase.commit",
        "name": "Write Commit Message",
        "description": "撰寫符合 conventional commit 格式的 commit message",
        "verification": [
          "commit_message_written == true",
          "commit_message_format_valid == true"
        ],
        "severity": "critical",
        "depends_on": [
  
        ],
        "linked_gates": [
          "gate.commit.commit_message_valid"
        ]
      },
      {
        "id": "obligation.commit.verify_pre_commit_hooks",
        "phase": "phase.commit",
        "name": "Verify Pre-commit Hooks",
        "description": "確認 pre-commit hook 執行成功，無 blocking errors",
        "verification": [
          "pre_commit_hooks_passed == true"
        ],
        "severity": "critical",
        "depends_on": [
          "obligation.commit.write_commit_message"
        ],
        "linked_gates": [
          "gate.commit.pre_commit_hooks_passed"
        ]
      },
      {
        "id": "obligation.push.verify_push_success",
        "phase": "phase.push",
        "name": "Verify Push Success",
        "description": "確認 git push 成功，遠端已接收所有 commit",
        "verification": [
          "push_exit_code == 0",
          "remote_ref_updated == true"
        ],
        "severity": "critical",
        "depends_on": [
  
        ],
        "linked_gates": [
          "gate.push.push_successful"
        ]
      },
      {
        "id": "obligation.readback.verify_clean_tree",
        "phase": "phase.readback",
        "name": "Verify Clean Working Tree",
        "description": "執行 git status 確認工作樹乾淨，無未追蹤或未提交的變更",
        "verification": [
          "working_tree_clean == true"
        ],
        "severity": "critical",
        "depends_on": [
          "obligation.push.verify_push_success"
        ],
        "linked_gates": [
          "gate.readback.working_tree_clean"
        ]
      },
      {
        "id": "obligation.readback.confirm_sync",
        "phase": "phase.readback",
        "name": "Confirm Remote Sync",
        "description": "確認本地與遠端同步（git log 比對）",
        "verification": [
          "local_remote_synced == true"
        ],
        "severity": "high",
        "depends_on": [
          "obligation.readback.verify_clean_tree"
        ],
        "linked_gates": [
          "gate.readback.remote_synced"
        ]
      },
      {
        "id": "obligation.finalize.close_ledger",
        "phase": "phase.finalize",
        "name": "Close Obligation Ledger",
        "description": "關閉本輪的 obligation ledger，清理已完成的 obligation 記錄",
        "verification": [
          "all_obligations_fulfilled == true",
          "ledger_closed == true"
        ],
        "severity": "critical",
        "depends_on": [
          "obligation.readback.confirm_sync"
        ],
        "linked_gates": [
          "gate.finalize.all_obligations_closed"
        ]
      },
      {
        "id": "obligation.finalize.run_knowledge_checkpoint",
        "phase": "phase.finalize",
        "name": "Run Knowledge Update Checkpoint",
        "description": "執行 governance/lifecycle/knowledge-update-flow.md 的知識更新檢查",
        "verification": [
          "knowledge_checkpoint_completed == true"
        ],
        "severity": "high",
        "depends_on": [
          "obligation.finalize.close_ledger"
        ],
        "linked_gates": [
          "gate.finalize.knowledge_updated"
        ]
      },
      {
        "id": "obligation.finalize.verify_output_governance",
        "phase": "phase.finalize",
        "name": "Verify Output Governance",
        "description": "最終 output governance 驗證：語言一致性、去敏狀態",
        "verification": [
          "final_output_governance_verified == true",
          "gate.output.final_language_check == passed",
          "gate.output.final_sanitization_check == passed"
        ],
        "severity": "critical",
        "depends_on": [
          "obligation.finalize.run_knowledge_checkpoint"
        ],
        "linked_gates": [
          "gate.output.final_language_check",
          "gate.output.final_sanitization_check"
        ]
      },
      {
        "id": "obligation.recovery.diagnose_root_cause",
        "phase": "phase.recovery",
        "name": "Diagnose Root Cause",
        "description": "診斷 state inconsistency 的根本原因",
        "verification": [
          "root_cause_identified == true"
        ],
        "severity": "critical",
        "depends_on": [
  
        ],
        "linked_gates": [
          "gate.recovery.repair_strategy_defined"
        ]
      },
      {
        "id": "obligation.recovery.apply_repair",
        "phase": "phase.recovery",
        "name": "Apply Repair",
        "description": "執行修復策略，修復 state/obligation/phase 不一致",
        "verification": [
          "repair_applied == true"
        ],
        "severity": "critical",
        "depends_on": [
          "obligation.recovery.diagnose_root_cause"
        ],
        "linked_gates": [
  
        ]
      },
      {
        "id": "obligation.recovery.verify_repair",
        "phase": "phase.recovery",
        "name": "Verify Repair",
        "description": "驗證修復是否成功，state 是否已一致",
        "verification": [
          "repair_verified == true",
          "state_consistent == true"
        ],
        "severity": "critical",
        "depends_on": [
          "obligation.recovery.apply_repair"
        ],
        "linked_gates": [
  
        ]
      }
    ],
    "metadata": {
      "priority": "P0",
      "confidence": "high",
      "compatibility_state": "new",
      "last_checked": "2026-05-15"
    }
  }

  # Source: runtime/output-governance/governance-gates.yaml
  OUTPUT_GOVERNANCE_GOVERNANCE_GATES =   {
    "governance_gates_version": "v1",
    "status": "active",
    "owner_layer": "runtime/output-governance",
    "parent": "runtime/gates/blocking-gates.yaml",
    "description": "Output governance blocking gates. 定義 validation 與 finalize phase 中 必須通過的語言與文件輸出檢查點。這些 gate 是 blocking-gates.yaml 的子集， 專注於輸出品質控制。\n",
    "gates": [
      {
        "id": "gate.output.language_consistency",
        "phase": "phase.validation",
        "name": "Language Consistency",
        "description": "所有輸出語言與使用者最後一次提問的語言一致",
        "severity": "critical",
        "check": {
          "type": "verification",
          "target": "runtime/output-governance/language-policy.yaml",
          "verification": "output_language_matches_user_language == true"
        },
        "failure_action": "block_transition",
        "failure_message": "無法通過 validation：輸出語言與使用者語言不一致，需修正後重新驗證",
        "canonical_source": [
          "enforcement/neutral-language.md",
          ".roomodes"
        ]
      },
      {
        "id": "gate.output.no_default_english",
        "phase": "phase.validation",
        "name": "No Default English",
        "description": "System prompt / custom instructions 中沒有 'Default to English' 固定設定",
        "severity": "high",
        "check": {
          "type": "audit",
          "target": ".roomodes, CLAUDE.md, .cursor/rules/*.mdc",
          "verification": "no_default_english_in_config == true"
        },
        "failure_action": "warn + suggest_fix",
        "failure_message": "警告：偵測到 'Default to English' 固定設定，建議改為軟性偏好",
        "canonical_source": [
          "ai-tools/agent/roo.md",
          "ai-tools/agent/claude.md",
          "ai-tools/agent/cursor.md"
        ]
      },
      {
        "id": "gate.output.sanitization_check",
        "phase": "phase.validation",
        "name": "Sanitization Check",
        "description": "所有輸出已去敏：無 token、密鑰、raw request、真實路徑、project incident 細節",
        "severity": "critical",
        "check": {
          "type": "audit",
          "target": "runtime/output-governance/output-rules.yaml",
          "verification": "all_output_sanitized == true"
        },
        "failure_action": "block_transition",
        "failure_message": "無法通過 validation：輸出包含未去敏的敏感資訊，需修正後重新驗證",
        "canonical_source": [
          "enforcement/sanitization.md"
        ]
      },
      {
        "id": "gate.output.placeholder_usage",
        "phase": "phase.validation",
        "name": "Placeholder Usage",
        "description": "可重用文件使用 <AI_SKILL_REPO>、<PROJECT_ROOT>、<WORKSPACE> 等占位符",
        "severity": "high",
        "check": {
          "type": "audit",
          "target": "reusable files",
          "verification": "placeholders_used_instead_of_real_paths == true"
        },
        "failure_action": "warn + suggest_fix",
        "failure_message": "警告：可重用文件包含真實路徑，建議改用占位符",
        "canonical_source": [
          "enforcement/sanitization.md"
        ]
      },
      {
        "id": "gate.output.tool_neutrality",
        "phase": "phase.validation",
        "name": "Tool Neutrality",
        "description": "可重用文件未包含工具專屬措辭（除非在允許位置）",
        "severity": "high",
        "check": {
          "type": "audit",
          "target": "enforcement/, workflow/, analysis/, intelligence/, reusable skills",
          "verification": "tool_specific_terms_only_in_allowed_locations == true"
        },
        "failure_action": "warn + suggest_fix",
        "failure_message": "警告：可重用文件包含工具專屬措辭，建議移至 ai-tools/ 或使用通用詞彙",
        "canonical_source": [
          "enforcement/tool-neutral-documentation.md"
        ]
      },
      {
        "id": "gate.output.format_compliance",
        "phase": "phase.validation",
        "name": "Format Compliance",
        "description": "輸出格式符合 markdown link、table、code block 規則",
        "severity": "medium",
        "check": {
          "type": "audit",
          "target": "all output",
          "verification": "output_format_compliant == true"
        },
        "failure_action": "warn + suggest_fix",
        "failure_message": "建議：輸出格式不完全符合規範，建議修正",
        "canonical_source": [
          "enforcement/neutral-language.md"
        ]
      },
      {
        "id": "gate.output.neutral_language",
        "phase": "phase.validation",
        "name": "Neutral Language",
        "description": "可重用文件使用中性描述，無高風險詞在標題/檔名/slug",
        "severity": "high",
        "check": {
          "type": "audit",
          "target": "enforcement/, workflow/, analysis/, intelligence/, reusable skills",
          "verification": "neutral_language_used == true AND no_high_risk_words_in_title == true"
        },
        "failure_action": "warn + suggest_fix",
        "failure_message": "警告：可重用文件包含高風險詞或非中性描述，建議改用中性詞",
        "canonical_source": [
          "enforcement/neutral-language.md"
        ]
      },
      {
        "id": "gate.output.final_language_check",
        "phase": "phase.finalize",
        "name": "Final Language Check",
        "description": "最終輸出前再次確認語言一致性",
        "severity": "critical",
        "check": {
          "type": "verification",
          "target": "attempt_completion result",
          "verification": "final_output_language_matches_user_language == true"
        },
        "failure_action": "block_finalize",
        "failure_message": "無法 finalize：最終輸出語言與使用者語言不一致",
        "canonical_source": [
          ".roomodes",
          "language-policy.yaml"
        ]
      },
      {
        "id": "gate.output.final_sanitization_check",
        "phase": "phase.finalize",
        "name": "Final Sanitization Check",
        "description": "最終輸出前再次確認去敏狀態",
        "severity": "critical",
        "check": {
          "type": "audit",
          "target": "attempt_completion result",
          "verification": "final_output_sanitized == true"
        },
        "failure_action": "block_finalize",
        "failure_message": "無法 finalize：最終輸出包含未去敏的敏感資訊",
        "canonical_source": [
          "enforcement/sanitization.md"
        ]
      }
    ]
  }

  # Source: runtime/output-governance/language-policy.yaml
  OUTPUT_GOVERNANCE_LANGUAGE_POLICY =   {
    "language_policy_version": "v1",
    "status": "active",
    "owner_layer": "runtime/output-governance",
    "canonical_sources": [
      "enforcement/neutral-language.md",
      ".roomodes",
      "ai-tools/agent/roo.md",
      "ai-tools/agent/claude.md",
      "ai-tools/agent/cursor.md"
    ],
    "description": "定義所有 agent 輸出必須遵守的語言規則。核心原則是「強制跟隨使用者語言」， 不存在預設語言。所有輸出（技術分析、表格、commit message、程式碼註解） 都必須與使用者當前語言一致。\n",
    "core_rules": [
      {
        "id": "lang.follow_user_language",
        "priority": "highest",
        "rule": "使用者用什麼語言提問，所有輸出就必須用什麼語言。 使用者第一次提問的語言即為本次對話語言。\n",
        "applies_to": [
          "all_output",
          "attempt_completion",
          "technical_analysis",
          "table_fields",
          "section_titles",
          "commit_messages",
          "code_comments"
        ],
        "exceptions": [
          "路徑、指令、環境變數、API 欄位、程式符號、固定 protocol 名稱、工具名稱",
          "翻譯後會降低精準度的專有名詞（如 \"AI-native Cognitive Execution System\"、\"reference-first\"）",
          "prompt template 內容（但 template 外的用途說明仍須跟隨使用者語言）"
        ]
      },
      {
        "id": "lang.no_default_language",
        "priority": "highest",
        "rule": "不存在「Default to English」。不存在預設語言，只有「使用者當前語言」。 禁止出現「Default to English」思維。\n",
        "applies_to": [
          "system_prompt",
          "custom_instructions",
          "all_output"
        ]
      },
      {
        "id": "lang.no_self_switch",
        "priority": "highest",
        "rule": "禁止自行切換語言。即使分析技術內容、寫 commit message、產生表格， 也必須與使用者當前語言一致。\n",
        "applies_to": [
          "all_output"
        ]
      },
      {
        "id": "lang.follow_user_switch",
        "priority": "high",
        "rule": "如果使用者切換語言，跟隨其切換。\n",
        "applies_to": [
          "all_output"
        ]
      },
      {
        "id": "lang.pre_attempt_completion_check",
        "priority": "high",
        "rule": "attempt_completion 前必須檢查：確認結果的語言與使用者最後一次提問的語言一致。\n",
        "applies_to": [
          "attempt_completion"
        ]
      }
    ],
    "anti_drift": [
      {
        "id": "drift.zh_consistency",
        "condition": "使用者用中文",
        "rule": "所有輸出（包含技術分析、表格欄位、章節標題、commit message、程式碼註解） 都必須使用中文。禁止中英混雜成段落。\n",
        "applies_to": [
          "all_output"
        ]
      },
      {
        "id": "drift.ja_consistency",
        "condition": "使用者用日文",
        "rule": "所有輸出都必須使用日文。\n",
        "applies_to": [
          "all_output"
        ]
      },
      {
        "id": "drift.en_consistency",
        "condition": "使用者用英文",
        "rule": "所有輸出都必須使用英文。\n",
        "applies_to": [
          "all_output"
        ]
      },
      {
        "id": "drift.neutral_language",
        "rule": "可重用文件正文預設使用繁體中文。英文保留給路徑、指令、環境變數、 API 欄位、程式符號與必要專有名詞。\n",
        "applies_to": [
          "enforcement/",
          "workflow/",
          "analysis/",
          "intelligence/",
          "reusable_skills",
          "feedback_lessons"
        ],
        "canonical_source": "enforcement/neutral-language.md"
      },
      {
        "id": "drift.no_high_risk_words_in_title",
        "rule": "標題、檔名、slug、索引列與 one-line summary 使用中性詞。 避免把特定品牌、產品、成人、仇恨、暴力、賭博、金融濫用、憑證、 繞過、破解、攻擊、盜用等高風險詞放在可重用文件標題或檔名。\n",
        "applies_to": [
          "enforcement/",
          "workflow/",
          "analysis/",
          "intelligence/",
          "reusable_skills"
        ],
        "canonical_source": "enforcement/neutral-language.md"
      }
    ],
    "tool_overrides": [
      {
        "tool": "roo_code",
        "config_location": ".roomodes (customInstructions) + SQLite globalStorage",
        "setup_steps": [
          "在 .roomodes 每個 mode 的 customInstructions 中加入語言強制規則（見 core_rules）",
          "修改 VS Code globalStorage SQLite 中 RooVeterinaryInc.roo-cline 的 language 欄位為 'zh-CN'"
        ],
        "note": ".roomodes 中的 customInstructions 會覆蓋 Roo Code Extension 設定面板中的全域 Custom Instructions。 SQLite 的 language 欄位（預設 \"en\"）會影響 system prompt 建構，需一併修改。\n",
        "reference": [
          "ai-tools/agent/roo.md",
          "intelligence/ide/vscode-extension-global-state.md"
        ]
      },
      {
        "tool": "claude_code",
        "config_location": "CLAUDE.md",
        "setup_steps": [
          "在 CLAUDE.md 中加入語言強制規則（見 core_rules）",
          "加入語言一致性強制規則段落"
        ],
        "note": "Claude Code 沒有 SQLite 全域資料庫，語言行為完全由 CLAUDE.md 中的 Custom Instructions 控制。\n",
        "reference": [
          "ai-tools/agent/claude.md",
          "CLAUDE.md"
        ]
      },
      {
        "tool": "cursor",
        "config_location": ".cursor/rules/*.mdc (alwaysApply: true)",
        "setup_steps": [
          "在 ~/.cursor/rules/ 的 alwaysApply: true 規則中加入語言強制規則（見 core_rules）",
          "加入語言一致性強制規則段落"
        ],
        "note": "Cursor 沒有公開的 SQLite 全域資料庫修改方式，語言偏好完全由 .cursor/rules/ 中的規則控制。 alwaysApply 規則會在每個 session 啟動時自動載入。\n",
        "reference": [
          "ai-tools/agent/cursor.md"
        ]
      }
    ],
    "validation": [
      {
        "id": "lang.validate_language_consistency",
        "phase": "phase.validation",
        "check": "確認所有輸出語言與使用者最後一次提問的語言一致",
        "failure_action": "block_transition",
        "failure_message": "語言一致性檢查失敗：輸出語言與使用者語言不一致"
      },
      {
        "id": "lang.validate_no_default_english",
        "phase": "phase.validation",
        "check": "確認 system prompt / custom instructions 中沒有 'Default to English' 固定設定",
        "failure_action": "warn + suggest_fix",
        "failure_message": "警告：偵測到 'Default to English' 固定設定，建議改為軟性偏好"
      }
    ]
  }

  # Source: runtime/output-governance/output-rules.yaml
  OUTPUT_GOVERNANCE_OUTPUT_RULES =   {
    "output_rules_version": "v1",
    "status": "active",
    "owner_layer": "runtime/output-governance",
    "canonical_sources": [
      "enforcement/sanitization.md",
      "enforcement/tool-neutral-documentation.md",
      "enforcement/neutral-language.md"
    ],
    "description": "定義所有可重用文件的輸出格式、內容邊界、去敏規則與工具中立性要求。 適用範圍：enforcement/、workflow/、analysis/、intelligence/、reusable skills、feedback lessons。\n",
    "format_rules": [
      {
        "id": "fmt.markdown_links",
        "rule": "所有回應中對 `language construct` 或 filename 的引用必須使用 clickable link 格式： [`filename`](relative/file/path.ext:line) 或 [`language.declaration()`](relative/file/path.ext:line)。 line 對 `syntax` 為必填，對 filename 為選填。\n",
        "applies_to": [
          "all_responses",
          "attempt_completion"
        ]
      },
      {
        "id": "fmt.table_format",
        "rule": "表格欄位優先使用繁體中文（除非欄位內容為固定英文術語）。 表格前後需有空行。\n",
        "applies_to": [
          "all_documents"
        ],
        "canonical_source": "enforcement/neutral-language.md"
      },
      {
        "id": "fmt.code_blocks",
        "rule": "Code block 可保留原始語言（英文程式碼、指令、環境變數）。 Code block 外的說明文字須跟隨使用者語言。\n",
        "applies_to": [
          "all_documents"
        ]
      },
      {
        "id": "fmt.no_mixed_language_paragraphs",
        "rule": "文件正文不要中英混雜成段落。除非文件本身是工具輸入範本、prompt、 程式註解或外部規格引用，否則標題、摘要、說明段落與表格欄位 優先使用繁體中文。\n",
        "applies_to": [
          "enforcement/",
          "workflow/",
          "analysis/",
          "intelligence/",
          "reusable_skills",
          "feedback_lessons"
        ],
        "canonical_source": "enforcement/neutral-language.md"
      },
      {
        "id": "fmt.english_retention",
        "rule": "下列內容可保留英文：檔名、路徑、命令、環境變數、程式符號、 API 欄位、固定 protocol 名稱、工具名稱，以及翻譯後會降低精準度的 專有名詞（如 \"AI-native Cognitive Execution System\"、\"reference-first\"、 \"Default Bootstrap\"、\"rule-weight\"、\"source-of-truth\"）。\n",
        "applies_to": [
          "all_documents"
        ],
        "canonical_source": "enforcement/neutral-language.md"
      },
      {
        "id": "fmt.prompt_template_language",
        "rule": "Prompt template 可保留英文，方便直接貼給 agent。 Template 外的用途說明、前後文與驗證要求仍以繁體中文描述。\n",
        "applies_to": [
          "reusable_skills",
          "templates"
        ],
        "canonical_source": "enforcement/neutral-language.md"
      }
    ],
    "content_boundary": [
      {
        "id": "boundary.tool_neutral",
        "rule": "可重用規則、skills、templates、索引與 feedback lessons 預設應保持工具中立。 使用通用詞彙（agent、AI tool、tool-specific adapter、local tool mirror、project tool config）。 除非規則確實只適用於單一工具，否則不要讓特定工具聽起來像必要條件。\n",
        "applies_to": [
          "enforcement/",
          "workflow/",
          "analysis/",
          "intelligence/",
          "reusable_skills",
          "feedback_lessons"
        ],
        "canonical_source": "enforcement/tool-neutral-documentation.md"
      },
      {
        "id": "boundary.tool_specific_allowed_locations",
        "rule": "工具名稱、路徑與 UI 操作只可出現在： - ai-tools/<tool>.md - tools/adapters/<tool>.md - 工具設定檔（.cursor/rules/*.mdc、.roomodes、CLAUDE.md） - 工具專屬 scripts 或 script docs - Project-local tool files\n",
        "applies_to": [
          "enforcement/",
          "workflow/",
          "analysis/",
          "intelligence/",
          "reusable_skills"
        ],
        "canonical_source": "enforcement/tool-neutral-documentation.md"
      },
      {
        "id": "boundary.no_tool_specific_in_generic",
        "rule": "除非章節明確討論工具整合，否則下列位置避免工具專屬措辭： - Root README.md - enforcement/README.md 索引摘要 - workflow/<domain>/、analysis/<domain>/、intelligence/<domain>/ 下的所有文件 - Skill templates 與 skills/ADDING_SKILLS.md - Feedback lessons 與可重用 checklists\n",
        "applies_to": [
          "root_readme",
          "enforcement_readme",
          "workflow",
          "analysis",
          "intelligence",
          "skill_templates",
          "feedback_lessons"
        ],
        "canonical_source": "enforcement/tool-neutral-documentation.md"
      },
      {
        "id": "boundary.reusable_guidance",
        "rule": "可重用文件不得包含 project-specific 的具體 app/project 名稱、endpoint、 host、payload fragment、sample ID、class/test 名稱、live run 結果或環境 quirks。 這些依 reusable-guidance-boundary.md 留在專案文件。\n",
        "applies_to": [
          "reusable_skills",
          "feedback_lessons",
          "enforcement/"
        ],
        "canonical_source": "enforcement/reusable-guidance-boundary.md"
      },
      {
        "id": "boundary.neutral_description",
        "rule": "使用功能、角色、資料流、狀態、風險類別描述，不使用嘲諷、煽動、獵奇、 誇大或道德判斷語。授權語境：分析、測試、逆向、抓包、工具、SDK 等內容 都要明確放在授權、合規、測試、文件化或相容性語境中。\n",
        "applies_to": [
          "enforcement/",
          "workflow/",
          "analysis/",
          "intelligence/",
          "reusable_skills",
          "feedback_lessons"
        ],
        "canonical_source": "enforcement/neutral-language.md"
      },
      {
        "id": "boundary.no_conclusion_change",
        "rule": "中性化不能把未授權行為寫成授權行為，也不能把安全風險寫成已解決。 若有風險，仍要明確描述控制、限制與驗證方法。\n",
        "applies_to": [
          "all_documents"
        ],
        "canonical_source": "enforcement/neutral-language.md"
      }
    ],
    "sanitization": [
      {
        "id": "sanit.no_auth_tokens",
        "rule": "可重用文件不得包含完整 Authorization / token、session cookie、可識別特定使用者的實體裝置識別",
        "severity": "critical",
        "applies_to": [
          "reusable_skills",
          "feedback_lessons",
          "enforcement/"
        ],
        "canonical_source": "enforcement/sanitization.md"
      },
      {
        "id": "sanit.no_crypto_keys",
        "rule": "可重用文件不得包含 AES/HMAC/簽名密鑰（除非教學用、合成、可公開測試向量）",
        "severity": "critical",
        "applies_to": [
          "reusable_skills",
          "feedback_lessons",
          "enforcement/"
        ],
        "canonical_source": "enforcement/sanitization.md"
      },
      {
        "id": "sanit.no_raw_requests",
        "rule": "可重用文件不得包含未去敏的 raw request/response",
        "severity": "high",
        "applies_to": [
          "reusable_skills",
          "feedback_lessons",
          "enforcement/"
        ],
        "canonical_source": "enforcement/sanitization.md"
      },
      {
        "id": "sanit.no_local_paths",
        "rule": "可重用文件不得包含本機真實絕對路徑、使用者帳號名稱、私用工作目錄、git clone 實體路徑",
        "severity": "high",
        "applies_to": [
          "reusable_skills",
          "feedback_lessons",
          "enforcement/"
        ],
        "canonical_source": "enforcement/sanitization.md"
      },
      {
        "id": "sanit.no_project_incident_details",
        "rule": "可重用文件不得包含 project incident 的具體 app/project 名稱、endpoint、host、payload fragment",
        "severity": "high",
        "applies_to": [
          "reusable_skills",
          "feedback_lessons",
          "enforcement/"
        ],
        "canonical_source": "enforcement/sanitization.md"
      },
      {
        "id": "sanit.use_placeholders",
        "rule": "一律改用 <AI_SKILL_REPO>、<PROJECT_ROOT>、<WORKSPACE> 等占位符",
        "applies_to": [
          "reusable_skills",
          "feedback_lessons",
          "enforcement/"
        ],
        "canonical_source": "enforcement/sanitization.md"
      },
      {
        "id": "sanit.original_terms_allowed",
        "rule": "可以在專案私有文件或證據附件中保留最小必要原詞，但必須同時滿足： 1. 原詞是重現、合約、UI 對照、法規或錯誤碼判讀所必需 2. 文件明確標示授權範圍與去敏狀態 3. 不會被複製到 reusable skill、enforcement rule、lesson、模板或公開範例 4. 可用中性別名在索引、摘要與 cross-link 中引用\n",
        "applies_to": [
          "project_private_documents",
          "evidence_attachments"
        ],
        "canonical_source": "enforcement/neutral-language.md"
      }
    ],
    "validation": [
      {
        "id": "output.validate_sanitization",
        "phase": "phase.validation",
        "check": "確認所有輸出已去敏：無 token、密鑰、raw request、真實路徑、project incident 細節",
        "failure_action": "block_transition",
        "failure_message": "去敏檢查失敗：輸出包含未去敏的敏感資訊"
      },
      {
        "id": "output.validate_tool_neutrality",
        "phase": "phase.validation",
        "check": "確認可重用文件未包含工具專屬措辭（除非在允許位置）",
        "failure_action": "warn + suggest_fix",
        "failure_message": "警告：可重用文件包含工具專屬措辭，建議移至 ai-tools/ 或使用通用詞彙"
      },
      {
        "id": "output.validate_format",
        "phase": "phase.validation",
        "check": "確認輸出格式符合 markdown link、table、code block 規則",
        "failure_action": "warn + suggest_fix",
        "failure_message": "警告：輸出格式不符合規範"
      },
      {
        "id": "output.validate_neutral_language",
        "phase": "phase.validation",
        "check": "確認可重用文件使用中性描述，無高風險詞在標題/檔名/slug",
        "failure_action": "warn + suggest_fix",
        "failure_message": "警告：可重用文件包含高風險詞或非中性描述"
      }
    ]
  }

  # Source: runtime/phases/phase-machine.yaml
  PHASES_PHASE_MACHINE =   {
    "phase_machine_version": "v1",
    "status": "active",
    "owner_layer": "runtime/phases",
    "description": "AI session execution phase state machine. 定義從 bootstrap 到 close-loop 的完整 phase 轉換圖。每個 phase 包含 entry_conditions、allowed/forbidden actions、 blocking gates、obligations 與 phase transition 規則。 Agent 必須在每次 action 前重新讀取此檔案，確認目前 phase 的 allowed actions。\n",
    "phases": [
      {
        "id": "phase.bootstrap",
        "name": "Bootstrap",
        "description": "Session 啟動階段。載入最小必要知識（CORE_BOOTSTRAP.md 3 rules）， 建立 session context，初始化 runtime state。\n",
        "entry_conditions": [
          {
            "session_start": true
          }
        ],
        "allowed_actions": [
          "read_runtime_config",
          "read_bootstrap",
          "check_git_status",
          "init_obligation_ledger"
        ],
        "forbidden_actions": [
          "execute_workflow",
          "commit",
          "push",
          "finalize"
        ],
        "blocking_gates": [
          "gate.bootstrap.core_rules_loaded",
          "gate.bootstrap.layout_loaded"
        ],
        "obligations": [
          "obligation.bootstrap.init_ledger"
        ],
        "next_phase": "phase.checkpoint",
        "phase_transition_triggers": [
          {
            "condition": "core_rules_loaded == true AND layout_loaded == true",
            "auto_transition": true
          }
        ]
      },
      {
        "id": "phase.checkpoint",
        "name": "Checkpoint",
        "description": "每輪工作開始前的檢查點。確認目前 state、讀取 obligation ledger、 檢查 blocking gates、執行知識更新檢查（快速路徑： 查詢 runtime.db generated_surfaces 表）、決定本輪要做什麼。\n",
        "entry_conditions": [
          {
            "bootstrap_complete": true
          },
          {
            "OR": [
              {
                "session_start": true
              },
              {
                "previous_phase_complete": true
              }
            ]
          }
        ],
        "allowed_actions": [
          "read_obligation_ledger",
          "read_phase_machine",
          "read_blocking_gates",
          "read_knowledge_update_phases",
          "read_system_upgrade_governance",
          "run_capability_discovery",
          "diff_review",
          "update_knowledge",
          "plan_next_actions"
        ],
        "forbidden_actions": [
          "finalize",
          "commit",
          "push"
        ],
        "blocking_gates": [
          "gate.checkpoint.obligation_ledger_loaded",
          "gate.checkpoint.blocking_gates_checked",
          "gate.checkpoint.knowledge_update_checked",
          "gate.checkpoint.system_upgrade_governance_checked",
          "gate.checkpoint.capability_discovery_completed"
        ],
        "obligations": [
          "obligation.checkpoint.review_pending_obligations",
          "obligation.checkpoint.check_knowledge_update_needed",
          "obligation.checkpoint.check_system_upgrade_governance",
          "obligation.checkpoint.run_capability_discovery"
        ],
        "next_phase": "phase.execution",
        "phase_transition_triggers": [
          {
            "condition": "obligation_ledger_loaded == true AND blocking_gates_checked == true AND knowledge_update_checked == true AND system_upgrade_governance_checked == true AND capability_discovery_completed == true",
            "auto_transition": false,
            "requires_confirmation": true
          }
        ]
      },
      {
        "id": "phase.execution",
        "name": "Execution",
        "description": "主要執行階段。執行本輪的 allowed actions，產出對應的 artifacts。 此階段是唯一允許修改檔案和執行命令的 phase。\n",
        "entry_conditions": [
          {
            "checkpoint_complete": true
          },
          {
            "blocking_gates_passed": true
          }
        ],
        "allowed_actions": [
          "read_file",
          "write_file",
          "apply_diff",
          "execute_command",
          "search_files",
          "read_workflow_prose",
          "read_domain_knowledge"
        ],
        "forbidden_actions": [
          "commit",
          "push",
          "finalize",
          "transition_phase"
        ],
        "blocking_gates": [
          "gate.execution.goal_defined",
          "gate.execution.scope_clear"
        ],
        "obligations": [
          "obligation.execution.log_actions",
          "obligation.execution.track_artifacts"
        ],
        "next_phase": "phase.validation",
        "phase_transition_triggers": [
          {
            "condition": "execution_complete == true",
            "auto_transition": false,
            "requires_confirmation": true
          }
        ]
      },
      {
        "id": "phase.validation",
        "name": "Validation",
        "description": "驗證階段。檢查 execution 的產出是否符合 gates、obligations 是否完成、 linked updates 是否完整、output governance 是否通過。 此階段不允許修改檔案。\n",
        "entry_conditions": [
          {
            "execution_complete": true
          }
        ],
        "allowed_actions": [
          "read_file",
          "diff_review",
          "run_validator",
          "check_linked_updates",
          "verify_obligations",
          "check_output_governance"
        ],
        "forbidden_actions": [
          "write_file",
          "apply_diff",
          "execute_command",
          "commit",
          "push",
          "finalize"
        ],
        "blocking_gates": [
          "gate.validation.all_obligations_met",
          "gate.validation.linked_updates_complete",
          "gate.validation.artifacts_complete",
          "gate.validation.no_forbidden_actions_used",
          "gate.output.language_consistency",
          "gate.output.sanitization_check",
          "gate.output.tool_neutrality",
          "gate.output.neutral_language"
        ],
        "obligations": [
          "obligation.validation.run_knowledge_runtime_validator",
          "obligation.validation.check_linked_updates",
          "obligation.validation.check_output_governance"
        ],
        "next_phase": "phase.commit",
        "phase_transition_triggers": [
          {
            "condition": "all_gates_passed == true",
            "auto_transition": false,
            "requires_confirmation": true
          },
          {
            "condition": "all_gates_passed == false",
            "auto_transition": false,
            "requires_confirmation": true,
            "target_phase": "phase.execution"
          }
        ]
      },
      {
        "id": "phase.commit",
        "name": "Commit",
        "description": "提交階段。執行 git add、git commit，確保所有變更被正確記錄。 此階段是 writeback transaction 的一部分。\n",
        "entry_conditions": [
          {
            "validation_passed": true
          }
        ],
        "allowed_actions": [
          "git_add",
          "git_commit",
          "read_file",
          "diff_review"
        ],
        "forbidden_actions": [
          "write_file",
          "apply_diff",
          "execute_command",
          "finalize"
        ],
        "blocking_gates": [
          "gate.commit.staging_complete",
          "gate.commit.commit_message_valid",
          "gate.commit.pre_commit_hooks_passed"
        ],
        "obligations": [
          "obligation.commit.write_commit_message",
          "obligation.commit.verify_pre_commit_hooks"
        ],
        "next_phase": "phase.push",
        "phase_transition_triggers": [
          {
            "condition": "commit_successful == true",
            "auto_transition": false,
            "requires_confirmation": true
          }
        ]
      },
      {
        "id": "phase.push",
        "name": "Push",
        "description": "推送階段。執行 git push，將本地 commit 推送到遠端。\n",
        "entry_conditions": [
          {
            "commit_successful": true
          }
        ],
        "allowed_actions": [
          "git_push",
          "read_file"
        ],
        "forbidden_actions": [
          "write_file",
          "apply_diff",
          "execute_command",
          "finalize"
        ],
        "blocking_gates": [
          "gate.push.remote_reachable",
          "gate.push.push_successful"
        ],
        "obligations": [
          "obligation.push.verify_push_success"
        ],
        "next_phase": "phase.readback",
        "phase_transition_triggers": [
          {
            "condition": "push_successful == true",
            "auto_transition": false,
            "requires_confirmation": true
          }
        ]
      },
      {
        "id": "phase.readback",
        "name": "Readback",
        "description": "讀回階段。push 完成後執行 git status 確認工作樹乾淨， 並讀回遠端狀態確認同步成功。這是 writeback transaction 的最後一步。\n",
        "entry_conditions": [
          {
            "push_successful": true
          }
        ],
        "allowed_actions": [
          "git_status",
          "git_log",
          "read_file",
          "diff_review"
        ],
        "forbidden_actions": [
          "write_file",
          "apply_diff",
          "execute_command",
          "finalize"
        ],
        "blocking_gates": [
          "gate.readback.working_tree_clean",
          "gate.readback.remote_synced"
        ],
        "obligations": [
          "obligation.readback.verify_clean_tree",
          "obligation.readback.confirm_sync"
        ],
        "next_phase": "phase.finalize",
        "phase_transition_triggers": [
          {
            "condition": "working_tree_clean == true AND remote_synced == true",
            "auto_transition": false,
            "requires_confirmation": true
          }
        ]
      },
      {
        "id": "phase.finalize",
        "name": "Finalize",
        "description": "最終階段。關閉 writeback transaction，清理 obligation ledger， 標記本輪完成。此階段後 session 可結束或回到 checkpoint 開始下一輪。 執行最終 output governance 檢查（語言一致性、去敏狀態）。\n",
        "entry_conditions": [
          {
            "readback_complete": true
          }
        ],
        "allowed_actions": [
          "close_obligation_ledger",
          "update_knowledge",
          "summarize_completion",
          "check_output_governance"
        ],
        "forbidden_actions": [
          "write_file",
          "apply_diff",
          "execute_command",
          "commit",
          "push"
        ],
        "blocking_gates": [
          "gate.finalize.all_obligations_closed",
          "gate.finalize.knowledge_updated",
          "gate.output.final_language_check",
          "gate.output.final_sanitization_check"
        ],
        "obligations": [
          "obligation.finalize.close_ledger",
          "obligation.finalize.run_knowledge_checkpoint",
          "obligation.finalize.verify_output_governance"
        ],
        "next_phase": "phase.checkpoint",
        "phase_transition_triggers": [
          {
            "condition": "finalize_complete == true",
            "auto_transition": false,
            "requires_confirmation": true
          },
          {
            "condition": "session_end == true",
            "auto_transition": false,
            "requires_confirmation": true,
            "target_phase": nil
          }
        ]
      },
      {
        "id": "phase.recovery",
        "name": "Recovery",
        "description": "恢復階段。當偵測到 state inconsistency、phase drift、或 obligation 遺失時， 進入此階段執行修復。此階段有最高的 action 權限。\n",
        "entry_conditions": [
          {
            "OR": [
              {
                "state_inconsistency_detected": true
              },
              {
                "phase_drift_detected": true
              },
              {
                "obligation_missing": true
              },
              {
                "generated_surface_stale": true
              }
            ]
          }
        ],
        "allowed_actions": [
          "read_file",
          "write_file",
          "apply_diff",
          "execute_command",
          "run_validator",
          "rebuild_obligation_ledger",
          "reconcile_phase"
        ],
        "forbidden_actions": [
          "commit",
          "push",
          "finalize"
        ],
        "blocking_gates": [
          "gate.recovery.repair_strategy_defined"
        ],
        "obligations": [
          "obligation.recovery.diagnose_root_cause",
          "obligation.recovery.apply_repair",
          "obligation.recovery.verify_repair"
        ],
        "next_phase": "phase.checkpoint",
        "phase_transition_triggers": [
          {
            "condition": "repair_complete == true AND verified == true",
            "auto_transition": false,
            "requires_confirmation": true
          }
        ]
      }
    ],
    "phase_transition_rules": [
      {
        "description": "禁止跳過 validation 直接 commit",
        "from": "phase.execution",
        "to": "phase.commit",
        "blocked": true,
        "reason": "Execution 完成後必須先經過 Validation，不能直接 commit"
      },
      {
        "description": "禁止跳過 commit 直接 push",
        "from": "phase.validation",
        "to": "phase.push",
        "blocked": true,
        "reason": "Validation 完成後必須先 commit，不能直接 push"
      },
      {
        "description": "禁止從 execution 直接 finalize",
        "from": "phase.execution",
        "to": "phase.finalize",
        "blocked": true,
        "reason": "Execution 完成後必須經過 validation → commit → push → readback 才能 finalize"
      },
      {
        "description": "Recovery 可從任何 phase 進入",
        "from": "*",
        "to": "phase.recovery",
        "blocked": false,
        "reason": "State inconsistency 可發生在任何 phase，recovery 必須能隨時介入"
      },
      {
        "description": "Recovery 完成後必須回到 checkpoint",
        "from": "phase.recovery",
        "to": "phase.checkpoint",
        "blocked": false,
        "reason": "修復完成後需要重新確認 state"
      }
    ],
    "metadata": {
      "priority": "P0",
      "confidence": "high",
      "compatibility_state": "new",
      "last_checked": "2026-05-15"
    }
  }

  # Source: runtime/pipeline/context-flow.yaml
  PIPELINE_CONTEXT_FLOW =   {
    "context_flow_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/pipeline",
    "description": "Progressive context expansion flow. 定義 context 從最輕量的 summary 開始， 只有在 summary 不足時才逐步擴展到 module summary → detailed source → raw source。 每個層級有明確的 entry condition、estimated tokens、exit condition 與 cache policy。\n",
    "expansion_levels": [
      {
        "id": "summary",
        "name": "Summary Level",
        "description": "最輕量的 context 層。只載入 summary（≤500 tokens）， 提供任務所需的高層次概覽。\n",
        "estimated_tokens": 500,
        "cache_policy": {
          "ttl": "session",
          "revalidate_on": "task_change"
        },
        "entry_conditions": [
          {
            "routing_complete": true
          },
          {
            "skill_identified": true
          }
        ],
        "sources": [
          {
            "type": "skill_summary",
            "path_template": "knowledge/summaries/{skill-id}.md",
            "description": "Skill summary（≤500 tokens）",
            "condition": "relevance_score >= 0.5"
          },
          {
            "type": "rule_summary",
            "path_template": "knowledge/summaries/{rule-id}.md",
            "description": "Rule summary（≤300 tokens）",
            "condition": "activation_rule_matched == true"
          }
        ],
        "exit_conditions": [
          {
            "summary_sufficient": true
          },
          {
            "OR summary_insufficient": true
          }
        ],
        "output": [
          {
            "summary_context": "Summary-level understanding of skill/rules"
          },
          {
            "next_level": "module_summary (if insufficient)"
          }
        ]
      },
      {
        "id": "module_summary",
        "name": "Module Summary Level",
        "description": "中量 context 層。載入 skill 的 module-level summary（如 README.md）， 提供模組結構與關鍵入口點。\n",
        "estimated_tokens": 1500,
        "cache_policy": {
          "ttl": "session",
          "revalidate_on": "task_change"
        },
        "entry_conditions": [
          {
            "summary_insufficient": true
          },
          {
            "module_summary_available": true
          }
        ],
        "sources": [
          {
            "type": "skill_readme",
            "path_template": "skills/{skill-id}/README.md",
            "description": "Skill README（~800 tokens）",
            "condition": "exists"
          },
          {
            "type": "module_readme",
            "path_template": "{module}/README.md",
            "description": "Module README（~500 tokens）",
            "condition": "cross_module_reference == true"
          }
        ],
        "exit_conditions": [
          {
            "module_summary_sufficient": true
          },
          {
            "OR detailed_source_needed": true
          }
        ],
        "output": [
          {
            "module_context": "Module-level understanding of structure"
          },
          {
            "next_level": "detailed_source (if insufficient)"
          }
        ]
      },
      {
        "id": "detailed_source",
        "name": "Detailed Source Level",
        "description": "重 context 層。載入完整 source file（如 SKILL.md、full rule file）， 提供執行所需的完整細節。\n",
        "estimated_tokens": 4500,
        "cache_policy": {
          "ttl": "task",
          "revalidate_on": "tool_call"
        },
        "entry_conditions": [
          {
            "module_summary_insufficient": true
          },
          {
            "detailed_source_available": true
          }
        ],
        "sources": [
          {
            "type": "skill_entrypoint",
            "path_template": "skills/{skill-id}/SKILL.md",
            "description": "Skill full source（~3500 tokens）",
            "condition": "summary_insufficient == true"
          },
          {
            "type": "rule_full",
            "path_template": "enforcement/{rule-id}.md",
            "description": "Rule full source（~2000 tokens）",
            "condition": "activation_rule_matched == true AND summary_insufficient == true"
          },
          {
            "type": "architecture_doc",
            "path_template": "architecture/{doc-id}.md",
            "description": "Architecture document（~3000 tokens）",
            "condition": "architecture_reference_needed == true"
          }
        ],
        "exit_conditions": [
          {
            "detailed_source_sufficient": true
          },
          {
            "OR raw_source_needed": true
          }
        ],
        "output": [
          {
            "detailed_context": "Full source-level understanding"
          },
          {
            "next_level": "raw_source (if insufficient)"
          }
        ]
      },
      {
        "id": "raw_source",
        "name": "Raw Source Level",
        "description": "最重 context 層。載入原始未壓縮的 source（如完整 feedback history、原始 log）， 僅在 debug 或 deep analysis 時使用。\n",
        "estimated_tokens": 10000,
        "cache_policy": {
          "ttl": "task",
          "revalidate_on": "tool_call",
          "max_age": 3
        },
        "entry_conditions": [
          {
            "detailed_source_insufficient": true
          },
          {
            "raw_source_available": true
          },
          {
            "raw_source_justified": true
          }
        ],
        "sources": [
          {
            "type": "feedback_history",
            "path_template": "skills/{skill-id}/feedback_history/{category}/{filename}",
            "description": "Full feedback history（~5000 tokens）",
            "condition": "debug_needed == true"
          },
          {
            "type": "raw_log",
            "path_template": "{log-path}",
            "description": "Raw log or output（~8000 tokens）",
            "condition": "deep_analysis_needed == true"
          },
          {
            "type": "full_technique",
            "path_template": "skills/{skill-id}/techniques/{technique}/README.md",
            "description": "Full technique documentation（~4000 tokens）",
            "condition": "technique_detail_needed == true"
          }
        ],
        "exit_conditions": [
          {
            "raw_source_sufficient": true
          },
          {
            "OR resource_exhausted": true
          }
        ],
        "output": [
          {
            "raw_context": "Raw source-level understanding"
          },
          {
            "fallback": "summary (if resource exhausted)"
          }
        ]
      }
    ],
    "rules": [
      {
        "id": "no-skip-levels",
        "description": "Context 必須依序擴展，不可跳過層級。 例如：不可直接從 summary 跳到 detailed source。\n",
        "severity": "error",
        "rationale": "跳過層級會失去 module summary 提供的結構理解， 可能導致載入不必要的完整 source。\n"
      },
      {
        "id": "cache-before-load",
        "description": "載入任何 source 前，先檢查 cache policy。 如果仍在 TTL 內且 task 未變更，使用 cached context。\n",
        "severity": "warning",
        "rationale": "避免重複載入相同 context，浪費 token。\n"
      },
      {
        "id": "justify-raw-source",
        "description": "進入 raw source level 前必須有明確理由（debug、deep analysis）。 不可因「summary 不足」就直接跳到 raw source。\n",
        "severity": "error",
        "rationale": "Raw source 是最高成本層級（~10000 tokens）， 必須有明確 justification 才能使用。\n"
      },
      {
        "id": "resource-exhaustion-fallback",
        "description": "如果 token budget 不足以載入下一層級， 回退到 summary level 並標記 context 為 degraded。\n",
        "severity": "warning",
        "rationale": "確保即使在資源受限的情況下， Agent 仍有 summary-level 的理解。\n"
      }
    ]
  }

  # Source: runtime/pipeline/guard-chain.yaml
  PIPELINE_GUARD_CHAIN =   {
    "guard_chain_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/pipeline",
    "description": "Guard execution order per pipeline stage. 定義每個 stage 中 guards 的執行順序、 檢查頻率與中斷行為。Guards 按 severity 與依賴關係排序， 確保高優先級 guards 先執行，低優先級 guards 在資源允許時才執行。\n",
    "stages": [
      {
        "id": "bootstrap",
        "name": "Bootstrap Stage",
        "description": "Session 啟動階段的 guard chain",
        "guards": [
          {
            "id": "token_budget",
            "order": 1,
            "severity": "critical",
            "description": "檢查 bootstrap token budget（2000 tokens）",
            "check": {
              "frequency": "once",
              "at": "after_init"
            },
            "action_on_violation": "halt_bootstrap",
            "source": "runtime/budget/token-budget.yaml"
          },
          {
            "id": "context_growth",
            "order": 2,
            "severity": "warning",
            "description": "檢查 bootstrap 是否載入過多 context",
            "check": {
              "frequency": "once",
              "at": "after_load",
              "threshold": "30%"
            },
            "action_on_violation": "warn + suggest_prune",
            "source": "runtime/guards/circuit-breaker.yaml"
          }
        ]
      },
      {
        "id": "routing",
        "name": "Routing Stage",
        "description": "任務路由階段的 guard chain",
        "guards": [
          {
            "id": "token_budget",
            "order": 1,
            "severity": "critical",
            "description": "檢查 routing token budget（2500 tokens）",
            "check": {
              "frequency": "once",
              "at": "after_routing"
            },
            "action_on_violation": "halt_routing + fallback_to_summary",
            "source": "runtime/budget/token-budget.yaml"
          },
          {
            "id": "context_growth",
            "order": 2,
            "severity": "warning",
            "description": "檢查 routing 是否載入過多 skill/rules",
            "check": {
              "frequency": "once",
              "at": "after_skill_load",
              "threshold": "30%"
            },
            "action_on_violation": "warn + limit_skill_load",
            "source": "runtime/guards/circuit-breaker.yaml"
          },
          {
            "id": "conflict_rules",
            "order": 3,
            "severity": "warning",
            "description": "檢查已載入的 rules 是否有衝突",
            "check": {
              "frequency": "once",
              "at": "after_activation"
            },
            "action_on_violation": "warn + suggest_resolution",
            "source": "runtime/guards/circuit-breaker.yaml"
          }
        ]
      },
      {
        "id": "execution",
        "name": "Execution Stage",
        "description": "任務執行階段的 guard chain（最複雜，最多 guards）",
        "guards": [
          {
            "id": "token_budget",
            "order": 1,
            "severity": "critical",
            "description": "檢查 execution token budget（100000 tokens）",
            "check": {
              "frequency": "per_tool_output",
              "thresholds": {
                "warning_at": 70000,
                "hard_stop_at": 90000
              }
            },
            "action_on_violation": {
              "warning": "warn + suggest_compression",
              "hard_stop": "halt_execution + force_close_loop"
            },
            "source": "runtime/budget/token-budget.yaml"
          },
          {
            "id": "recursive_depth",
            "order": 2,
            "severity": "critical",
            "description": "檢查遞迴深度是否超過限制",
            "check": {
              "frequency": "per_tool_call",
              "threshold": 4
            },
            "action_on_violation": "halt_execution + force_close_loop",
            "source": "runtime/guards/circuit-breaker.yaml"
          },
          {
            "id": "tool_calls",
            "order": 3,
            "severity": "critical",
            "description": "檢查工具呼叫次數是否超過限制",
            "check": {
              "frequency": "per_tool_call",
              "thresholds": {
                "warn_at": 15,
                "hard_stop_at": 20
              }
            },
            "action_on_violation": {
              "warning": "warn + suggest_decomposition",
              "hard_stop": "halt_execution + suggest_new_session"
            },
            "source": "runtime/guards/circuit-breaker.yaml"
          },
          {
            "id": "context_growth",
            "order": 4,
            "severity": "high",
            "description": "檢查 context 增長是否超過限制",
            "check": {
              "frequency": "per_task",
              "thresholds": {
                "warn_at": "30%",
                "hard_stop_at": "80%"
              }
            },
            "action_on_violation": {
              "warning": "warn + suggest_prune",
              "hard_stop": "force_prune"
            },
            "source": "runtime/guards/circuit-breaker.yaml"
          },
          {
            "id": "hallucination_risk",
            "order": 5,
            "severity": "high",
            "description": "檢查 hallucination risk factors",
            "check": {
              "frequency": "per_task",
              "factors": [
                "no_canonical_source",
                "speculative_reasoning",
                "conflicting_information",
                "missing_dependency"
              ],
              "threshold": 0.7
            },
            "action_on_violation": "halt + suggest_source_read",
            "source": "runtime/guards/circuit-breaker.yaml"
          },
          {
            "id": "conversation_length",
            "order": 6,
            "severity": "high",
            "description": "檢查對話長度是否超過限制",
            "check": {
              "frequency": "per_task",
              "threshold": 50
            },
            "action_on_violation": "warn + suggest_new_session",
            "source": "runtime/guards/context-pollution.yaml"
          },
          {
            "id": "conflict_rules",
            "order": 7,
            "severity": "medium",
            "description": "檢查執行期間是否有規則衝突",
            "check": {
              "frequency": "per_task"
            },
            "action_on_violation": "warn + suggest_resolution",
            "source": "runtime/guards/circuit-breaker.yaml"
          },
          {
            "id": "repetitive_edits",
            "order": 8,
            "severity": "medium",
            "description": "檢查是否對同一檔案反覆修改",
            "check": {
              "frequency": "per_edit",
              "threshold": 5
            },
            "action_on_violation": "warn + suggest_review",
            "source": "runtime/guards/context-pollution.yaml"
          },
          {
            "id": "module_count",
            "order": 9,
            "severity": "medium",
            "description": "檢查載入的模組數量",
            "check": {
              "frequency": "per_task",
              "threshold": 20
            },
            "action_on_violation": "warn + suggest_prune",
            "source": "runtime/guards/context-pollution.yaml"
          },
          {
            "id": "cross_reference_depth",
            "order": 10,
            "severity": "medium",
            "description": "檢查依賴鏈深度",
            "check": {
              "frequency": "per_task",
              "threshold": 5
            },
            "action_on_violation": "warn + suggest_simplify",
            "source": "runtime/guards/context-pollution.yaml"
          },
          {
            "id": "token_utilization",
            "order": 11,
            "severity": "medium",
            "description": "檢查 token 使用率",
            "check": {
              "frequency": "per_task",
              "threshold": "85%"
            },
            "action_on_violation": "warn + suggest_archive",
            "source": "runtime/guards/context-pollution.yaml"
          }
        ]
      },
      {
        "id": "close-loop",
        "name": "Close-loop Stage",
        "description": "Session 閉環階段的 guard chain",
        "guards": [
          {
            "id": "token_budget",
            "order": 1,
            "severity": "critical",
            "description": "檢查 close-loop token budget（1000 tokens）",
            "check": {
              "frequency": "once",
              "at": "before_write"
            },
            "action_on_violation": "force_minimal_write",
            "source": "runtime/budget/token-budget.yaml"
          }
        ]
      }
    ],
    "rules": [
      {
        "id": "ordered-execution",
        "description": "Guards 必須按 order 順序執行。高 order 的 guard 必須在低 order 的 guard 執行完成後才能執行。\n",
        "severity": "error",
        "rationale": "Token budget 必須在 recursive depth 之前檢查， 因為 recursive depth 檢查需要 token 來執行。\n"
      },
      {
        "id": "layered-violation",
        "description": "如果 critical guard 違規，直接中斷 stage，不繼續執行 lower-order guards。 如果 high/medium guard 違規，記錄警告但繼續執行。\n",
        "severity": "warning",
        "rationale": "Critical guards（token_budget、recursive_depth、tool_calls） 直接影響系統穩定性，違規時必須立即中斷。\n"
      },
      {
        "id": "frequency-respect",
        "description": "Guards 必須按定義的 frequency 執行。 per_tool_call 的 guard 不可只在 per_task 時檢查。\n",
        "severity": "error",
        "rationale": "降低檢查頻率可能導致 guard 在違規發生後才被觸發。\n"
      },
      {
        "id": "cross-stage-escalation",
        "description": "如果某個 guard 在 routing stage 已違規（如 conflict_rules）， 在 execution stage 也必須重新檢查，不可假設 routing stage 已處理。\n",
        "severity": "warning",
        "rationale": "Execution stage 可能載入新的 rules，導致新的衝突。\n"
      }
    ]
  }

  # Source: runtime/pipeline/relevance-engine.yaml
  PIPELINE_RELEVANCE_ENGINE =   {
    "relevance_engine_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/pipeline",
    "description": "Skill Relevance Engine. 在 routing stage 執行，對 routing-registry.yaml 中的每個 route 計算 relevance score（0.0-1.0）。Score ≥ 0.5 的 skill 才會被載入 summary。 Score < 0.5 的 skill 被跳過，避免浪費 token 在不相關的 skill 上。\n",
    "scoring": {
      "method": "weighted_combination",
      "description": "Relevance score 由三個子分數加權組合而成： - trigger_match_score（0.5）：基於 user intent 與 skill triggers 的匹配度 - domain_match_score（0.3）：基於 task domain 與 skill domains 的匹配度 - weight_score（0.2）：基於 route 的基礎 weight（來自 routing-registry.yaml）\n最終分數 = trigger_match × 0.5 + domain_match × 0.3 + weight × 0.2\n如果 skill 與已選 skill 有 conflicts，最終分數 × 0.5（penalty）。\n",
      "threshold": 0.5,
      "components": [
        {
          "id": "trigger_match_score",
          "weight": 0.5,
          "description": "計算 user intent 中的關鍵詞與 skill triggers 的匹配度。 匹配度 = matched_triggers / total_triggers（但上限 1.0）。 如果 user intent 包含 trigger 的精確匹配，該 trigger 計為 2 分。\n",
          "scoring_logic": "matched_count = count(intent_keywords ∩ skill.triggers) exact_match_count = count(exact_match(intent_keywords, skill.triggers)) raw_score = (matched_count + exact_match_count) / len(skill.triggers) return min(raw_score, 1.0)\n"
        },
        {
          "id": "domain_match_score",
          "weight": 0.3,
          "description": "計算 task domain 與 skill domains 的匹配度。 如果 task domain 與 skill domain 完全匹配，得分 1.0。 如果 task domain 是 skill domain 的子領域，得分 0.7。 如果 task domain 與 skill domain 無關，得分 0.0。\n",
          "scoring_logic": "if task_domain in skill.domains:\n  return 1.0\nelif is_subdomain(task_domain, skill.domains):\n  return 0.7\nelse:\n  return 0.0\n"
        },
        {
          "id": "weight_score",
          "weight": 0.2,
          "description": "直接使用 routing-registry.yaml 中的 weight 值（0.0-1.0）。 這是 skill 的基礎重要性，不受 task 影響。\n",
          "scoring_logic": "return skill.weight\n"
        }
      ],
      "penalties": [
        {
          "id": "conflict_penalty",
          "description": "如果 skill 與已選 skill 有 conflicts，最終分數 × 0.5。 這確保衝突的 skill 不會同時被載入。\n",
          "condition": "has_conflict_with_selected_skill(skill)",
          "multiplier": 0.5
        },
        {
          "id": "dependency_missing_penalty",
          "description": "如果 skill 的 dependencies 不存在或無法載入，最終分數 × 0.8。 這確保依賴不完整的 skill 不會被優先載入。\n",
          "condition": "dependency_check_failed(skill)",
          "multiplier": 0.8
        }
      ]
    },
    "execution": {
      "steps": [
        {
          "id": "parse-intent",
          "action": "analyze",
          "input": "user_message",
          "output": {
            "keywords": [
              "extracted keywords"
            ],
            "domain": "identified domain"
          },
          "estimated_tokens": 200
        },
        {
          "id": "score-all-skills",
          "action": "compute",
          "input": [
            "routing-registry.yaml",
            "parsed_intent"
          ],
          "output": {
            "scored_skills": [
              {
                "skill_id": "skill-id",
                "trigger_match_score": "0.0-1.0",
                "domain_match_score": "0.0-1.0",
                "weight_score": "0.0-1.0",
                "raw_score": "0.0-1.0",
                "penalties": [
                  "conflict",
                  "dependency_missing"
                ],
                "final_score": "0.0-1.0"
              }
            ]
          },
          "estimated_tokens": 200
        },
        {
          "id": "select-skills",
          "action": "filter",
          "input": "scored_skills",
          "condition": "final_score >= 0.5",
          "output": {
            "selected_skills": [
              {
                "skill_id": "skill-id",
                "final_score": "0.0-1.0",
                "summary_path": "knowledge/summaries/{skill-id}.md"
              }
            ]
          },
          "estimated_tokens": 100
        },
        {
          "id": "skip-irrelevant",
          "action": "skip",
          "input": "scored_skills",
          "condition": "final_score < 0.5",
          "output": {
            "skipped_skills": [
              {
                "skill_id": "skill-id",
                "final_score": "0.0-1.0",
                "reason": "below_threshold"
              }
            ]
          },
          "estimated_tokens": 50
        }
      ]
    },
    "examples": [
      {
        "scenario": "APK traffic analysis task",
        "user_intent": "分析這個 APK 的網路流量，看看有沒有可疑的連線",
        "parsed": {
          "keywords": [
            "apk",
            "traffic",
            "network",
            "analysis"
          ],
          "domain": "mobile-security"
        },
        "results": [
          {
            "skill": "apk-analysis",
            "trigger_match": 0.8,
            "domain_match": 1.0,
            "weight": 0.9,
            "raw_score": "0.8 × 0.5 + 1.0 × 0.3 + 0.9 × 0.2 = 0.88",
            "penalties": [
  
            ],
            "final_score": 0.88
          },
        ]
      },
      {
        "scenario": "iOS API design review",
        "user_intent": "幫我 review 這個 iOS app 的 API 設計",
        "parsed": {
          "keywords": [
            "ios",
            "api",
            "design",
            "review"
          ],
          "domain": "mobile-development"
        },
        "results": [
          {
            "skill": "apk-analysis",
            "trigger_match": 0.0,
            "domain_match": 0.0,
            "weight": 0.9,
            "raw_score": "0.0 × 0.5 + 0.0 × 0.3 + 0.9 × 0.2 = 0.18",
            "penalties": [

            ],
            "final_score": 0.18
          }
        ]
      },
      {
        "scenario": "Conflicting skills scenario",
        "user_intent": "分析這個 APK 的架構，同時給開發建議",
        "parsed": {
          "keywords": [
            "apk",
            "architecture",
            "development",
            "analysis"
          ],
          "domain": "mobile-security"
        },
        "results": [
          {
            "skill": "apk-analysis",
            "trigger_match": 0.6,
            "domain_match": 1.0,
            "weight": 0.9,
            "raw_score": "0.6 × 0.5 + 1.0 × 0.3 + 0.9 × 0.2 = 0.78",
            "penalties": [
  
            ],
            "final_score": 0.78
          },
        ]
      }
    ],
    "rules": [
      {
        "id": "threshold-gate",
        "description": "Final score < 0.5 的 skill 必須被跳過，不可載入其 summary 或 source。 即使 Agent 認為該 skill 可能相關，也必須遵守 threshold。\n",
        "severity": "error",
        "rationale": "Threshold 0.5 確保只有「明顯相關」的 skill 被載入， 避免「可能相關」的 skill 浪費 token。\n"
      },
      {
        "id": "conflict-penalty-mandatory",
        "description": "如果 skill 與已選 skill 有 conflicts，conflict_penalty（×0.5） 必須被應用。不可因 Agent 判斷而跳過 penalty。\n",
        "severity": "error",
        "rationale": "Conflicts 定義在 routing-registry.yaml 中，代表已知的目標衝突。 忽略 conflict penalty 可能導致 context pollution。\n"
      },
      {
        "id": "score-transparency",
        "description": "Agent 在決定載入或跳過 skill 時，必須記錄每個子分數與最終分數。 這使得 scoring 決策可追溯、可調試。\n",
        "severity": "warning",
        "rationale": "透明的 scoring 記錄有助於調試 relevance engine 的行為， 並在必要時調整 weight 或 threshold。\n"
      },
      {
        "id": "no-manual-override",
        "description": "Agent 不可手動覆蓋 relevance score 或強制載入 score < 0.5 的 skill。 如果 Agent 認為某個 skill 應該被載入但 score < 0.5， 應記錄 feedback 供後續調整 weight 或 triggers。\n",
        "severity": "error",
        "rationale": "手動覆蓋會破壞 relevance engine 的可靠性， 導致 token 浪費在不相關的 skill 上。\n"
      }
    ]
  }

  # Source: runtime/pipeline/session-lifecycle.yaml
  PIPELINE_SESSION_LIFECYCLE =   {
    "session_lifecycle_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/pipeline",
    "description": "Session lifecycle stages for AI runtime. 每個 session 依序經過四個階段， 每個階段有明確的 token budget allocation、guard 檢查點與產出物。\n",
    "stages": [
      {
        "id": "bootstrap",
        "name": "Bootstrap Stage",
        "description": "Session 啟動階段。載入最小必要知識，建立 session context。 目標：在 ~800 tokens 內完成啟動，不載入任何非必要知識。\n",
        "entry_conditions": [
          {
            "session_start": true
          }
        ],
        "execution": {
          "steps": [
            {
              "id": "load-core-bootstrap",
              "action": "read_file",
              "target": "CORE_BOOTSTRAP.md",
              "description": "載入 3 條核心規則（rule-weight, dependency-reading, goal-ledger）",
              "estimated_tokens": 800
            },
            {
              "id": "load-layout",
              "action": "read_file",
              "target": "README.md",
              "description": "載入 OS layout 與 quickstart",
              "estimated_tokens": 300
            },
            {
              "id": "init-token-budget",
              "action": "check_budget",
              "target": "runtime/budget/token-budget.yaml",
              "description": "初始化 session token budget"
            },
            {
              "id": "init-memory",
              "action": "init_memory",
              "target": "memory/working/",
              "description": "初始化 working memory"
            }
          ]
        },
        "token_budget": {
          "allocation": 2000,
          "warning_at": 1400,
          "hard_stop_at": 1800
        },
        "guards_active": [
          "token_budget",
          "context_growth"
        ],
        "exit_conditions": [
          {
            "core_bootstrap_loaded": true
          },
          {
            "layout_loaded": true
          },
          {
            "token_budget_initialized": true
          }
        ],
        "output": [
          {
            "session_context": "Core rules + OS layout loaded"
          },
          {
            "working_memory": "Initialized empty"
          }
        ]
      },
      {
        "id": "routing",
        "name": "Routing Stage",
        "description": "任務路由階段。根據 user intent 找到對應 skill，決定哪些知識需要載入。 目標：在 ~1500 tokens 內完成路由，只載入 summary 層。\n",
        "entry_conditions": [
          {
            "bootstrap_complete": true
          },
          {
            "user_intent_received": true
          }
        ],
        "execution": {
          "steps": [
            {
              "id": "parse-intent",
              "action": "analyze",
              "target": "user_message",
              "description": "解析 user intent，提取關鍵詞與領域",
              "estimated_tokens": 200
            },
            {
              "id": "query-skill-index",
              "action": "read_file",
              "target": "routing-registry.yaml",
              "description": "查詢 routing-registry.yaml 找到對應 route",
              "estimated_tokens": 500
            },
            {
              "id": "run-relevance-engine",
              "action": "score",
              "target": "routing-registry.yaml",
              "description": "使用 relevance engine 計算每個 route 的相關性分數",
              "estimated_tokens": 200
            },
            {
              "id": "load-skill-summary",
              "action": "read_file",
              "target": "knowledge/summaries/{skill-id}.md",
              "description": "載入對應 skill 的 summary（300-500 tokens）",
              "estimated_tokens": 500,
              "condition": "relevance_score >= 0.5"
            },
            {
              "id": "check-activation-rules",
              "action": "read_file",
              "target": "runtime/router/activation-rules.yaml",
              "description": "檢查哪些 lazy-load rules 需要 activate",
              "estimated_tokens": 300
            },
            {
              "id": "load-activated-rules",
              "action": "read_file",
              "target": "enforcement/{rule-id}.md",
              "description": "載入 activated lazy-load rules",
              "estimated_tokens": 800,
              "condition": "activation_rule_matched == true"
            }
          ]
        },
        "token_budget": {
          "allocation": 2500,
          "warning_at": 1750,
          "hard_stop_at": 2250
        },
        "guards_active": [
          "token_budget",
          "context_growth",
          "conflict_rules"
        ],
        "exit_conditions": [
          {
            "skill_identified": true
          },
          {
            "summary_loaded": true
          },
          {
            "activation_rules_checked": true
          }
        ],
        "output": [
          {
            "selected_skill": "skill-id"
          },
          {
            "relevance_score": "0.0-1.0"
          },
          {
            "activated_rules": [
              "rule-id-1",
              "rule-id-2"
            ]
          },
          {
            "summary_context": "Skill summary loaded"
          }
        ]
      },
      {
        "id": "execution",
        "name": "Execution Stage",
        "description": "任務執行階段。執行實際工作，使用工具、載入完整 source、產生輸出。 目標：在 token budget 內完成任務，動態調整 compression level。\n",
        "entry_conditions": [
          {
            "routing_complete": true
          },
          {
            "skill_selected": true
          }
        ],
        "execution": {
          "steps": [
            {
              "id": "load-full-source",
              "action": "read_file",
              "target": "{skill-entrypoint}",
              "description": "需要時才載入完整 source",
              "estimated_tokens": 4500,
              "condition": "summary_insufficient == true"
            },
            {
              "id": "execute-tool-calls",
              "action": "tool_loop",
              "description": "執行工具呼叫，每個工具呼叫前檢查 circuit breaker",
              "max_calls": 20
            },
            {
              "id": "check-health",
              "action": "evaluate",
              "target": "context_health",
              "description": "每個 task 完成後檢查 context health score",
              "frequency": "per_task"
            },
            {
              "id": "adjust-compression",
              "action": "adjust",
              "target": "compression_level",
              "description": "根據 context health 調整 compression level",
              "condition": "health_score < 0.75"
            },
            {
              "id": "check-budget",
              "action": "check_budget",
              "target": "token_budget",
              "description": "每個工具輸出後檢查 token budget",
              "frequency": "per_tool_output"
            }
          ]
        },
        "token_budget": {
          "allocation": 100000,
          "warning_at": 70000,
          "hard_stop_at": 90000
        },
        "guards_active": [
          "token_budget",
          "recursive_depth",
          "tool_calls",
          "context_growth",
          "hallucination_risk",
          "conflict_rules",
          "conversation_length",
          "repetitive_edits",
          "module_count",
          "cross_reference_depth",
          "token_utilization"
        ],
        "exit_conditions": [
          {
            "task_complete": true
          },
          {
            "OR budget_hard_stop_reached": true
          },
          {
            "OR pollution_critical": true
          }
        ],
        "output": [
          {
            "task_result": "完成的工作產出"
          },
          {
            "health_score": "0.0-1.0"
          },
          {
            "token_used": "number"
          },
          {
            "compression_level": "raw|summary|structured|minimal"
          }
        ]
      },
      {
        "id": "close-loop",
        "name": "Close-loop Stage",
        "description": "任務閉環階段。寫入 memory、記錄 decision、清理 context。 目標：在 ~1000 tokens 內完成閉環，確保知識不遺失。\n",
        "entry_conditions": [
          {
            "execution_complete": true
          },
          {
            "OR budget_hard_stop_reached": true
          },
          {
            "OR pollution_critical": true
          }
        ],
        "execution": {
          "steps": [
            {
              "id": "write-working-memory",
              "action": "write",
              "target": "memory/working/",
              "description": "寫入 session working memory",
              "estimated_tokens": 300
            },
            {
              "id": "write-session-summary",
              "action": "write",
              "target": "memory/summary/",
              "description": "寫入壓縮的 session summary（≤500 tokens）",
              "estimated_tokens": 500,
              "condition": "task_complete == true"
            },
            {
              "id": "record-decision",
              "action": "write",
              "target": "decisions/",
              "description": "記錄架構決策（如有）",
              "estimated_tokens": 200,
              "condition": "architecture_decision_made == true"
            },
            {
              "id": "auto-archive",
              "action": "archive",
              "target": "memory/working/session-archive-{timestamp}.md",
              "description": "自動歸檔 session（pollution critical 時強制）",
              "condition": "pollution_critical == true"
            },
            {
              "id": "suggest-new-session",
              "action": "suggest",
              "description": "建議新 session（budget hard stop 或 pollution critical 時）",
              "condition": "budget_hard_stop_reached == true OR pollution_critical == true"
            }
          ]
        },
        "token_budget": {
          "allocation": 1000,
          "warning_at": 700,
          "hard_stop_at": 900
        },
        "guards_active": [
          "token_budget"
        ],
        "exit_conditions": [
          {
            "memory_written": true
          },
          {
            "summary_written": true
          },
          {
            "OR archived": true
          }
        ],
        "output": [
          {
            "session_summary": "Compressed session history"
          },
          {
            "decisions": [
              "ADR notes if any"
            ]
          },
          {
            "archive_path": "path if archived"
          },
          {
            "next_action": "continue | new_session"
          }
        ]
      }
    ]
  }

  # Source: runtime/prompt-artifacts/artifact-templates.yaml
  PROMPT_ARTIFACTS_ARTIFACT_TEMPLATES =   {
    "prompt_artifact_templates_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/prompt-artifacts",
    "description": "Task-specific prompt artifact templates. 每個 template 定義一個完整的 prompt 結構， 包含 sections、composition rules、與 token budget。Agent 在 routing stage 完成後， 根據 task type 載入對應模板，依序組裝各 section 成為 execution stage 的 prompt。\n",
    "templates": [
      {
        "task_type": "apk-analysis",
        "name": "APK Analysis Prompt Artifact",
        "description": "授權 APK 流量分析的完整 prompt 結構。包含分析範圍、執行流程、 工程智慧（路由選擇、證據優先）、分析方法、產出閘門與輸出格式。\n",
        "skill_id": "apk-analysis",
        "triggers": [
          "apk",
          "mitm",
          "flutter",
          "frida",
          "traffic-analysis",
          "reverse-engineering"
        ],
        "estimated_total_tokens": 4500,
        "token_budget": {
          "warning_at": 3500,
          "hard_stop_at": 4000
        },
        "sections": [
          {
            "id": "task-context",
            "description": "任務上下文。包含授權確認、scope 定義、APK 版本、裝置/模擬器資訊、 允許的操作範圍。\n",
            "source": "inline",
            "content": "## Task Context\n- Authorization: [confirmed / pending]\n- Scope: [APK traffic analysis / runtime behavior / specific feature]\n- APK Version: [version]\n- Device/Emulator: [device info]\n- Allowed Actions: [capture / hook / decrypt / document]\n",
            "estimated_tokens": 200,
            "required": true
          },
          {
            "id": "workflow-steps",
            "description": "APK 分析執行流程。引用 workflow/apk-analysis/ 中的執行步驟。\n",
            "source": "workflow/apk-analysis/execution-flow.md",
            "composition_rule": "workflow-steps",
            "estimated_tokens": 1200,
            "required": true
          },
          {
            "id": "intelligence-atoms",
            "description": "APK 分析相關的工程智慧 atoms。根據 task 類型選擇性載入。\n",
            "source": "intelligence/engineering/analytical-reasoning/",
            "composition_rule": "intelligence-atoms",
            "estimated_tokens": 800,
            "required": false,
            "load_strategy": "on_condition",
            "load_condition": "task_requires_judgment == true"
          },
          {
            "id": "analysis-methods",
            "description": "APK 分析的具體方法。引用 analysis/apk/ 中的技術內容。\n",
            "source": "analysis/apk/",
            "composition_rule": "analysis-methods",
            "estimated_tokens": 1000,
            "required": false,
            "load_strategy": "on_condition",
            "load_condition": "technique_detail_needed == true"
          },
          {
            "id": "artifact-gates",
            "description": "產出閘門定義。引用 workflow/apk-analysis/artifact-gates.md。\n",
            "source": "workflow/apk-analysis/artifact-gates.md",
            "estimated_tokens": 800,
            "required": true
          },
          {
            "id": "output-format",
            "description": "輸出格式要求。定義最終產出的結構與品質標準。\n",
            "source": "inline",
            "content": "## Output Format\n- UI Architecture Map (page-level)\n- API Catalog (endpoint-level)\n- Domain/Runtime Baseline\n- Feature Reconstruction Handoff (if applicable)\n- Unknowns & Open Questions\n- Sanitization Checklist\n",
            "estimated_tokens": 300,
            "required": true
          }
        ]
      },
      {
        "task_type": "repo-analysis",
        "name": "Repository Analysis Prompt Artifact",
        "description": "Repository 分析的完整 prompt 結構。包含分析類型、方法、工程智慧與輸出格式。\n",
        "skill_id": "repo-analysis",
        "triggers": [
          "repo-analysis",
          "codebase-analysis",
          "onboarding",
          "tech-debt",
          "migration-impact"
        ],
        "estimated_total_tokens": 3500,
        "token_budget": {
          "warning_at": 2800,
          "hard_stop_at": 3200
        },
        "sections": [
          {
            "id": "task-context",
            "description": "分析任務上下文。包含分析類型、目標 repository、分析目標。\n",
            "source": "inline",
            "content": "## Task Context\n- Analysis Type: [onboarding / deep-analysis / migration-impact / tech-debt / security-audit]\n- Target Repository: [name / URL]\n- Analysis Goal: [specific questions to answer]\n- Constraints: [time / depth / focus areas]\n",
            "estimated_tokens": 200,
            "required": true
          },
          {
            "id": "analysis-workflow",
            "description": "Repository 分析執行流程。引用 workflow/repo-analysis/。\n",
            "source": "workflow/repo-analysis/",
            "composition_rule": "workflow-steps",
            "estimated_tokens": 1000,
            "required": true
          },
          {
            "id": "analysis-methods",
            "description": "Repository 分析的具體方法。引用 analysis/repo/。\n",
            "source": "analysis/repo/",
            "composition_rule": "analysis-methods",
            "estimated_tokens": 800,
            "required": false,
            "load_strategy": "on_condition",
            "load_condition": "deep_analysis == true"
          },
          {
            "id": "intelligence-atoms",
            "description": "架構與領域相關的工程智慧 atoms。\n",
            "source": "intelligence/engineering/architecture/",
            "composition_rule": "intelligence-atoms",
            "estimated_tokens": 600,
            "required": false,
            "load_strategy": "on_condition",
            "load_condition": "architecture_assessment_needed == true"
          },
          {
            "id": "output-format",
            "description": "分析結果輸出格式。\n",
            "source": "inline",
            "content": "## Output Format\n- Repository Mental Model: [high-level structure summary]\n- Key Findings: [list of findings with evidence references]\n- Risk Assessment: [risks identified]\n- Recommendations: [actionable next steps]\n- Open Questions: [unresolved items]\n",
            "estimated_tokens": 300,
            "required": true
          }
        ]
      },
      {
        "task_type": "travel-planning",
        "name": "Travel Planning Prompt Artifact",
        "description": "旅遊規劃的完整 prompt 結構。包含規劃範圍、執行流程、領域智慧與輸出格式。\n",
        "skill_id": "travel-planning",
        "triggers": [
          "travel",
          "itinerary",
          "route",
          "transportation",
          "accommodation",
          "sightseeing",
          "trip",
          "vacation"
        ],
        "estimated_total_tokens": 3500,
        "token_budget": {
          "warning_at": 2800,
          "hard_stop_at": 3200
        },
        "sections": [
          {
            "id": "task-context",
            "description": "規劃任務上下文。包含目的地、日期、人數、預算與旅遊風格。\n",
            "source": "inline",
            "content": "## Task Context\n- Destination: [city / region / multi-city]\n- Dates: [start / end / duration]\n- Travelers: [adults / children / seniors]\n- Budget: [total / daily]\n- Style: [relaxed / packed / adventure / culture / food]\n- Special Requirements: [dietary / mobility / interests]\n",
            "estimated_tokens": 200,
            "required": true
          },
          {
            "id": "planning-workflow",
            "description": "旅遊規劃執行流程。引用 workflow/travel-planning/。\n",
            "source": "workflow/travel-planning/",
            "composition_rule": "workflow-steps",
            "estimated_tokens": 1200,
            "required": true
          },
          {
            "id": "intelligence-atoms",
            "description": "旅遊領域智慧 atoms。根據目的地與規劃類型選擇性載入。\n",
            "source": "intelligence/travel/",
            "composition_rule": "intelligence-atoms",
            "estimated_tokens": 500,
            "required": false,
            "load_strategy": "on_condition",
            "load_condition": "domain_specific_heuristics_available == true"
          },
          {
            "id": "output-format",
            "description": "規劃結果輸出格式。\n",
            "source": "inline",
            "content": "## Output Format\n- Daily Itinerary: [day-by-day plan with time slots]\n- Transportation: [routes / costs / alternatives]\n- Accommodation: [options with pros/cons]\n- Budget Breakdown: [estimated costs by category]\n- Backup Plans: [alternatives for key items]\n- Preparation Checklist: [packing / documents / bookings]\n",
            "estimated_tokens": 400,
            "required": true
          }
        ]
      },
      {
        "task_type": "repo-governance",
        "name": "Repository Governance Prompt Artifact",
        "description": "AI-native Knowledge OS 本身維護與治理的完整 prompt 結構。 包含治理範圍、生命週期、驗證閘門與輸出格式。\n",
        "skill_id": "repo-governance",
        "triggers": [
          "governance",
          "migration",
          "architecture-upgrade",
          "repo-maintenance",
          "cleanup",
          "deprecation",
          "lifecycle",
          "validation",
          "metadata",
          "schema"
        ],
        "estimated_total_tokens": 3000,
        "token_budget": {
          "warning_at": 2400,
          "hard_stop_at": 2700
        },
        "sections": [
          {
            "id": "task-context",
            "description": "治理任務上下文。包含治理類型、範圍與目標。\n",
            "source": "inline",
            "content": "## Task Context\n- Governance Type: [migration / upgrade / cleanup / validation / deprecation]\n- Scope: [layer / directory / file pattern]\n- Goal: [specific outcome]\n- Constraints: [backward-compatibility / timeline / dependencies]\n",
            "estimated_tokens": 200,
            "required": true
          },
          {
            "id": "governance-workflow",
            "description": "治理執行流程。引用 governance/ 中的生命週期與驗證規則。\n",
            "source": "governance/",
            "composition_rule": "workflow-steps",
            "estimated_tokens": 1000,
            "required": true
          },
          {
            "id": "lifecycle-rules",
            "description": "知識生命週期規則。引用 governance/lifecycle/。\n",
            "source": "governance/lifecycle/",
            "estimated_tokens": 600,
            "required": true
          },
          {
            "id": "validation-gates",
            "description": "驗證閘門定義。引用 governance/validation/。\n",
            "source": "governance/validation/",
            "estimated_tokens": 600,
            "required": true
          },
          {
            "id": "output-format",
            "description": "治理結果輸出格式。\n",
            "source": "inline",
            "content": "## Output Format\n- Changes Made: [list of changes with file paths]\n- Validation Results: [pass / fail per gate]\n- Migration Status: [completed / pending / blocked]\n- Linked Updates Required: [files that need sync]\n- Open Issues: [unresolved items]\n",
            "estimated_tokens": 300,
            "required": true
          }
        ]
      },
      {
        "task_type": "knowledge-navigation",
        "name": "Knowledge Navigation Prompt Artifact",
        "description": "知識導航與查詢的完整 prompt 結構。包含查詢類型、索引、摘要與輸出格式。\n",
        "skill_id": "knowledge-navigation",
        "triggers": [
          "navigation",
          "routing",
          "index",
          "summary",
          "knowledge-graph",
          "context-loading",
          "runtime"
        ],
        "estimated_total_tokens": 2500,
        "token_budget": {
          "warning_at": 2000,
          "hard_stop_at": 2250
        },
        "sections": [
          {
            "id": "task-context",
            "description": "導航任務上下文。包含查詢類型與目標。\n",
            "source": "inline",
            "content": "## Task Context\n- Query Type: [find-atom / find-skill / find-summary / trace-relationship]\n- Query Target: [keyword / topic / file pattern]\n- Depth: [summary / full / graph]\n",
            "estimated_tokens": 150,
            "required": true
          },
          {
            "id": "navigation-index",
            "description": "知識導航索引。引用 knowledge/indexes/。\n",
            "source": "knowledge/indexes/",
            "estimated_tokens": 800,
            "required": true
          },
          {
            "id": "knowledge-summary",
            "description": "知識摘要。引用 knowledge/summaries/。\n",
            "source": "knowledge/summaries/",
            "composition_rule": "knowledge-summary",
            "estimated_tokens": 600,
            "required": false,
            "load_strategy": "on_condition",
            "load_condition": "summary_needed == true"
          },
          {
            "id": "output-format",
            "description": "導航結果輸出格式。\n",
            "source": "inline",
            "content": "## Output Format\n- Found Atoms: [list with paths and status]\n- Relevance Score: [per atom]\n- Relationships: [graph edges if applicable]\n- Recommended Load Order: [priority-ordered list]\n",
            "estimated_tokens": 300,
            "required": true
          }
        ]
      },
      {
        "task_type": "feedback-promotion",
        "name": "Feedback & Promotion Prompt Artifact",
        "description": "Feedback lesson 管理與 promotion 的完整 prompt 結構。 包含 lesson 內容、promotion 目標、驗證規則與輸出格式。\n",
        "skill_id": "feedback-promotion",
        "triggers": [
          "feedback",
          "lesson",
          "promotion",
          "failure-learning",
          "extraction",
          "refinement"
        ],
        "estimated_total_tokens": 3000,
        "token_budget": {
          "warning_at": 2400,
          "hard_stop_at": 2700
        },
        "sections": [
          {
            "id": "task-context",
            "description": "Feedback 任務上下文。包含操作類型與 lesson 範圍。\n",
            "source": "inline",
            "content": "## Task Context\n- Operation: [create-lesson / promote / extract / refine / replay]\n- Lesson Source: [feedback_history / failure-patterns / conversation]\n- Target Layer: [intelligence / workflow / analysis / enforcement]\n",
            "estimated_tokens": 200,
            "required": true
          },
          {
            "id": "promotion-workflow",
            "description": "Promotion 執行流程。引用 feedback/promotion/ 與 feedback/pipeline/。\n",
            "source": "feedback/promotion/",
            "composition_rule": "workflow-steps",
            "estimated_tokens": 1000,
            "required": true
          },
          {
            "id": "lesson-format",
            "description": "Lesson 格式要求。引用 enforcement/feedback-lessons.md。\n",
            "source": "enforcement/feedback-lessons.md",
            "estimated_tokens": 800,
            "required": true
          },
          {
            "id": "validation-rules",
            "description": "Promotion 驗證規則。引用 governance/validation/。\n",
            "source": "governance/validation/",
            "estimated_tokens": 600,
            "required": true
          },
          {
            "id": "output-format",
            "description": "Feedback 操作結果輸出格式。\n",
            "source": "inline",
            "content": "## Output Format\n- Lesson ID: [generated or existing]\n- Promotion Target: [layer / file path]\n- Validation Results: [pass / fail per gate]\n- Linked Updates: [files that need sync]\n- Status: [created / promoted / rejected / reworked]\n",
            "estimated_tokens": 300,
            "required": true
          }
        ]
      }
    ]
  }

  # Source: runtime/prompt-artifacts/composition-rules.yaml
  PROMPT_ARTIFACTS_COMPOSITION_RULES =   {
    "prompt_composition_rules_version": "v1",
    "status": "candidate",
    "owner_layer": "runtime/prompt-artifacts",
    "description": "Composition rules for prompt artifact generation. 定義如何從各層動態選擇 內容嵌入 prompt artifact。每個規則包含 selection criteria、priority order、 token budget 與 conflict resolution。\n",
    "rules": [
      {
        "rule_id": "workflow-steps",
        "name": "Workflow Steps Selector",
        "description": "根據 task type 從 workflow/ 選擇對應的執行步驟。 每個 task type 對應一個 workflow 子目錄。\n",
        "selection": {
          "method": "task_type_mapping",
          "mapping": {
            "apk-analysis": "workflow/apk-analysis/",
            "repo-analysis": "workflow/repo-analysis/",
            "travel-planning": "workflow/travel-planning/",
            "repo-governance": "governance/",
            "feedback-promotion": "feedback/promotion/"
          },
          "fallback": "workflow/README.md"
        },
        "loading": {
          "strategy": "read_summary_first",
          "summary_source": "{workflow_dir}/README.md",
          "full_source": "{workflow_dir}/{execution-flow}.md",
          "expand_condition": "summary_insufficient == true"
        },
        "estimated_tokens": {
          "summary": 400,
          "full": 1200
        },
        "priority": "P1"
      },
      {
        "rule_id": "intelligence-atoms",
        "name": "Intelligence Atoms Selector",
        "description": "根據 task type 與 task context 從 intelligence/ 選擇相關 atoms。 使用 keyword matching 與 domain mapping 決定哪些 atoms 需要嵌入。\n",
        "selection": {
          "method": "keyword_and_domain_matching",
          "domain_mapping": {
            "apk-analysis": {
              "primary_domain": "engineering/apk-analysis",
              "secondary_domains": [
                "engineering/heuristics",
                "engineering/failure"
              ],
              "keywords": [
                "traffic-analysis",
                "reverse-engineering",
                "flutter",
                "mitm"
              ]
            },
            "repo-analysis": {
              "primary_domain": "engineering/architecture",
              "secondary_domains": [
                "engineering/domain",
                "engineering/anti-patterns"
              ],
              "keywords": [
                "architecture",
                "code-quality",
                "tech-debt"
              ]
            },
            "travel-planning": {
              "primary_domain": "travel",
              "secondary_domains": [
  
              ],
              "keywords": [
                "itinerary",
                "transportation",
                "accommodation"
              ]
            },
            "repo-governance": {
              "primary_domain": "engineering",
              "secondary_domains": [
                "engineering/architecture",
                "engineering/tradeoffs"
              ],
              "keywords": [
                "architecture",
                "migration",
                "lifecycle"
              ]
            },
            "feedback-promotion": {
              "primary_domain": "engineering/heuristics",
              "secondary_domains": [
                "engineering/failure"
              ],
              "keywords": [
                "learning",
                "pattern",
                "heuristic"
              ]
            }
          },
          "max_atoms": 3,
          "min_relevance_score": 0.5
        },
        "loading": {
          "strategy": "inline_summary",
          "format": "### Relevant Intelligence Atoms\n{atom_list}\nEach atom: principle (1 sentence) → when to apply → when not to apply\n",
          "expand_condition": "atom_detail_needed == true"
        },
        "estimated_tokens": {
          "summary_per_atom": 200,
          "full_per_atom": 600
        },
        "priority": "P2"
      },
      {
        "rule_id": "analysis-methods",
        "name": "Analysis Methods Selector",
        "description": "根據 task type 從 analysis/ 選擇相關的分析方法。 僅在需要深入技術細節時載入。\n",
        "selection": {
          "method": "task_type_mapping",
          "mapping": {
            "apk-analysis": "analysis/apk/",
            "repo-analysis": "analysis/repo/"
          },
          "fallback": "analysis/README.md"
        },
        "loading": {
          "strategy": "read_summary_first",
          "summary_source": "{analysis_dir}/README.md",
          "full_source": "{analysis_dir}/{method-file}.md",
          "expand_condition": "technique_detail_needed == true"
        },
        "estimated_tokens": {
          "summary": 300,
          "full": 800
        },
        "priority": "P3",
        "load_condition": "deep_technical_work == true"
      },
      {
        "rule_id": "knowledge-summary",
        "name": "Knowledge Summary Selector",
        "description": "根據 skill_id 從 knowledge/summaries/ 選擇對應的知識摘要。 在 navigation 與 routing 場景中使用。\n",
        "selection": {
          "method": "skill_id_mapping",
          "mapping": {
            "apk-analysis": "knowledge/summaries/apk-analysis-pilot.md",
            "travel-planning": "knowledge/summaries/travel-planning.md",
            "repo-governance": "knowledge/summaries/repo-governance.md",
            "knowledge-navigation": "knowledge/summaries/knowledge-navigation.md",
            "feedback-promotion": "knowledge/summaries/feedback-promotion.md"
          },
          "fallback": "knowledge/summaries/README.md"
        },
        "loading": {
          "strategy": "inline_full",
          "format": "### Knowledge Summary\n{summary_content}\n"
        },
        "estimated_tokens": 500,
        "priority": "P2"
      }
    ],
    "section_ordering": {
      "description": "定義 artifact sections 的預設順序與相依關係。 某些 sections 必須在其他 sections 之前載入。\n",
      "default_order": [
        "task-context",
        "workflow-steps",
        "intelligence-atoms",
        "analysis-methods",
        "artifact-gates",
        "output-format"
      ],
      "dependencies": {
        "intelligence-atoms": {
          "depends_on": "task-context",
          "reason": "需要 task context 來決定哪些 atoms 相關"
        },
        "analysis-methods": {
          "depends_on": "workflow-steps",
          "reason": "需要 workflow steps 來決定哪些方法需要"
        },
        "artifact-gates": {
          "depends_on": "workflow-steps",
          "reason": "需要 workflow steps 來決定哪些 gates 適用"
        }
      }
    },
    "token_budget_rules": {
      "description": "每個 section 的 token budget 分配規則。 如果總 token 超過 hard_stop，依 priority 降序裁減。\n",
      "allocation_strategy": "priority_based_culling",
      "culling_order": [
        "analysis-methods",
        "intelligence-atoms",
        "knowledge-summary",
        "artifact-gates",
        "workflow-steps",
        "task-context",
        "output-format"
      ],
      "culling_method": "switch_to_summary",
      "culling_description": "裁減時不直接移除 section，而是從 full source 切換到 summary 模式。 只有當 summary 模式仍超過 budget 時才移除 section。\n"
    },
    "conflict_resolution": {
      "description": "當多個 composition rules 指向同一個檔案或內容時的處理規則。\n",
      "conflict_rules": [
        {
          "scenario": "同一檔案被多個 sections 引用",
          "action": "只載入一次，共享內容"
        },
        {
          "scenario": "intelligence atom 同時被多個 rules 選中",
          "action": "只嵌入一次，在第一次出現的 section 中"
        },
        {
          "scenario": "workflow step 與 analysis method 內容重疊",
          "action": "workflow step 優先（P1），analysis method 引用 workflow 中的內容"
        }
      ]
    }
  }

  # Source: runtime/recovery/obligation-rebuild.yaml
  RECOVERY_OBLIGATION_REBUILD =   {
    "version": 1.0,
    "status": "candidate",
    "owner_layer": "runtime/recovery",
    "last_updated": "2026-05-15",
    "rebuild_procedures": [
      {
        "id": "rebuild.obligation.full_rebuild",
        "name": "Full Obligation Ledger Rebuild",
        "description": "從頭重建完整的 obligation ledger。 適用於 ledger 完全遺失或嚴重損毀的情況。\n",
        "trigger_conditions": [
          "obligation-ledger.yaml 無法讀取",
          "所有 obligations 狀態皆為 unknown",
          "recovery phase 首次進入"
        ],
        "procedure": [
          {
            "step": 1,
            "action": "LOAD_PHASE_MACHINE",
            "description": "載入 phase-machine.yaml 取得所有 phase 定義",
            "command": "讀取 runtime/phases/phase-machine.yaml"
          },
          {
            "step": 2,
            "action": "ITERATE_ALL_PHASES",
            "description": "遍歷所有 9 個 phase",
            "command": "依 phase 順序（bootstrap→checkpoint→execution→validation→commit→push→readback→finalize→recovery）"
          },
          {
            "step": 3,
            "action": "EXTRACT_OBLIGATIONS_PER_PHASE",
            "description": "從每個 phase 的 obligations 欄位提取義務",
            "command": "讀取 phase-machine.yaml 中各 phase 的 obligations 列表"
          },
          {
            "step": 4,
            "action": "MAP_TO_LEDGER_FORMAT",
            "description": "將 obligations 轉換為 ledger 格式",
            "command": "對每個 obligation 填入：id、phase、name、description、verification_criteria、severity、depends_on、linked_gates"
          },
          {
            "step": 5,
            "action": "SET_INITIAL_STATE",
            "description": "設定初始狀態（全部為 pending）",
            "command": "所有 obligations 初始狀態設為 pending"
          },
          {
            "step": 6,
            "action": "VERIFY_LEDGER_INTEGRITY",
            "description": "驗證重建後的 ledger 完整性",
            "command": "檢查是否有遺漏的 phase 或 obligation"
          }
        ],
        "verify": [
          "所有 9 個 phase 皆有對應的 obligations",
          "每個 obligation 有完整的 id、phase、name、verification_criteria",
          "severity 皆為有效值（critical/high/medium）",
          "depends_on 與 linked_gates 的引用可解析"
        ]
      },
      {
        "id": "rebuild.obligation.phase_partial",
        "name": "Partial Phase Obligation Rebuild",
        "description": "只重建特定 phase 的 obligations。 適用於單一 phase 的 obligation 狀態損毀。\n",
        "trigger_conditions": [
          "特定 phase 的 obligations 狀態遺失",
          "phase transition 後新 phase 的 obligations 未初始化"
        ],
        "procedure": [
          {
            "step": 1,
            "action": "IDENTIFY_TARGET_PHASE",
            "description": "確認需要重建的 phase",
            "command": "讀取目前 current_phase 或指定 phase"
          },
          {
            "step": 2,
            "action": "LOAD_PHASE_OBLIGATIONS",
            "description": "載入該 phase 的 obligations",
            "command": "從 phase-machine.yaml 中該 phase 的 obligations 欄位提取"
          },
          {
            "step": 3,
            "action": "CROSS_REFERENCE_LEDGER",
            "description": "對照 obligation-ledger.yaml 確認完整定義",
            "command": "從 obligation-ledger.yaml 中過濾 phase 對應的 obligations"
          },
          {
            "step": 4,
            "action": "SET_PHASE_OBLIGATIONS",
            "description": "設定該 phase 的 obligations 狀態",
            "command": "所有 obligations 設為 pending，等待執行"
          },
          {
            "step": 5,
            "action": "CHECK_DEPENDENCIES",
            "description": "檢查 dependencies 是否已滿足",
            "command": "對每個 obligation 檢查 depends_on 中的上游 obligation 是否已完成"
          }
        ],
        "verify": [
          "目標 phase 的所有 obligations 已重建",
          "dependencies 鏈完整無斷裂",
          "linked_gates 可解析"
        ]
      },
      {
        "id": "rebuild.obligation.dependency_chain",
        "name": "Dependency Chain Verification",
        "description": "驗證並修復 obligation 之間的 dependency 鏈。 確保沒有 circular dependency 或 dangling reference。\n",
        "trigger_conditions": [
          "obligation 的 depends_on 引用不存在的 obligation",
          "偵測到 circular dependency",
          "linked_gates 引用不存在的 gate"
        ],
        "procedure": [
          {
            "step": 1,
            "action": "BUILD_DEPENDENCY_GRAPH",
            "description": "建立完整的 dependency 圖",
            "command": "從 obligation-ledger.yaml 中提取所有 depends_on 關係"
          },
          {
            "step": 2,
            "action": "DETECT_CIRCULAR",
            "description": "檢測 circular dependency",
            "command": "使用拓撲排序檢查是否有環"
          },
          {
            "step": 3,
            "action": "DETECT_DANGLING",
            "description": "檢測 dangling reference",
            "command": "檢查所有 depends_on 與 linked_gates 的目標是否存在"
          },
          {
            "step": 4,
            "action": "REPAIR_BROKEN_LINKS",
            "description": "修復斷裂的引用",
            "command": "移除或修正無法解析的 depends_on/linked_gates"
          },
          {
            "step": 5,
            "action": "REBUILD_ORDERED_LIST",
            "description": "重建依賴順序的 obligation 列表",
            "command": "依拓撲排序重建執行順序"
          }
        ],
        "verify": [
          "無 circular dependency",
          "所有 depends_on 與 linked_gates 可解析",
          "拓撲排序可產生合法的執行順序"
        ]
      }
    ]
  }

  # Source: runtime/recovery/phase-reconciliation.yaml
  RECOVERY_PHASE_RECONCILIATION =   {
    "version": 1.0,
    "status": "candidate",
    "owner_layer": "runtime/recovery",
    "last_updated": "2026-05-15",
    "reconciliation_procedures": [
      {
        "id": "reconcile.phase.forward_sync",
        "name": "Forward Phase Sync",
        "description": "當 phase_history 顯示已完成的 phase 但 current_phase 落後時的 reconciliation。 例如：validation 已完成但 current_phase 仍停留在 execution。\n",
        "detection": {
          "signal": "phase_history 的最後 phase 與 current_phase 不一致",
          "check": "比對 phase_history 的最後記錄與 current_phase"
        },
        "procedure": [
          {
            "step": 1,
            "action": "COMPARE_HISTORY_VS_CURRENT",
            "description": "比對 phase_history 與 current_phase",
            "command": "phase_history.last vs current_phase"
          },
          {
            "step": 2,
            "action": "IDENTIFY_MISSING_TRANSITIONS",
            "description": "找出缺少的 phase transitions",
            "command": "依 phase 順序找出從 current_phase 到 phase_history.last 之間缺少的 transitions"
          },
          {
            "step": 3,
            "action": "FAST_FORWARD_PHASE",
            "description": "快速前進到正確的 phase",
            "command": "依序執行每個缺少 phase 的 entry_conditions（不重複執行 obligations）"
          },
          {
            "step": 4,
            "action": "UPDATE_CURRENT_PHASE",
            "description": "更新 current_phase",
            "command": "設定 current_phase = phase_history.last"
          },
          {
            "step": 5,
            "action": "RELOAD_PHASE_STATE",
            "description": "重新載入正確 phase 的 state",
            "command": "載入 phase-machine.yaml 中該 phase 的 allowed_actions、blocking_gates、obligations"
          }
        ],
        "verify": [
          "current_phase == phase_history.last",
          "current_phase 的 allowed_actions 與 blocking gates 已載入",
          "無遺漏的 phase transition"
        ]
      },
      {
        "id": "reconcile.phase.backward_sync",
        "name": "Backward Phase Sync",
        "description": "當 current_phase 超前於實際完成狀態時的 reconciliation。 例如：current_phase 為 commit 但 validation 的 obligations 尚未完成。\n",
        "detection": {
          "signal": "current_phase 超前於 obligations 完成狀態",
          "check": "檢查 current_phase 之前的 phase 是否有未完成的 obligations"
        },
        "procedure": [
          {
            "step": 1,
            "action": "IDENTIFY_PRECEDING_PHASE",
            "description": "找出 current_phase 的前一個 phase",
            "command": "依 phase 順序找出 current_phase 的前一個 phase"
          },
          {
            "step": 2,
            "action": "CHECK_PRECEDING_OBLIGATIONS",
            "description": "檢查前一個 phase 的 obligations 完成狀態",
            "command": "讀取 obligation-ledger.yaml 中前一個 phase 的 obligations 狀態"
          },
          {
            "step": 3,
            "action": "ROLLBACK_IF_INCOMPLETE",
            "description": "如果前一個 phase 未完成，回滾到該 phase",
            "command": "設定 current_phase = 前一個 phase，重新載入 state"
          },
          {
            "step": 4,
            "action": "COMPLETE_MISSED_OBLIGATIONS",
            "description": "補完成遺漏的 obligations",
            "command": "執行前一個 phase 的 pending obligations"
          },
          {
            "step": 5,
            "action": "RE_ATTEMPT_TRANSITION",
            "description": "重新嘗試 phase transition",
            "command": "依 phase_transition_rules 重新 transition 到 current_phase"
          }
        ],
        "verify": [
          "所有前置 phase 的 obligations 已完成",
          "phase_transition_rules 中 from→to 合法",
          "current_phase 的 blocking gates 已通過"
        ]
      },
      {
        "id": "reconcile.phase.gate_obligation_mismatch",
        "name": "Gate-Obligation State Mismatch",
        "description": "當 blocking gate 狀態與 obligation 狀態不一致時的 reconciliation。 例如：gate 顯示 passed 但對應的 obligation 顯示 pending。\n",
        "detection": {
          "signal": "gate 狀態與 linked obligation 狀態不一致",
          "check": "比對 blocking-gates.yaml 的 gate 結果與 obligation-ledger.yaml 中 linked_gates 的 obligation 狀態"
        },
        "procedure": [
          {
            "step": 1,
            "action": "LIST_GATE_OBLIGATION_PAIRS",
            "description": "列出所有 gate-obligation 配對",
            "command": "從 blocking-gates.yaml 中提取 linked_obligations，對照 obligation-ledger.yaml"
          },
          {
            "step": 2,
            "action": "IDENTIFY_MISMATCHES",
            "description": "找出狀態不一致的配對",
            "command": "比對 gate 的 passed/failed 與 obligation 的 completed/pending"
          },
          {
            "step": 3,
            "action": "RESOLVE_EACH_MISMATCH",
            "description": "逐項解決不一致",
            "command": "如果 gate passed 但 obligation pending → 重新驗證 obligation；如果 gate failed 但 obligation completed → 重新檢查 gate"
          },
          {
            "step": 4,
            "action": "SYNC_STATES",
            "description": "同步兩邊的狀態",
            "command": "更新不一致的 gate 或 obligation 狀態"
          }
        ],
        "verify": [
          "所有 linked gate-obligation 配對狀態一致",
          "無 gate passed 但 obligation pending 的情況",
          "無 gate failed 但 obligation completed 的情況"
        ]
      },
      {
        "id": "reconcile.phase.full_state_reconciliation",
        "name": "Full State Reconciliation",
        "description": "當多個 phase 狀態元件不一致時的完整 reconciliation。 同時檢查：current_phase、phase_history、blocking_gates、obligations 四者的一致性。\n",
        "trigger_conditions": [
          "recovery phase 進入時",
          "偵測到 2 個以上不一致",
          "手動觸發完整 reconciliation"
        ],
        "procedure": [
          {
            "step": 1,
            "action": "SNAPSHOT_ALL_STATE",
            "description": "快照所有 phase state",
            "command": "記錄 current_phase、phase_history、所有 gate 狀態、所有 obligation 狀態"
          },
          {
            "step": 2,
            "action": "RUN_ALL_DETECTIONS",
            "description": "執行所有不一致檢測",
            "command": "依序執行 forward_sync、backward_sync、gate_obligation_mismatch 的 detection"
          },
          {
            "step": 3,
            "action": "PRIORITIZE_ISSUES",
            "description": "依嚴重度排序不一致問題",
            "command": "critical（phase 順序錯誤）> high（gate-obligation 不一致）> medium（state 過期）"
          },
          {
            "step": 4,
            "action": "RESOLVE_IN_ORDER",
            "description": "依序解決所有不一致",
            "command": "從最高嚴重度開始，依 reconciliation procedure 解決"
          },
          {
            "step": 5,
            "action": "VERIFY_CONSISTENCY",
            "description": "驗證所有 state 一致",
            "command": "重新執行所有 detection 確認無不一致"
          }
        ],
        "verify": [
          "current_phase 與 phase_history 一致",
          "所有 blocking gates 狀態與 obligations 狀態一致",
          "無遺漏的 phase transition",
          "所有 phase 的 obligations 狀態正確"
        ]
      }
    ]
  }

  # Source: runtime/recovery/recovery-strategies.yaml
  RECOVERY_RECOVERY_STRATEGIES =   {
    "version": 1.0,
    "status": "candidate",
    "owner_layer": "runtime/recovery",
    "last_updated": "2026-05-20",
    "recovery_state_machine": {
      "states": [
        "escalation",
        "recovery",
        "rediscovery",
        "replan",
        "execution"
      ],
      "transitions": [
        {
          "from": "escalation",
          "to": "recovery",
          "condition": "mismatch_escalation guard emits suspend_execution or enter_recovery",
          "required_actions": [
            "suspend_assumption",
            "capture_contradicting_evidence"
          ]
        },
        {
          "from": "recovery",
          "to": "rediscovery",
          "condition": "old assumption is marked invalid and trigger class is named",
          "required_actions": [
            "reload_workflow_primary_source",
            "reload_owner_docs"
          ]
        },
        {
          "from": "rediscovery",
          "to": "replan",
          "condition": "required source-of-truth is read or marked not applicable/source missing",
          "required_actions": [
            "rebuild_execution_graph",
            "define_validation_checkpoint"
          ]
        },
        {
          "from": "replan",
          "to": "execution",
          "condition": "new execution graph has goal, route, dependencies, checkpoint, and validation",
          "required_actions": [
            "explain_failure",
            "resume_with_new_plan"
          ]
        }
      ]
    },
    "recovery_levels": [
      {
        "level": "L1",
        "name": "simple_retry",
        "meaning": "single transient tool or command failure",
        "default_action": "retry_once"
      },
      {
        "level": "L2",
        "name": "reload_local_workflow",
        "meaning": "local rule or workflow may be stale",
        "default_action": "reload_current_workflow"
      },
      {
        "level": "L3",
        "name": "reload_source_of_truth",
        "meaning": "repeated failure or assumption drift indicates source context is insufficient",
        "default_action": "reload_owner_docs_and_primary_sources"
      },
      {
        "level": "L4",
        "name": "rebuild_execution_graph",
        "meaning": "user contradiction, evidence conflict, source-of-truth miss, or automation drift",
        "default_action": "enter_recovery_and_rebuild_execution_graph"
      },
      {
        "level": "L5",
        "name": "assumption_collapse_rediscovery",
        "meaning": "task framing or route selection may be wrong",
        "default_action": "return_to_discovery"
      }
    ],
    "required_actions": [
      {
        "id": "suspend_assumption",
        "description": "write the old assumption and mark it invalid before continuing"
      },
      {
        "id": "reload_workflow_primary_source",
        "description": "reload the selected workflow primary source plus artifact gates"
      },
      {
        "id": "reload_owner_docs",
        "description": "reload owner docs, UI map, API catalog, runtime source, or project contract as applicable"
      },
      {
        "id": "rebuild_execution_graph",
        "description": "rebuild goal -> route -> dependencies -> source-of-truth -> checkpoint -> validation"
      },
      {
        "id": "explain_failure",
        "description": "explain assumption, contradiction, new source-of-truth, revised plan, and validation"
      }
    ],
    "output_schema": {
      "required_fields": [
        "trigger_class",
        "level",
        "old_assumption",
        "contradicting_evidence",
        "source_of_truth_reloaded",
        "new_execution_graph",
        "next_action",
        "validation"
      ],
      "validation": [
        "trigger_class is one of repeated_failure, user_contradiction, evidence_conflict, source_of_truth_miss, automation_drift, assumption_drift",
        "level is L1-L5",
        "source_of_truth_reloaded lists read, not_applicable, or source_missing items",
        "new_execution_graph contains goal, route, dependencies, checkpoint, and validation"
      ]
    },
    "strategies": [
      {
        "id": "strategy.mismatch_escalation_recovery",
        "name": "Mismatch Escalation Recovery",
        "description": "當 execution 中段出現 repeated failure、user contradiction、evidence conflict、source-of-truth miss 或 automation drift 時，停止局部 patch 並重建 execution graph。",
        "applicable_gates": [
          "runtime.guard.mismatch_escalation",
          "enforcement.escalation-policy"
        ],
        "detection": {
          "signal": "mismatch_escalation guard emits suspend_execution, enter_recovery, or rediscovery_required",
          "check": "檢查 runtime/pipeline/guard-chain.yaml execution stage 與 runtime/guards/circuit-breaker.yaml mismatch_escalation trigger"
        },
        "state_transitions": [
          "escalation",
          "recovery",
          "rediscovery",
          "replan",
          "execution"
        ],
        "repair_steps": [
          {
            "action": "SUSPEND_ASSUMPTION",
            "description": "停止沿用舊 mental model，寫出舊假設與反證",
            "command": "依 enforcement/escalation-policy.md 的 Required Recovery Output 填寫 old_assumption 與 contradicting_evidence"
          },
          {
            "action": "RELOAD_SOURCE_OF_TRUTH",
            "description": "重讀 workflow primary source、artifact gates、owner docs 或 runtime source",
            "command": "列出 read / not_applicable / source_missing 的 source-of-truth reload set"
          },
          {
            "action": "REBUILD_EXECUTION_GRAPH",
            "description": "重建 goal、route、dependencies、checkpoint 與 validation",
            "command": "產出 new_execution_graph 並標記下一個可安全執行的 action"
          },
          {
            "action": "EXPLAIN_AND_RESUME",
            "description": "說明 failure 與新路線後恢復 execution",
            "command": "輸出 recovery schema，然後依新 execution graph 繼續"
          }
        ],
        "verify_steps": [
          "舊假設已明確 suspend",
          "required source-of-truth 已讀或標記 not_applicable/source_missing",
          "new_execution_graph 包含 goal、route、dependencies、checkpoint、validation",
          "下一步不再沿用被反證的 route 或 automation"
        ],
        "escalation": {
          "condition": "rediscovery 後仍無法建立可信 execution graph",
          "action": "升級 L5 assumption collapse，回到 discovery 並請使用者確認 task framing"
        }
      },
      {
        "id": "strategy.phase_drift",
        "name": "Phase Drift Recovery",
        "description": "當 agent 跳過某個 phase 或 phase 順序錯誤時的修復策略。 例如：直接從 execution 嘗試 commit（違反 execution→validation→commit 規則）。\n",
        "applicable_gates": [
          "gate.phase_bootstrap.core_rules_loaded",
          "gate.phase_bootstrap.layout_loaded",
          "gate.phase_checkpoint.obligation_ledger_loaded",
          "gate.phase_checkpoint.blocking_gates_checked"
        ],
        "detection": {
          "signal": "agent 嘗試執行不屬於目前 phase 的 allowed_actions",
          "check": "比對目前 phase 的 allowed_actions 與 agent 的下一步意圖"
        },
        "repair_steps": [
          {
            "action": "IDENTIFY_CURRENT_PHASE",
            "description": "確認目前 phase 與 blocking gate",
            "command": "讀取 runtime/phases/phase-machine.yaml 中目前 phase 的 allowed_actions"
          },
          {
            "action": "ROLLBACK_TO_CORRECT_PHASE",
            "description": "回到正確的 phase 重新執行",
            "command": "更新 phase state 到上一個合法 phase"
          },
          {
            "action": "REPLAY_MISSED_STEPS",
            "description": "補執行被跳過的 phase 步驟",
            "command": "依 phase-machine.yaml 的 phase_transition_rules 重新依序執行"
          },
          {
            "action": "VERIFY_PHASE_SEQUENCE",
            "description": "確認 phase 順序已恢復正確",
            "command": "檢查 phase_transition_rules 中 from→to 是否合法"
          }
        ],
        "verify_steps": [
          "目前 phase 的 allowed_actions 包含 agent 的下一步意圖",
          "所有 blocking gates 已通過",
          "phase_transition_rules 中 from→to 未被 blocked"
        ],
        "escalation": {
          "condition": "修復後仍無法通過 blocking gate",
          "action": "記錄 failure pattern 到 enforcement/failure-patterns/，通知使用者"
        }
      },
      {
        "id": "strategy.obligation_missed",
        "name": "Missed Obligation Recovery",
        "description": "當 agent 遺漏某個 obligation 時的修復策略。 例如：commit 前未執行 linked_updates_complete 檢查。\n",
        "applicable_gates": [
          "gate.phase_execution.all_obligations_met",
          "gate.phase_validation.linked_updates_complete",
          "gate.phase_validation.artifacts_complete",
          "gate.phase_commit.commit_message_valid",
          "gate.phase_commit.pre_commit_hooks_passed",
          "gate.phase_finalize.all_obligations_closed",
          "gate.phase_finalize.knowledge_updated"
        ],
        "detection": {
          "signal": "blocking gate 因 obligation 未滿足而阻斷",
          "check": "查詢 runtime/obligations/obligation-ledger.yaml 中本 phase 的 pending obligations"
        },
        "repair_steps": [
          {
            "action": "LIST_PENDING_OBLIGATIONS",
            "description": "列出本 phase 所有未滿足的 obligations",
            "command": "讀取 obligation-ledger.yaml 中 phase 對應的 obligations，過濾未完成者"
          },
          {
            "action": "EXECUTE_EACH_OBLIGATION",
            "description": "依序執行每個 pending obligation 的 verification_criteria",
            "command": "對每個 pending obligation 執行其 verification_criteria 定義的檢查"
          },
          {
            "action": "RE_VERIFY_GATES",
            "description": "重新檢查 blocking gates",
            "command": "執行 blocking-gates.yaml 中本 phase 的 gates，確認 severity 狀態"
          },
          {
            "action": "UPDATE_OBLIGATION_STATE",
            "description": "更新 obligation 狀態為 completed",
            "command": "記錄 obligation 已完成，允許 phase transition"
          }
        ],
        "verify_steps": [
          "本 phase 所有 obligations 的 verification_criteria 已滿足",
          "linked_gates 中關聯的 blocking gates 已通過"
        ],
        "escalation": {
          "condition": "obligation 無法滿足（例如缺少必要 artifact）",
          "action": "記錄為 failure，進入 recovery phase 重新規劃"
        }
      },
      {
        "id": "strategy.gate_blocked",
        "name": "Blocking Gate Recovery",
        "description": "當 blocking gate 因條件不滿足而阻斷 phase transition 時的通用修復策略。 適用於任何 severity=critical 或 severity=high 的 gate。\n",
        "applicable_gates": [
          "*"
        ],
        "detection": {
          "signal": "blocking gate 檢查回傳 failed",
          "check": "讀取 blocking-gates.yaml 中該 gate 的 failure_message 與 failure_action"
        },
        "repair_steps": [
          {
            "action": "READ_GATE_DETAILS",
            "description": "讀取阻斷 gate 的完整定義",
            "command": "查詢 blocking-gates.yaml 中 gate 的 check.type、check.verification、failure_action"
          },
          {
            "action": "EXECUTE_FAILURE_ACTION",
            "description": "執行 gate 定義的 failure_action",
            "command": "依 failure_action 執行（例如：re-verify、re-plan、re-execute）"
          },
          {
            "action": "RE_RUN_GATE_CHECK",
            "description": "重新執行 gate 檢查",
            "command": "依 check.verification 重新驗證"
          },
          {
            "action": "CONFIRM_GATE_PASSED",
            "description": "確認 gate 已通過",
            "command": "檢查 gate 狀態是否變為 passed"
          }
        ],
        "verify_steps": [
          "gate 狀態為 passed",
          "phase transition 不再被此 gate 阻斷"
        ],
        "escalation": {
          "condition": "gate 連續 3 次修復仍失敗",
          "action": "進入 recovery phase，記錄 failure pattern，通知使用者"
        }
      },
      {
        "id": "strategy.transaction_incomplete",
        "name": "Incomplete Transaction Recovery",
        "description": "當 Ai-skill writeback transaction 未正確關閉時的修復策略。 例如：commit 後未 push、push 後未讀回確認。\n",
        "applicable_gates": [
          "gate.phase_push.remote_reachable",
          "gate.phase_push.push_successful",
          "gate.phase_readback.working_tree_clean",
          "gate.phase_readback.remote_synced"
        ],
        "detection": {
          "signal": "git status 顯示 dirty tree 或 ahead/behind",
          "check": "執行 git status --short --branch 檢查"
        },
        "repair_steps": [
          {
            "action": "CHECK_GIT_STATUS",
            "description": "檢查 git 狀態",
            "command": "git status --short --branch"
          },
          {
            "action": "STAGE_UNCOMMITTED",
            "description": "將未 stage 的檔案加入",
            "command": "git add -A"
          },
          {
            "action": "COMMIT_WITH_MESSAGE",
            "description": "提交變更",
            "command": "git commit -m '<descriptive message>'"
          },
          {
            "action": "PUSH_TO_REMOTE",
            "description": "推送至遠端",
            "command": "git push"
          },
          {
            "action": "VERIFY_CLEAN",
            "description": "確認 working tree clean",
            "command": "git status --short --branch 應顯示 clean"
          }
        ],
        "verify_steps": [
          "git status 顯示 clean",
          "branch 沒有 ahead/behind",
          "所有必要檔案已 commit 並 push"
        ],
        "escalation": {
          "condition": "push 失敗（remote rejected / network error）",
          "action": "記錄錯誤，通知使用者手動處理"
        }
      },
      {
        "id": "strategy.knowledge_stale",
        "name": "Stale Knowledge Recovery",
        "description": "當 knowledge 或 runtime surface 與 canonical source 不同步時的修復策略。 例如：prose 已修改但 generated YAML 未更新。\n",
        "applicable_gates": [
          "gate.phase_finalize.knowledge_updated",
          "gate.phase_readback.remote_synced"
        ],
        "detection": {
          "signal": "pre-commit hook 或 validation gate 檢測到 prose/YAML 不一致",
          "check": "比較 prose source 與 generated YAML 的修改時間"
        },
        "repair_steps": [
          {
            "action": "IDENTIFY_STALE_SURFACES",
            "description": "找出不同步的 generated surfaces",
            "command": "查詢 runtime.db 確認 generated_surfaces 表的最新記錄是否與 prose source 一致"
          },
          {
            "action": "RE_RUN_COMPILER",
            "description": "重新執行 compiler 更新 runtime.db",
            "command": "執行 ruby runtime/compiler/compiler-engine.rb"
          },
          {
            "action": "VERIFY_SYNC",
            "description": "驗證 runtime.db 已更新",
            "command": "執行 ruby scripts/validate-runtime-db.rb"
          },
          {
            "action": "COMMIT_SYNCED_SURFACES",
            "description": "提交同步後的變更（含 compiler 更新）",
            "command": "git add -A && git commit"
          }
        ],
        "verify_steps": [
          "runtime.db 的 generated_surfaces 表包含最新記錄",
          "validate-runtime-db.rb 回傳 exit 0"
        ],
        "escalation": {
          "condition": "compiler 不存在或無法執行",
          "action": "手動同步 generated YAML，記錄 compiler 建置需求到 P2"
        }
      },
      {
        "id": "strategy.unknown_failure",
        "name": "Unknown Failure Recovery",
        "description": "當 failure pattern 不屬於上述任何已知策略時的通用修復策略。 先診斷、再嘗試修復、最後記錄為新的 failure pattern。\n",
        "applicable_gates": [
          "*"
        ],
        "detection": {
          "signal": "無法匹配任何已知 strategy 的 applicable_gates",
          "check": "執行通用診斷流程"
        },
        "repair_steps": [
          {
            "action": "DIAGNOSE",
            "description": "收集診斷資訊",
            "command": "檢查：目前 phase、blocking gate、obligation 狀態、git status、最近操作"
          },
          {
            "action": "CLASSIFY",
            "description": "分類 failure type",
            "command": "比對 enforcement/failure-patterns/ 中的已知 pattern"
          },
          {
            "action": "ATTEMPT_REPAIR",
            "description": "嘗試最佳猜測修復",
            "command": "根據診斷結果選擇最接近的 strategy 嘗試"
          },
          {
            "action": "RECORD_PATTERN",
            "description": "記錄新的 failure pattern",
            "command": "寫入 enforcement/failure-patterns/ 供後續參考"
          }
        ],
        "verify_steps": [
          "blocking gate 已通過或已進入 recovery phase",
          "failure pattern 已記錄"
        ],
        "escalation": {
          "condition": "所有修復嘗試均失敗",
          "action": "通知使用者，提供完整診斷資訊"
        }
      }
    ]
  }

  # Source: runtime/recovery/state-repair.yaml
  RECOVERY_STATE_REPAIR =   {
    "version": 1.0,
    "status": "candidate",
    "owner_layer": "runtime/recovery",
    "last_updated": "2026-05-15",
    "repair_procedures": [
      {
        "id": "repair.state.current_phase_corrupted",
        "name": "Current Phase State Repair",
        "description": "當 current_phase 遺失或與實際狀態不一致時的修復程序。 例如：recovery 完成後 current_phase 仍停留在 recovery 而非 checkpoint。\n",
        "detection": {
          "signal": "current_phase 為 nil 或不符合 phase_transition_rules",
          "check": "比對 phase_history 的最後合法 phase 與 current_phase"
        },
        "procedure": [
          {
            "step": 1,
            "action": "READ_PHASE_HISTORY",
            "description": "讀取 phase_history 找出最後一個合法 phase",
            "command": "檢查 conversation history 或 phase log 中的最後一個 completed phase"
          },
          {
            "step": 2,
            "action": "DETERMINE_CORRECT_PHASE",
            "description": "根據 phase_transition_rules 決定正確的 current_phase",
            "command": "從最後合法 phase 依 phase_transition_rules 推斷應在的 phase"
          },
          {
            "step": 3,
            "action": "RESET_CURRENT_PHASE",
            "description": "重置 current_phase 到正確值",
            "command": "設定 current_phase = 正確 phase"
          },
          {
            "step": 4,
            "action": "RELOAD_PHASE_STATE",
            "description": "重新載入該 phase 的 allowed_actions、forbidden_actions、blocking_gates",
            "command": "讀取 phase-machine.yaml 中該 phase 的定義"
          }
        ],
        "verify": [
          "current_phase 不為 nil",
          "current_phase 的 allowed_actions 與目前操作一致",
          "phase_transition_rules 中 from=current_phase 的規則可解析"
        ]
      },
      {
        "id": "repair.state.blocking_gates_stale",
        "name": "Blocking Gates State Repair",
        "description": "當 blocking gates 狀態過期或未正確初始化時的修復程序。 例如：phase 已變更但 blocking gates 仍停留在上一個 phase 的狀態。\n",
        "detection": {
          "signal": "blocking gates 狀態與目前 phase 不匹配",
          "check": "比對 blocking-gates.yaml 中目前 phase 的 gates 與實際檢查結果"
        },
        "procedure": [
          {
            "step": 1,
            "action": "IDENTIFY_CURRENT_PHASE_GATES",
            "description": "找出目前 phase 應檢查的所有 gates",
            "command": "讀取 blocking-gates.yaml 中 phase 對應的 gates"
          },
          {
            "step": 2,
            "action": "RESET_GATE_STATES",
            "description": "重置所有 gate 狀態為 unchecked",
            "command": "清除舊的 gate 檢查結果"
          },
          {
            "step": 3,
            "action": "RE_RUN_GATES",
            "description": "重新執行所有 gates 的檢查",
            "command": "依 blocking-gates.yaml 的 check.verification 逐項執行"
          },
          {
            "step": 4,
            "action": "RECORD_GATE_RESULTS",
            "description": "記錄新的 gate 檢查結果",
            "command": "記錄每個 gate 的 passed/failed 狀態"
          }
        ],
        "verify": [
          "所有目前 phase 的 gates 已檢查",
          "無未檢查的 critical/high gates",
          "gate 結果與目前 phase 一致"
        ]
      },
      {
        "id": "repair.state.obligation_ledger_desync",
        "name": "Obligation Ledger Desync Repair",
        "description": "當 obligation ledger 狀態與實際完成情況不一致時的修復程序。 例如：obligation 已執行但 ledger 未更新為 completed。\n",
        "detection": {
          "signal": "obligation 狀態與實際執行情況不符",
          "check": "比對 obligation-ledger.yaml 的 verification_criteria 與實際檢查結果"
        },
        "procedure": [
          {
            "step": 1,
            "action": "LIST_ALL_OBLIGATIONS",
            "description": "列出目前 phase 的所有 obligations",
            "command": "讀取 obligation-ledger.yaml 中 phase 對應的 obligations"
          },
          {
            "step": 2,
            "action": "VERIFY_EACH_OBLIGATION",
            "description": "逐項驗證每個 obligation 的完成狀態",
            "command": "對每個 obligation 執行 verification_criteria"
          },
          {
            "step": 3,
            "action": "UPDATE_LEDGER_STATE",
            "description": "更新 ledger 中每個 obligation 的狀態",
            "command": "標記已完成的 obligation 為 completed，未完成的為 pending"
          },
          {
            "step": 4,
            "action": "RE_CHECK_GATES",
            "description": "重新檢查 linked gates",
            "command": "對 linked_gates 中關聯的 blocking gates 重新檢查"
          }
        ],
        "verify": [
          "所有 obligations 的狀態與實際一致",
          "linked_gates 的狀態已同步更新",
          "無 false-positive completed obligations"
        ]
      },
      {
        "id": "repair.state.phase_history_gap",
        "name": "Phase History Gap Repair",
        "description": "當 phase_history 有遺漏（gap）時的修復程序。 例如：直接從 execution 跳到 finalize，缺少 validation→commit→push→readback。\n",
        "detection": {
          "signal": "phase_history 中缺少中間 phase",
          "check": "比對 phase_history 的 phase 順序與 phase-machine.yaml 的 phase_transition_rules"
        },
        "procedure": [
          {
            "step": 1,
            "action": "IDENTIFY_GAP",
            "description": "找出 phase_history 中缺少的 phase",
            "command": "依 phase-machine.yaml 的 phase 順序比對 phase_history"
          },
          {
            "step": 2,
            "action": "ASSESS_GAP_IMPACT",
            "description": "評估 gap 的影響範圍",
            "command": "檢查缺少的 phase 是否有必須執行的 obligations 或 blocking gates"
          },
          {
            "step": 3,
            "action": "BACKFILL_GAP",
            "description": "補執行缺少的 phase",
            "command": "依 phase-machine.yaml 執行缺少 phase 的 entry_conditions 與 obligations"
          },
          {
            "step": 4,
            "action": "REBUILD_PHASE_HISTORY",
            "description": "重建完整的 phase_history",
            "command": "記錄補執行後的完整 phase 順序"
          }
        ],
        "verify": [
          "phase_history 中無 gap",
          "所有 phase 的 obligations 已滿足",
          "phase_transition_rules 中所有 from→to 皆合法"
        ]
      },
      {
        "id": "repair.state.allowed_actions_mismatch",
        "name": "Allowed Actions Mismatch Repair",
        "description": "當 agent 嘗試的 action 不屬於目前 phase 的 allowed_actions 時的修復程序。 例如：在 execution phase 嘗試執行 commit。\n",
        "detection": {
          "signal": "agent action 不在目前 phase 的 allowed_actions 中",
          "check": "比對 agent 的下一步意圖與 phase-machine.yaml 中目前 phase 的 allowed_actions"
        },
        "procedure": [
          {
            "step": 1,
            "action": "IDENTIFY_FORBIDDEN_ACTION",
            "description": "確認哪個 action 被禁止",
            "command": "查詢 phase-machine.yaml 中目前 phase 的 forbidden_actions"
          },
          {
            "step": 2,
            "action": "CHECK_PHASE_TRANSITION_RULES",
            "description": "檢查 phase_transition_rules 確認是否可先 transition",
            "command": "查詢 phase_transition_rules 中是否有合法路徑可到達目標 action 的 phase"
          },
          {
            "step": 3,
            "action": "PERFORM_LEGAL_TRANSITION",
            "description": "如果可 transition，先執行合法 phase transition",
            "command": "依 phase_transition_rules 依序 transition 到目標 phase"
          },
          {
            "step": 4,
            "action": "BLOCK_IF_ILLEGAL",
            "description": "如果 transition 不合法，阻斷並通知",
            "command": "回傳 failure_message，記錄 forbidden action 嘗試"
          }
        ],
        "verify": [
          "agent 的下一步 action 在目前 phase 的 allowed_actions 中",
          "無 forbidden action 被執行",
          "phase_transition_rules 未被違反"
        ]
      }
    ]
  }

  # Source: runtime/router/activation-rules.yaml
  ROUTER_ACTIVATION_RULES =   {
    "activation_rules_version": "v2",
    "status": "validated",
    "owner_layer": "runtime/router",
    "description": "Runtime activation rules for lazy-load enforcement rules. Core Bootstrap rules (rule-weight, dependency-reading, conversation-goal-ledger) 永遠 preload。 其餘 rules 依本檔定義的條件 activate。\n",
    "core_bootstrap": [
      "rule-weight",
      "dependency-reading",
      "conversation-goal-ledger"
    ],
    "rules": [
      {
        "rule_id": "enforcement.linked-updates",
        "description": "連動更新 — 改一處時哪些文件必須同步",
        "activation": {
          "when": [
            {
              "type": "file_change",
              "pattern": "**/*.md",
              "count": ">=2"
            },
            {
              "type": "task_intent",
              "matches": [
                "architecture-update",
                "migration",
                "refactor",
                "rename",
                "restructure"
              ]
            },
            {
              "type": "user_signal",
              "matches": [
                "連動",
                "同步更新",
                "linked",
                "影響其他文件"
              ]
            }
          ]
        },
        "load": {
          "strategy": "lazy",
          "priority": "P1",
          "estimated_tokens": 800,
          "source": "enforcement/linked-updates.md"
        }
      },
      {
        "rule_id": "enforcement.failure-learning-system",
        "description": "失效學習系統 — 將重複 agent 錯誤轉成 durable prevention",
        "activation": {
          "when": [
            {
              "type": "user_signal",
              "matches": [
                "失誤",
                "漏讀",
                "忘記",
                "錯誤",
                "failure",
                "miss",
                "漏做",
                "閉環不完整"
              ]
            },
            {
              "type": "validation_gap",
              "matches": [
                "close-loop",
                "commit",
                "push",
                "sync",
                "dirty",
                "ahead",
                "behind"
              ]
            },
            {
              "type": "task_intent",
              "matches": [
                "debug",
                "troubleshoot",
                "fix-error"
              ]
            }
          ]
        },
        "load": {
          "strategy": "lazy",
          "priority": "P1",
          "estimated_tokens": 1200,
          "source": "enforcement/failure-learning-system.md"
        }
      },
      {
        "rule_id": "enforcement.decision-efficiency",
        "description": "決策效率 — 以最高收益路線做決策",
        "activation": {
          "when": [
            {
              "type": "task_complexity",
              "routes": ">=3"
            },
            {
              "type": "user_signal",
              "matches": [
                "選擇",
                "路線",
                "priority",
                "先做哪個",
                "哪條路",
                "比較"
              ]
            },
            {
              "type": "task_intent",
              "matches": [
                "route-selection",
                "cost-benefit",
                "trade-off"
              ]
            }
          ]
        },
        "load": {
          "strategy": "lazy",
          "priority": "P2",
          "estimated_tokens": 600,
          "source": "enforcement/decision-efficiency.md"
        }
      },
      {
        "rule_id": "enforcement.tool-neutral-documentation",
        "description": "工具中立文件 — 可重用文件保持工具中立",
        "activation": {
          "when": [
            {
              "type": "task_intent",
              "matches": [
                "write-documentation",
                "create-template",
                "update-readme",
                "create-enforcement-rule"
              ]
            },
            {
              "type": "file_change",
              "pattern": "enforcement/**"
            },
            {
              "type": "file_change",
              "pattern": "skills/**/SKILL.md"
            }
          ]
        },
        "load": {
          "strategy": "lazy",
          "priority": "P2",
          "estimated_tokens": 500,
          "source": "enforcement/tool-neutral-documentation.md"
        }
      },
      {
        "rule_id": "enforcement.document-todo-list",
        "description": "文件 TODO 表 — 文件前段 TODO 與 goal/章節連結",
        "activation": {
          "when": [
            {
              "type": "task_intent",
              "matches": [
                "update-document",
                "complete-document",
                "review-document"
              ]
            },
            {
              "type": "file_has_todo",
              "pattern": "**/*.md"
            }
          ]
        },
        "load": {
          "strategy": "lazy",
          "priority": "P2",
          "estimated_tokens": 300,
          "source": "enforcement/document-todo-list.md"
        }
      },
      {
        "rule_id": "enforcement.goal-action-validation",
        "description": "目標、執行、驗證 — 重要工作單元的驗證流程",
        "activation": {
          "when": [
            {
              "type": "task_intent",
              "matches": [
                "critical-change",
                "destructive-action",
                "production-deploy",
                "security-review"
              ]
            },
            {
              "type": "user_signal",
              "matches": [
                "驗證",
                "validate",
                "確認",
                "confirm"
              ]
            }
          ]
        },
        "load": {
          "strategy": "lazy",
          "priority": "P2",
          "estimated_tokens": 500,
          "source": "enforcement/goal-action-validation.md"
        }
      },
      {
        "rule_id": "enforcement.neutral-language",
        "description": "中性低爭議用語 — 文件語言一致性",
        "activation": {
          "when": [
            {
              "type": "task_intent",
              "matches": [
                "write-documentation",
                "review-documentation",
                "translate"
              ]
            },
            {
              "type": "file_change",
              "pattern": "enforcement/**"
            },
            {
              "type": "file_change",
              "pattern": "skills/**/README.md"
            }
          ]
        },
        "load": {
          "strategy": "lazy",
          "priority": "P2",
          "estimated_tokens": 300,
          "source": "enforcement/neutral-language.md"
        }
      },
      {
        "rule_id": "enforcement.sanitization",
        "description": "去敏與占位符 — 什麼不可寫進可重用文件",
        "activation": {
          "when": [
            {
              "type": "task_intent",
              "matches": [
                "write-feedback",
                "create-lesson",
                "write-feedback-history"
              ]
            },
            {
              "type": "file_change",
              "pattern": "**/feedback_history/**"
            }
          ]
        },
        "load": {
          "strategy": "lazy",
          "priority": "P1",
          "estimated_tokens": 400,
          "source": "enforcement/sanitization.md"
        }
      },
      {
        "rule_id": "enforcement.authorization-scope",
        "description": "授權與範圍 — 僅在授權範圍內分析",
        "activation": {
          "when": [
            {
              "type": "task_intent",
              "matches": [
                "security-analysis",
                "penetration-test",
                "vulnerability-scan",
                "external-system-access"
              ]
            },
            {
              "type": "user_signal",
              "matches": [
                "授權",
                "authorization",
                "scope",
                "邊界"
              ]
            }
          ]
        },
        "load": {
          "strategy": "lazy",
          "priority": "P0",
          "estimated_tokens": 300,
          "source": "enforcement/authorization-scope.md"
        }
      },
      {
        "rule_id": "enforcement.cross-skill-references",
        "description": "Cross-skill references — 一個 skill 引用另一個 skill 的規範",
        "activation": {
          "when": [
            {
              "type": "task_intent",
              "matches": [
                "cross-skill",
                "multi-skill",
                "skill-integration"
              ]
            },
            {
              "type": "file_change",
              "pattern": "skills/**",
              "count": ">=2"
            }
          ]
        },
        "load": {
          "strategy": "lazy",
          "priority": "P2",
          "estimated_tokens": 400,
          "source": "enforcement/cross-skill-references.md"
        }
      },
      {
        "rule_id": "enforcement.feedback-lessons",
        "description": "Feedback 與技巧條目 — 檔名規則、模板、agent 行為",
        "activation": {
          "when": [
            {
              "type": "task_intent",
              "matches": [
                "write-feedback",
                "create-lesson",
                "promote-lesson",
                "review-feedback"
              ]
            },
            {
              "type": "file_change",
              "pattern": "**/feedback_history/**"
            }
          ]
        },
        "load": {
          "strategy": "lazy",
          "priority": "P2",
          "estimated_tokens": 600,
          "source": "enforcement/feedback-lessons.md"
        }
      },
      {
        "rule_id": "enforcement.content-layering",
        "description": "內容分層 — 共用規則／技巧／業務專案各自放哪",
        "activation": {
          "when": [
            {
              "type": "task_intent",
              "matches": [
                "migration",
                "restructure",
                "content-organization"
              ]
            },
            {
              "type": "file_change",
              "pattern": "**/README.md",
              "count": ">=3"
            }
          ]
        },
        "load": {
          "strategy": "lazy",
          "priority": "P2",
          "estimated_tokens": 300,
          "source": "enforcement/content-layering.md"
        }
      },
      {
        "rule_id": "enforcement.reusable-guidance-boundary",
        "description": "可重用規則與專案證據邊界",
        "activation": {
          "when": [
            {
              "type": "task_intent",
              "matches": [
                "write-feedback",
                "create-lesson",
                "promote-to-enforcement-rule"
              ]
            },
            {
              "type": "user_signal",
              "matches": [
                "泛化",
                "可重用",
                "reusable",
                "抽象化"
              ]
            }
          ]
        },
        "load": {
          "strategy": "lazy",
          "priority": "P2",
          "estimated_tokens": 400,
          "source": "enforcement/reusable-guidance-boundary.md"
        }
      }
    ]
  }

  # Source: runtime/scheduler/execution-queue.yaml
  SCHEDULER_EXECUTION_QUEUE =   {
    "version": 1.0,
    "status": "candidate",
    "owner_layer": "runtime/scheduler",
    "last_updated": "2026-05-15",
    "queue_structure": {
      "name": "Phase Execution Queue",
      "description": "每個 phase 有一個 execution queue，包含該 phase 的所有待執行項目。 Queue 是動態的：完成一個項目可能新增相依項目到 queue 中。\n",
      "fields": [
        {
          "name": "id",
          "type": "string",
          "description": "項目唯一識別碼",
          "format": "queue.{phase}.{type}.{seq}"
        },
        {
          "name": "phase",
          "type": "string",
          "description": "所屬 phase",
          "enum": [
            "bootstrap",
            "checkpoint",
            "execution",
            "validation",
            "commit",
            "push",
            "readback",
            "finalize",
            "recovery"
          ]
        },
        {
          "name": "type",
          "type": "string",
          "description": "項目類型",
          "enum": [
            "gate",
            "obligation",
            "action",
            "repair",
            "reconciliation"
          ]
        },
        {
          "name": "ref",
          "type": "string",
          "description": "引用 blocking-gates.yaml、obligation-ledger.yaml 或 recovery 中的 id"
        },
        {
          "name": "priority",
          "type": "string",
          "description": "優先級",
          "enum": [
            "P0",
            "P1",
            "P2",
            "P3"
          ]
        },
        {
          "name": "status",
          "type": "string",
          "description": "執行狀態",
          "enum": [
            "pending",
            "running",
            "completed",
            "failed",
            "skipped"
          ]
        },
        {
          "name": "depends_on",
          "type": "list[string]",
          "description": "此項目依賴的其他 queue item id"
        },
        {
          "name": "estimated_cost",
          "type": "object",
          "description": "估計的資源消耗",
          "properties": {
            "tokens": "integer",
            "steps": "integer"
          }
        },
        {
          "name": "deadline",
          "type": "integer",
          "description": "應在多少步驟內完成（null = 無期限）"
        },
        {
          "name": "result",
          "type": "object",
          "description": "執行結果",
          "properties": {
            "status": "string",
            "message": "string",
            "timestamp": "string"
          }
        }
      ]
    },
    "queue_management_rules": [
      {
        "id": "queue.rule.init_on_phase_entry",
        "name": "Initialize Queue on Phase Entry",
        "description": "進入新 phase 時初始化 execution queue",
        "rule": "載入 phase-machine.yaml 中該 phase 的 blocking_gates 與 obligations， 加上 blocking-gates.yaml 與 obligation-ledger.yaml 的完整定義， 建立初始 queue。\n",
        "trigger": "phase transition"
      },
      {
        "id": "queue.rule.sort_by_priority",
        "name": "Sort by Priority",
        "description": "queue 依 priority 降序排列",
        "rule": "P0 > P1 > P2 > P3，同 priority 內依 deadline 升序",
        "trigger": "每次 queue 變更"
      },
      {
        "id": "queue.rule.resolve_dependencies",
        "name": "Resolve Dependencies",
        "description": "有 depends_on 的項目必須在上游完成後才能執行",
        "rule": "如果 depends_on 中有未完成的項目，該項目設為 blocked",
        "trigger": "每次 queue 變更"
      },
      {
        "id": "queue.rule.dequeue_next",
        "name": "Dequeue Next Runnable",
        "description": "取出 queue 中 priority 最高且 dependencies 已滿足的項目",
        "rule": "過濾 status=pending 且所有 depends_on 皆為 completed 的項目，取 priority 最高者",
        "trigger": "目前項目完成時"
      },
      {
        "id": "queue.rule.handle_failure",
        "name": "Handle Item Failure",
        "description": "當項目執行失敗時的處理",
        "rule": "如果 severity=critical → 阻斷 phase transition，進入 recovery phase。 如果 severity=high → 阻斷 transition 但可進入 recovery。 如果 severity=medium → 記錄 warning，繼續執行。\n",
        "trigger": "項目 status=failed"
      },
      {
        "id": "queue.rule.recovery_override",
        "name": "Recovery Override",
        "description": "recovery phase 的 queue 優先於所有其他 queue",
        "rule": "當 current_phase=recovery 時，recovery queue 的項目提升至 P0",
        "trigger": "phase=recovery"
      }
    ],
    "queue_examples": [
      {
        "phase": "execution",
        "initial_queue": [
          {
            "id": "queue.execution.gate.01",
            "type": "gate",
            "ref": "gate.phase_execution.all_obligations_met",
            "priority": "P0",
            "status": "pending",
            "depends_on": [
  
            ],
            "estimated_cost": {
              "tokens": 100,
              "steps": 1
            }
          },
          {
            "id": "queue.execution.obligation.01",
            "type": "obligation",
            "ref": "obl.execution.log_actions",
            "priority": "P2",
            "status": "pending",
            "depends_on": [
  
            ],
            "estimated_cost": {
              "tokens": 50,
              "steps": 1
            }
          },
          {
            "id": "queue.execution.obligation.02",
            "type": "obligation",
            "ref": "obl.execution.track_artifacts",
            "priority": "P2",
            "status": "pending",
            "depends_on": [
  
            ],
            "estimated_cost": {
              "tokens": 80,
              "steps": 2
            }
          }
        ],
        "execution_order": [
          "queue.execution.gate.01 (P0, no dependencies)",
          "queue.execution.obligation.01 (P2, no dependencies)",
          "queue.execution.obligation.02 (P2, no dependencies)"
        ]
      },
      {
        "phase": "validation",
        "initial_queue": [
          {
            "id": "queue.validation.gate.01",
            "type": "gate",
            "ref": "gate.phase_validation.linked_updates_complete",
            "priority": "P0",
            "status": "pending",
            "depends_on": [
  
            ],
            "estimated_cost": {
              "tokens": 100,
              "steps": 1
            }
          },
          {
            "id": "queue.validation.gate.02",
            "type": "gate",
            "ref": "gate.phase_validation.artifacts_complete",
            "priority": "P0",
            "status": "pending",
            "depends_on": [
  
            ],
            "estimated_cost": {
              "tokens": 100,
              "steps": 1
            }
          },
          {
            "id": "queue.validation.obligation.01",
            "type": "obligation",
            "ref": "obl.validation.run_knowledge_runtime_validator",
            "priority": "P1",
            "status": "pending",
            "depends_on": [
  
            ],
            "estimated_cost": {
              "tokens": 200,
              "steps": 3
            }
          },
          {
            "id": "queue.validation.obligation.02",
            "type": "obligation",
            "ref": "obl.validation.check_linked_updates",
            "priority": "P1",
            "status": "pending",
            "depends_on": [
              "queue.validation.gate.01"
            ],
            "estimated_cost": {
              "tokens": 150,
              "steps": 2
            }
          }
        ],
        "execution_order": [
          "queue.validation.gate.01 (P0, no dependencies)",
          "queue.validation.gate.02 (P0, no dependencies)",
          "queue.validation.obligation.01 (P1, no dependencies)",
          "queue.validation.obligation.02 (P1, depends on gate.01)"
        ]
      },
      {
        "phase": "recovery",
        "initial_queue": [
          {
            "id": "queue.recovery.action.01",
            "type": "action",
            "ref": "strategy.unknown_failure",
            "priority": "P0",
            "status": "pending",
            "depends_on": [
  
            ],
            "estimated_cost": {
              "tokens": 300,
              "steps": 4
            }
          },
          {
            "id": "queue.recovery.repair.01",
            "type": "repair",
            "ref": "repair.state.current_phase_corrupted",
            "priority": "P0",
            "status": "pending",
            "depends_on": [
              "queue.recovery.action.01"
            ],
            "estimated_cost": {
              "tokens": 200,
              "steps": 4
            }
          },
          {
            "id": "queue.recovery.reconciliation.01",
            "type": "reconciliation",
            "ref": "reconcile.phase.full_state_reconciliation",
            "priority": "P1",
            "status": "pending",
            "depends_on": [
              "queue.recovery.repair.01"
            ],
            "estimated_cost": {
              "tokens": 250,
              "steps": 5
            }
          }
        ],
        "execution_order": [
          "queue.recovery.action.01 (P0, no dependencies)",
          "queue.recovery.repair.01 (P0, depends on action.01)",
          "queue.recovery.reconciliation.01 (P1, depends on repair.01)"
        ]
      }
    ]
  }

  # Source: runtime/scheduler/priority-scheduler.yaml
  SCHEDULER_PRIORITY_SCHEDULER =   {
    "version": 1.0,
    "status": "candidate",
    "owner_layer": "runtime/scheduler",
    "last_updated": "2026-05-15",
    "scheduling_policy": {
      "name": "Default Priority Policy",
      "description": "預設排程策略：先依 priority 降序，再依 deadline 升序， 最後依 blocking_dependencies 決定拓撲順序。\n",
      "sort_order": [
        {
          "field": "priority",
          "direction": "desc"
        },
        {
          "field": "deadline",
          "direction": "asc"
        },
        {
          "field": "blocking_dependencies",
          "direction": "topological"
        }
      ]
    },
    "priority_levels": [
      {
        "level": "P0",
        "name": "Critical",
        "description": "必須立即執行，阻斷所有其他工作",
        "color": "red",
        "examples": [
          "blocking gate severity=critical 的檢查",
          "recovery phase 的修復程序",
          "security 相關的 obligation"
        ],
        "preempts": true
      },
      {
        "level": "P1",
        "name": "High",
        "description": "應在目前 phase 內完成",
        "color": "orange",
        "examples": [
          "blocking gate severity=high 的檢查",
          "phase transition 前的 obligations",
          "commit/push 相關的 obligations"
        ],
        "preempts": false
      },
      {
        "level": "P2",
        "name": "Medium",
        "description": "應在目前 session 內完成",
        "color": "yellow",
        "examples": [
          "knowledge update checkpoint",
          "非 blocking 的 obligations",
          "artifacts 整理"
        ],
        "preempts": false
      },
      {
        "level": "P3",
        "name": "Low",
        "description": "可在有空時執行",
        "color": "green",
        "examples": [
          "非緊急的 linked updates",
          "選擇性的 validation",
          "文件整理"
        ],
        "preempts": false
      }
    ],
    "scheduling_rules": [
      {
        "id": "sched.rule.blocking_first",
        "name": "Blocking Gates First",
        "description": "所有 blocking gates 必須在 obligations 之前檢查",
        "rule": "先執行 blocking-gates.yaml 中 severity=critical 或 severity=high 的 gates，再執行 obligations",
        "priority": "P0",
        "applies_to": [
          "phase entry",
          "phase transition"
        ]
      },
      {
        "id": "sched.rule.dependency_before_dependent",
        "name": "Dependency Before Dependent",
        "description": "有 depends_on 的 obligation 必須在上游完成後才能執行",
        "rule": "依 obligation-ledger.yaml 的 depends_on 建立拓撲排序",
        "priority": "P0",
        "applies_to": [
          "obligation execution"
        ]
      },
      {
        "id": "sched.rule.deadline_aware",
        "name": "Deadline-Aware Scheduling",
        "description": "接近 deadline 的項目優先執行",
        "rule": "deadline 在 5 步驟內的項目提升一級 priority",
        "priority": "P1",
        "applies_to": [
          "all scheduled items"
        ]
      },
      {
        "id": "sched.rule.phase_transition_gate",
        "name": "Phase Transition Gate",
        "description": "phase transition 前必須先通過所有 blocking gates",
        "rule": "transition 前執行 blocking-gates.yaml 中目標 phase 的所有 gates",
        "priority": "P0",
        "applies_to": [
          "phase transition"
        ]
      },
      {
        "id": "sched.rule.recovery_preempts_all",
        "name": "Recovery Preempts All",
        "description": "recovery phase 的工作優先於所有其他 phase 的工作",
        "rule": "當 current_phase=recovery 時，所有 recovery 相關工作提升至 P0",
        "priority": "P0",
        "applies_to": [
          "recovery phase"
        ]
      },
      {
        "id": "sched.rule.token_budget_aware",
        "name": "Token Budget Awareness",
        "description": "在 token 預算內優先執行高價值項目",
        "rule": "當 token 使用率 > 70% 時，只執行 P0 和 P1 項目",
        "priority": "P2",
        "applies_to": [
          "token management"
        ]
      }
    ],
    "scheduling_examples": [
      {
        "scenario": "Phase: execution → validation transition",
        "items": [
          {
            "name": "Check all_obligations_met gate",
            "priority": "P0",
            "reason": "blocking gate severity=critical"
          },
          {
            "name": "Check artifacts_complete gate",
            "priority": "P0",
            "reason": "blocking gate severity=high"
          },
          {
            "name": "Run knowledge runtime validator",
            "priority": "P1",
            "reason": "phase obligation"
          },
          {
            "name": "Check linked updates",
            "priority": "P1",
            "reason": "phase obligation"
          },
          {
            "name": "Log actions",
            "priority": "P2",
            "reason": "non-blocking obligation"
          }
        ],
        "execution_order": [
          "Check all_obligations_met gate",
          "Check artifacts_complete gate",
          "Run knowledge runtime validator",
          "Check linked updates",
          "Log actions"
        ]
      },
      {
        "scenario": "Recovery phase entry",
        "items": [
          {
            "name": "Diagnose root cause",
            "priority": "P0",
            "reason": "recovery preempts all"
          },
          {
            "name": "Apply repair",
            "priority": "P0",
            "reason": "recovery preempts all"
          },
          {
            "name": "Verify repair",
            "priority": "P0",
            "reason": "recovery preempts all"
          },
          {
            "name": "Run full state reconciliation",
            "priority": "P1",
            "reason": "post-recovery verification"
          }
        ],
        "execution_order": [
          "Diagnose root cause",
          "Apply repair",
          "Verify repair",
          "Run full state reconciliation"
        ]
      }
    ]
  }

  # Source: runtime/transactions/transaction-machine.yaml
  TRANSACTIONS_TRANSACTION_MACHINE =   {
    "version": 1.0,
    "status": "candidate",
    "owner_layer": "runtime/transactions",
    "last_updated": "2026-05-15",
    "transaction_states": [
      {
        "state": "closed",
        "description": "無 active transaction，可開始新 transaction",
        "entry_conditions": [
  
        ],
        "allowed_actions": [
          "START_TRANSACTION"
        ],
        "transitions": [
          {
            "to": "open",
            "trigger": "agent 第一次寫入 Ai-skill 檔案",
            "condition": "寫入路徑在 <AI_SKILL_REPO> 內"
          }
        ]
      },
      {
        "state": "open",
        "description": "transaction 已開始，正在累積變更",
        "entry_conditions": [
          "canonical repo 已定位（git rev-parse --show-toplevel）",
          "touched paths 已記錄"
        ],
        "allowed_actions": [
          "WRITE_FILE",
          "MODIFY_FILE",
          "DELETE_FILE",
          "CHECK_GIT_STATUS",
          "RECORD_TOUCHED_PATH"
        ],
        "forbidden_actions": [
          "COMMIT",
          "PUSH",
          "CLOSE_TRANSACTION"
        ],
        "blocking_gates": [
          "gate.phase_commit.staging_complete"
        ],
        "transitions": [
          {
            "to": "staging",
            "trigger": "所有變更已完成，準備 stage",
            "condition": "所有 touched paths 已記錄"
          }
        ]
      },
      {
        "state": "staging",
        "description": "變更已準備 stage，正在進行去敏檢查",
        "entry_conditions": [
          "所有 touched paths 已記錄"
        ],
        "allowed_actions": [
          "GIT_ADD",
          "RUN_SANITIZATION_CHECK",
          "CHECK_LINKED_UPDATES"
        ],
        "forbidden_actions": [
          "WRITE_NEW_FILES",
          "CLOSE_TRANSACTION"
        ],
        "blocking_gates": [
          "gate.phase_commit.staging_complete"
        ],
        "transitions": [
          {
            "to": "commit_ready",
            "trigger": "去敏檢查通過",
            "condition": "sanitization.md 檢查無問題"
          },
          {
            "to": "open",
            "trigger": "去敏檢查失敗",
            "condition": "發現需修正的內容"
          }
        ]
      },
      {
        "state": "commit_ready",
        "description": "變更已 stage，準備 commit",
        "entry_conditions": [
          "git add 已完成",
          "去敏檢查已通過",
          "linked updates 已檢查"
        ],
        "allowed_actions": [
          "GIT_COMMIT",
          "WRITE_COMMIT_MESSAGE",
          "VERIFY_PRE_COMMIT_HOOKS"
        ],
        "forbidden_actions": [
          "WRITE_NEW_FILES",
          "MODIFY_STAGED_FILES"
        ],
        "blocking_gates": [
          "gate.phase_commit.commit_message_valid",
          "gate.phase_commit.pre_commit_hooks_passed"
        ],
        "transitions": [
          {
            "to": "committed",
            "trigger": "commit 成功",
            "condition": "git commit 成功，pre-commit hooks 通過"
          },
          {
            "to": "staging",
            "trigger": "commit 失敗",
            "condition": "pre-commit hooks 失敗或 commit rejected"
          }
        ]
      },
      {
        "state": "committed",
        "description": "變更已 commit，準備 push",
        "entry_conditions": [
          "git commit 成功",
          "commit message 符合規範"
        ],
        "allowed_actions": [
          "GIT_PUSH",
          "CHECK_REMOTE"
        ],
        "forbidden_actions": [
          "MODIFY_COMMITTED_FILES"
        ],
        "blocking_gates": [
          "gate.phase_push.remote_reachable"
        ],
        "transitions": [
          {
            "to": "pushed",
            "trigger": "push 成功",
            "condition": "git push 成功"
          },
          {
            "to": "commit_ready",
            "trigger": "push 失敗",
            "condition": "remote rejected 或 network error"
          }
        ]
      },
      {
        "state": "pushed",
        "description": "變更已 push，準備讀回確認",
        "entry_conditions": [
          "git push 成功"
        ],
        "allowed_actions": [
          "VERIFY_CLEAN_TREE",
          "CHECK_REMOTE_SYNC",
          "REREAD_UPDATED_FILES"
        ],
        "forbidden_actions": [
          "NEW_WRITES"
        ],
        "blocking_gates": [
          "gate.phase_readback.working_tree_clean",
          "gate.phase_readback.remote_synced"
        ],
        "transitions": [
          {
            "to": "verified",
            "trigger": "讀回確認完成",
            "condition": "git status clean，branch 無 ahead/behind"
          },
          {
            "to": "committed",
            "trigger": "讀回發現不一致",
            "condition": "working tree dirty 或 branch ahead/behind"
          }
        ]
      },
      {
        "state": "verified",
        "description": "transaction 已驗證完成，可關閉",
        "entry_conditions": [
          "git status clean",
          "branch 無 ahead/behind",
          "所有必要檔案已 reread"
        ],
        "allowed_actions": [
          "CLOSE_TRANSACTION",
          "UPDATE_OBLIGATION_STATE"
        ],
        "transitions": [
          {
            "to": "closed",
            "trigger": "CLOSE_TRANSACTION",
            "condition": "所有 close conditions 已滿足"
          }
        ]
      }
    ],
    "transaction_rules": [
      {
        "id": "tx.rule.lock_check",
        "name": "Transaction Lock Check",
        "description": "開始 transaction 前必須檢查是否有 active lock",
        "rule": "檢查是否有其他 agent/user 正在操作；若有，不得自動 commit/push",
        "severity": "critical"
      },
      {
        "id": "tx.rule.canonical_first",
        "name": "Canonical Source First",
        "description": "必須先在 canonical repo 修改，不得先改 mirror copy",
        "rule": "禁止先修改工具部署/mirror 路徑再回頭找 source",
        "severity": "critical"
      },
      {
        "id": "tx.rule.sanitize_before_commit",
        "name": "Sanitize Before Commit",
        "description": "commit 前必須執行去敏檢查",
        "rule": "依 sanitization.md 檢查所有新增/修改的可重用文件",
        "severity": "high"
      },
      {
        "id": "tx.rule.linked_updates_check",
        "name": "Linked Updates Check",
        "description": "commit 前必須檢查必要的 linked updates",
        "rule": "必要的 linked updates 已同步或明確寫出不適用理由",
        "severity": "high"
      },
      {
        "id": "tx.rule.reread_after_push",
        "name": "Reread After Push",
        "description": "push 後必須重新讀取更新過的入口、主要依賴、索引與 promotion target",
        "rule": "push 完成後重新讀取 CORE_BOOTSTRAP.md、routing-registry.yaml",
        "severity": "high"
      },
      {
        "id": "tx.rule.final_clean_check",
        "name": "Final Clean Check",
        "description": "transaction 關閉前必須確認 git status clean",
        "rule": "最後一次 git status --short --branch 顯示 clean，branch 無 ahead/behind",
        "severity": "critical"
      }
    ],
    "transaction_templates": [
      {
        "id": "tx.template.skill_update",
        "name": "Skill Update Transaction",
        "description": "修改 skill 或 enforcement rule 時的標準 transaction",
        "typical_steps": [
          "在 canonical repo 修改 source",
          "記錄 touched paths",
          "執行去敏檢查",
          "檢查 linked updates",
          "git add",
          "git commit",
          "git push",
          "重新讀取更新檔案",
          "確認 git status clean",
          "關閉 transaction"
        ]
      },
      {
        "id": "tx.template.new_skill",
        "name": "New Skill Creation Transaction",
        "description": "建立新 skill 時的標準 transaction",
        "typical_steps": [
          "在 canonical repo 建立新檔案",
          "更新 routing-registry.yaml（如需要）",
          "記錄 touched paths",
          "執行去敏檢查",
          "檢查 linked updates",
          "git add",
          "git commit",
          "git push",
          "重新讀取更新檔案",
          "確認 git status clean",
          "關閉 transaction"
        ]
      },
      {
        "id": "tx.template.feedback_lesson",
        "name": "Feedback Lesson Transaction",
        "description": "寫入 feedback lesson 時的標準 transaction",
        "typical_steps": [
          "在 feedback/history/ 下建立 lesson",
          "執行去敏檢查（移除 project-specific evidence）",
          "更新 promotion target（如需要）",
          "記錄 touched paths",
          "git add",
          "git commit",
          "git push",
          "重新讀取更新檔案",
          "確認 git status clean",
          "關閉 transaction"
        ]
      }
    ]
  }

  # Source: runtime/transactions/transaction-templates.yaml
  TRANSACTIONS_TRANSACTION_TEMPLATES =   {
    "version": 1.0,
    "status": "candidate",
    "owner_layer": "runtime/transactions",
    "last_updated": "2026-05-15",
    "templates": [
      {
        "id": "template.enforcement_rule_update",
        "name": "Enforcement Rule Update",
        "description": "修改 enforcement/ 下的 enforcement rule",
        "applicable_scenarios": [
          "修改 dependency-reading.md",
          "修改 goal-action-validation.md",
          "修改 failure-learning-system.md",
          "新增 enforcement/failure-patterns/ 下的 pattern"
        ],
        "pre_checks": [
          "確認修改不違反 rule-weight 的 P0/P1/P2/P3 權重體系",
          "確認修改不破壞 reference-first 原則",
          "確認修改不破壞 activation model"
        ],
        "steps": [
          {
            "step": 1,
            "action": "LOCATE_CANONICAL_SOURCE",
            "description": "確認 canonical repo 路徑",
            "command": "git rev-parse --show-toplevel"
          },
          {
            "step": 2,
            "action": "MAKE_CHANGES",
            "description": "在 canonical source 中修改",
            "command": "編輯 enforcement/ 下對應檔案"
          },
          {
            "step": 3,
            "action": "CHECK_LINKED_UPDATES",
            "description": "檢查需要連動更新的檔案",
            "command": "依 enforcement/linked-updates.md 檢查"
          },
          {
            "step": 4,
            "action": "RUN_SANITIZATION",
            "description": "執行去敏檢查",
            "command": "依 enforcement/sanitization.md 檢查"
          },
          {
            "step": 5,
            "action": "UPDATE_ROUTING_REGISTRY",
            "description": "更新 routing-registry.yaml（如需要）",
            "command": "檢查是否需要更新 knowledge/runtime/routing-registry.yaml"
          },
          {
            "step": 6,
            "action": "GIT_ADD",
            "description": "stage 所有變更",
            "command": "git add -A"
          },
          {
            "step": 7,
            "action": "GIT_COMMIT",
            "description": "提交變更",
            "command": "git commit -m 'rule: ...'"
          },
          {
            "step": 8,
            "action": "GIT_PUSH",
            "description": "推送至遠端",
            "command": "git push"
          },
          {
            "step": 9,
            "action": "REREAD",
            "description": "重新讀取更新檔案",
            "command": "讀取 CORE_BOOTSTRAP.md、routing-registry.yaml"
          },
          {
            "step": 10,
            "action": "VERIFY_CLEAN",
            "description": "確認 git status clean",
            "command": "git status --short --branch"
          }
        ],
        "post_checks": [
          "git status clean",
          "branch 無 ahead/behind",
          "所有 linked updates 已完成"
        ]
      },
      {
        "id": "template.runtime_component_update",
        "name": "Runtime Component Update",
        "description": "新增或修改 runtime/ 下的元件",
        "applicable_scenarios": [
          "新增 phase-machine.yaml",
          "新增 obligation-ledger.yaml",
          "新增 blocking-gates.yaml",
          "新增 recovery/ 下的檔案",
          "新增 scheduler/ 下的檔案",
          "新增 transactions/ 下的檔案"
        ],
        "pre_checks": [
          "確認新元件與現有 phase-machine.yaml 的 phase 定義一致",
          "確認新元件的 id 不與現有 routing-registry.yaml 衝突",
          "確認新元件遵循 generated YAML 格式（如適用）"
        ],
        "steps": [
          {
            "step": 1,
            "action": "CREATE_COMPONENT",
            "description": "建立新 runtime 元件",
            "command": "在 runtime/ 下建立對應 YAML 檔案"
          },
          {
            "step": 2,
            "action": "UPDATE_ROUTING_REGISTRY",
            "description": "在 routing-registry.yaml 新增 route record",
            "command": "新增 route.runtime.{component} 記錄"
          },
          {
            "step": 3,
            "action": "UPDATE_CORE_BOOTSTRAP",
            "description": "更新 CORE_BOOTSTRAP.md（如需要）",
            "command": "如果元件需在 bootstrap 階段載入，更新啟動流程"
          },
          {
            "step": 4,
            "action": "CHECK_INTEGRATION",
            "description": "檢查與現有元件的整合",
            "command": "確認新元件的 references 可解析"
          },
          {
            "step": 5,
            "action": "GIT_ADD",
            "description": "stage 所有變更",
            "command": "git add -A"
          },
          {
            "step": 6,
            "action": "GIT_COMMIT",
            "description": "提交變更",
            "command": "git commit -m 'feat(runtime): ...'"
          },
          {
            "step": 7,
            "action": "GIT_PUSH",
            "description": "推送至遠端",
            "command": "git push"
          },
          {
            "step": 8,
            "action": "REREAD",
            "description": "重新讀取更新檔案",
            "command": "讀取 CORE_BOOTSTRAP.md、routing-registry.yaml"
          },
          {
            "step": 9,
            "action": "VERIFY_CLEAN",
            "description": "確認 git status clean",
            "command": "git status --short --branch"
          }
        ],
        "post_checks": [
          "git status clean",
          "branch 無 ahead/behind",
          "routing-registry.yaml 包含新 route record",
          "CORE_BOOTSTRAP.md 已更新（如需要）"
        ]
      },
      {
        "id": "template.workflow_update",
        "name": "Workflow Update",
        "description": "修改 workflow/ 下的 prose 文件",
        "applicable_scenarios": [
          "修改 workflow/apk-analysis/execution-flow.md",
          "修改 workflow/software-delivery/development-process.md",
          "修改 workflow/travel-planning/execution-flow.md"
        ],
        "pre_checks": [
          "確認修改不破壞 workflow 的 reference-first 原則",
          "確認修改與對應的 artifact-gates.md 一致"
        ],
        "steps": [
          {
            "step": 1,
            "action": "LOCATE_CANONICAL_SOURCE",
            "description": "確認 canonical repo 路徑",
            "command": "git rev-parse --show-toplevel"
          },
          {
            "step": 2,
            "action": "MAKE_CHANGES",
            "description": "修改 workflow prose",
            "command": "編輯 workflow/ 下對應檔案"
          },
          {
            "step": 3,
            "action": "CHECK_GENERATED_SURFACES",
            "description": "檢查是否有對應的 generated surface（runtime.db）需要更新",
            "command": "查詢 runtime.db 確認對應的 generated_surfaces 記錄是否存在"
          },
          {
            "step": 4,
            "action": "UPDATE_GENERATED_YAML",
            "description": "更新 generated YAML（如存在）",
            "command": "執行 compiler 或手動同步"
          },
          {
            "step": 5,
            "action": "GIT_ADD",
            "description": "stage 所有變更",
            "command": "git add -A"
          },
          {
            "step": 6,
            "action": "GIT_COMMIT",
            "description": "提交變更",
            "command": "git commit -m 'workflow: ...'"
          },
          {
            "step": 7,
            "action": "GIT_PUSH",
            "description": "推送至遠端",
            "command": "git push"
          },
          {
            "step": 8,
            "action": "REREAD",
            "description": "重新讀取更新檔案",
            "command": "讀取更新後的 workflow 檔案"
          },
          {
            "step": 9,
            "action": "VERIFY_CLEAN",
            "description": "確認 git status clean",
            "command": "git status --short --branch"
          }
        ],
        "post_checks": [
          "git status clean",
          "branch 無 ahead/behind",
          "prose 與 generated YAML 一致（如適用）"
        ]
      },
      {
        "id": "template.feedback_lesson_write",
        "name": "Feedback Lesson Write",
        "description": "寫入 feedback lesson 到 feedback/history/",
        "applicable_scenarios": [
          "記錄新的 feedback lesson",
          "更新現有 feedback lesson"
        ],
        "pre_checks": [
          "確認 lesson 已去敏（無 project-specific evidence、絕對路徑、使用者名稱）",
          "確認 lesson 已分類（domain/category）"
        ],
        "steps": [
          {
            "step": 1,
            "action": "DETERMINE_PATH",
            "description": "決定 lesson 存放路徑",
            "command": "feedback/history/<domain>/<category>/<date>-<title>.md"
          },
          {
            "step": 2,
            "action": "WRITE_LESSON",
            "description": "寫入 lesson 內容",
            "command": "依 feedback-lessons.md 格式撰寫"
          },
          {
            "step": 3,
            "action": "RUN_SANITIZATION",
            "description": "執行去敏檢查",
            "command": "依 sanitization.md 檢查"
          },
          {
            "step": 4,
            "action": "CHECK_PROMOTION",
            "description": "檢查是否需要 promotion",
            "command": "依 intelligence-extraction-pipeline.md 判斷"
          },
          {
            "step": 5,
            "action": "GIT_ADD",
            "description": "stage 所有變更",
            "command": "git add -A"
          },
          {
            "step": 6,
            "action": "GIT_COMMIT",
            "description": "提交變更",
            "command": "git commit -m 'feedback: ...'"
          },
          {
            "step": 7,
            "action": "GIT_PUSH",
            "description": "推送至遠端",
            "command": "git push"
          },
          {
            "step": 8,
            "action": "REREAD",
            "description": "重新讀取更新檔案",
            "command": "讀取 feedback/README.md、promotion targets"
          },
          {
            "step": 9,
            "action": "VERIFY_CLEAN",
            "description": "確認 git status clean",
            "command": "git status --short --branch"
          }
        ],
        "post_checks": [
          "git status clean",
          "branch 無 ahead/behind",
          "lesson 已去敏",
          "promotion target 已更新（如需要）"
        ]
      },
      {
        "id": "template.multi_file_change",
        "name": "Multi-File Change",
        "description": "同時修改多個層級的檔案（例如：修改 workflow 同時更新 enforcement）",
        "applicable_scenarios": [
          "修改 workflow 同時更新對應的 enforcement rules",
          "新增 skill 同時更新 routing-registry.yaml",
          "重構目錄結構同時更新所有 cross-references"
        ],
        "pre_checks": [
          "列出所有受影響的檔案",
          "確認變更範圍不超出任務 scope",
          "確認無 circular dependency between changes"
        ],
        "steps": [
          {
            "step": 1,
            "action": "LIST_ALL_FILES",
            "description": "列出所有需要修改的檔案",
            "command": "依 linked-updates.md 找出所有連動檔案"
          },
          {
            "step": 2,
            "action": "ORDER_CHANGES",
            "description": "依相依關係排序修改順序",
            "command": "先修改 canonical source，再修改 references"
          },
          {
            "step": 3,
            "action": "MAKE_ALL_CHANGES",
            "description": "依序修改所有檔案",
            "command": "從相依關係的 root 開始修改"
          },
          {
            "step": 4,
            "action": "VERIFY_CROSS_REFERENCES",
            "description": "驗證 cross-references 可解析",
            "command": "檢查所有檔案間的 link 是否正確"
          },
          {
            "step": 5,
            "action": "RUN_SANITIZATION",
            "description": "執行去敏檢查",
            "command": "依 sanitization.md 檢查所有修改檔案"
          },
          {
            "step": 6,
            "action": "GIT_ADD",
            "description": "stage 所有變更",
            "command": "git add -A"
          },
          {
            "step": 7,
            "action": "GIT_COMMIT",
            "description": "提交變更",
            "command": "git commit -m '...'"
          },
          {
            "step": 8,
            "action": "GIT_PUSH",
            "description": "推送至遠端",
            "command": "git push"
          },
          {
            "step": 9,
            "action": "REREAD",
            "description": "重新讀取所有更新檔案",
            "command": "讀取所有修改過的主要檔案"
          },
          {
            "step": 10,
            "action": "VERIFY_CLEAN",
            "description": "確認 git status clean",
            "command": "git status --short --branch"
          }
        ],
        "post_checks": [
          "git status clean",
          "branch 無 ahead/behind",
          "所有 cross-references 可解析",
          "所有 linked updates 已完成"
        ]
      }
    ]
  }

end