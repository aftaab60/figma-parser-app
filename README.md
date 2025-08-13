## Development

### Project Structure

```
figma-parser-app/
├── frontend/           # React TypeScript app
│   ├── docker/        # Docker configuration
│   │   ├── Dockerfile
│   │   ├── docker-compose.yml
│   │   └── docker-compose.fullstack.yml
│   ├── src/
│   │   ├── components/ # Reusable UI components
│   │   ├── pages/      # Page components
│   │   ├── types/      # TypeScript interfaces
│   │   └── utils/      # API utilities
│   ├── Makefile       # Frontend commands
│   └── package.json
├── backend/           # Go API server
│   ├── docker/        # Docker configuration
│   │   ├── Dockerfile
│   │   ├── docker-compose.yml
│   │   ├── docker_postgresql.env
│   │   ├── init-db.sql
│   │   └── conf/      # Database configuration
│   ├── handler/       # HTTP handlers
│   ├── internal/      # Internal packages
│   │   ├── db_manager/     # Database connection
│   │   ├── errors/         # Error handling
│   │   └── figma_manager/  # Figma API client
│   ├── middlewares/   # Custom middleware
│   ├── models/        # Data models
│   ├── repositories/  # Data access layer
│   ├── services/      # Business logic
│   ├── Makefile      # Backend commands
│   ├── go.mod
│   └── main.go
├── Makefile          # Root orchestration commands
└── README.md
```

### API Endpoints

- `POST /parse-figma-file` - Parse a Figma file (requires token)
- `GET /figma-files/:id` - Get file details with components/instances

### Execution in local

- Get whole project (db, backend, frontend) up and running: `make start`
- Stop all comonents: `make stop`
- Start backend service only: `cd backend; make start`
- Stop backend service only: `cd backend; make stop`
- Start frontend service only: `cd frontend; make start`
- Stop backend service only: `cd frontend; make stop`

### Example curl

````curl --location 'localhost:3000/parse-figma-file' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer figd_f9W-TF4VkCFtuVTjkQ_ypyi7HGHAONAH5tTUY2BY' \
--data '{
    "figma_file_url": "https://www.figma.com/design/DNCLfE7Tf8A0mudOOLZUYx/Locofy.ai.test?t=94QjDVDsuj5cHJEU-1"
}'```
````

### Local execution result
<img width="1272" height="784" alt="image" src="https://github.com/user-attachments/assets/907d2fdb-c79c-4ff1-808c-cab8b1341773" />

<img width="1146" height="803" alt="image" src="https://github.com/user-attachments/assets/b90d239b-353c-415c-af7d-835d65d516bf" />
