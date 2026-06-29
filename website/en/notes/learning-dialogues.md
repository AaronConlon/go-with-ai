# Learning Dialogues

This page mirrors the Chinese learning-dialogue log at a summary level.

The detailed dialogue record remains in Chinese because it preserves the learner's real questions and the coaching context. English notes should capture durable conclusions when a question becomes broadly reusable.

## Recording Rule

For each durable question, keep:

- The simplified question.
- The coaching explanation.
- A reusable conclusion.
- The next action.

## 2026-06-29 Go if syntax variants

Go `if` has three forms: basic condition (`if x > 0 {}`), if-else with `else` on the same line as `}`, and init-statement form (`if err := fn(); err != nil {}`). Key rules: no parentheses, braces always required, init-statement variables are scoped to the if block. See [Go if Syntax](/en/notes/stage-1-startup/go-if-syntax) for details.

## 2026-06-28 Go slice vs array

Go has two distinct collection types: `[]T` (slice, dynamic) and `[N]T` (array, fixed). Slices are the default. Arrays copy the full value when passed to a function; slices copy a small descriptor (24 bytes). This is a key difference from JavaScript, where there is no fixed-length array type.

See [Go Slice vs Array](/en/notes/stage-1-startup/go-slice-vs-array) for details.

## Current Source of Truth

The detailed Chinese version is available at [Learning Dialogue Records](/notes/learning-dialogues).

