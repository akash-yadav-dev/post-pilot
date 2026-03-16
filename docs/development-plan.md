# Post Pilot Development Plan

## Objective
Build and operate a production-grade social publishing platform with two deployable services:
- API service for authentication, scheduling, and post management.
- Worker service for asynchronous publishing and retries.

## Phase 0 - Foundation and Standards
### Goals
- Define engineering standards and operational expectations before feature work.

### Tasks
- Confirm module and workspace baseline (`go.mod`, `go.work`) and enforce single root module usage.
- Establish coding standards for package boundaries (`apps/api/internal`, `apps/worker/internal`, `packages/*`).
- Define branch strategy (`main`, `develop`, feature branches) and PR checklist.
- Configure static checks in CI: `go fmt`, `go vet`, `go test`.
- Add conventional commit style and changelog generation process.

### Deliverables
- CI pipeline with required checks.
- Repository standards document.
- Baseline project bootstrap that compiles end-to-end.

## Phase 1 - Configuration and Secrets Management
### Goals
- Make runtime configuration explicit, validated, and environment-specific.

### Tasks
- Introduce strict config validation for API and worker startup.
- Split required vs optional environment variables.
- Add `.env.example` templates for local/dev/staging/prod.
- Move secrets to secret store strategy (Kubernetes Secret, Vault, or cloud secret manager).
- Add startup checks that fail fast on invalid or missing critical configuration.

### Deliverables
- Typed config loaders for API and worker.
- Documented environment variables and defaults.
- Secret rotation and ownership runbook.

## Phase 2 - Data Layer and Migrations
### Goals
- Establish robust persistence and schema evolution.

### Tasks
- Implement PostgreSQL connection pool and health checks.
- Add migration runner and migration naming policy.
- Create core tables: users, auth_tokens, posts, schedules, publish_attempts, dead_letter_jobs.
- Add indexes for scheduled dispatch and user/content lookup patterns.
- Define repository interfaces for users/posts/schedules.

### Deliverables
- Production-ready DB package with retry/backoff.
- Initial migration set and rollback scripts.
- Data access layer tests.

## Phase 3 - Identity and Access
### Goals
- Secure API access and enforce authorization boundaries.

### Tasks
- Implement user registration, login, refresh token flow.
- Add password hashing and password policy enforcement.
- Add JWT middleware with expiry and rotation strategy.
- Define role model (user/admin) and protected routes.
- Add brute force protection and lockout policy.

### Deliverables
- Auth handlers/services/repositories.
- Auth middleware with unit tests.
- Threat model notes for auth attack surfaces.

## Phase 4 - Post and Scheduling Domain
### Goals
- Build core post lifecycle and scheduling capabilities.

### Tasks
- Implement CRUD for drafts and scheduled posts.
- Add validation rules per social network constraints.
- Implement scheduling service and queue producer.
- Add timezone handling and idempotency keys for enqueue operations.
- Add API pagination, filtering, and sort support.

### Deliverables
- Stable post APIs.
- Scheduler integration with queue abstraction.
- Domain tests for schedule transitions.

## Phase 5 - Queueing and Worker Reliability
### Goals
- Ensure durable, observable async processing.

### Tasks
- Finalize queue abstraction and Redis-backed implementation.
- Implement worker loop with visibility timeout and retry policy.
- Add dead-letter queue and replay tooling.
- Implement idempotent publish job execution.
- Add backoff strategy with jitter for transient failures.

### Deliverables
- Worker service with graceful shutdown.
- Retry/dead-letter flow with operational commands.
- Integration tests for queue and worker behavior.

## Phase 6 - Social Provider Integrations
### Goals
- Integrate external APIs safely and consistently.

### Tasks
- Implement provider clients (Twitter, LinkedIn, Mastodon, Bluesky) with unified interface.
- Add token storage/encryption and refresh lifecycle.
- Add request signing, timeout, retries, and circuit breaking.
- Normalize provider error mapping for consistent handling.
- Add dry-run mode for non-production validation.

### Deliverables
- Provider adapters and publish services.
- Sandbox tests for each social network.
- Compatibility matrix (supported post types and limits).

## Phase 7 - API Contract and Developer Experience
### Goals
- Make APIs predictable and maintainable.

### Tasks
- Finalize versioned API routes and response envelope standards.
- Publish OpenAPI spec and generated client SDK baseline.
- Add request validation middleware and consistent error codes.
- Add local developer scripts for seed/migrate/test/dev-start.
- Add Makefile or task runner for common workflows.

### Deliverables
- API contract docs and examples.
- Developer onboarding workflow under 15 minutes.
- Contract tests for critical routes.

## Phase 8 - Observability and Operational Readiness
### Goals
- Support production troubleshooting and performance optimization.

### Tasks
- Add structured logging with correlation/request IDs.
- Expose metrics (HTTP latency, queue depth, job failures, publish success ratio).
- Add distributed tracing for API and worker jobs.
- Define SLOs and alert thresholds.
- Add health/readiness probes and dependency checks.

### Deliverables
- Dashboards and alert rules.
- Service-level objectives document.
- Incident response playbooks.

## Phase 9 - Security Hardening
### Goals
- Reduce security risk before broad rollout.

### Tasks
- Add dependency scanning and SAST in CI.
- Add input sanitization and strict payload validation.
- Enforce TLS and secure headers at ingress/API.
- Add audit logs for sensitive operations.
- Run penetration test and remediate findings.

### Deliverables
- Security checklist with completion evidence.
- Pen-test report and mitigation log.
- Compliance posture notes (as required by business context).

## Phase 10 - Deployment, Scaling, and Release Management
### Goals
- Achieve repeatable and safe deployments.

### Tasks
- Productionize Dockerfiles (multi-stage builds, minimal runtime image).
- Finalize Kubernetes manifests and resource sizing.
- Add rolling/blue-green deployment strategy.
- Add backup/restore and disaster recovery testing.
- Introduce release gates and canary verification.

### Deliverables
- CI/CD deployment pipeline.
- Verified rollback procedure.
- Capacity and performance test report.

## Phase 11 - Quality Engineering and Automation
### Goals
- Prevent regressions and improve release confidence.

### Tasks
- Build test pyramid: unit, integration, contract, end-to-end.
- Add deterministic test fixtures and fake provider stubs.
- Add load tests for scheduling and worker throughput.
- Add flaky test detection and quarantine workflow.
- Enforce minimum coverage thresholds for critical packages.

### Deliverables
- Automated quality gates.
- Performance baseline and trend tracking.
- Test strategy document.

## Phase 12 - Productization and Growth
### Goals
- Prepare platform for broader user adoption and feature expansion.

### Tasks
- Add billing/plan enforcement if needed.
- Add multi-tenant controls and usage quotas.
- Add analytics events for product insights.
- Add admin tooling for support operations.
- Build roadmap for advanced features (AI captions, campaign templates, approval workflows).

### Deliverables
- Scalable product architecture direction.
- Feature roadmap and dependency graph.
- Operational support model.

## Cross-Phase Governance
- Definition of done: code, tests, docs, metrics, and security checks complete.
- Weekly architecture review and risk tracking.
- Monthly reliability review against SLOs.
- Release retrospective and action item ownership.

## Suggested Milestone Cadence
- Milestone 1: Phases 0-3 complete (platform secure and bootstrapped).
- Milestone 2: Phases 4-6 complete (core product workflow and provider publishing).
- Milestone 3: Phases 7-10 complete (operationally production-ready).
- Milestone 4: Phases 11-12 complete (scale and product growth).
