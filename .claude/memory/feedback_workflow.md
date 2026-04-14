---
name: Workflow and commit rules
description: How the user wants commits, PRs, and in-progress work handled
type: feedback
originSessionId: 24a106a6-d389-4c1f-b68e-f7633a0973cd
---
Do not create commits unless the user explicitly asks in the current conversation.

**Why:** CLAUDE.md explicitly states this. The user commits manually.

**How to apply:** Leave all changes in the working tree. Never run `git commit` proactively. Never push or open PRs unless asked.
