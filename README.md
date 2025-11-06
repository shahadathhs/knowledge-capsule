# knowledge-capsule-api

## Personal Knowledge Capsule API

**Description:**
A REST API that lets users store, retrieve, and manage **small “knowledge capsules”** — short text entries, quotes, links, or even personal thoughts — organized by **topics and tags**.
Think of it as your personal lightweight “digital second brain,” but purely backend.

## Core Requirements

### 1. **Entities**

#### a. User

* `id` (uuid)
* `name`
* `email`
* `password_hash`
* `created_at`
* `updated_at`

#### b. Capsule

* `id` (uuid)
* `user_id`
* `title`
* `content`
* `topic`
* `tags` (array of strings)
* `is_private` (bool)
* `created_at`
* `updated_at`

#### c. Topic

* `id` (uuid)
* `name`
* `description`
* `created_at`
* `updated_at`

## 🔧 Endpoints

### **Auth**

* `POST /api/auth/register` → Create a new user
* `POST /api/auth/login` → Return a JWT (no third-party libs if possible; use `crypto/hmac` and `encoding/base64`)

### **Capsules**

* `GET /api/capsules` → List user’s capsules
* `POST /api/capsules` → Create new capsule
* `GET /api/capsules/{id}` → Get capsule details
* `PUT /api/capsules/{id}` → Update capsule
* `DELETE /api/capsules/{id}` → Delete capsule

### **Topics**

* `GET /api/topics` → List topics
* `POST /api/topics` → Add new topic
* `GET /api/topics/{id}` → Get topic with its capsules

### **Search**

* `GET /api/search?q=keyword` → Search capsules by title/content/tags

## Technical Requirements

* Use only Go’s **standard library** (`net/http`, `encoding/json`, `crypto`, `os`, etc.)
* Use a **local JSON file or boltDB** (for simplicity) for data persistence
* Implement **middleware manually**, e.g.:

  * Logging
  * Authentication (JWT)
  * Panic recovery
* Graceful shutdown with `context.WithTimeout`
* Clean folder structure (e.g., `/handlers`, `/models`, `/middleware`, `/store`)

## Folder Structure

```
knowledge-capsule-api/
│
├── main.go
├── handlers/
│   ├── auth.go
│   ├── capsule.go
│   ├── topic.go
│
├── middleware/
│   ├── auth.go
│   ├── logger.go
│
├── models/
│   ├── capsule.go
│   ├── topic.go
│   ├── user.go
│
├── store/
│   ├── store.go  (JSON or BoltDB persistence)
│
└── utils/
    ├── jwt.go
    ├── hash.go
```

## Future Advanced Features

1. **Export capsules as Markdown** → `GET /api/export`
   Generate `.md` file dynamically using Go templates.

2. **Rate Limiting Middleware** → custom in-memory request counter with reset.

3. **Versioned API** → `/api/v1/...` to simulate production API design.

4. **CLI Tool (bonus)** → Add a small Go CLI that interacts with the API using `net/http`.

5. **Encrypted Capsules** → Store `is_private` capsules encrypted using AES before saving to disk.

## Learning Goals

* Deep dive into **Go’s `net/http`** without frameworks like Gin/Fiber
* Build **manual JWT** authentication
* Learn about **clean architecture** in raw Go
* Handle **middleware and request routing** yourself
* Understand **data persistence** with JSON or local DB
* Practice **structuring production-like Go apps**
