# Usage Guide for Project Template

This document outlines how to use this project as a template for new applications.

## 1. Clone the Repository

```bash
git clone <repository-url>
cd <project-name>
```

## 2. Initialize Git (if starting a new project)

If you cloned this template to start a new project, you might want to remove the `.git` directory and re-initialize:

```bash
rm -rf .git
git init
git add .
git commit -m "Initial commit from template"
```

## 3. Frontend Setup

Navigate to the `frontend` directory and install dependencies:

```bash
cd frontend
npm install
```

To run the frontend development server:

```bash
npm run dev
```

To build the frontend for production:

```bash
npm run build
```

## 4. Backend Setup

Navigate to the `server` directory and install Go modules:

```bash
cd server
go mod tidy
```

To run the backend server:

```bash
go run main.go
```

To build the backend executable:

```bash
go build -o server.exe main.go
```

## 5. Customization

- **Frontend**: Start developing your Vue.js application in `frontend/src/`.
- **Backend**: Implement your Go API logic in `server/`.
- **Configuration**: Update `.env.example` and create your `.env` files for both frontend and backend.

## 6. Cleaning Up

Before starting your new project, ensure you remove any remaining placeholder or example code that is not relevant to your new application.