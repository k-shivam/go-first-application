Flow A:
main.go------>config.load()----> loads populates config struct first

main.go ------> Logger.New(cfg.Env) ------> Initializes slog (JSON or text)

main.go---------> database.NewPostgresPool(ctx, cfg.DB)------> connects to Postgres via pgxpool --> executes pool.Ping(ctx)

main.go ----------->Dependency Initiation(Bottom-Up Injection)
|->Repositories: NewUserRepo,NewProjectRepo,NewTaskRepository
|->Services: NewAuthServic, NewProjectService, NewTaskService
|->Handlers: NewAuthHandlr, NewProjectHandlr, NewTaskHandlr

main.go ------>Router Wiring (chi.NewRouter())
|---> Global Middleware: RequestID, RealIP, Logger, Recoverer
|---> Public Routes: /api/v1/auth/*
|---> Protected Sub-router: /api/v1/* wrapped in middleware.Authenticate(cfg.JWT)

main.go --------> Server Setup & Graceful Shutdown:
|->Launches server.ListenAndServe() in a background goroutine
|->Blocks main thread on signal.Notify(SIGINT, SIGTERM)
|->On signal -> executes server.Shutdown(ctxWithTimeout) to drain inflight requests

FLow B:

client:----> HTTP POST /api/v1/auth/register {"email": "...", "password": "...", "full_name": "..."}
[chi Router] ────────────────► Routes to AuthHandler.Register

[AuthHandler.Register]
├── 1. json.NewDecoder().Decode(&req)
            |----If malformed JSON -> returns 400 Bad Request
├── 2. validator.ValidateStruct(req)
            |---If fields invalid -> returns 422 Unprocessable Entity with FieldErrors
└── 3. Calls authService.Register(ctx, req)

[AuthService.Register]
├── 1. bcrypt.GenerateFromPassword(req.Password, DefaultCost)
├── 2. Constructs domain.User struct with new uuid.New() and UTC timestamps
├── 3. Calls userRepo.Create(ctx, user)

[UserRepository.Create]
├── Executes `INSERT INTO users (...) VALUES ($1, $2, ...)`
└── If Postgres returns error code 23505 (unique violation): Returns domain.ErrConflict

├── 4. If ErrConflict returned -> Service returns to Handler -> Handler writes 409 Conflict
├── 5. jwt.NewWithClaims(...) -> Signs token using cfg.Secret└── 6. Returns domain.AuthResponse{Token: token, User: user}

[AuthHandler.Register]
└── response.JSON(w, http.StatusCreated, res) -> Writes 201 Created with JSON envelope

