# 🌐 **Knowledge Capsule API**

### ⚡ A Lightweight, Go-Powered Knowledge Management Backend

**Knowledge Capsule API** is a fast, simple, Go-based backend that allows you to create, store, search, and organize **“knowledge capsules”** — bite-sized learning notes categorized by topics and tags.
Perfect for personal knowledge bases, team learning platforms, or lightweight documentation systems.

📌 **Live API & Swagger Docs:**
👉 [https://knowledge-capsule-api.onrender.com/docs/index.html](https://knowledge-capsule-api.onrender.com/docs/index.html)

## ✨ **Features**

* 🔐 **User Authentication** – Secure JWT-based login & registration
* 🧠 **Capsule Management** – Create, read, and organize knowledge entries
* 🗂️ **Topic Organization** – Categorize capsules using topics
* 🔍 **Powerful Search** – Search capsules by title or content
* 🏷️ **Tagging System** – Add tags for deeper filtering
* 💾 **File-based Storage** – JSON storage, no DB required — ultra simple setup

## 🧰 **Tech Stack**

* 🏎️ **Go (1.23+)**
* 📦 **Docker & Docker Compose**
* 🔁 **Air (Live Reload)**
* 🛠️ **Makefile** for workflow automation
* ⚙️ **Lefthook** for Git hooks

## 🚀 Getting Started

### 1️⃣ **Clone the Repository**

```bash
git clone https://github.com/shahadathhs/knowledge-capsule-api.git
cd knowledge-capsule-api
```

### 2️⃣ **Environment Setup**

Create `.env` file:

```bash
PORT=8080
GO_ENV=development
JWT_SECRET=your_super_secret_key_here
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

* 📥 **GET** `/api/topics` – Fetch topics
* ➕ **POST** `/api/topics` – Create topic

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

**GET** `/api/capsules`

## 🔍 **Search Capsules**

**GET** `/api/search?q=<query>`

## ❤️‍🩹 **Health Check**

**GET** `/health`
✔ Confirms server is alive

## 🧱 **Project Structure**

```
knowledge-capsule-api/
├── app/
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # Auth, logger, etc.
│   ├── models/         # Data models
│   └── store/          # JSON-based storage
├── pkg/
│   ├── config/         # Configuration loading
│   └── utils/          # Helpers
├── data/               # JSON data store
├── scripts/            # Helper scripts
├── Dockerfile
├── Dockerfile.dev
├── compose.yaml
├── Makefile
└── main.go
```

## 🛠️ **Development Commands**

* 📘 `make help` – See all commands
* ▶️ `make run` – Run locally
* 🔨 `make build-local` – Build binary
* ✨ `make fmt` – Format code
* 🔍 `make vet` – Static analysis
* 🧹 `make tidy` – Cleanup modules
* 🧪 `make test` – Run tests
* 🔐 `make g-jwt` – Generate JWT secret
