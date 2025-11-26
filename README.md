# 📦 InventoryPulse

A modern, real-time inventory management system built with Go (Gin) backend and Svelte frontend. Features JWT authentication, WebSocket updates, product history tracking, and a beautiful glassmorphism UI.

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)
![Svelte](https://img.shields.io/badge/Svelte-4.0+-FF3E00?style=flat-square&logo=svelte)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-4169E1?style=flat-square&logo=postgresql)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

## ✨ Features

- **🔐 JWT Authentication** - Secure login with access and refresh tokens
- **👥 Role-Based Access Control** - Admin and Client roles
- **📦 Product Management** - Full CRUD operations with multiple categories support
- **📁 Category Management** - Organize products into categories (many-to-many)
- **📜 Product History** - Track price and stock changes over time
- **🔍 Unified Search** - Search products and categories in one endpoint
- **⚡ Real-Time Updates** - WebSocket-powered live data synchronization
- **🎨 Modern UI** - Glassmorphism design with Svelte
- **📊 Dashboard** - Overview stats and inventory value
- **🔄 Auto-Migration** - Database schema managed automatically

## 🗄️ Database Schema

```
┌─────────────────┐       ┌──────────────────────┐       ┌─────────────────┐
│     users       │       │  product_categories  │       │   categories    │
├─────────────────┤       ├──────────────────────┤       ├─────────────────┤
│ id (PK)         │       │ product_id (PK, FK)  │───────│ id (PK)         │
│ email (unique)  │       │ category_id (PK, FK) │       │ name (unique)   │
│ password_hash   │       └──────────────────────┘       │ description     │
│ role            │                 │                    │ created_at      │
│ created_at      │                 │                    │ updated_at      │
│ updated_at      │                 │                    │ deleted_at      │
│ deleted_at      │       ┌─────────┴─────────┐          └─────────────────┘
└─────────────────┘       │                   │
                          ▼                   │
                ┌─────────────────┐           │
                │    products     │           │
                ├─────────────────┤           │
                │ id (PK)         │───────────┘
                │ name            │
                │ description     │
                │ sku (unique)    │
                │ stock           │
                │ price           │
                │ category_id(FK) │
                │ created_at      │
                │ updated_at      │
                │ deleted_at      │
                └────────┬────────┘
                         │
                         │
                         ▼
                ┌─────────────────┐
                │ product_history │
                ├─────────────────┤
                │ id (PK)         │
                │ product_id (FK) │
                │ price           │
                │ stock           │
                │ changed_at      │
                └─────────────────┘
```

### Models

| Table | Fields |
|-------|--------|
| **products** | id, name, description, sku, price, stock, category_id, created_at, updated_at |
| **categories** | id, name, description, created_at, updated_at |
| **product_categories** | product_id, category_id |
| **product_history** | id, product_id, price, stock, changed_at |
| **users** | id, email, password_hash, role, created_at, updated_at |

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
- Docker & Docker Compose

### 1. Clone and Setup

```bash
git clone https://github.com/yourusername/inventorypulse.git
cd inventorypulse
```

### 2. Configure Environment

Create a `.env` file with the following configuration:

```env
# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=inventoryuser
DB_PASSWORD=inventorypass
DB_NAME=inventorypulse
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY_HOURS=24
JWT_REFRESH_EXPIRY_HOURS=168

# Admin User (created on first run)
ADMIN_EMAIL=admin@inventorypulse.com
ADMIN_PASSWORD=admin123
```

### 3. Start Database with Docker Compose

```bash
docker-compose up -d
```

Or manually:

```bash
docker run -d \
  --name inventorypulse_db \
  -e POSTGRES_USER=inventoryuser \
  -e POSTGRES_PASSWORD=inventorypass \
  -e POSTGRES_DB=inventorypulse \
  -p 5432:5432 \
  postgres:16-alpine
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

## 📚 API Documentation

### Swagger UI

Access the interactive API documentation at: `http://localhost:8080/swagger/index.html`

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
| GET | `/api/categories` | List categories (paginated) | Required |
| GET | `/api/categories/:id` | Get category by ID | Required |
| POST | `/api/categories` | Create category | Admin |
| PUT | `/api/categories/:id` | Update category | Admin |
| DELETE | `/api/categories/:id` | Delete category | Admin |

### Products

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/products` | List products (paginated, filterable) | Required |
| GET | `/api/products/:id` | Get product by ID | Required |
| POST | `/api/products` | Create product | Admin |
| PUT | `/api/products/:id` | Update product | Admin |
| DELETE | `/api/products/:id` | Delete product | Admin |
| PATCH | `/api/products/:id/stock` | Update stock | Admin |
| GET | `/api/products/:id/history` | Get product price/stock history | Required |

#### Product History Query Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `start` | string | Start date filter (YYYY-MM-DD) |
| `end` | string | End date filter (YYYY-MM-DD) |
| `page` | int | Page number (default: 1) |
| `page_size` | int | Items per page (default: 10) |

### Search

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/search` | Unified search for products and categories | Required |

#### Search Query Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `q` | string | Search query (required) |
| `type` | string | Filter by type: `product`, `category`, or empty for both |
| `page` | int | Page number (default: 1) |
| `page_size` | int | Items per page (default: 10) |

**Example:**
```bash
# Search all
GET /api/search?q=laptop

# Search only products
GET /api/search?q=laptop&type=product

# Search only categories
GET /api/search?q=electronics&type=category
```

## 🔌 WebSocket Documentation

### Connection

Connect to the WebSocket endpoint:

```
ws://localhost:8080/ws
```

### Events

The WebSocket server broadcasts the following events in real-time:

#### Product Events

| Event | Description | Payload |
|-------|-------------|---------|
| `product.created` | New product added | Product object |
| `product.updated` | Product modified | Product object |
| `product.deleted` | Product removed | `{ "id": <product_id> }` |
| `stock.updated` | Stock quantity changed | Product object |

#### Category Events

| Event | Description | Payload |
|-------|-------------|---------|
| `category.created` | New category added | Category object |
| `category.updated` | Category modified | Category object |
| `category.deleted` | Category removed | `{ "id": <category_id> }` |

### Message Format

```json
{
  "type": "product.created",
  "payload": {
    "id": 1,
    "name": "Product Name",
    "description": "Description",
    "sku": "SKU-001",
    "stock": 100,
    "price": 29.99,
    "category_id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### JavaScript Example

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
  console.log('Connected to WebSocket');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Event:', message.type);
  console.log('Payload:', message.payload);

  switch (message.type) {
    case 'product.created':
      // Handle new product
      break;
    case 'product.updated':
    case 'stock.updated':
      // Handle product update
      break;
    case 'product.deleted':
      // Handle product deletion
      break;
    case 'category.created':
    case 'category.updated':
    case 'category.deleted':
      // Handle category changes
      break;
  }
};

ws.onclose = () => {
  console.log('Disconnected from WebSocket');
};
```

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
├── docs/                 # Swagger documentation
├── frontend/             # Svelte frontend application
│   ├── src/
│   │   ├── lib/          # Components, stores, API client
│   │   └── routes/       # Page components
│   └── ...
├── docker-compose.yml
├── Makefile
└── README.md
```

## 🎨 Design Decisions

### Architecture

1. **Clean Architecture**: The project follows clean architecture principles with clear separation between layers:
   - **Handlers**: HTTP request handling and validation
   - **Services**: Business logic
   - **Repositories**: Database operations

2. **Dependency Injection**: All dependencies are injected through constructors, making the code testable and maintainable.

### Database

1. **GORM with AutoMigrate**: Automatic schema management for rapid development.

2. **Soft Deletes**: Products and categories use soft deletes (`deleted_at`) to preserve data integrity and allow recovery.

3. **Product History**: A separate table tracks all price and stock changes for auditing and analytics.

4. **Many-to-Many Categories**: Products can belong to multiple categories through the `product_categories` junction table.

### Authentication

1. **JWT with Refresh Tokens**: Implements secure authentication with short-lived access tokens and long-lived refresh tokens.

2. **Role-Based Access Control**: Two roles implemented:
   - **Admin**: Full CRUD access to all resources
   - **Client**: Read-only access to products and categories

### Real-Time Updates

1. **WebSocket Hub Pattern**: A central hub manages all WebSocket connections and broadcasts events efficiently.

2. **Event-Driven Updates**: All CRUD operations emit WebSocket events, keeping connected clients in sync.

## 🚢 Deployment

### Using Docker Compose (Recommended)

1. Build the backend:
```bash
go build -o bin/inventorypulse ./cmd/api
```

2. Start all services:
```bash
docker-compose up -d
```

### Manual Deployment

1. **Database**: Set up a PostgreSQL 16+ instance

2. **Backend**:
```bash
# Build
go build -o inventorypulse ./cmd/api

# Set environment variables
export DB_HOST=your-db-host
export DB_PORT=5432
export DB_USER=your-user
export DB_PASSWORD=your-password
export DB_NAME=inventorypulse
export JWT_SECRET=your-production-secret
export GIN_MODE=release

# Run
./inventorypulse
```

3. **Frontend**:
```bash
cd frontend
npm run build
# Serve the dist/ folder with nginx or any static server
```

### Production Checklist

- [ ] Set `GIN_MODE=release`
- [ ] Use a strong `JWT_SECRET`
- [ ] Enable HTTPS/TLS
- [ ] Configure proper database credentials
- [ ] Set up database backups
- [ ] Configure CORS for your domain
- [ ] Use a reverse proxy (nginx, Caddy)

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
| `DB_SSLMODE` | disable | Database SSL mode |
| `JWT_SECRET` | (required) | JWT signing secret |
| `JWT_EXPIRY_HOURS` | 24 | Access token expiry |
| `JWT_REFRESH_EXPIRY_HOURS` | 168 | Refresh token expiry |
| `ADMIN_EMAIL` | admin@inventorypulse.com | Initial admin email |
| `ADMIN_PASSWORD` | admin123 | Initial admin password |

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
