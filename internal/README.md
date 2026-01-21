# 🧩 Adonis Platform One — Internal Architecture Guide
This directory contains the **core backend implementation** of Adonis Platform One (A1), including the runtime engine, agent orchestration logic, execution workflows, AI/LLM integrations, financial components, and security modules.

The code under `internal/` is **not** exposed as a public library.  
It is used exclusively by the A1 platform services.

---

# 🏗️ 1. Directory Structure

```
internal/
 ├── runtime/           # Core execution engine (Go)
 ├── orchestration/     # Agent scheduler, queues, state machine
 ├── ai/                # LLM integrations, embeddings, reasoning modules
 ├── inference/         # Worker adapters (Groq, OpenAI, Anthropic, Local)
 ├── agents/            # Agent roles, behaviors, policies
 ├── events/            # Event bus, domain events, audit hooks
 ├── financial/         # Escrow, credits, ledger, payments
 ├── compliance/        # Rules, audit logs, policy engine
 ├── storage/           # Database interfaces (Postgres) + caches (Redis)
 ├── messaging/         # Message queues & job dispatchers
 ├── security/          # AuthN, AuthZ, signatures, encryption
 ├── utils/             # Shared utilities
 └── config/            # Runtime configuration models
```

---

# ⚙️ 2. Core Runtime

The A1 Runtime is implemented in **Go** for performance and safety.

### Responsibilities:
- Manage task lifecycle  
- Dispatch tasks to agents  
- Execute workflows with deterministic control  
- Track state transitions  
- Handle retries, failures, rollbacks  
- Integrate with AI reasoning and validation layers  

### Key Components:
- `TaskManager`  
- `WorkflowEngine`  
- `Scheduler`  
- `StateMachine`  
- `ExecutionContext`  

---

# 🤖 3. Agent Orchestration Layer

This layer coordinates:

- **Execution Agents** (Go routines)
- **AI-Supervised Agents** (Python / LLM reasoning)
- **Domain-Specific Agents** (WASM modules)
- **Financial Agents** (escrow verification, settlement)
- **Audit Agents** (compliance verification)

Agents interact through a **strict policy-based system**:

✓ Role-based behavior  
✓ Input/output schema validation  
✓ Policy-driven restrictions  
✓ Deterministic execution when required  

---

# 🧠 4. AI & LLM Integration (Python Layer)

This is where A1 becomes unique.

The `ai/` module includes:

- LLM orchestration (Groq / OpenAI / Anthropic / Local models)
- Prompt templates & structured outputs
- Chain-of-thought supervisor (private)
- Structured action planners
- Knowledge embeddings
- AI validation layer (verifies correctness of tasks)

Python modules in `worker/` call these components via:

- gRPC
- REST
- Message queues

---

# ⚡ 5. Inference Workers

Workers handle:

- AI model inference
- Skill execution (Python functions, WASM modules)
- Embedding generation
- Tool usage (web search, scraping, analysis)
- IoT/Edge device integration (future scope)

Workers are **stateless** and **scalable**.

A default worker is included:

```
cmd/worker/worker.py
```

---

# 💰 6. Financial Layer

The `financial/` module implements:

- **Escrow Engine**
- **ADON Credits**
- **Ledger entries**
- **Transaction signing**
- **Fee calculation**
- **Task-based payouts**
- **Refund & dispute logic**

These features also support:

- A1 Marketplace
- Enterprise billing
- OEM licensing

This part is important for **NIW** because it demonstrates:

> “A platform solving critical problems in financial transparency, digital execution trust, and AI-governed transactions.”

---

# 📜 7. Compliance, Audit, and Rules Engine

The `compliance/` module ensures:

- AI outputs are validated
- Execution steps are logged
- Domain-specific rules are applied
- Disputes can be resolved
- All actions are auditable

This is one of the **core innovations** that distinguishes A1 from other agent platforms.

---

# 💾 8. Storage Layer

The `storage/` module implements:

- Postgres driver + migrations
- Redis caching
- Job queues
- Distributed locks
- Search indexes (future)

All database access is abstracted through:

```
storage/repositories/
```

---

# 📬 9. Messaging Layer

Implements inter-service messaging:

- Job queues
- Event dispatching
- Worker scheduling
- Stream processing

Supports either:

- Redis Streams (**MVP**)
- NATS or Kafka (future enterprise versions)

---

# 🔐 10. Security Model

A1 follows:

- Zero-trust principles  
- Signed workflows  
- Permissioned agent operations  
- Key-based authentication  
- Encrypted storage for sensitive data  

Files in `security/` include:

- JWT manager  
- API key generator  
- Signature verifier  
- Encryption helpers  

---

# 🧪 11. Testing Framework

Tests live under:

```
/test
```

Coverage includes:

- Runtime tests
- Agent behavior tests
- AI validation tests
- Integration tests
- Financial model tests

---

# 🏁 Status

> Internal architecture is actively evolving.  
> All modules are placeholders until MVP-1 implementation.

