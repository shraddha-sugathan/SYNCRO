# Syncro Billing Engine

Syncro is a lightweight, full-stack subscription billing engine built for high performance and rapid deployment. It provides Role-Based Access Control (RBAC), automated recurring billing, and a unified administrative dashboard.

## Architecture Diagram

```mermaid
graph TD
    subgraph Frontend
        Vue[Vue.js 3 SPA]
    end

    subgraph Go Server
        Router[Chi Router]
        Auth[JWT Middleware]
        API[Handlers: Auth, Products, Subs, Invoices]
        Worker[Background Billing Engine]
    end

    subgraph Database
        PG[(PostgreSQL)]
    end

    Vue -->|HTTP REST + JWT| Router
    Router --> Auth
    Auth --> API
    API -->|Read/Write| PG
    Worker -->|Sweeps expiries every 60s| PG