# BitwiseLearn V2 — Backend Rewrite Plan (Cost-Optimized, Autoscaling 100 → 10,000 concurrent)

**Status:** Phases 0-4 implemented in `apps/go-server` (Echo + Mongo driver), verified against live MongoDB/RabbitMQ/Piston. Phase 5 (Azure Container Apps Bicep + GitHub Actions CI/CD) written in `apps/go-server/infra/` and `.github/workflows/deploy-go-server.yml`, not yet deployed to a live subscription.
**Decision recorded:** Full rewrite in Go; RabbitMQ kept (not Azure Service Bus); storage moved to Azure Blob (not S3/Cloudinary).

---

## 1. Stack selection (and why)

The current backend is a Python **FastAPI monolith + MongoDB + RabbitMQ + Piston** (~5,600 LOC, 14 routers, 26 models). The workload is **IO-bound**; the only CPU-heavy work (code execution) is already offloaded to Piston.

For a rewrite whose primary goal is **minimum cost of operation under autoscaling**, the deciding metrics are: memory per replica (Azure Container Apps bills memory-GiB-seconds), cold-start time (enables scale-to-zero without latency penalty), and requests/sec per vCPU (fewer replicas at peak).

### Chosen stack: **Go (Golang)**

| Layer | Choice | Why this one |
|---|---|---|
| Language | **Go 1.22+** | ~15–40 MB idle RAM vs ~150–250 MB for Python; ~15 MB static-binary image; sub-100 ms cold start; goroutines give massive IO concurrency per core → lowest Azure cost per request |
| HTTP framework | **Echo** (or **Chi**) on `net/http` | Standard-library-compatible, HTTP/2, mature middleware, safe under autoscaling. (Fiber/fasthttp is ~10–15% faster but non-stdlib; use only if you want max throughput) |
| Database | **MongoDB** via `mongo-go-driver` | **Keep Mongo** — the data is document-shaped and the whole app already uses it. **No data migration, no Postgres.** |
| Auth | `golang-jwt/jwt/v5` | Re-implement the **exact same JWT scheme** so existing tokens and the frontend keep working unchanged |
| Validation | `go-playground/validator/v10` | Replaces Pydantic request validation |
| Cache/session | `redis/go-redis/v9` → **Azure Cache for Redis** | Cache read-heavy lists, rate-limit counters, optional sessions |
| Async jobs | `rabbitmq/amqp091-go` **or** Azure Service Bus SDK | Keep the existing queue/worker pattern; worker is a separate Go binary |
| Object storage | `Azure/azure-sdk-for-go` (azblob) | **Azure Blob** replaces S3/boto3 (aligns billing); keep S3 SDK only if you must |
| Email | SendGrid / Azure Communication Services SDK | **Never blocking** — always via the worker |
| Excel bulk upload | `qax-os/excelize` (streaming) | Replaces openpyxl; streams rows, low memory |
| Code execution | `net/http` client + `errgroup` | Reuse one pooled client; run test cases **in parallel** (goroutines), not serially |
| Config | stdlib `os` + `caarlos0/env` (or Viper) | 12-factor env config, same vars as today |
| Container | distroless/scratch, multi-stage | ~15 MB image, fast pulls → faster scale-out |

**Alternative if team velocity / one-language-with-frontend matters more than absolute cost:** **TypeScript + Fastify (or NestJS)** — shares the language with the Next.js frontend (shared DTO types), larger ecosystem, faster to write. Costs slightly more RAM (~80–150 MB/replica) than Go but still autoscales cheaply. *Go remains the lower-cost choice; this is the trade-off if you optimize for developer throughput instead.*

> Honest note: a rewrite is the most expensive and highest-risk path to the cost goal — the same 100→10k target is reachable by fixing the existing FastAPI defects. This plan proceeds with the rewrite as requested, structured to de-risk it (strangler migration, API parity, no data migration).

---

## 2. Target architecture

```
                         INTERNET
                            │
              ┌─────────────────────────────┐
              │ Azure Front Door + WAF + CDN│
              └─────────────────────────────┘
                            │
        ┌──────────────────────────────────────────┐
        │        AZURE CONTAINER APPS (KEDA)        │
        ├──────────────────────────────────────────┤
        │  api        (Go, Echo)      min2  max25   │  HTTP-concurrency scaler
        │  worker     (Go)            min0  max10   │  queue-length scaler
        │  code-exec  (Go + Piston)   min1  max15   │  CPU scaler, capped
        │  frontend   (Next.js)       min1  maxN    │
        └──────────────────────────────────────────┘
             │            │              │
   ┌─────────┘            │              └──────────┐
   ▼                      ▼                         ▼
MongoDB (Atlas/      Azure Cache for         Azure Blob Storage
 Cosmos serverless)   Redis                   (files, reports)
   │
   ▼
RabbitMQ / Azure Service Bus  ──►  Email (SendGrid / Azure Comm Services)
```

