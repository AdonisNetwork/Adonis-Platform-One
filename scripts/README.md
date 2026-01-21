# 🛠️ scripts — Development & Automation Tools for Adonis Platform One

This directory contains **automation scripts** used for development, testing, deployment, CI/CD, and infrastructure operations. These scripts help maintain consistent workflows across cloud environments, local development, and future Edge/IoT runtimes.

The structure is designed for MVP readiness and future enterprise scaling.

---

# 📁 1. Directory Structure

```
scripts/
 ├── dev/                # Local development helpers
 │    ├── run_api.sh
 │    ├── run_worker.sh
 │    ├── lint.sh
 │    └── format.sh
 │
 ├── db/                 # DB utilities
 │    ├── migrate.sh
 │    ├── seed.sh
 │    └── reset.sh
 │
 ├── docker/             # Docker & Compose tools
 │    ├── build.sh
 │    ├── up.sh
 │    ├── down.sh
 │    └── logs.sh
 │
 ├── ci/                 # CI/CD tools (GitHub Actions / GitLab / Jenkins)
 │    ├── test.sh
 │    ├── security_scan.sh
 │    ├── build_release.sh
 │    └── version_bump.sh
 │
 ├── utils/              # General-purpose utilities
 │    ├── env_check.sh
 │    ├── health_check.sh
 │    └── cleanup.sh
 │
 └── README.md           # This documentation
```

---

# ⚙️ 2. Development Tools (Local)

Scripts under `scripts/dev/` help contributors run the platform easily:

### ```run_api.sh```
Runs the Go API service with hot-reload (if using air):

```bash
#!/bin/bash
air api
```

### ```run_worker.sh```
Starts the Python worker:

```bash
#!/bin/bash
python3 cmd/worker/worker.py
```

### ```lint.sh```
Runs all linters (Go, Python):

```bash
#!/bin/bash
golangci-lint run
flake8 .
```

---

# 🛢️ 3. Database Utilities

### ```migrate.sh```
Apply SQL migrations:

```bash
#!/bin/bash
psql "$A1_DB" -f infra/migrations.sql
```

### ```reset.sh```
Drop + recreate the development database.

---

# 🐳 4. Docker / Compose Scripts

This ensures consistent environment setup:

### ```build.sh```
```bash
#!/bin/bash
docker compose -f infra/docker-compose.yml build
```

### ```up.sh```
Bring everything online:

```bash
docker compose -f infra/docker-compose.yml up -d
```

### ```down.sh```
```bash
docker compose -f infra/docker-compose.yml down
```

---

# 🔐 5. CI/CD Automation

These help with release management and security:

### ```security_scan.sh```
Runs Trivy + secret scan:

```bash
trivy fs .
gitleaks detect
```

### ```version_bump.sh```
Auto-update v0.x → v0.x+1

---

# 🧪 6. Test Automation

```bash
./scripts/ci/test.sh
```

Runs:

- Go unit tests
- Python unit tests
- API contract tests (future)

---

# 🌐 7. Environment Check Tools

Useful for deployment to cloud or on-prem:

```bash
./scripts/utils/env_check.sh
```

Checks for:

- Redis availability
- Postgres connection
- Required environment variables
- Worker queue connectivity

---

# 📦 8. Why This Matters (NIW + Investors)

A strong automation folder shows:

✓ Engineering maturity  
✓ Deployment readiness  
✓ Reliability & reproducibility  
✓ Enterprise compliance mindset  
✓ Ability to scale multi-domain systems  

این دقیقاً همان چیزی است که NIW آفیسرها و اینوستورها می‌خواهند ببینند.

---

# 📌 Status

> Scripts are placeholders until full automation is implemented in MVP-1.

