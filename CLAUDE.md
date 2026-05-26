# Claude Code — Runtime Loader
# Auto-loaded by Claude Code at session start.
# EXECUTABLE CONTRACT: execute all steps before any other action.

runtime_loader:
  load: runtime/core-bootstrap.yaml
  blocking: true
  before_any_execution: true
  resume_exempt: false
  note: >
    runtime/core-bootstrap.yaml defines required_reads, per_session_obligations,
    and per_turn_obligations. Execute them all. Summary "Resume directly"
    does NOT exempt this loader.