**Stateless by design** (critical for clean autoscaling): JWT auth carries identity, no server-side session; any shared state (rate-limit counters, caches) lives in Redis. Replicas are interchangeable → KEDA can add/remove them freely.

---

## 3. Module mapping (FastAPI → Go packages)

Repository layout for the new service (`apps/go-server/`):

```
apps/go-server/
├── cmd/
│   ├── api/main.go            # API entrypoint (Echo)
│   └── worker/main.go         # Queue consumer entrypoint
├── internal/
│   ├── config/                # env config  (was config.py)
│   ├── auth/                  # JWT, OTP, password (bcrypt), reset tokens (was utils/ + middleware/)
│   ├── middleware/            # JWT guard + role guards (was middleware/auth.py)
│   ├── db/                    # mongo client, index sync (was config.py connect_to_mongo)
│   ├── models/                # 26 structs + bson tags (was models/*.py)
│   ├── handlers/              # one file per current router (see table)
│   ├── services/              # email, blob, redis, piston, queue (was services/*)
│   └── response/              # api_response equivalent — MUST match current JSON shape
└── Dockerfile
```

| Current router (LOC) | New Go handler package | Migration order | Notes |
|---|---|---|---|
| `dsa_problem.py` (406) | `handlers/problem` | **1** | Read-heavy public lists; fix N+1 with `$in`; cache in Redis |
| `course.py` (981) | `handlers/course` | **2** | Largest; batch the per-section queries |
| `auth.py` (282) | `handlers/auth` | **3** | Keep JWT/bcrypt scheme byte-compatible |
| `student/teacher/institution/vendor/batch/admin` (≈900) | `handlers/identity` | 4 | CRUD; email creation → **publish to worker** |
| `assessment.py` (485) | `handlers/assessment` | 5 | |
| `code_runner.py` (223) | `handlers/code` (+ code-exec app) | 6 | Parallelize test cases; pooled Piston client |
| `report.py` (266) | `handlers/report` | 7 | Already well-optimized; mirror the `$in`+gather pattern |
| `bulk_upload.py` (295) | `handlers/bulkupload` | 8 | Streaming excelize; `insertMany` |
| `contact.py` (21) | `handlers/contact` | 9 | Email via worker |

**API parity rule:** the `internal/response` package must reproduce the exact JSON envelope of the current `utils/api_response.py` (`statusCode`, `message`, `data`, `error`) so the **Next.js frontend needs zero changes**.

---

## 4. Migration strategy — Strangler Fig (no big-bang)

1. **Stand up Go `api` alongside** the FastAPI monolith, same MongoDB, same JWT secret.
2. **Route incrementally at the edge** (Front Door / gateway path rules): move one endpoint group at a time to Go, starting with read-only public routes (lowest risk), ending with writes/reports/code-exec.
3. **Contract-test each module** against the current API responses before cutover (golden-response tests; the existing `pytest` suite in `apps/python-server/tests/` defines expected behavior to port into Go table tests).
4. **No data migration** — both stacks read/write the same Mongo collections during transition.
5. **Decommission FastAPI** only after all groups are migrated and parity-verified.

This keeps the system live throughout and lets you roll back any single endpoint group instantly.

---

## 5. Code-level optimizations baked into the rewrite

Carried over from the earlier review, made structural so they can't regress:
- **No blocking IO on request path** — email + blob uploads go through the worker / async SDKs (Go has no event loop to block, but offloading email keeps p99 low).
- **MongoDB indexes** created at startup for every foreign-key field queried (`problem_id`, `section_id`, `course_id`, `batch_id`, `assessment_id`, etc.) — eliminates collection scans, lowers DB tier.
- **Connection pooling** for Mongo, Redis, and the Piston HTTP client (single shared client, not per-call).
- **Parallel code-exec** with `errgroup` + a bounded semaphore so a submission's test cases run concurrently without overwhelming Piston.
- **Pagination everywhere** (`skip/limit`) on list endpoints; no unbounded `find_all().to_list()`.
- **Redis caching** with short TTLs on the read-mostly catalogs (problems, courses, leaderboards).

---

## 6. Azure Container Apps autoscaling spec

**Plan:** Consumption (serverless). Bills per vCPU-second + memory-GiB-second only while replicas run. Stateless replicas + KEDA = scale on real demand.

