# Mobile-First Economics

Cost model: mobile and constrained-context defects are cheaper to solve before component structure and interaction hierarchy solidify.

## Cost Curve

```text
Design stage       1x
Implementation     5x
QA                20x
Production       100x
```

The multipliers are illustrative. The reusable claim is directional: late responsive fixes are expensive because they often require structure, density, content priority, interaction, and validation changes.

## Economic Rule

When a user-visible UI has a declared mobile or narrow-mobile contract:

- desktop-only evidence leaves unpriced risk;
- fixed/sticky behavior should be validated before QA handoff;
- responsive evidence should be collected while layout structure is still cheap to change;
- deferral should name the unsupported or deferred render context explicitly.

## Boundary

This model does not mandate mobile-first for an explicitly desktop-only product. It requires the contract to be explicit so the cost is visible.
