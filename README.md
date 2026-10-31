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
    
    
Tech Stack
Backend: Go (1.20+) with chi router

Frontend: Vue.js 3 (Composition API)

Database: PostgreSQL

Security: Google OAuth, JWT (HMAC-SHA256), bcrypt

System Workflow
Authentication: Users can register standard accounts or authenticate via Google. New users default to the customer role.

State Management: The UI state (like active subscriptions) persists via localStorage to prevent duplicate checkout flows on page refreshes.

Automated Billing: A dedicated Goroutine polls the subscriptions table every 60 seconds. When a subscription period ends, it automatically generates a new pending invoice and extends the active billing period.

Admin Oversight: Administrators have access to real-time MRR (Monthly Recurring Revenue) calculations and comprehensive customer activity reports.

How to Run
Ensure PostgreSQL is running.

Start the Go server (the database schema auto-migrates on startup):

Bash
go run main.go
Open http://localhost:8080 in your browser.

