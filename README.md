# 📦 InventoryPulse

A modern, real-time inventory management system built with Go (Gin) backend and Svelte frontend. Features JWT authentication, WebSocket updates, and a beautiful glassmorphism UI.

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)
![Svelte](https://img.shields.io/badge/Svelte-4.0+-FF3E00?style=flat-square&logo=svelte)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-4169E1?style=flat-square&logo=postgresql)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

## ✨ Features

- **🔐 JWT Authentication** - Secure login with access and refresh tokens
- **👥 Role-Based Access Control** - Admin and Viewer roles
- **📦 Product Management** - Full CRUD operations with categories
- **📁 Category Management** - Organize products into categories
- **⚡ Real-Time Updates** - WebSocket-powered live data synchronization
- **🎨 Modern UI** - Glassmorphism design with Svelte
- **📊 Dashboard** - Overview stats and inventory value
- **🔄 Auto-Migration** - Database schema managed automatically

## 🛠️ Tech Stack

### Backend
- **Framework**: [Gin](https://gin-gonic.com/) (Go)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt)
- **WebSocket**: gorilla/websocket
- **Documentation**: Swagger (swaggo)

### Frontend
- **Framework**: [Svelte](https://svelte.dev/) + [Vite](https://vitejs.dev/)
- **Styling**: Custom CSS with Glassmorphism theme
- **State Management**: Svelte Stores

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- Node.js 18+
- Docker (for PostgreSQL)

### 1. Clone and Setup

```bash
git clone https://github.com/yourusername/inventorypulse.git
cd inventorypulse
```

### 2. Start PostgreSQL

```bash
docker run -d \
  --name inventorypulse_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=inventorypulse \
  -p 5432:5432 \
  postgres:16-alpine
```

### 3. Configure Environment

```bash
cp .env.example .env
# Edit .env if needed
```

### 4. Run Backend

```bash
go mod download
go run ./cmd/api
```

The API will start at `http://localhost:8080`

### 5. Run Frontend

```bash
cd frontend
npm install
npm run dev
```

The frontend will start at `http://localhost:5173`

### 6. Login

Default admin credentials:
- **Email**: `admin@inventorypulse.com`
- **Password**: `admin123`

## 📚 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/auth/login` | User login | Public |
| POST | `/api/auth/refresh` | Refresh token | Public |
| GET | `/api/auth/me` | Get current user | Required |
| POST | `/api/auth/register` | Create user | Admin |

### Categories

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/categories` | List categories | Required |
| GET | `/api/categories/:id` | Get category | Required |
| POST | `/api/categories` | Create category | Admin |
| PUT | `/api/categories/:id` | Update category | Admin |
| DELETE | `/api/categories/:id` | Delete category | Admin |

### Products

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/products` | List products | Required |
| GET | `/api/products/:id` | Get product | Required |
| POST | `/api/products` | Create product | Admin |
| PUT | `/api/products/:id` | Update product | Admin |
| DELETE | `/api/products/:id` | Delete product | Admin |
| PATCH | `/api/products/:id/stock` | Update stock | Admin |

### WebSocket

| Endpoint | Description |
|----------|-------------|
| `ws://localhost:8080/ws` | Real-time updates |

**WebSocket Events:**
- `product.created` - New product added
- `product.updated` - Product modified
- `product.deleted` - Product removed
- `stock.updated` - Stock quantity changed

## 🏗️ Project Structure

```
inventorypulse/
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/           # Configuration loading
│   ├── domain/models/    # Data models and DTOs
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # Auth, CORS middleware
│   ├── repository/       # Database operations
│   └── service/          # Business logic
├── pkg/
│   ├── database/         # DB connection, migrations, seeder
│   ├── jwt/              # JWT utilities
│   └── websocket/        # WebSocket hub and handlers
├── frontend/             # Svelte frontend application
│   ├── src/
│   │   ├── lib/          # Components, stores, API client
│   │   └── routes/       # Page components
│   └── ...
├── docker-compose.yml
├── Makefile
└── README.md
```

## 🔧 Development

### Make Commands

```bash
make run          # Run the application
make build        # Build binary
make test         # Run tests
make docker-up    # Start PostgreSQL
make docker-down  # Stop PostgreSQL
make swagger      # Generate Swagger docs
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | 8080 | API server port |
| `GIN_MODE` | debug | Gin mode (debug/release) |
| `DB_HOST` | localhost | Database host |
| `DB_PORT` | 5432 | Database port |
| `DB_USER` | postgres | Database user |
| `DB_PASSWORD` | postgres | Database password |
| `DB_NAME` | inventorypulse | Database name |
| `JWT_SECRET` | (required) | JWT signing secret |
| `JWT_EXPIRY_HOURS` | 24 | Access token expiry |
| `JWT_REFRESH_EXPIRY_HOURS` | 168 | Refresh token expiry |

## 📝 License

MIT License - see [LICENSE](LICENSE) for details.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

Built with ❤️ using Go and Svelte
