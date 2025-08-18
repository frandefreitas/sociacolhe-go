# Explicação linha a linha (arquivos-chave)

> Para manter legível, detalhamos **cada instrução** dos arquivos mais centrais.
> Os demais seguem os mesmos padrões e possuem comentários no código.

---

## cmd/server/main.go
1. `package main` — ponto de entrada da aplicação.
2-6. Imports de `log`, `os` e módulos internos `config`, `db`, `routes`.
8. `func main()` — função principal.
9. `cfg := config.Load()` — carrega variáveis de ambiente.
10. `db.Connect(cfg.DatabaseURL)` — abre conexão com Postgres.
13. `r := routes.SetupRouter(cfg)` — monta o roteador HTTP (Gin).
14. Loga porta.
15-18. `r.Run(":"+cfg.Port)` — inicia servidor e encerra se houver erro.

## internal/config/config.go
1. `package config` — pacote de configuração.
3-6. Imports de `log` e `os` (ler env e debugar).
8-13. `type Config` — estrutura de config.
15-23. `Load()` — lê envs com `get`, define defaults e loga resumo.
25-31. `get(k, def string)` — utilitário para pegar env com fallback.

## internal/routes/routes.go
1. `package routes` — agrupador de rotas/handlers.
3-8. Imports `net/http`, `gin`, `config`, `handlers`, `middleware`, `ws`.
10. `SetupRouter(cfg config.Config)` — constrói `*gin.Engine`.
11. `r := gin.Default()` — cria engine com middlewares padrão.
12. `jwt := middleware.JWTAuth{Secret: cfg.JWTSecret}` — emite/valida JWT.
15. `r.GET("/health", ...)` — rota de saúde (200 OK).
18-21. Autenticação: POST `/auth/register` e `/auth/login`.
24-25. Admin: PUT `/admin/users/:userID/approve` com guard `jwt.Require("ADMIN")`.
28-30. Requests: cria (`POST /requests`, papel PATIENT) e lista abertas (`GET /requests/open`).
33-34. Triagem manual: `PUT /requests/:requestID/triage` (PSYCHOLOGIST/ADMIN).
37-38. Atribuição: aceitar atendimento (PSYCHOLOGIST/STUDENT).
41-42. Sessões: registrar sessão.
45-46. Feedbacks: endpoint aberto (permite anônimo).
49-52. Chat: cria sala e abre WebSocket por sala.
54. `return r` — devolve engine.

## internal/middleware/auth.go
1. Pacote `middleware`.
3-8. Imports HTTP, `strings`, `time`, Gin e JWT.
10. `type JWTAuth` — guarda a `Secret` da assinatura.
12-16. `Claims` — payload: `uid`, `role` e `expires`.
18-26. `Issue` — cria token HS256 (expira em 72h).
28-51. `Require(roles...)` — middleware:
- Lê `Authorization: Bearer ...`.
- Valida assinatura e expiração.
- Se `roles` informado, verifica papel autorizado.
- Injeta `uid` e `role` no contexto.

## internal/handlers/auth.go
1. Pacote `handlers`.
3-11. Imports: HTTP, validação, UUID, DB, JWT e bcrypt.
13-16. `AuthHandler` — dependências (JWT + flag `AllowSignup`).
18-26. `registerReq` — payload do cadastro.
28. `loginReq` — payload do login.
30-43. `Register` — cria usuário:
- Bloqueia se signup desabilitado.
- Valida JSON, gera hash de senha.
- Seta `IsApproved = true` para `PATIENT`.
- Persiste no banco e retorna `id` + `isApproved`.
45-55. `Login` — autentica:
- Busca por email, confere senha, exige aprovação (não para paciente).
- Emite JWT com `uid` e `role`.

## internal/handlers/request.go
1. Pacote `handlers` + imports.
9. `RequestHandler` — marcador.
11-14. `createReq` — campos da solicitação.
16-28. `Create` — cria `ServiceRequest`:
- Valida JSON; roda `AutoTriage` para `area` e `urgency`.
- Se não anônimo, preenche `patientId` com `uid` do token.
- `status = OPEN`.
30-36. `ListOpen` — lista OPEN/IN_TRIAGE em ordem de urgência e data.

## internal/services/triage.go
1. Pacote `services`.
3. Import `strings`.
6-16. `AutoTriage` — heurística MVP:
- Urgência 5 se conter termos de autoagressão.
- Detecta área por palavras-chave (Ansiedade, Depressão, Violência).
- Default: `Geral`.

## internal/handlers/triage.go
1-2. Pacote/imports.
10-14. `triageReq` — campos editáveis na triagem manual.
16-28. `Manual` — atualiza `area`/`urgency` e seta `status = IN_TRIAGE`.

## internal/handlers/assignment.go
1-2. Pacote/imports.
10-21. `Accept` — profissional/estudante aceita atendimento:
- Cria `Assignment` vinculando `requestId` e `assigneeId`.
- Atualiza `status = ASSIGNED`.

## internal/handlers/session.go
1-2. Pacote/imports.
10-21. `sessionReq` — dados da sessão (data, duração, tipo, notas, CRP supervisor).
23-36. `Create` — persiste sessão vinculada ao profissional e (opcional) a request/estudante.

## internal/handlers/feedback.go
1-2. Pacote/imports.
10-19. `feedbackReq` — notas 1..5 e comentário/flag.
21-30. `Create` — grava feedback para `requestId` (sem auth para permitir anônimo).

## internal/ws/hub.go
1-2. Pacote/imports.
8. `upgrader` — permite CORS amplo (MVP).
10. `Hub` — sala -> conjunto de conexões.
12. `NewHub` — construtor.
14-24. `Join` — sobe WS, registra conexão e faz broadcast de mensagens.

## internal/models/models.go
- Define tipos e structs com `uuid` e timestamps.
- Enums são representados como `string` no Go, e mapeados às types do Postgres via GORM tags.

## migrations/001_init.sql
- Cria enums `user_role` e `request_status`.
- Cria tabelas com chaves e FKs necessárias para MVP.

---

### Observações
- Handlers seguem padrão: valida payload -> acessa `db.DB` (GORM) -> monta resposta JSON.
- O RBAC é garantido por `middleware.JWTAuth.Require(...)` nas rotas.
- O chat WS é **broadcast simples** por sala; produção: persistência de mensagens, rate-limit e autenticação.
