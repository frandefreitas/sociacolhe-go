# Sociacolhe API (Go)

API mínima viável para conectar psicólogos/estudantes a pacientes (anônimos ou identificados), com triagem, atribuição, sessões, feedback e chat anônimo via WebSocket.

## Rodando rápido

### Requisitos
- Go 1.22+
- PostgreSQL 16+ (ou use Docker)

### Com Docker (recomendado)
```bash
docker compose up --build
# API: http://localhost:8080/health
```
As migrações rodam automaticamente via `migrations/001_init.sql`.

### Local (sem Docker)
1. Crie o banco `sociacolhe` e rode o SQL de `migrations/001_init.sql`.
2. Copie `.env.example` para variáveis de ambiente (ou exporte no shell).
3. Rode:
```bash
go mod tidy
go run ./cmd/server
```

## Endpoints principais
- `POST /auth/register` — cadastro (Paciente auto-aprovado)
- `POST /auth/login` — login (JWT)
- `PUT /admin/users/:userID/approve` — aprovar usuário (ADMIN)
- `POST /requests` — criar solicitação (PATIENT)
- `GET /requests/open` — listar solicitações abertas (PSYCHOLOGIST/STUDENT/ADMIN)
- `PUT /requests/:requestID/triage` — triagem manual (PSYCHOLOGIST/ADMIN)
- `POST /requests/:requestID/assignments` — aceitar atendimento (PSYCHOLOGIST/STUDENT)
- `POST /sessions` — registrar sessão (PSYCHOLOGIST/STUDENT)
- `POST /feedbacks` — criar feedback (aberto)
- `POST /chat/rooms` — criar sala de chat (aberto)
- `GET /ws/chat/:roomID` — conectar via WebSocket

## IMPORTANTE (Segurança/Ética)
- **Notas de sessão** em claro (MVP). Em produção, criptografe campo-a-campo (KMS).
- **Logs**: evite logs de dados sensíveis.
- **RBAC** simples por papel. Refine conforme necessidade.
- **Auditoria**: adicione trilhas de auditoria em produção.
# sociacolhe-go
