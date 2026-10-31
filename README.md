# Syncro Billing Engine

Syncro is a lightweight, full-stack subscription billing engine built for high performance and rapid deployment. It features Role-Based Access Control (RBAC), automated background billing, Google OAuth integration, and a unified administrative dashboard.

## Architecture Diagram

```mermaid
graph TD
    subgraph Frontend
        Vue[Vue.js 3 SPA]
        GAuth[Google Identity Services]
    end

    subgraph Go Server
        Router[Chi Router]
        Auth[JWT Middleware]
        API[Handlers: Auth, Products, Subs, Invoices, Admin]
        Worker[Background Billing Engine]
    end

    subgraph Database
        PG[(PostgreSQL)]
    end

    Vue -->|HTTP REST + JWT| Router
    GAuth -->|OAuth Token| Vue
    Router --> Auth
    Auth --> API
    API -->|Read/Write| PG
    Worker -->|Sweeps expiries every 60s| PG
...
---

# 🚀 SYNCRO – Features, Architecture & Technical Overview

SYNCRO is a lightweight, full-stack subscription billing engine designed to simplify subscription management, automate recurring billing, and provide financial insights through a unified platform.

It brings together product management, customer subscriptions, invoice tracking, automated billing, and revenue analytics to streamline subscription-based business operations.

## ✨ Key Features

### 1. Authentication & Authorization

* User registration and login.
* JWT-based authentication.
* Role-Based Access Control (RBAC).
* Separate access permissions for administrators and customers.
* Google OAuth integration.

### 2. Product Management

* Create and manage subscription products.
* Define product pricing and billing intervals.
* Support for monthly and yearly billing plans.
* Retrieve available products through REST APIs.

### 3. Subscription Management

* Create and manage customer subscriptions.
* Support recurring monthly and yearly billing.
* Track subscription status and renewal cycles.
* Automated background billing and renewal processing.

### 4. Invoice Management

* Retrieve customer-specific invoices.
* Allow administrators to access all invoices.
* Track invoice amounts, payment status, and due dates.
* Maintain invoice creation and update timestamps.

### 5. Revenue Analytics

* Calculate Monthly Recurring Revenue (MRR).
* Normalize yearly subscription prices into monthly revenue.
* Generate revenue reports based on active subscriptions.
* Restrict financial reporting to authorized administrators.

### 6. Administrative Dashboard

* Centralized interface for managing subscription operations.
* Product and customer management.
* Access to invoice and revenue information.

## 🛠️ Technology Stack

| Component             | Technology                    |
| --------------------- | ----------------------------- |
| Frontend              | Vue.js, JavaScript, HTML, CSS |
| Backend               | Go (Golang)                   |
| Database              | PostgreSQL                    |
| API Architecture      | REST APIs                     |
| Authentication        | JWT, Google OAuth             |
| Routing               | Chi Router                    |
| Background Processing | Go routines                   |
| Version Control       | Git and GitHub                |

## 🏗️ System Architecture

SYNCRO follows a modular client-server architecture.

* **Frontend Layer:** Provides interfaces for customers and administrators.
* **Backend Layer:** Handles business logic, authentication, authorization, subscriptions, billing, and API requests.
* **Database Layer:** Stores user, product, subscription, and invoice information.
* **Background Processing Layer:** Handles recurring billing and subscription renewal operations.

## 🔗 API Overview

The backend exposes RESTful endpoints for product management, invoice retrieval, and revenue reporting.

| HTTP Method | Endpoint             | Description                             |
| ----------- | -------------------- | --------------------------------------- |
| POST        | `/api/products`      | Create a subscription product           |
| GET         | `/api/products`      | Retrieve available products             |
| GET         | `/api/invoices`      | Retrieve invoices based on user role    |
| GET         | `/api/revenue`       | Retrieve revenue information            |
| GET         | `/api/admin/revenue` | Retrieve administrative revenue reports |

### Access Control

* Product creation is restricted to administrators.
* Product listing is publicly accessible.
* Customers can retrieve only their own invoices.
* Administrators can retrieve all invoices.
* Revenue reporting endpoints require administrator authorization.

## 📂 Project Structure

```text
SYNCRO/
│
├── auth/             # Authentication and authorization
├── handlers/         # API request handlers
├── models/           # Data models
├── workers/          # Background billing engine
├── frontend/         # Vue.js frontend application
├── main.go           # Backend entry point
└── README.md
```

## ⚙️ Local Development

### Prerequisites

* Go
* Node.js and npm
* PostgreSQL
* Git

### Clone the Repository

```bash
git clone https://github.com/shraddha-sugathan/SYNCRO.git
cd SYNCRO
```

### Run the Backend

Configure PostgreSQL and the required environment variables before starting the backend.

```bash
go mod download
go run .
```

The backend is configured to run on:

`http://localhost:8080`

### Run the Frontend

Open a separate terminal:

```bash
cd frontend
npm install
npm run dev
```

Open the local URL displayed by Vite, usually:

`http://localhost:5173`

**Note:** The frontend and backend must both be running for complete application functionality. Ensure the database and backend environment configuration match the project's setup.

## 🔐 Security Considerations

* JWT-based authentication for protected resources.
* Role-based authorization for administrative operations.
* Customer-specific invoice access.
* Environment-based configuration for sensitive credentials.
* Database-backed storage of application data.

## 🔮 Future Enhancements

* Payment gateway integration.
* Email notifications for billing and renewals.
* Advanced revenue analytics.
* Subscription usage tracking.
* Enhanced payment history and reporting.

---

**SYNCRO – Simplifying Subscriptions. Streamlining Billing.**
