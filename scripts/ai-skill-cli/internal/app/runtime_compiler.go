package app

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const goRuntimeCompilerVersion = "2.0.0"

type compilerMappingFile struct {
	SourceTargetMapping []compilerMapping `yaml:"source_target_mapping" json:"source_target_mapping"`
	CompilationRules    []map[string]any  `yaml:"compilation_rules" json:"compilation_rules"`
	CompilerWorkflow    map[string]any    `yaml:"compiler_workflow" json:"compiler_workflow"`
}

type compilerMapping struct {
	Source      string `yaml:"source" json:"source"`
	Target      string `yaml:"target" json:"target"`
	CompileRule string `yaml:"compile_rule" json:"compile_rule"`
}

type runtimeConfigMapping struct {
	rel     string
	table   string
	idCol   string
	listKey string
	idKeys  []string
}

func buildNativeRuntimeDBFromSources(repo string, outputDB string) Check {
	docs, err := loadRuntimeCanonicalDocuments(repo)
	if err != nil {
		return Check{Name: "runtime_compile_native", Status: "failed", Message: err.Error()}
	}
	tempDB := outputDB
	if filepath.Clean(outputDB) == filepath.Clean(filepath.Join(repo, "runtime", "runtime.db")) {
		tempDB = outputDB + ".tmp"
	}
	_ = os.Remove(tempDB)
	db, err := sql.Open("sqlite", tempDB)
	if err != nil {
		return Check{Name: "runtime_compile_native", Status: "failed", Message: err.Error()}
	}
	defer db.Close()
	if err := createGoRuntimeSchema(db); err != nil {
		return Check{Name: "runtime_compile_native", Status: "failed", Message: err.Error()}
	}
	if err := insertRuntimeConfigDocuments(db, docs); err != nil {
		return Check{Name: "runtime_compile_native", Status: "failed", Message: err.Error()}
	}
	if err := compileStructuredRuntimeSources(repo, db, docs); err != nil {
		return Check{Name: "runtime_compile_native", Status: "failed", Message: err.Error()}
	}
	if err := compileProseRuntimeSources(repo, db, docs); err != nil {
		return Check{Name: "runtime_compile_native", Status: "failed", Message: err.Error()}
	}
	if _, err := db.Exec("VACUUM"); err != nil {
		return Check{Name: "runtime_compile_native", Status: "failed", Message: err.Error()}
	}
	if err := db.Close(); err != nil {
		return Check{Name: "runtime_compile_native", Status: "failed", Message: err.Error()}
	}
	if tempDB != outputDB {
		if err := os.Rename(tempDB, outputDB); err != nil {
			return Check{Name: "runtime_compile_native", Status: "failed", Message: err.Error()}
		}
	}
	check := nativeRuntimeDBValidation(outputDB)
	if check.Status != "ok" {
		check.Name = "runtime_compile_native"
		return check
	}
	return Check{Name: "runtime_compile_native", Status: "ok", Message: outputDB}
}

