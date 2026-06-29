# Stage 5: AI Integration

## Goal

Generate stable summaries with an LLM API.

## Tasks

- Define summary data structures.
- Call the LLM API from backend code.
- Validate structured output.
- Add retry and rate-limit handling.
- Keep secrets out of source control.

## Key Concepts

- Treat model output as untrusted external data.
- Prefer structured output when downstream code depends on fields.
- Retry only bounded, likely transient failures.
- Keep prompts, schemas, and error handling reviewable.

