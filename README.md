# 🌐 **Knowledge Capsule**

### ⚡ A Lightweight, Go-Powered Knowledge Management Backend

**Knowledge Capsule** is a fast, simple, Go-based backend that allows you to create, store, search, and organize **“knowledge capsules”** — bite-sized learning notes categorized by topics and tags.
Perfect for personal knowledge bases, team learning platforms, or lightweight documentation systems.

## ✨ **Features**

* 🔐 **User Authentication** – Secure JWT-based login & registration
* 🧠 **Capsule Management** – Create, read, and organize knowledge entries
* 🗂️ **Topic Organization** – Categorize capsules using topics
* 🔍 **Powerful Search** – Search capsules by title or content
* 🏷️ **Tagging System** – Add tags for deeper filtering
* 💾 **PostgreSQL + GORM** – Persistent database storage
* 👤 **RBAC** – Roles: user, admin, superadmin (role assignment by admin/superadmin)
* 👥 **User Management** – Profile, avatar, list users (admin), admin team (superadmin)
* 🔍 **Global Search** – Admin-only search across users, topics, capsules
* 📋 **Filtering** – Query params on GET endpoints (topic, tags, q, is_private, role)
* 💬 **Real-time Chat** – Fully WebSocket-based (send messages, fetch history over socket)
* 📂 **File Uploads** – Upload and serve files locally

## 🧰 **Tech Stack**

* 🏎️ **Go (1.25+)**
* 🐘 **PostgreSQL** – Database
* 📦 **Docker & Docker Compose**
* 🔁 **Air (Live Reload)**
* 🛠️ **Makefile** for workflow automation
* ⚙️ **Lefthook** for Git hooks

## 🚀 Getting Started

### 1️⃣ **Clone the Repository**

```bash
git clone https://github.com/shahadathhs/knowledge-capsule.git
cd knowledge-capsule
```

### 2️⃣ **Environment Setup**

Copy `.env.example` to `.env` and fill in values:

```bash
cp .env.example .env
```

Required variables:

| Variable | Description |
|----------|-------------|
| `PORT` | Server port (default: 8080) |
| `GO_ENV` | `development` or `production` |
| `JWT_SECRET` | Secret for JWT signing |
| `DATABASE_URL` | PostgreSQL connection string |
| `POSTGRES_USER` | DB user (for Docker Compose) |
| `POSTGRES_PASSWORD` | DB password |
| `POSTGRES_DB` | Database name |
| `SUPERADMIN_EMAIL` | (Optional) Superadmin email – creates/updates on startup |
| `SUPERADMIN_PASSWORD` | (Optional) Superadmin password |
| `SUPERADMIN_NAME` | (Optional) Superadmin display name |

💡 Generate JWT secret: `make g-jwt`

## 🐳 Run Using Docker (Recommended)

Docker Compose starts **PostgreSQL** + **API** together. The database is included in both `dev` and `prod` profiles.

### ▶️ Development Mode (with Live Reload)

```bash
make up-dev
```

