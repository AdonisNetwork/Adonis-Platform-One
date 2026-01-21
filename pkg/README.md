# 📦 pkg — Public Reusable Modules for Adonis Platform One

This directory contains **public-facing Go modules and libraries** intended for reuse outside of the core platform. Unlike `internal/`, code placed under `pkg/` is considered stable, versioned, and consumable by external projects, extensions, OEM partners, and SDKs.

This layout follows the established conventions of large-scale Go service platforms (Kubernetes, Temporal, HashiCorp, etc).

---

# 🎯 Purpose of `pkg/`

The `pkg/` directory exists to:

✔ Enable modular reuse of domain logic  
✔ Provide building blocks for extensions and integrations  
✔ Allow OEM and Enterprise customers to embed A1 modules  
✔ Support SDK and language bindings  
✔ Enable public packages without exposing the entire runtime  

---

# 🏗️ 1. Planned Module Structure

The module structure is designed for growth:

```
pkg/
 ├── sdk/              # SDKs for external integrations
 ├── client/           # API clients for Go and other languages
 ├── schema/           # Task, agent, and financial schema models
 ├── protocol/         # Agent protocol definitions & message types
 ├── adapters/         # Connectors & data adapters
 ├── workflows/        # Reusable standardized execution workflows
 ├── plugins/          # Plugin interfaces for domain-specific modules
 ├── crypto/           # Signing, keys, hashing utilities
 ├── credits/          # ADON Credits & Billing utilities
 ├── compliance/       # Rule bundles for compliance extensions
 └── utils/            # Lightweight helper libraries
```

---

# 🧱 2. Key Concepts

### **📌 Public vs Private Code**

| Location | Visibility | Purpose |
|---|---|---|
| `internal/` | Private | Core runtime, orchestration, AI pipeline |
| `pkg/` | Public | Reusable components, SDKs, OEM modules |

---

# 🧩 3. Example Use Cases

These modules may be imported by:

### ✔ AI/Agent Developers
to implement new agent roles or workflows

### ✔ Enterprise Integrators
to integrate A1 with:

- ERP systems
- Financial systems
- Industrial stacks
- Healthcare data pipelines (for Nevacoin use case)

### ✔ OEM Hardware Partners
(e.g., IoT + Medical Devices)

> e.g., integrating a **QRMA health analyzer** with A1 workflows

### ✔ Blockchain Ecosystem
(e.g., ADON credits + Smart settlement future roadmap)

---

# 🔌 4. Future SDK Targets

Planned SDK language targets:

| Language | Status |
|---|---|
| Go | Native |
| Python | Planned |
| TypeScript | Planned |
| Rust | Planned |
| C (WASM) | Future |
| Swift/Kotlin | Mobile bridge (future) |

This supports both:

🟦 Cloud execution  
🟩 Edge/IoT execution (QRMA / Medical devices)

---

# 📜 5. Stability & Version Guarantees

Packages in `pkg/` will follow semantic versioning:

```
v0.x  → unstable / MVP
v1.x  → stable API
v2.x+ → enterprise / OEM stability
```

---

# 🧃 6. Sample Imports

Example usage from Go:

```go
import (
    "github.com/adonisnetwork/a1/pkg/schema"
    "github.com/adonisnetwork/a1/pkg/protocol"
)
```

---

# 🛑 Important Notes

- `pkg/` **does not** contain core orchestration logic  
- `pkg/` is safe for embedding in **OEM devices**
- `pkg/` is aligned with **Enterprise product strategy**
- `pkg/` is required for **NIW demonstration → National Economic Interest**

---

# 📈 NIW Relevance Justification

This structure supports:

✓ Technological Merit  
✓ Substantial Economic Impact  
✓ National Interest through:

> “Creation of a modular AI execution platform enabling secure AI-driven digital labor pipelines, verifiable task execution, and hardware integration pathways across healthcare, finance, and cybersecurity domains.”

---

# 🔐 Licensing Implications

Code here is compatible with:

- `AOL-1.0` (Non-commercial open license)
- `ACLA-1.0` (Commercial license)
- `AOLA-1.0` (OEM license)

All three are supported.

---

# 📦 Status

> `pkg/` is currently in scaffolding stage.  
> Modules will be populated incrementally during MVP execution.

