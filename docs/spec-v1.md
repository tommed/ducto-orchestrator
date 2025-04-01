<!--suppress HtmlDeprecatedAttribute -->
<p align="right">
    <a href="https://github.com/tommed" title="See Project Ducto">
        <img src="../assets/ducto-logo-small.png" alt="A part of Project Ducto"/>
    </a>
</p>

# Ducto Orchestrator — Specification (v1)

## Purpose
Ducto-Orchestrator is a lightweight, embeddable, and extensible data stream processor designed to execute transformation pipelines using the Ducto-DSL as its core transformation language.

It is intended for:
- Streamlining JSON-like event transformations 
- Running on local CLI, GCP, AWS, Azure, and other environments 
- Future integration of feature flag evaluation and routing 
- Being lightweight and dependency-lean

---

## 🟣 Core Pipeline

```
[input] -> [pre-processors] -> [Ducto-DSL Transformer] -> [post-processors] -> [output]
```

### Inputs:
- stdin (CLI)
- GCP Pub/Sub message 
- HTTP Event (GCP)
- Future: SQS, S3, EventArc, EventBridge, Kafka

### Outputs:
- stdout (CLI)
- GCP Pub/Sub message
- Future: GCS, BigQuery, SNS, SQS, HTTP callback

### Optional:
- Attribute Projection (input fields projected to PubSub/SNS attributes)
- Error Collection (inject @dsl_errors)
- Debug Fields (inject @dsl_debug)

---
## 🟣 Planned Processors
- ducto-dsl (always present)
- Feature Flag Processor (planned)
- Enrichment Processors (planned)
- Debug / Logging Processor (optional)

---
## 🟣 Runtime Configuration
- Program file (DSL)
- Attribute Projection Rules 
- Input Source Configuration (pubsub, stdin, http, etc.)
- Output Target Configuration (pubsub, stdout, etc.)

---
## 🟣 Operator Compliance
- Fully compatible with ducto-dsl operators defined in spec-v1.md 
- Supports evolving DSL versions ("version": 1 mandatory)

---
## 🟣 Future Extensions
- Feature Flag integration (ducto-featureflag)
- Attribute Projection Language (simple DSL for selecting attribute fields)
- Cloud Specific Optimizations 
- Telemetry & Metrics 
- Secure Configuration (env, secret managers, IAM)
- Playground (WebAssembly + Vue3)

---
## Notes:
- The orchestrator is not tied to GCP; AWS, Azure, and local deployments are equally supported 
- Operators execute exactly as specified by ducto-dsl 
- Processors are composable and can be disabled or re-ordered for specialized editions 
- Avoids linking cloud SDKs unless strictly needed by the build target 
- Modular and OpenTelemetry ready

This document will evolve into versioned specifications as we proceed.