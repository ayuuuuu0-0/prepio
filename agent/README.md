# Agent Context

This file is an index, not a rulebook. **The rules live only in `.ai/`** so they cannot drift apart.

## Read order

1. [.ai/PRODUCT.MD](../.ai/PRODUCT.MD): what we are building and the lesson UX contract
2. [.ai/ARCHITECTURE.MD](../.ai/ARCHITECTURE.MD): domain model, ownership, grading rules, standards
3. [.ai/CONTENT_SYSTEM.MD](../.ai/CONTENT_SYSTEM.MD): lessons, steps, skills, mastery, content-as-code
4. [.ai/EXECUTION.MD](../.ai/EXECUTION.MD): the current phase, task order, gates, out-of-scope

If a proposed change conflicts with any of them, stop and explain the conflict. If a document is wrong, follow the Change Process in EXECUTION.MD and update the documents in the same change.

## Current phase

See "CURRENT PHASE" in `.ai/EXECUTION.MD`.

## Practical rules for edits

- Start from the smallest concrete file or symbol that owns the behavior
- Prefer the nearest test or call site that can disconfirm the current hypothesis
- After the first substantive edit, validate the touched slice before broadening scope
- Remove whatever your change makes dead (code, routes, env vars, flags, tests, docs)
- Keep the root README factual and concise