👉 API at **[http://localhost:8081](http://localhost:8081)** · PostgreSQL on `localhost:5432`

### ▶️ Production Mode

```bash
make up
```

👉 API at **[http://localhost:8080](http://localhost:8080)** · PostgreSQL on `localhost:5432`

### 🐘 Database Only (for local dev without full compose)

If you run the API locally (`make run`) and want PostgreSQL in Docker:

```bash
make db        # Start PostgreSQL
make down-db   # Stop PostgreSQL
```

### ⏹️ Stop Containers

```bash
make down-dev   # dev
make down       # prod
```

## 🖥️ Run Locally (Without Docker)

Ensure PostgreSQL is running (e.g. `make db` or your own instance) and `DATABASE_URL` is set in `.env`.

Install dependencies:

```bash
make install
```

Start server with live reload:

```bash
make run
```

Or build & run binary:

```bash
make build-local
./tmp/server
```

## 🧪 **Test Chat UI**

**GET** `/test-ws` — WebSocket chat test page (same origin as API, no CORS issues)

## 📘 **API Documentation**

Swagger UI at `/docs/index.html`

**Protected endpoints:** Click **Authorize**, enter `Bearer <your-jwt-token>` (get token from POST `/api/auth/login`), then **Authorize** again. All subsequent requests will include the token.

## 🔐 **Authentication Endpoints**

### ➕ Register:

**POST** `/api/auth/register`
Body:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

### 🔑 Login:

**POST** `/api/auth/login`

## 👤 **User & Profile** (Requires JWT)

* 📥 **GET** `/api/users/me` – Current user profile (id, name, email, role, avatar_url)
* ✏️ **PATCH** `/api/users/me` – Update name, avatar_url

## 🗂️ **Topic Management** (Requires JWT)

* 📥 **GET** `/api/topics?page=1&limit=20&q=<search>` – Fetch topics (paginated, filterable)
* ➕ **POST** `/api/topics` – Create topic
* ✏️ **PUT** `/api/topics/{id}` – Update topic
* 🗑️ **DELETE** `/api/topics/{id}` – Delete topic

## 🧠 **Capsule Management** (Requires JWT)

### ➕ Create Capsule

**POST** `/api/capsules`

```json
{
  "title": "Interfaces in Go",
  "content": "Interfaces are named collections of method signatures...",
  "topic": "Golang",
  "tags": ["programming", "go"],
  "is_private": false
}
```

### 📥 Get Capsules

**GET** `/api/capsules?page=1&limit=20&topic=&tags=&q=&is_private=` (all query params optional)

### ✏️ Update Capsule

**PUT** `/api/capsules/{id}`

### 🗑️ Delete Capsule

**DELETE** `/api/capsules/{id}`

## 🔍 **Search & Filter**

**GET endpoints support search + filter** via query params (`q`, `page`, `limit`, etc.):
- `GET /api/capsules?q=&topic=&tags=&is_private=` – Search/filter capsules (owner only)
- `GET /api/topics?q=` – Search/filter topics
- `GET /api/users?q=&role=` – Search/filter users (admin only)

**GET** `/api/admin/search?q=<query>&limit=10` – **Global search** (admin only): searches users, topics, and capsules in one request

## 👥 **Admin** (Admin/Superadmin)

* 📥 **GET** `/api/users` – List users (admin, paginated: `q`, `role`, `page`, `limit`)
* 📥 **GET** `/api/users/{id}` – Get user by ID (admin)
* 📥 **GET** `/api/admin/admins` – List admins (superadmin only)
* ✏️ **POST** `/api/admin/users/{id}/role` – Set user role (superadmin only): `{"role":"user|admin|superadmin"}`

## ❤️‍🩹 **Health Check**

**GET** `/health`
✔ Confirms server is alive

## 💬 **Chat & Uploads** (Requires JWT)

### 🔌 WebSocket Chat (Fully socket-based)
**GET** `/ws/chat` — Connect with `?token=<jwt>`
* **Send message:** `{ "type": "send", "payload": { "receiver_id": "...", "content": "...", "type": "text" } }`
* **Get history:** `{ "type": "get_history", "payload": { "user_id": "...", "page": 1, "limit": 20 } }`
* **Server responses:** `{ "type": "message"|"history"|"error", "payload": {...} }`

### 📤 Upload File
**POST** `/api/upload`
* Body: `multipart/form-data` with `file` field.

### 📂 Serve File
**GET** `/uploads/:filename`

## 🧱 **Project Structure**

```
knowledge-capsule/
├── app/
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # Auth, logger, etc.
│   ├── models/         # Data models
│   └── store/          # GORM stores
├── pkg/
│   ├── config/         # Configuration loading
│   ├── db/             # PostgreSQL connection
│   └── utils/          # Helpers
├── web/                # Frontend assets (Chat UI)
├── docs/               # Swagger API docs
├── uploads/            # Uploaded files
├── scripts/            # Helper scripts
├── Dockerfile
├── Dockerfile.dev
├── compose.yaml
├── Makefile
└── main.go
```

## 🛠️ **Development Commands**

* 📘 `make help` – See all commands
* ▶️ `make run` – Run locally with live reload
* 🔨 `make build-local` – Build binary
* 🐘 `make db` – Start PostgreSQL (for local dev)
* ✨ `make fmt` – Format code
* 🔍 `make vet` – Static analysis
* 🧹 `make tidy` – Cleanup modules
* 📝 `make swagger` – Generate API docs
* 🔐 `make g-jwt` – Generate JWT secret
