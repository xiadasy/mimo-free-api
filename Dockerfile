# Build frontend
FROM node:20-slim AS frontend
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# Build backend
FROM golang:1.22-alpine AS backend
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
COPY --from=frontend /app/static ./static
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o mimo-gateway .

# Runtime
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=backend /app/mimo-gateway .
RUN printf '#!/bin/sh\ntest -f config.json && exec ./mimo-gateway\nprintf "{\\"port\\":\\"%%s\\",\\"api_key\\":\\"%%s\\",\\"default_model\\":\\"%%s\\",\\"accounts\\":[{\\"id\\":\\"account-1\\",\\"service_token\\":\\"%%s\\",\\"user_id\\":\\"%%s\\",\\"ph\\":\\"%%s\\",\\"active\\":true}]}\\n" "$${PORT:-8080}" "$${MIMO_API_KEY:-sk-mimo}" "$${MIMO_DEFAULT_MODEL:-mimo-v2.5-pro}" "$$MIMO_SERVICE_TOKEN" "$$MIMO_USER_ID" "$$MIMO_PH" > config.json\nexec ./mimo-gateway\n' > /app/start.sh && chmod +x /app/start.sh
EXPOSE 8080
CMD ["/app/start.sh"]
