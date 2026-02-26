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

Create `.env` file:

```bash
PORT=8080
GO_ENV=development
JWT_SECRET=your_super_secret_key_here
DATABASE_URL=postgres://user:pass@localhost:5432/knowledge?sslmode=disable
```

💡 Generate secret automatically:
`make g-jwt`

## 🐳 Run Using Docker (Recommended)

### ▶️ Development Mode (with Live Reload)

```bash
make up-dev
```

👉 Runs at: **[http://localhost:8081](http://localhost:8081)**

### ▶️ Production Mode

```bash
make up
```

👉 Runs at: **[http://localhost:8080](http://localhost:8080)**

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

Swagger docs available at:
`/docs/index.html`

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

## 🗂️ **Topic Management** (Requires JWT)

* 📥 **GET** `/api/topics?page=1&limit=20` – Fetch topics (paginated)
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

**GET** `/api/capsules?page=1&limit=20`

### ✏️ Update Capsule

**PUT** `/api/capsules/{id}`

### 🗑️ Delete Capsule

**DELETE** `/api/capsules/{id}`

## 🔍 **Search Capsules**

**GET** `/api/search?q=<query>&page=1&limit=20`

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
