# =========================
# 🗄️ DATABASE CONFIGURATION
# =========================

# PostgreSQL connection string (DSN)
# Cấu trúc: postgres://USER:PASSWORD@HOST:PORT/DBNAME?sslmode=disable
DB_URL=postgres://postgres:123456@localhost:5432/my_database?sslmode=disable

# =========================
# 🌐 SERVER CONFIGURATION
# =========================
# Port để chạy server Gin
PORT=3000

# =========================
# ⚙️ ENVIRONMENT SETTINGS
# =========================
# Chế độ chạy: development | production
GIN_MODE=development
# =========================
# LẤY SECRET KEY
# =========================
# chạy lệnh trên cmd: openssl rand -hex 32
# để lấy secret_key
# =========================