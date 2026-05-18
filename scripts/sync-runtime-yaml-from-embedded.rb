#!/usr/bin/env ruby
# frozen_string_literal: true
# Restore canonical runtime/**/*.yaml from EmbeddedRuntimeData (one-way sync for missing files).

require "fileutils"
require "yaml"

ROOT = File.expand_path("..", __dir__)
require File.join(ROOT, "runtime/compiler/embedded_data.rb")

MAPPINGS = {
  "runtime/context/ttl-policy.yaml" => EmbeddedRuntimeData::CONTEXT_TTL_POLICY,
  "runtime/budget/token-budget.yaml" => EmbeddedRuntimeData::BUDGET_TOKEN_BUDGET,
  "runtime/guards/circuit-breaker.yaml" => EmbeddedRuntimeData::GUARDS_CIRCUIT_BREAKER,
  "runtime/guards/context-pollution.yaml" => EmbeddedRuntimeData::GUARDS_CONTEXT_POLLUTION,
  "runtime/health/context-health-score.yaml" => EmbeddedRuntimeData::HEALTH_CONTEXT_HEALTH_SCORE,
  "runtime/intelligence/intelligence-routing.yaml" => EmbeddedRuntimeData::INTELLIGENCE_INTELLIGENCE_ROUTING,
  "runtime/pipeline/session-lifecycle.yaml" => EmbeddedRuntimeData::PIPELINE_SESSION_LIFECYCLE,
  "runtime/pipeline/context-flow.yaml" => EmbeddedRuntimeData::PIPELINE_CONTEXT_FLOW,
  "runtime/pipeline/guard-chain.yaml" => EmbeddedRuntimeData::PIPELINE_GUARD_CHAIN,
  "runtime/pipeline/relevance-engine.yaml" => EmbeddedRuntimeData::PIPELINE_RELEVANCE_ENGINE,
  "runtime/prompt-artifacts/artifact-templates.yaml" => EmbeddedRuntimeData::PROMPT_ARTIFACTS_ARTIFACT_TEMPLATES,
  "runtime/prompt-artifacts/composition-rules.yaml" => EmbeddedRuntimeData::PROMPT_ARTIFACTS_COMPOSITION_RULES,
  "runtime/output-governance/language-policy.yaml" => EmbeddedRuntimeData::OUTPUT_GOVERNANCE_LANGUAGE_POLICY,
  "runtime/output-governance/output-rules.yaml" => EmbeddedRuntimeData::OUTPUT_GOVERNANCE_OUTPUT_RULES,
  "runtime/output-governance/governance-gates.yaml" => EmbeddedRuntimeData::OUTPUT_GOVERNANCE_GOVERNANCE_GATES,
  "runtime/distributed/distributed-locks.yaml" => EmbeddedRuntimeData::DISTRIBUTED_DISTRIBUTED_LOCKS,
  "runtime/distributed/multi-agent-coordination.yaml" => EmbeddedRuntimeData::DISTRIBUTED_MULTI_AGENT_COORDINATION,
  "runtime/distributed/async-job-lifecycle.yaml" => EmbeddedRuntimeData::DISTRIBUTED_ASYNC_JOB_LIFECYCLE,
  "runtime/router/activation-rules.yaml" => EmbeddedRuntimeData::ROUTER_ACTIVATION_RULES
}.freeze

README_SNIPPETS = {
  "runtime/pipeline" => "# Runtime Pipeline\n\nCanonical YAML: `session-lifecycle.yaml`, `context-flow.yaml`, `guard-chain.yaml`, `relevance-engine.yaml`.\n",
  "runtime/prompt-artifacts" => "# Prompt Artifacts\n\nCanonical YAML: `artifact-templates.yaml`, `composition-rules.yaml`.\n",
  "runtime/output-governance" => "# Output Governance\n\nCanonical YAML: `language-policy.yaml`, `output-rules.yaml`, `governance-gates.yaml`.\n",
  "runtime/distributed" => "# Distributed Runtime\n\nCanonical YAML: `distributed-locks.yaml`, `multi-agent-coordination.yaml`, `async-job-lifecycle.yaml`.\n",
  "runtime/context" => "# Context TTL\n\nCanonical YAML: `ttl-policy.yaml`. Prose: `governance/lifecycle/context-ttl-philosophy.md`.\n",
  "runtime/guards" => "# Runtime Guards\n\nCanonical YAML: `circuit-breaker.yaml`, `context-pollution.yaml`.\n",
  "runtime/intelligence" => "# Runtime Intelligence Routing\n\nCanonical YAML: `intelligence-routing.yaml`. Domain knowledge: `intelligence/`.\n",
  "runtime/router" => "# Runtime Router\n\nActivation: `activation-rules.yaml`, `activation-table.md`, `activation-engine.rb`. Routing prose: `governance/lifecycle/routing-philosophy.md`.\n"
}.freeze

MAPPINGS.each do |rel, data|
  path = File.join(ROOT, rel)
  FileUtils.mkdir_p(File.dirname(path))
  File.write(path, data.to_yaml(line_width: -1))
  puts "  ✓ #{rel}"
end

README_SNIPPETS.each do |dir, body|
  readme = File.join(ROOT, dir, "README.md")
  next if File.exist?(readme)

  FileUtils.mkdir_p(File.dirname(readme))
  File.write(readme, body)
  puts "  ✓ #{dir}/README.md"
end

puts "Sync complete (#{MAPPINGS.size} YAML files)"
