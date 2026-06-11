# Responsive Cost Curve

Cost model: responsive bugs compound across layout, interaction, validation, and product decisions.

## Cost Drivers

- Layout: fixed widths and duplicated shell constants spread through components.
- Interaction: sticky and fixed controls can hide actions or drift from content.
- Validation: each missing render context creates another evidence gap.
- Product: late mobile issues can force content priority and density decisions, not just CSS edits.

## Escalation Pressure

Responsive issues become economically risky when the affected surface includes:

- primary navigation or submit controls;
- identity, entitlement, payment, or irreversible actions;
- fixed, sticky, modal, sheet, or full-height scroll-root behavior;
- defect reports tied to viewport switching, reload, safe area, orientation, or browser chrome.

## Deferral Rule

Deferring a responsive issue is acceptable only when the release posture, unsupported context, owner, and follow-up evidence plan are explicit.
