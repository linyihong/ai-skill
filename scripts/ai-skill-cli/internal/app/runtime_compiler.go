package app

import (
	"database/sql"
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
	SourceTargetMapping []compilerMapping `yaml:"source_target_mapping"`
	CompilationRules    []map[string]any  `yaml:"compilation_rules"`
	CompilerWorkflow    map[string]any    `yaml:"compiler_workflow"`
}

type compilerMapping struct {
	Source      string `yaml:"source"`
	Target      string `yaml:"target"`
	CompileRule string `yaml:"compile_rule"`
}

func buildNativeRuntimeDBFromSources(repo string, outputDB string) Check {
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
	if err := compileStructuredRuntimeSources(repo, db); err != nil {
		return Check{Name: "runtime_compile_native", Status: "failed", Message: err.Error()}
	}
	if err := compileProseRuntimeSources(repo, db); err != nil {
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
		`CREATE TABLE activation_rules (rule_id TEXT PRIMARY KEY, description TEXT, activation_when TEXT, load_strategy TEXT DEFAULT 'lazy', load_priority TEXT DEFAULT 'P2', load_estimated_tokens INTEGER DEFAULT 0, load_source TEXT, metadata TEXT, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE core_bootstrap_rules (rule_id TEXT PRIMARY KEY, ordinal INTEGER DEFAULT 0);`,
		`CREATE TABLE discovery_checkpoints (id INTEGER PRIMARY KEY AUTOINCREMENT, phase TEXT NOT NULL, trigger TEXT NOT NULL, description TEXT, discovery_targets TEXT, metadata TEXT);`,
		`CREATE TABLE discovery_search_strategy (id INTEGER PRIMARY KEY AUTOINCREMENT, priority_order TEXT, fallback TEXT, min_confidence_threshold TEXT);`,
		`CREATE TABLE compiler_metadata (key TEXT PRIMARY KEY, value TEXT NOT NULL);`,
		`CREATE TABLE generated_surfaces (id INTEGER PRIMARY KEY AUTOINCREMENT, source_path TEXT NOT NULL, target_key TEXT NOT NULL, compile_rule TEXT NOT NULL, compiled_at TEXT NOT NULL, compiler_version TEXT NOT NULL, status TEXT NOT NULL, data TEXT NOT NULL, UNIQUE(source_path, target_key));`,
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
		`CREATE TABLE activation_rules_mirror (id INTEGER PRIMARY KEY AUTOINCREMENT, rule_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE transaction_templates_ext (id INTEGER PRIMARY KEY AUTOINCREMENT, template_name TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE distributed_locks (id INTEGER PRIMARY KEY AUTOINCREMENT, lock_name TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE multi_agent_coordination (id INTEGER PRIMARY KEY AUTOINCREMENT, rule_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE async_job_lifecycle (id INTEGER PRIMARY KEY AUTOINCREMENT, state TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
		`CREATE TABLE capability_checkpoints (id INTEGER PRIMARY KEY AUTOINCREMENT, checkpoint_id TEXT, content TEXT NOT NULL, created_at TEXT DEFAULT (datetime('now')), updated_at TEXT DEFAULT (datetime('now')));`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}

func compileStructuredRuntimeSources(repo string, db *sql.DB) error {
	phaseMachine, err := readRuntimeYAMLMap(repo, "runtime/phases/phase-machine.yaml")
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

	obligationLedger, err := readRuntimeYAMLMap(repo, "runtime/obligations/obligation-ledger.yaml")
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

	blockingGates, err := readRuntimeYAMLMap(repo, "runtime/gates/blocking-gates.yaml")
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

	transactionMachine, err := readRuntimeYAMLMap(repo, "runtime/transactions/transaction-machine.yaml")
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

	activationRules, err := readRuntimeYAMLMap(repo, "runtime/router/activation-rules.yaml")
	if err != nil {
		return err
	}
	for index, rule := range runtimeSlice(activationRules["core_bootstrap"]) {
		if _, err := db.Exec(`INSERT INTO core_bootstrap_rules (rule_id, ordinal) VALUES (?, ?)`, runtimeString(rule), index); err != nil {
			return err
		}
	}
	for _, rule := range runtimeSliceOfMaps(activationRules["rules"]) {
		activation := runtimeMap(rule["activation"])
		load := runtimeMap(rule["load"])
		if _, err := db.Exec(`INSERT INTO activation_rules (rule_id, description, activation_when, load_strategy, load_priority, load_estimated_tokens, load_source, metadata) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			runtimeString(rule["rule_id"]), runtimeString(rule["description"]), runtimeJSON(activation["when"]), runtimeDefaultString(load["strategy"], "lazy"), runtimeDefaultString(load["priority"], "P2"), runtimeInt(load["estimated_tokens"]), runtimeString(load["source"]), runtimeJSON(rule["metadata"]),
		); err != nil {
			return err
		}
	}

	discoveryCheckpoints, err := readRuntimeYAMLMap(repo, "runtime/discovery/capability-checkpoints.yaml")
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

	configs := []struct {
		rel     string
		table   string
		idCol   string
		listKey string
		idKeys  []string
	}{
		{"runtime/budget/token-budget.yaml", "runtime_budget", "model_name", "per_model", []string{"name", "model"}},
		{"runtime/context/ttl-policy.yaml", "context_ttl_policy", "ttl_type", "rules", []string{"name", "id"}},
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
		{"runtime/scheduler/execution-queue.yaml", "execution_queue", "queue_name", "queue_structure", []string{"name", "id"}},
		{"runtime/scheduler/priority-scheduler.yaml", "priority_scheduler", "priority_level", "levels", []string{"name", "level"}},
		{"runtime/router/activation-rules.yaml", "activation_rules_mirror", "rule_id", "rules", []string{"rule_id", "name"}},
		{"runtime/transactions/transaction-machine.yaml", "transaction_templates_ext", "template_name", "transaction_templates", []string{"id", "name"}},
	}
	for _, config := range configs {
		if err := insertRuntimeConfigRows(repo, db, config.rel, config.table, config.idCol, config.listKey, config.idKeys); err != nil {
			return err
		}
	}
	if err := insertTransactionTemplateRows(repo, db); err != nil {
		return err
	}
	if err := insertCompilerRuleRows(repo, db); err != nil {
		return err
	}
	return nil
}

func insertRuntimeConfigRows(repo string, db *sql.DB, rel string, table string, idCol string, listKey string, idKeys []string) error {
	config, err := readRuntimeYAMLMap(repo, rel)
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

func insertTransactionTemplateRows(repo string, db *sql.DB) error {
	config, err := readRuntimeYAMLMap(repo, "runtime/transactions/transaction-templates.yaml")
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

func insertCompilerRuleRows(repo string, db *sql.DB) error {
	var config compilerMappingFile
	if err := readRuntimeYAML(repo, "runtime/compiler/compiler-rules.yaml", &config); err != nil {
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

func compileProseRuntimeSources(repo string, db *sql.DB) error {
	var config compilerMappingFile
	if err := readRuntimeYAML(repo, "runtime/compiler/compiler-rules.yaml", &config); err != nil {
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
	return nil
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
