\# Syncro Billing Engine



Syncro is a lightweight, full-stack subscription billing engine built for high performance and rapid deployment. It features Role-Based Access Control (RBAC), automated background billing, Google OAuth integration, and a unified administrative dashboard.



\## Architecture Diagram



```mermaid

graph TD

&#x20;   subgraph Frontend

&#x20;       Vue\[Vue.js 3 SPA]

&#x20;       GAuth\[Google Identity Services]

&#x20;   end



&#x20;   subgraph Go Server

&#x20;       Router\[Chi Router]

&#x20;       Auth\[JWT Middleware]

&#x20;       API\[Handlers: Auth, Products, Subs, Invoices, Admin]

&#x20;       Worker\[Background Billing Engine]

&#x20;   end



&#x20;   subgraph Database

&#x20;       PG\[(PostgreSQL)]

&#x20;   end



&#x20;   Vue -->|HTTP REST + JWT| Router

&#x20;   GAuth -->|OAuth Token| Vue

&#x20;   Router --> Auth

&#x20;   Auth --> API

&#x20;   API -->|Read/Write| PG

&#x20;   Worker -->|Sweeps expiries every 60s| PG

```



\## Core Features



\* Product Management

\* Invoice Management

\* Revenue Reporting

\* Authentication and Authorization

\* Subscription Billing

\* Background Processing



\## Technology Stack



| Component        | Technology |

| ---------------- | ---------- |

| Backend          | Go         |

| Frontend         | Vue.js 3   |

| Database         | PostgreSQL |

| Routing          | Chi Router |

| Authentication   | JWT        |

| API Architecture | REST       |



\## Product and Finance APIs



| Method | Endpoint             | Description                   |

| ------ | -------------------- | ----------------------------- |

| POST   | `/api/products`      | Create a product (Admin only) |

| GET    | `/api/products`      | Retrieve products             |

| GET    | `/api/invoices`      | Retrieve invoices             |

| GET    | `/api/revenue`       | Revenue report                |

| GET    | `/api/admin/revenue` | Admin revenue report          |



\## Access Control



\* Customers can view their own invoices.

\* Administrators can access invoices across users.

\* Product creation is restricted to administrators.

\* Revenue reporting is restricted to administrators.



\## Project Structure



\* `auth/` – Authentication and authorization

\* `frontend/` – Frontend application

\* `handlers/` – API handlers

\* `models/` – Data models

\* `workers/` – Background processing

\* `main.go` – Backend entry point

\* `syncro.sql` – Database schema



\## Future Enhancements



\* Additional payment gateway integrations

\* Enhanced revenue analytics

\* Expanded subscription management

\* Improved billing and notification workflows