func createGoRuntimeSchema(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE phases (id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT, entry_conditions TEXT, allowed_actions TEXT, forbidden_actions TEXT, blocking_gates TEXT, obligations TEXT, next_phase TEXT, phase_transition_triggers TEXT, metadata TEXT, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE phase_transitions (id INTEGER PRIMARY KEY AUTOINCREMENT, from_phase TEXT NOT NULL, to_phase TEXT NOT NULL, trigger TEXT, condition TEXT, rule_type TEXT DEFAULT 'normal', description TEXT);`,
		`CREATE TABLE obligations (id TEXT PRIMARY KEY, phase TEXT NOT NULL, name TEXT NOT NULL, description TEXT, verification TEXT, severity TEXT DEFAULT 'high', depends_on TEXT, linked_gates TEXT, metadata TEXT, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE gates (id TEXT PRIMARY KEY, phase TEXT NOT NULL, name TEXT NOT NULL, description TEXT, severity TEXT DEFAULT 'high', check_type TEXT, check_target TEXT, check_verification TEXT, failure_action TEXT, failure_message TEXT, metadata TEXT, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE transaction_states (state TEXT PRIMARY KEY, description TEXT, entry_conditions TEXT, allowed_actions TEXT, forbidden_actions TEXT, blocking_gates TEXT, metadata TEXT, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE transaction_transitions (id INTEGER PRIMARY KEY AUTOINCREMENT, from_state TEXT NOT NULL, to_state TEXT NOT NULL, trigger TEXT, condition TEXT, description TEXT);`,
		`CREATE TABLE transaction_rules (id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT, rule TEXT, severity TEXT DEFAULT 'high');`,
		`CREATE TABLE transaction_templates (id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT, typical_steps TEXT, content TEXT DEFAULT '{}');`,
		`CREATE TABLE compiler_rules (id INTEGER PRIMARY KEY AUTOINCREMENT, rule_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE core_bootstrap_rules (rule_id TEXT PRIMARY KEY, ordinal INTEGER DEFAULT 0);`,
		`CREATE TABLE discovery_checkpoints (id INTEGER PRIMARY KEY AUTOINCREMENT, phase TEXT NOT NULL, trigger TEXT NOT NULL, description TEXT, discovery_targets TEXT, metadata TEXT);`,
		`CREATE TABLE discovery_search_strategy (id INTEGER PRIMARY KEY AUTOINCREMENT, priority_order TEXT, fallback TEXT, min_confidence_threshold TEXT);`,
		`CREATE TABLE decision_recording (id INTEGER PRIMARY KEY AUTOINCREMENT, section TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE runtime_config_documents (logical_id TEXT PRIMARY KEY, owner_layer TEXT, status TEXT, schema_version TEXT, content_json TEXT NOT NULL, checksum TEXT NOT NULL, updated_at TEXT NOT NULL DEFAULT (datetime('now')));`,
		`CREATE TABLE runtime_config_projections (id INTEGER PRIMARY KEY AUTOINCREMENT, logical_id TEXT NOT NULL, target_table TEXT NOT NULL, row_key TEXT NOT NULL, checksum TEXT NOT NULL, updated_at TEXT NOT NULL DEFAULT (datetime('now')), UNIQUE(logical_id, target_table, row_key));`,
		`CREATE TABLE runtime_source_files (source_path TEXT PRIMARY KEY, source_kind TEXT NOT NULL DEFAULT 'db', target_table TEXT NOT NULL, compile_rule TEXT NOT NULL, compiled_at TEXT NOT NULL, compiler_version TEXT NOT NULL, status TEXT NOT NULL);`,
		`CREATE TABLE compiler_metadata (key TEXT PRIMARY KEY, value TEXT NOT NULL);`,
		`CREATE TABLE generated_surfaces (id INTEGER PRIMARY KEY AUTOINCREMENT, source_path TEXT NOT NULL, target_key TEXT NOT NULL, compile_rule TEXT NOT NULL, compiled_at TEXT NOT NULL, compiler_version TEXT NOT NULL, status TEXT NOT NULL, data TEXT NOT NULL, UNIQUE(source_path, target_key));`,
		`CREATE TABLE repository_topology (subtree TEXT PRIMARY KEY, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE derived_forbidden_tokens (token TEXT NOT NULL, canonical_token TEXT NOT NULL, owning_project_id TEXT NOT NULL, source_metadata_path TEXT NOT NULL, suggested_placeholder TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), PRIMARY KEY (token, owning_project_id, source_metadata_path));`,
		`CREATE TABLE sanitization_patterns (category TEXT PRIMARY KEY, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE runtime_budget (id INTEGER PRIMARY KEY AUTOINCREMENT, model_name TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE context_ttl_policy (id INTEGER PRIMARY KEY AUTOINCREMENT, ttl_type TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE circuit_breaker (id INTEGER PRIMARY KEY AUTOINCREMENT, guard_name TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE context_pollution (id INTEGER PRIMARY KEY AUTOINCREMENT, signal_name TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE context_health_score (id INTEGER PRIMARY KEY AUTOINCREMENT, dimension TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE intelligence_routing (id INTEGER PRIMARY KEY AUTOINCREMENT, rule_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE obligation_ledger (id INTEGER PRIMARY KEY AUTOINCREMENT, obligation_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE language_policy (id INTEGER PRIMARY KEY AUTOINCREMENT, section TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE output_rules (id INTEGER PRIMARY KEY AUTOINCREMENT, section TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE governance_gates (id INTEGER PRIMARY KEY AUTOINCREMENT, gate_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE blocking_gates (id INTEGER PRIMARY KEY AUTOINCREMENT, gate_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE phase_machine (id INTEGER PRIMARY KEY AUTOINCREMENT, phase_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE pipeline_context_flow (id INTEGER PRIMARY KEY AUTOINCREMENT, level TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE guard_chain (id INTEGER PRIMARY KEY AUTOINCREMENT, stage TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE relevance_engine (id INTEGER PRIMARY KEY AUTOINCREMENT, component TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE session_lifecycle (id INTEGER PRIMARY KEY AUTOINCREMENT, stage TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE prompt_artifact_templates (id INTEGER PRIMARY KEY AUTOINCREMENT, template_name TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE prompt_composition_rules (id INTEGER PRIMARY KEY AUTOINCREMENT, rule_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE recovery_strategies (id INTEGER PRIMARY KEY AUTOINCREMENT, strategy_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE state_repair (id INTEGER PRIMARY KEY AUTOINCREMENT, procedure_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE obligation_rebuild (id INTEGER PRIMARY KEY AUTOINCREMENT, procedure_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE phase_reconciliation (id INTEGER PRIMARY KEY AUTOINCREMENT, procedure_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE execution_queue (id INTEGER PRIMARY KEY AUTOINCREMENT, queue_name TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE priority_scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, priority_level TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE transaction_templates_ext (id INTEGER PRIMARY KEY AUTOINCREMENT, template_name TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE distributed_locks (id INTEGER PRIMARY KEY AUTOINCREMENT, lock_name TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE multi_agent_coordination (id INTEGER PRIMARY KEY AUTOINCREMENT, rule_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE async_job_lifecycle (id INTEGER PRIMARY KEY AUTOINCREMENT, state TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE capability_checkpoints (id INTEGER PRIMARY KEY AUTOINCREMENT, checkpoint_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE cognitive_modes (id INTEGER PRIMARY KEY, task_id TEXT, execution_mode TEXT, context_mode TEXT, governance_mode TEXT, memory_mode TEXT, resolved_at TEXT, source TEXT);`,
		`CREATE TABLE discovery_signals (id INTEGER PRIMARY KEY AUTOINCREMENT, signal_name TEXT NOT NULL, signal_type TEXT NOT NULL, pattern TEXT, execution_mode TEXT, context_mode TEXT, governance_mode TEXT, memory_mode TEXT, priority INTEGER DEFAULT 0, description TEXT);`,
		// discovery_proposals: per-task ephemeral state for the Workflow
		// Activation Discovery Bridge (plan 2026-06-06-1700, Phase A). Raw
		// data table, NOT projected from a runtime YAML — config lives in
		// runtime/discovery-bridge.yaml (projected to
		// generated_surfaces[runtime.discovery.config]), but the proposal
		// rows themselves are per-task state with TTL 24h.
		`CREATE TABLE discovery_proposals (id INTEGER PRIMARY KEY AUTOINCREMENT, task_hash TEXT NOT NULL, route_candidates_json TEXT NOT NULL, signal_snapshot_json TEXT NOT NULL, scoring_version TEXT NOT NULL, current_best_confidence REAL NOT NULL, status TEXT NOT NULL, miss_reason TEXT, created_at TEXT NOT NULL, updated_at TEXT NOT NULL, expires_at TEXT NOT NULL);`,
		`CREATE INDEX idx_discovery_proposals_task_hash ON discovery_proposals(task_hash);`,
		`CREATE INDEX idx_discovery_proposals_status ON discovery_proposals(status);`,
		`CREATE INDEX idx_discovery_proposals_expires_at ON discovery_proposals(expires_at);`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}

func compileStructuredRuntimeSources(repo string, db *sql.DB, docs map[string]map[string]any) error {
	phaseMachine, err := runtimeCanonicalDocument(docs, "runtime/phases/phase-machine.yaml")
	if err != nil {
		return err
	}
	for _, phase := range runtimeSliceOfMaps(phaseMachine["phases"]) {
		if _, err := db.Exec(`INSERT INTO phases (id, name, description, entry_conditions, allowed_actions, forbidden_actions, blocking_gates, obligations, next_phase, phase_transition_triggers, metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			runtimeString(phase["id"]), runtimeString(phase["name"]), runtimeString(phase["description"]), runtimeJSON(phase["entry_conditions"]), runtimeJSON(phase["allowed_actions"]), runtimeJSON(phase["forbidden_actions"]), runtimeJSON(phase["blocking_gates"]), runtimeJSON(phase["obligations"]), runtimeString(phase["next_phase"]), runtimeJSON(phase["phase_transition_triggers"]), runtimeJSON(phase["metadata"]),
		); err != nil {
			return err
		}
	}
	for _, transition := range runtimeSliceOfMaps(phaseMachine["phase_transition_rules"]) {
		if _, err := db.Exec(`INSERT INTO phase_transitions (from_phase, to_phase, trigger, condition, rule_type, description) VALUES (?, ?, ?, ?, ?, ?)`,
			runtimeString(transition["from"]), runtimeString(transition["to"]), runtimeString(transition["trigger"]), runtimeString(transition["condition"]), runtimeDefaultString(transition["type"], "normal"), runtimeString(transition["description"]),
		); err != nil {
			return err
		}
	}

	obligationLedger, err := runtimeCanonicalDocument(docs, "runtime/obligations/obligation-ledger.yaml")
	if err != nil {
		return err
	}
	for _, obligation := range runtimeSliceOfMaps(obligationLedger["obligations"]) {
		if _, err := db.Exec(`INSERT INTO obligations (id, phase, name, description, verification, severity, depends_on, linked_gates, metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			runtimeString(obligation["id"]), runtimeString(obligation["phase"]), runtimeString(obligation["name"]), runtimeString(obligation["description"]), runtimeJSON(obligation["verification"]), runtimeDefaultString(obligation["severity"], "high"), runtimeJSON(obligation["depends_on"]), runtimeJSON(obligation["linked_gates"]), runtimeJSON(obligation["metadata"]),
		); err != nil {
			return err
		}
	}

	blockingGates, err := runtimeCanonicalDocument(docs, "runtime/gates/blocking-gates.yaml")
	if err != nil {
		return err
	}
	for _, gate := range runtimeSliceOfMaps(blockingGates["gates"]) {
		check := runtimeMap(gate["check"])
		if _, err := db.Exec(`INSERT INTO gates (id, phase, name, description, severity, check_type, check_target, check_verification, failure_action, failure_message, metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			runtimeString(gate["id"]), runtimeString(gate["phase"]), runtimeString(gate["name"]), runtimeString(gate["description"]), runtimeDefaultString(gate["severity"], "high"), runtimeString(check["type"]), runtimeString(check["target"]), runtimeString(check["verification"]), runtimeString(gate["failure_action"]), runtimeString(gate["failure_message"]), runtimeJSON(gate["metadata"]),
		); err != nil {
			return err
		}
	}

	transactionMachine, err := runtimeCanonicalDocument(docs, "runtime/transactions/transaction-machine.yaml")
	if err != nil {
		return err
	}
	for _, state := range runtimeSliceOfMaps(transactionMachine["transaction_states"]) {
		stateID := runtimeString(state["state"])
		if _, err := db.Exec(`INSERT INTO transaction_states (state, description, entry_conditions, allowed_actions, forbidden_actions, blocking_gates, metadata) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			stateID, runtimeString(state["description"]), runtimeJSON(state["entry_conditions"]), runtimeJSON(state["allowed_actions"]), runtimeJSON(state["forbidden_actions"]), runtimeJSON(state["blocking_gates"]), runtimeJSON(state["metadata"]),
		); err != nil {
			return err
		}
		for _, transition := range runtimeSliceOfMaps(state["transitions"]) {
			if _, err := db.Exec(`INSERT INTO transaction_transitions (from_state, to_state, trigger, condition, description) VALUES (?, ?, ?, ?, ?)`,
				stateID, runtimeString(transition["to"]), runtimeString(transition["trigger"]), runtimeString(transition["condition"]), runtimeString(transition["description"]),
			); err != nil {
				return err
			}
		}
	}
	for _, rule := range runtimeSliceOfMaps(transactionMachine["transaction_rules"]) {
		if _, err := db.Exec(`INSERT INTO transaction_rules (id, name, description, rule, severity) VALUES (?, ?, ?, ?, ?)`,
			runtimeString(rule["id"]), runtimeString(rule["name"]), runtimeString(rule["description"]), runtimeString(rule["rule"]), runtimeDefaultString(rule["severity"], "high"),
		); err != nil {
			return err
		}
	}
	for _, template := range runtimeSliceOfMaps(transactionMachine["transaction_templates"]) {
		if _, err := db.Exec(`INSERT OR REPLACE INTO transaction_templates (id, name, description, typical_steps) VALUES (?, ?, ?, ?)`,
			runtimeString(template["id"]), runtimeString(template["name"]), runtimeString(template["description"]), runtimeJSON(template["typical_steps"]),
		); err != nil {
			return err
		}
	}

	activationRules, err := runtimeCanonicalDocument(docs, "runtime/router/activation-rules.yaml")
	if err != nil {
		return err
	}
	for index, rule := range runtimeSlice(activationRules["core_bootstrap"]) {
		if _, err := db.Exec(`INSERT INTO core_bootstrap_rules (rule_id, ordinal) VALUES (?, ?)`, runtimeString(rule), index); err != nil {
			return err
		}
	}
	if err := insertRuntimeSourceFile(db, "runtime/router/activation-rules.yaml", "core_bootstrap_rules", "core_bootstrap_order_config"); err != nil {
		return err
	}
	discoveryCheckpoints, err := runtimeCanonicalDocument(docs, "runtime/discovery/capability-checkpoints.yaml")
	if err != nil {
		return err
	}
	for _, checkpoint := range runtimeSliceOfMaps(discoveryCheckpoints["checkpoints"]) {
		if _, err := db.Exec(`INSERT INTO discovery_checkpoints (phase, trigger, description, discovery_targets, metadata) VALUES (?, ?, ?, ?, ?)`,
			runtimeString(checkpoint["phase"]), runtimeString(checkpoint["trigger"]), runtimeString(checkpoint["description"]), runtimeJSON(checkpoint["discovery_targets"]), runtimeJSON(checkpoint["metadata"]),
		); err != nil {
			return err
		}
	}
	searchStrategy := runtimeMap(discoveryCheckpoints["search_strategy"])
	if len(searchStrategy) > 0 {
		if _, err := db.Exec(`INSERT INTO discovery_search_strategy (priority_order, fallback, min_confidence_threshold) VALUES (?, ?, ?)`, runtimeJSON(searchStrategy["priority_order"]), runtimeJSON(searchStrategy["fallback"]), runtimeString(searchStrategy["min_confidence_threshold"])); err != nil {
			return err
		}
	}

	if _, err := db.Exec(`INSERT INTO compiler_metadata (key, value) VALUES ('compiler_version', ?)`, goRuntimeCompilerVersion); err != nil {
		return err
	}
	if _, err := db.Exec(`INSERT INTO compiler_metadata (key, value) VALUES ('compiled_at', ?)`, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return err
	}
	if _, err := db.Exec(`INSERT INTO compiler_metadata (key, value) VALUES ('schema_version', '1.0')`); err != nil {
		return err
	}

	for _, config := range runtimeConfigMappings() {
		if err := insertRuntimeConfigRows(db, docs, config.rel, config.table, config.idCol, config.listKey, config.idKeys); err != nil {
			return err
		}
		if config.rel == "runtime/budget/token-budget.yaml" {
			if err := insertTokenBudgetFrameworkRows(db, docs); err != nil {
				return err
			}
		}
		if err := insertRuntimeSourceFile(db, config.rel, config.table, "structured_runtime_config"); err != nil {
			return err
		}
	}
	if err := insertTransactionTemplateRows(db, docs); err != nil {
		return err
	}
	if err := insertRuntimeSourceFile(db, "runtime/transactions/transaction-templates.yaml", "transaction_templates", "transaction_templates_config"); err != nil {
		return err
	}
	if err := insertCompilerRuleRows(db, docs); err != nil {
		return err
	}
	if err := insertRuntimeSourceFile(db, "runtime/compiler/compiler-rules.yaml", "compiler_rules", "compiler_rules_config"); err != nil {
		return err
	}
	if err := compileDerivedForbiddenTokens(repo, db); err != nil {
		return err
	}
	// Phase 1C₁ (2026-06-09): repository topology now uses a custom
	// projection function because v2 schema's `subtrees:` + per-subtree
	// owner/purpose fields do not fit the tuple-format pipeline.
	// See scripts/ai-skill-cli/internal/app/repository_topology_compile.go.
	if err := compileRepositoryTopology(repo, db); err != nil {
		return err
	}
	return nil
}

func runtimeConfigMappings() []runtimeConfigMapping {
	return []runtimeConfigMapping{
		{"runtime/budget/token-budget.yaml", "runtime_budget", "model_name", "per_model", []string{"name", "model"}},
		{"runtime/context/ttl-policy.yaml", "context_ttl_policy", "ttl_type", "rules", []string{"name", "id"}},
		{"runtime/constitution/decision-recording.yaml", "decision_recording", "section", "", nil},
		{"runtime/discovery/capability-checkpoints.yaml", "capability_checkpoints", "checkpoint_id", "checkpoints", []string{"phase", "name"}},
		{"runtime/distributed/distributed-locks.yaml", "distributed_locks", "lock_name", "locks", []string{"name", "id", "state"}},
		{"runtime/distributed/multi-agent-coordination.yaml", "multi_agent_coordination", "rule_id", "coordination_rules", []string{"name", "id", "rule_id"}},
		{"runtime/distributed/async-job-lifecycle.yaml", "async_job_lifecycle", "state", "states", []string{"name", "state"}},
		{"runtime/gates/blocking-gates.yaml", "blocking_gates", "gate_id", "gates", []string{"id", "name"}},
		{"runtime/guards/circuit-breaker.yaml", "circuit_breaker", "guard_name", "", nil},
		{"runtime/guards/context-pollution.yaml", "context_pollution", "signal_name", "signals", []string{"name", "id"}},
		{"runtime/health/context-health-score.yaml", "context_health_score", "dimension", "dimensions", []string{"name", "dimension"}},
		{"runtime/intelligence/intelligence-routing.yaml", "intelligence_routing", "rule_id", "routing_rules", []string{"name", "id", "rule_id"}},
		{"runtime/obligations/obligation-ledger.yaml", "obligation_ledger", "obligation_id", "obligations", []string{"id", "name"}},
		{"runtime/output-governance/language-policy.yaml", "language_policy", "section", "rules", []string{"name", "id"}},
		{"runtime/output-governance/output-rules.yaml", "output_rules", "section", "rules", []string{"name", "id"}},
		{"runtime/output-governance/governance-gates.yaml", "governance_gates", "gate_id", "gates", []string{"id", "name"}},
		{"runtime/phases/phase-machine.yaml", "phase_machine", "phase_id", "phases", []string{"id", "name"}},
		{"runtime/pipeline/context-flow.yaml", "pipeline_context_flow", "level", "levels", []string{"name", "level"}},
		{"runtime/pipeline/guard-chain.yaml", "guard_chain", "stage", "stages", []string{"name", "stage"}},
		{"runtime/pipeline/relevance-engine.yaml", "relevance_engine", "component", "scoring.components", []string{"id", "name"}},
		{"runtime/pipeline/session-lifecycle.yaml", "session_lifecycle", "stage", "stages", []string{"name", "stage"}},
		{"runtime/prompt-artifacts/artifact-templates.yaml", "prompt_artifact_templates", "template_name", "templates", []string{"name", "id"}},
		{"runtime/prompt-artifacts/composition-rules.yaml", "prompt_composition_rules", "rule_id", "rules", []string{"rule_id", "name"}},
		{"runtime/recovery/recovery-strategies.yaml", "recovery_strategies", "strategy_id", "recovery_strategies", []string{"id", "name"}},
		{"runtime/recovery/state-repair.yaml", "state_repair", "procedure_id", "repair_procedures", []string{"name", "id"}},
		{"runtime/recovery/obligation-rebuild.yaml", "obligation_rebuild", "procedure_id", "rebuild_procedures", []string{"name", "id"}},
		{"runtime/recovery/phase-reconciliation.yaml", "phase_reconciliation", "procedure_id", "reconciliation_procedures", []string{"name", "id"}},
		// runtime/repository-topology.yaml moved out of tuple-driven projection in
		// Phase 1C₁ (2026-06-09). The v2 schema's `subtrees:` + per-subtree
		// owner/purpose/path fields do not fit the generic tuple format; the
		// canonical projection now flows through
		// repository_topology_compile.go::compileRepositoryTopology.
		{"runtime/sanitization-patterns.yaml", "sanitization_patterns", "category", "pattern_families", []string{"category", "id", "name"}},
		{"runtime/scheduler/execution-queue.yaml", "execution_queue", "queue_name", "queue_structure", []string{"name", "id"}},
		{"runtime/scheduler/priority-scheduler.yaml", "priority_scheduler", "priority_level", "levels", []string{"name", "level"}},
		{"runtime/transactions/transaction-machine.yaml", "transaction_templates_ext", "template_name", "transaction_templates", []string{"id", "name"}},
	}
}

func runtimeCanonicalDocumentPaths() []string {
	seen := map[string]bool{}
	paths := []string{}
	for _, rel := range []string{
		"runtime/phases/phase-machine.yaml",
		"runtime/obligations/obligation-ledger.yaml",
		"runtime/gates/blocking-gates.yaml",
		"runtime/transactions/transaction-machine.yaml",
		"runtime/router/activation-rules.yaml",
		"runtime/discovery/capability-checkpoints.yaml",
		"runtime/transactions/transaction-templates.yaml",
		"runtime/compiler/compiler-rules.yaml",
		// Phase 1C₁: repository-topology.yaml moved out of runtimeConfigMappings
		// because its v2 schema doesn't fit the tuple format; it still needs
		// to appear in runtime_config_documents for introspection and drift
		// checks, so it is listed explicitly here.
		"runtime/repository-topology.yaml",
	} {
		paths = append(paths, rel)
		seen[rel] = true
	}
	for _, mapping := range runtimeConfigMappings() {
		if !seen[mapping.rel] {
			paths = append(paths, mapping.rel)
			seen[mapping.rel] = true
		}
	}
	sort.Strings(paths)
	return paths
}

func insertRuntimeSourceFile(db *sql.DB, rel string, table string, compileRule string) error {
	if _, err := db.Exec(`INSERT OR REPLACE INTO runtime_source_files (source_path, source_kind, target_table, compile_rule, compiled_at, compiler_version, status) VALUES (?, 'db', ?, ?, datetime('now'), ?, 'synced')`, rel, table, compileRule, goRuntimeCompilerVersion); err != nil {
		return err
	}
	return insertRuntimeConfigProjection(db, rel, table, "__config__")
}

func insertTokenBudgetFrameworkRows(db *sql.DB, docs map[string]map[string]any) error {
	config, err := runtimeCanonicalDocument(docs, "runtime/budget/token-budget.yaml")
	if err != nil {
		return err
	}
	if defaultBudget := runtimeMap(config["default_budget"]); len(defaultBudget) > 0 {
		if _, err := db.Exec(`INSERT INTO runtime_budget (model_name, content) VALUES (?, ?)`, "default_budget", runtimeJSON(defaultBudget)); err != nil {
			return err
		}
	}
	layerBudget := runtimeMap(config["layer_budget"])
	layerNames := make([]string, 0, len(layerBudget))
	for name := range layerBudget {
		layerNames = append(layerNames, name)
	}
	sort.Strings(layerNames)
	for _, name := range layerNames {
		value := runtimeMap(layerBudget[name])
		entry := map[string]any{"layer": name}
		for key, nested := range value {
			entry[key] = nested
		}
		if _, err := db.Exec(`INSERT INTO runtime_budget (model_name, content) VALUES (?, ?)`, "layer:"+name, runtimeJSON(entry)); err != nil {
			return err
		}
	}
	for _, group := range []struct {
		key    string
		prefix string
	}{
		{key: "on_warning", prefix: "on_warning"},
		{key: "on_hard_stop", prefix: "on_hard_stop"},
	} {
		for index, action := range runtimeSliceOfMaps(config[group.key]) {
			actionName := runtimeDefaultString(action["action"], fmt.Sprintf("action_%d", index+1))
			if _, err := db.Exec(`INSERT INTO runtime_budget (model_name, content) VALUES (?, ?)`, fmt.Sprintf("%s:%02d:%s", group.prefix, index+1, actionName), runtimeJSON(action)); err != nil {
				return err
			}
		}
	}
	return nil
}

func insertRuntimeConfigRows(db *sql.DB, docs map[string]map[string]any, rel string, table string, idCol string, listKey string, idKeys []string) error {
	config, err := runtimeCanonicalDocument(docs, rel)
	if err != nil {
		return err
	}
	values := []any{}
	if listKey == "" {
		for key, value := range config {
			if key == "status" || key == "owner_layer" || key == "schema_version" || key == "source_path" || key == "summary" {
				continue
			}
			if nested, ok := value.(map[string]any); ok {
				entry := map[string]any{"name": key}
				for nestedKey, nestedValue := range nested {
					entry[nestedKey] = nestedValue
				}
				values = append(values, entry)
			}
		}
	} else {
		values = runtimeNestedValues(config, listKey)
	}
	for index, value := range values {
		id := runtimeIDFromValue(value, idKeys, fmt.Sprintf("entry_%d", index+1))
		if _, err := db.Exec(fmt.Sprintf("INSERT INTO %s (%s, content) VALUES (?, ?)", table, idCol), id, runtimeJSON(value)); err != nil {
			return err
		}
	}
	if _, err := db.Exec(fmt.Sprintf("INSERT INTO %s (%s, content) VALUES ('__config__', ?)", table, idCol), runtimeJSON(config)); err != nil {
		return err
	}
	return nil
}

func insertTransactionTemplateRows(db *sql.DB, docs map[string]map[string]any) error {
	config, err := runtimeCanonicalDocument(docs, "runtime/transactions/transaction-templates.yaml")
	if err != nil {
		return err
	}
	for _, template := range runtimeSliceOfMaps(config["templates"]) {
		id := runtimeDefaultString(template["id"], runtimeDefaultString(template["name"], "default"))
		if _, err := db.Exec(`INSERT OR REPLACE INTO transaction_templates (id, name, description, typical_steps, content) VALUES (?, ?, ?, ?, ?)`, id, runtimeDefaultString(template["name"], id), runtimeString(template["description"]), runtimeJSON(template["steps"]), runtimeJSON(template)); err != nil {
			return err
		}
	}
	if _, err := db.Exec(`INSERT OR REPLACE INTO transaction_templates (id, name, description, typical_steps, content) VALUES ('__config__', '__config__', 'Full YAML config', '[]', ?)`, runtimeJSON(config)); err != nil {
		return err
	}
	return nil
}

func insertCompilerRuleRows(db *sql.DB, docs map[string]map[string]any) error {
	var config compilerMappingFile
	if err := runtimeCanonicalDocumentStruct(docs, "runtime/compiler/compiler-rules.yaml", &config); err != nil {
		return err
	}
	for _, rule := range config.CompilationRules {
		id := runtimeIDFromValue(rule, []string{"id", "name"}, "rule")
		if _, err := db.Exec(`INSERT INTO compiler_rules (rule_id, content) VALUES (?, ?)`, id, runtimeJSON(rule)); err != nil {
			return err
		}
	}
	for index, mapping := range config.SourceTargetMapping {
		id := mapping.Source
		if id == "" {
			id = fmt.Sprintf("mapping_%d", index+1)
		}
		if _, err := db.Exec(`INSERT INTO compiler_rules (rule_id, content) VALUES (?, ?)`, id, runtimeJSON(mapping)); err != nil {
			return err
		}
	}
	if steps, ok := config.CompilerWorkflow["steps"]; ok {
		for _, step := range runtimeSliceOfMaps(steps) {
			id := runtimeIDFromValue(step, []string{"action", "id", "name"}, "step")
			if _, err := db.Exec(`INSERT INTO compiler_rules (rule_id, content) VALUES (?, ?)`, id, runtimeJSON(step)); err != nil {
				return err
			}
		}
	}
	if _, err := db.Exec(`INSERT INTO compiler_rules (rule_id, content) VALUES ('__config__', ?)`, runtimeJSON(config)); err != nil {
		return err
	}
	return nil
}

func compileProseRuntimeSources(repo string, db *sql.DB, docs map[string]map[string]any) error {
	var config compilerMappingFile
	if err := runtimeCanonicalDocumentStruct(docs, "runtime/compiler/compiler-rules.yaml", &config); err != nil {
		return err
	}
	for _, mapping := range config.SourceTargetMapping {
		paths, err := runtimeCompilerGlob(repo, mapping.Source)
		if err != nil {
			return err
		}
		if strings.Contains(mapping.Source, "plans/active/*.md") {
			plans := []map[string]any{}
			for _, path := range paths {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				name := strings.TrimSuffix(filepath.Base(path), ".md")
				title := ""
				if headings := runtimeHeadings(string(content), `(?m)^#\s+(.+)$`); len(headings) > 0 {
					title = runtimeString(headings[0]["name"])
				}
				plans = append(plans, map[string]any{"plan_id": name, "filename": name, "title": title})
			}
			if err := insertGeneratedSurface(db, "plans/active/*.md", "plans.index", mapping.CompileRule, map[string]any{"plans": plans}); err != nil {
				return err
			}
			continue
		}
		if strings.Contains(mapping.Source, "analysis/apk/workflows/*.md") {
			entries := []map[string]any{}
			for _, path := range paths {
				if filepath.Base(path) == "README.md" {
					continue
				}
				entry, err := compileAPKWorkflow(path)
				if err != nil {
					return err
				}
				entries = append(entries, entry)
			}
			if err := insertGeneratedSurface(db, "analysis/apk/workflows/*.md", "analysis.apk.workflow_phases", mapping.CompileRule, map[string]any{"workflows": entries}); err != nil {
				return err
			}
			continue
		}
		for _, path := range paths {
			rel, err := filepath.Rel(repo, path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			data, targetKey, err := compileProseFile(repo, path, mapping)
			if err != nil {
				return err
			}
			if err := insertGeneratedSurface(db, rel, targetKey, mapping.CompileRule, data); err != nil {
				return err
			}
		}
	}
	if err := compileExecutableYAMLContracts(repo, db); err != nil {
		return err
	}
	if err := compileDiscoverySignals(repo, db); err != nil {
		return err
	}
	if err := compileCognitiveModeEnforcementRules(db); err != nil {
		return err
	}
	if err := compileCognitiveModePOCSeed(db); err != nil {
		return err
	}
	if err := compileBootstrapEnforcementRules(db); err != nil {
		return err
	}
	return nil
}

func compileCognitiveModePOCSeed(db *sql.DB) error {
	// Seed the Phase 1 POC task record per plan acceptance criterion (d):
	// "至少 1 個 POC task 的 cognitive_modes row 寫入 runtime.db".
	// Since compile wipes the table, the seed must be re-inserted each compile
	// until behavioral runtime (Phase 4) populates rows from real task resolutions.
	_, err := db.Exec(
		`INSERT OR IGNORE INTO cognitive_modes (id, task_id, execution_mode, context_mode, governance_mode, memory_mode, resolved_at, source) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		1,
		"phase1.poc.cognitive_modes_bootstrap",
		"NORMAL",
		"SUMMARY_FIRST",
		"STANDARD",
		"NONE",
		"2026-05-22T00:00:00Z",
		"plan:2026-05-22-1629-runtime-cognitive-modes-system §Phase 1 POC seed (compile-time)",
	)
	return err
}

func compileBootstrapEnforcementRules(db *sql.DB) error {
	_, err := db.Exec(
		`INSERT OR IGNORE INTO obligations (id, phase, name, description, verification, severity, depends_on, linked_gates, metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"obligation.bootstrap.receipt_acknowledged",
		"phase.bootstrap",
		"Bootstrap Receipt Acknowledged",
		"Session start（含 resume from summary）first-turn 必須輸出 Bootstrap Receipt（rules=✓ phase=<id> obligations=<n> gates=<n>），證明已讀 CORE_BOOTSTRAP.md + README.md 並查過 runtime.db。",
		`["bootstrap_receipt_emitted == true","first_message_contains 'Bootstrap: rules=✓ phase='","resume_sessions_not_exempt"]`,
		"high",
		`["gate.bootstrap.core_rules_loaded","gate.bootstrap.layout_loaded"]`,
		`["gate.bootstrap.receipt_present"]`,
		`{"source":"CLAUDE.md + CORE_BOOTSTRAP.md","contract":"Bootstrap Receipt clause","applies_to":"new + resume sessions"}`,
	)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		`INSERT OR IGNORE INTO gates (id, phase, name, description, severity, check_type, check_target, check_verification, failure_action, failure_message, metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"gate.bootstrap.receipt_present",
		"phase.bootstrap",
		"Bootstrap Receipt Present",
		"執行任何非-Read 工具（Edit/Write/Bash/git/...）前，first user-facing message 必須含 Bootstrap Receipt。Resume from summary 不豁免。",
		"high",
		"verification",
		"first_message contains 'Bootstrap: rules=✓ phase='",
		"bootstrap_receipt_emitted == true",
		"block_execution",
		"Bootstrap Receipt 未輸出：請先讀 CORE_BOOTSTRAP.md + README.md，查 runtime.db 取得 phase/obligations/gates，然後在 first message 輸出 'Bootstrap: rules=✓ phase=<id> obligations=<n> gates=<n>'。Resume from summary 同樣適用。",
		`{"source":"CLAUDE.md","contract":"Bootstrap Receipt clause","failure_pattern":"bootstrap-bypass-on-resume","deferred":"behavioral first-turn detection"}`,
	)
	return err
}

func compileCognitiveModeEnforcementRules(db *sql.DB) error {
	_, err := db.Exec(
		`INSERT OR IGNORE INTO obligations (id, phase, name, description, verification, severity, depends_on, linked_gates, metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"obligation.execution.resolve_cognitive_mode",
		"phase.execution",
		"Resolve Cognitive Mode",
		"任務開始前必須解析 4 維 cognitive mode（execution/context/governance/memory），並記錄到 cognitive_modes 表或 final report。",
		`["cognitive_mode_resolved == true","cognitive_modes table has row OR final_report includes Cognitive Mode block"]`,
		"high",
		`[]`,
		`["gate.execution.cognitive_mode_resolved"]`,
		`{"source":"cognitive-modes-phase-integration","contract":"runtime/cognitive-modes.yaml"}`,
	)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		`INSERT OR IGNORE INTO gates (id, phase, name, description, severity, check_type, check_target, check_verification, failure_action, failure_message, metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"gate.execution.cognitive_mode_resolved",
		"phase.execution",
		"Cognitive Mode Resolved",
		"執行任何 close-loop 動作前，cognitive mode 必須已解析。NONE 解析等同於未解析（除非是 Phase D doc-only）。",
		"high",
		"verification",
		"cognitive_modes table OR final_report Cognitive Mode block",
		"cognitive_mode_resolved == true",
		"block_execution",
		"Cognitive mode 未解析：請先執行 mode discovery 並記錄 4 維 mode 值，再繼續 close-loop。",
		`{"source":"cognitive-modes-governance-integration","contract":"runtime/cognitive-modes-governance-integration.yaml","deferred":"behavioral pre-commit wiring"}`,
	)
	return err
}

func compileDiscoverySignals(repo string, db *sql.DB) error {
	yamlPath := filepath.Join(repo, "runtime", "cognitive-modes-discovery.yaml")
	content, err := os.ReadFile(yamlPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var data map[string]any
	if err := yaml.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("compile discovery signals: %w", err)
	}
	for _, sig := range runtimeSliceOfMaps(data["signals"]) {
		priority := 0
		if p, ok := sig["priority"].(int); ok {
			priority = p
		}
		if _, err := db.Exec(
			`INSERT INTO discovery_signals (signal_name, signal_type, pattern, execution_mode, context_mode, governance_mode, memory_mode, priority, description) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			runtimeString(sig["name"]),
			runtimeString(sig["signal_type"]),
			runtimeString(sig["pattern"]),
			runtimeNullableString(sig["execution_mode"]),
			runtimeNullableString(sig["context_mode"]),
			runtimeNullableString(sig["governance_mode"]),
			runtimeNullableString(sig["memory_mode"]),
			priority,
			runtimeString(sig["description"]),
		); err != nil {
			return err
		}
	}
	return nil
}

func compileExecutableYAMLContracts(repo string, db *sql.DB) error {
	sourceRoots := []string{
		"governance",
		"enforcement",
		"workflow",
		"ai-tools",
		filepath.ToSlash(filepath.Join("metadata", "rules")),
		"runtime",
	}
	return filepath.Walk(repo, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info == nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".yaml") {
			return nil
		}
		rel, err := filepath.Rel(repo, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		inContractRoot := false
		for _, root := range sourceRoots {
			if rel == root || strings.HasPrefix(rel, root+"/") {
				inContractRoot = true
				break
			}
		}
		if !inContractRoot {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var data map[string]any
		if err := yaml.Unmarshal(content, &data); err != nil {
			return fmt.Errorf("compile %s: %w", rel, err)
		}
		normalized := runtimeMap(runtimeNormalizeYAML(data))
		projection := runtimeMap(normalized["runtime_projection"])
		if !runtimeBool(projection["enabled"]) {
			return nil
		}
		targetKey := runtimeDefaultString(projection["target_key"], runtimeDefaultString(normalized["id"], rel))
		return insertGeneratedSurface(db, rel, targetKey, "executable YAML contract", normalized)
	})
}

func compileProseFile(repo string, path string, mapping compilerMapping) (map[string]any, string, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	content := string(contentBytes)
	rel, err := filepath.Rel(repo, path)
	if err != nil {
		return nil, "", err
	}
	rel = filepath.ToSlash(rel)
	domain := runtimeCompilerDomain(rel)
	switch {
	case strings.Contains(mapping.CompileRule, "execution-flow"):
		phases := runtimeCompileSections(content, `(?m)^##\s+\d+\.\s+.+$`)
		return map[string]any{"phases": phases}, "workflow." + domain + ".phases", nil
	case strings.Contains(mapping.CompileRule, "artifact gates"):
		return map[string]any{"compiled_from": rel, "artifacts": runtimeHeadings(content, `(?m)^##\s+\d+\.\s+(.+)$`), "verification_gates": runtimeHeadings(content, `(?m)^###\s+(.+)$`)}, "workflow." + domain + ".artifacts", nil
	case strings.Contains(mapping.CompileRule, "writeback transaction"):
		return map[string]any{"compiled_from": rel, "states": runtimeHeadings(content, `(?m)^###?\s+(.+)$`)}, "enforcement." + domain + ".transactions", nil
	case strings.Contains(mapping.CompileRule, "validation gate"):
		return map[string]any{"compiled_from": rel, "gates": runtimeNumberedLines(content), "scenarios": runtimeMarkdownTableRows(content)}, "enforcement." + domain + ".gates", nil
	case strings.Contains(mapping.CompileRule, "failure taxonomy"):
		return map[string]any{"compiled_from": rel, "failure_taxonomy": runtimeMarkdownTableRows(content), "loop_steps": runtimeNumberedLines(content)}, "enforcement." + domain + ".recovery", nil
	case strings.Contains(mapping.CompileRule, "neutral-language") || strings.Contains(mapping.CompileRule, "sanitization") || strings.Contains(mapping.CompileRule, "tool neutrality"):
		return map[string]any{"compiled_from": rel, "sections": runtimeHeadings(content, `(?m)^###?\s+(.+)$`)}, "governance." + domain, nil
	case strings.Contains(mapping.CompileRule, "knowledge-update-flow.md 的 11"):
		return map[string]any{"compiled_from": rel, "steps": runtimeHeadings(content, `(?m)^##\s+Step\s+\d+[：:]\s*(.+)$`)}, "governance.knowledge_update_flow", nil
	case strings.Contains(mapping.CompileRule, "classification"):
		return map[string]any{"compiled_from": rel, "classification_dimensions": runtimeMarkdownTableRows(content)}, "governance.classification_rules", nil
	case strings.Contains(mapping.CompileRule, "系統升級治理"):
		return map[string]any{"compiled_from": rel, "checklist_categories": runtimeHeadings(content, `(?m)^###\s+(.+)$`), "forced_rules": runtimeHeadings(content, `(?m)^###\s+規則\s+\d+[：:]\s*(.+)$`)}, "governance.system_upgrade", nil
	default:
		return map[string]any{"compiled_from": rel, "sections": runtimeHeadings(content, `(?m)^##+\s+(.+)$`)}, rel, nil
	}
}

func insertGeneratedSurface(db *sql.DB, sourcePath string, targetKey string, compileRule string, data map[string]any) error {
	_, err := db.Exec(`INSERT OR REPLACE INTO generated_surfaces (source_path, target_key, compile_rule, compiled_at, compiler_version, status, data) VALUES (?, ?, ?, datetime('now'), ?, 'synced', ?)`, sourcePath, targetKey, compileRule, goRuntimeCompilerVersion, runtimeJSON(data))
	return err
}

func runtimeCompileSections(content string, headingPattern string) []map[string]any {
	matches := regexp.MustCompile(headingPattern).FindAllStringIndex(content, -1)
	sections := []map[string]any{}
	for index, match := range matches {
		end := len(content)
		if index+1 < len(matches) {
			end = matches[index+1][0]
		}
		heading := strings.TrimSpace(strings.TrimPrefix(content[match[0]:match[1]], "##"))
		body := strings.TrimSpace(content[match[1]:end])
		section := map[string]any{"name": heading, "source_line": strings.Count(content[:match[0]], "\n") + 1}
		if actions := runtimeBulletLines(body); len(actions) > 0 {
			section["allowed_actions"] = actions
		}
		if tables := runtimeMarkdownTableRows(body); len(tables) > 0 {
			section["tables"] = tables
		}
		sections = append(sections, section)
	}
	return sections
}

func compileAPKWorkflow(path string) (map[string]any, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(contentBytes)
	steps := runtimeHeadings(content, `(?m)^##\s+步驟\s+\d+[：:]\s*(.+)$`)
	return map[string]any{"workflow_name": strings.TrimSuffix(filepath.Base(path), ".md"), "total_steps": len(steps), "steps": steps}, nil
}

func runtimeHeadings(content string, pattern string) []map[string]any {
	re := regexp.MustCompile(pattern)
	result := []map[string]any{}
	for _, match := range re.FindAllStringSubmatchIndex(content, -1) {
		value := strings.TrimSpace(content[match[2]:match[3]])
		result = append(result, map[string]any{"name": value, "source_line": strings.Count(content[:match[0]], "\n") + 1})
	}
	return result
}

func runtimeBulletLines(content string) []map[string]any {
	re := regexp.MustCompile(`(?m)^-\s+(.+)$`)
	result := []map[string]any{}
	for _, match := range re.FindAllStringSubmatch(content, -1) {
		result = append(result, map[string]any{"action": strings.TrimSpace(match[1]), "source": "bullet"})
	}
	return result
}

func runtimeNumberedLines(content string) []map[string]any {
	re := regexp.MustCompile(`(?m)^\d+\.\s+(.+)$`)
	result := []map[string]any{}
	for _, match := range re.FindAllStringSubmatch(content, -1) {
		result = append(result, map[string]any{"rule": strings.TrimSpace(match[1])})
	}
	return result
}

func runtimeMarkdownTableRows(content string) []map[string]any {
	result := []map[string]any{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") || strings.Contains(line, "---") {
			continue
		}
		cols := []string{}
		for _, col := range strings.Split(line, "|") {
			col = strings.TrimSpace(col)
			if col != "" {
				cols = append(cols, col)
			}
		}
		if len(cols) > 0 {
			result = append(result, map[string]any{"columns": cols})
		}
	}
	return result
}

func runtimeCompilerGlob(repo string, pattern string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(repo, filepath.FromSlash(pattern)))
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)
	return matches, nil
}

func runtimeCompilerDomain(sourcePath string) string {
	if match := regexp.MustCompile(`workflow/([^/]+)/`).FindStringSubmatch(sourcePath); len(match) > 1 {
		return match[1]
	}
	if match := regexp.MustCompile(`enforcement/(.+)\.md$`).FindStringSubmatch(sourcePath); len(match) > 1 {
		return strings.ReplaceAll(match[1], "-", "_")
	}
	return "unknown"
}

func loadRuntimeCanonicalDocuments(repo string) (map[string]map[string]any, error) {
	dbPath := filepath.Join(repo, "runtime", "runtime.db")
	if docs, err := readRuntimeCanonicalDocumentsFromDB(dbPath); err == nil && len(docs) > 0 {
		for _, rel := range runtimeCanonicalDocumentPaths() {
			if _, ok := docs[rel]; ok {
				continue
			}
			doc, readErr := readRuntimeYAMLMap(repo, rel)
			if readErr != nil {
				return nil, fmt.Errorf("load missing canonical runtime document %s: %w", rel, readErr)
			}
			docs[rel] = doc
		}
		return docs, nil
	}
	return importRuntimeCanonicalDocumentsFromYAML(repo)
}

func readRuntimeCanonicalDocumentsFromDB(dbPath string) (map[string]map[string]any, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var tableCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'runtime_config_documents'`).Scan(&tableCount); err != nil {
		return nil, err
	}
	if tableCount == 0 {
		return nil, fmt.Errorf("runtime_config_documents table missing")
	}
	rows, err := db.Query(`SELECT logical_id, content_json FROM runtime_config_documents ORDER BY logical_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	docs := map[string]map[string]any{}
	for rows.Next() {
		var logicalID string
		var content string
		if err := rows.Scan(&logicalID, &content); err != nil {
			return nil, err
		}
		var doc map[string]any
		if err := json.Unmarshal([]byte(content), &doc); err != nil {
			return nil, fmt.Errorf("%s canonical JSON: %w", logicalID, err)
		}
		docs[logicalID] = doc
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return nil, fmt.Errorf("runtime_config_documents is empty")
	}
	return docs, nil
}

func importRuntimeCanonicalDocumentsFromYAML(repo string) (map[string]map[string]any, error) {
	docs := map[string]map[string]any{}
	for _, rel := range runtimeCanonicalDocumentPaths() {
		doc, err := readRuntimeYAMLMap(repo, rel)
		if err != nil {
			return nil, fmt.Errorf("load canonical runtime document %s: %w", rel, err)
		}
		docs[rel] = doc
	}
	return docs, nil
}

func insertRuntimeConfigDocuments(db *sql.DB, docs map[string]map[string]any) error {
	for _, logicalID := range runtimeCanonicalDocumentPaths() {
		doc, err := runtimeCanonicalDocument(docs, logicalID)
		if err != nil {
			return err
		}
		content := runtimeJSON(doc)
		checksum := runtimeDocumentChecksum(content)
		if _, err := db.Exec(`INSERT OR REPLACE INTO runtime_config_documents (logical_id, owner_layer, status, schema_version, content_json, checksum, updated_at) VALUES (?, ?, ?, ?, ?, ?, datetime('now'))`,
			logicalID, runtimeString(doc["owner_layer"]), runtimeDefaultString(doc["status"], "active"), runtimeString(doc["schema_version"]), content, checksum,
		); err != nil {
			return err
		}
	}
	return nil
}

func insertRuntimeConfigProjection(db *sql.DB, logicalID string, targetTable string, rowKey string) error {
	var checksum string
	if err := db.QueryRow(`SELECT checksum FROM runtime_config_documents WHERE logical_id = ?`, logicalID).Scan(&checksum); err != nil {
		return err
	}
	_, err := db.Exec(`INSERT OR REPLACE INTO runtime_config_projections (logical_id, target_table, row_key, checksum, updated_at) VALUES (?, ?, ?, ?, datetime('now'))`, logicalID, targetTable, rowKey, checksum)
	return err
}

func runtimeCanonicalDocument(docs map[string]map[string]any, logicalID string) (map[string]any, error) {
	doc, ok := docs[logicalID]
	if !ok {
		return nil, fmt.Errorf("canonical runtime document missing: %s", logicalID)
	}
	return doc, nil
}

func runtimeCanonicalDocumentStruct(docs map[string]map[string]any, logicalID string, target any) error {
	doc, err := runtimeCanonicalDocument(docs, logicalID)
	if err != nil {
		return err
	}
	content, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, target)
}

func runtimeDocumentChecksum(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

func readRuntimeYAML(repo string, rel string, target any) error {
	content, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(rel)))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, target)
}

func readRuntimeYAMLMap(repo string, rel string) (map[string]any, error) {
	var result map[string]any
	if err := readRuntimeYAML(repo, rel, &result); err != nil {
		return nil, err
	}
	return runtimeMap(runtimeNormalizeYAML(result)), nil
}

func runtimeNormalizeYAML(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		result := map[string]any{}
		for key, nested := range typed {
			result[strings.TrimPrefix(key, ":")] = runtimeNormalizeYAML(nested)
		}
		return result
	case []any:
		result := make([]any, len(typed))
		for index, nested := range typed {
			result[index] = runtimeNormalizeYAML(nested)
		}
		return result
	default:
		return value
	}
}

func runtimeNestedValues(config map[string]any, path string) []any {
	current := any(config)
	for _, part := range strings.Split(path, ".") {
		current = runtimeMap(current)[part]
	}
	if list, ok := current.([]any); ok {
		return list
	}
	if current != nil {
		return []any{current}
	}
	return nil
}

func runtimeMap(value any) map[string]any {
	if value == nil {
		return map[string]any{}
	}
	if result, ok := value.(map[string]any); ok {
		return result
	}
	return map[string]any{}
}

func runtimeSlice(value any) []any {
	if result, ok := value.([]any); ok {
		return result
	}
	return nil
}

func runtimeSliceOfMaps(value any) []map[string]any {
	result := []map[string]any{}
	for _, item := range runtimeSlice(value) {
		if mapped, ok := item.(map[string]any); ok {
			result = append(result, mapped)
		}
	}
	return result
}

func runtimeIDFromValue(value any, keys []string, fallback string) string {
	mapped := runtimeMap(value)
	for _, key := range keys {
		if id := runtimeString(mapped[key]); id != "" {
			return id
		}
	}
	return fallback
}

func runtimeString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case fmt.Stringer:
		return typed.String()
	case nil:
		return ""
	default:
		return fmt.Sprint(typed)
	}
}

func runtimeNullableString(value any) any {
	if value == nil {
		return nil
	}
	s := runtimeString(value)
	if s == "" {
		return nil
	}
	return s
}

func runtimeDefaultString(value any, fallback string) string {
	if result := runtimeString(value); result != "" {
		return result
	}
	return fallback
}

func runtimeInt(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	default:
		return 0
	}
}

func runtimeBool(value any) bool {
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		return strings.EqualFold(typed, "true") || typed == "1" || strings.EqualFold(typed, "yes")
	default:
		return false
	}
}

func runtimeJSON(value any) string {
	if value == nil {
		return "{}"
	}
	content, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(content)
}
