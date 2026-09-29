# ADR 0004: Credential Storage and Secrets Management Strategy

## Status

Accepted

## Context

UniGo operates across macOS, Linux, and Windows environments, providing CLI workflows, self-update capabilities, and network interactions. Configuration values and sensitive operational tokens (such as proxy credentials, access tokens, and environment parameters) must be managed securely on local client filesystems.

Key challenges include:

1. **Local File Persistence:** Configuration files persisted to disk (`unigo.toml`) must maintain strict permission boundaries to prevent unauthorized inspection by non-privileged accounts on the same system.
2. **Crash Resilience and Concurrency:** Concurrently running processes or sudden process termination during configuration writes can truncate target files if written directly in place without atomic rename semantics.
3. **Cross-Platform Credential Stores:** Desktop operating systems provide native secret vaults (macOS Keychain, Windows Credential Manager, Linux Secret Service via D-Bus), while minimal CLI and containerized environments rely on environment variables and secure filesystem stores.

## Decision

We adopt a two-phase credential storage and protection architecture:

### Phase 1: Baseline Security & Filesystem Hardening (Current Implementation)

1. **Strict File Permissions:** Configuration files are explicitly created and maintained with `0600` permissions (`-rw-------`), ensuring access is restricted strictly to the current operating system user.
2. **Atomic Configuration Persistence:** All configuration updates use temporary file generation with atomic rename semantics, preventing partial writes and file corruption.
3. **Data Sanitization in Logs and Memory:** Credentials and authentication tokens are excluded from structured loggers (`zerolog`), debug outputs, and crash diagnostics.

### Phase 2: Platform-Native Keyring Integration (Evolutionary Roadmap)

1. **Abstract Secret Storage Interface:** Introduce a lightweight storage contract `pkg/secret.Store` exposing `Get(key string)`, `Set(key string, val []byte)`, and `Delete(key string)`.
2. **OS Keychain Providers:**
   - Under macOS, delegate storage to Apple Keychain Services.
   - Under Windows, delegate to Windows Credential Manager.
   - Under Linux desktop environments, communicate with the Secret Service D-Bus interface.
3. **Graceful Fallback:** If the system credential service is unavailable (e.g. headless Linux servers, Docker containers, or minimal CI runners), gracefully fall back to encrypted local storage derived from user-scoped machine keys or environment variables (`UNIGO_PROXY_PASSWORD`).

## Consequences

### Positive

- **Transparent Defense-in-Depth:** Filesystem permissions (`0600`) and sanitized log pipelines protect credentials against inspection on multi-user hosts.
- **Portability:** UniGo executes reliably in headless environments, CI pipelines, and all major desktop platforms without mandatory external C library dependencies.
- **Clear Evolution Path:** Developers and security auditors have an unambiguous architectural blueprint for native keyring integration.

### Negative

- Native OS keyring integrations require platform-specific API shims, which will be implemented incrementally in Phase 2 to preserve pure-Go cross-compilation guarantees (`CGO_ENABLED=0`).