### 6.1 `api` (Echo)
```yaml
# Container Apps - api
properties:
  configuration:
    ingress: { external: true, targetPort: 8080, transport: auto }  # HTTP/2
  template:
    containers:
      - name: api
        image: <acr>.azurecr.io/bitwise-api:${TAG}
        resources: { cpu: 0.25, memory: 0.5Gi }   # Go runs comfortably here
        probes:
          - type: Liveness  { httpGet: { path: /health, port: 8080 }, periodSeconds: 15 }
          - type: Readiness { httpGet: { path: /ready,  port: 8080 }, periodSeconds: 10 }  # checks Mongo
          - type: Startup   { httpGet: { path: /health, port: 8080 }, failureThreshold: 10 }
    scale:
      minReplicas: 2        # avoid cold start on the critical path; set 0 for max savings off-hours
      maxReplicas: 25
      rules:
        - name: http-concurrency
          http: { metadata: { concurrentRequests: "100" } }   # add a replica per ~100 in-flight reqs
        - name: cpu
          custom: { type: cpu, metadata: { type: Utilization, value: "70" } }
```

### 6.2 `worker` (queue consumer) — scale to zero
```yaml
    scale:
      minReplicas: 0        # nothing in queue = no cost
      maxReplicas: 10
      rules:
        - name: queue-depth
          custom:
            type: rabbitmq        # or azure-servicebus
            metadata: { queueName: jobs, mode: QueueLength, value: "20" }
            auth: [{ secretRef: rabbitmq-conn, triggerParameter: host }]
```

### 6.3 `code-exec` (+ Piston) — isolated, CPU-bound, capped
```yaml
    containers:
      - { name: code-exec, resources: { cpu: 1.0, memory: 2.0Gi } }   # compilers need headroom
    scale:
      minReplicas: 1
      maxReplicas: 15        # hard cap protects cost; queue submissions beyond this
      rules:
        - name: cpu
          custom: { type: cpu, metadata: { type: Utilization, value: "65" } }
```

### 6.4 Sizing rationale (why this hits 10k concurrent cheaply)
- 10,000 **concurrent users** ≈ **1,000–2,000 RPS** peak (with normal think-time), almost all light IO-bound calls.
- A 0.25-vCPU Go replica serves several hundred such RPS → **~6–12 `api` replicas at peak**, scaling back to 2 (or 0) when idle.
- The CPU bottleneck is **only** code execution, which is a separate app with its own cap — it never starves the API.
- Result: you pay for ~2 small replicas most of the time, burst to ~a dozen briefly, and the heavy compiler pool scales independently. This is the minimum-cost shape for the workload.

---

## 7. CI/CD & observability (unchanged-in-spirit, lighter)
- **GitHub Actions** → build static Go binary → push to **Azure Container Registry** → `az containerapp update` (revisions, no Helm/AKS needed at this scale).
- **Application Insights + Azure Monitor** for traces/metrics/alerts; add Prometheus/Grafana only if/when you outgrow it.

---

## 8. Phased rollout & estimated effort
| Phase | Scope | Concurrency target |
|---|---|---|
| 0 | Scaffold Go service, config/db/auth/response, Dockerfile, ACA env + autoscale | — |
| 1 | Migrate problem + course read endpoints behind edge routing | ~500 |
| 2 | Auth + identity CRUD + worker (email/jobs) | ~3,000 |
| 3 | Assessment, report, bulk upload | ~6,000 |
| 4 | Code-exec service + Piston isolation; decommission FastAPI | ~10,000 |

Each phase is independently shippable and reversible (strangler routing).

---

## 9. Risks & mitigations
- **Rewrite cost/time** — largest risk. Mitigated by strangler migration (ship value continuously, never a frozen big-bang).
- **API drift breaking the frontend** — mitigated by the strict `response` parity package + golden-response contract tests ported from the existing pytest suite.
- **JWT incompatibility** — mitigated by reusing the same secret + claim structure so existing sessions survive cutover.
- **Hidden business logic in the 981-line course router** — audit and port with tests before routing traffic.

---

## 10. Open decisions before coding starts
1. **Queue:** keep RabbitMQ, or switch to **Azure Service Bus** (managed, less to operate)?
2. **Storage:** move S3 → **Azure Blob**, or keep S3 during transition?
3. **`api` minReplicas:** `2` (no cold-start, slightly higher floor cost) or `0` off-hours (lowest cost, cold-start risk)?
4. **HTTP framework:** Echo/Chi (stdlib-safe, recommended) vs Fiber (max throughput)?
5. Confirm **Go** vs the **TypeScript/Fastify** alternative for the rewrite.
