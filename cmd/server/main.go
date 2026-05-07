package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"vniizht/internal/config"
	"vniizht/internal/db"
	"vniizht/internal/handler"
	"vniizht/internal/middleware"
	"vniizht/internal/repository"
	"vniizht/internal/service"
	"vniizht/internal/storage"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.RunMigrations(ctx, cfg.DSN()); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	pool, err := db.NewPool(ctx, cfg)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pool.Close()

	minioClient, err := storage.NewMinioClient(ctx, cfg)
	if err != nil {
		log.Fatalf("minio: %v", err)
	}

	userRepo := repository.NewUserRepo(pool)
	sessionRepo := repository.NewSessionRepo(pool)
	categoryRepo := repository.NewCategoryRepo(pool)
	tagRepo := repository.NewTagRepo(pool)
	questionRepo := repository.NewQuestionRepo(pool)
	answerRepo := repository.NewAnswerRepo(pool)
	commentRepo := repository.NewCommentRepo(pool)
	attachmentRepo := repository.NewAttachmentRepo(pool)
	searchRepo := repository.NewSearchRepo(pool)
	notificationRepo := repository.NewNotificationRepo(pool)

	authSvc := service.NewAuthService(userRepo, sessionRepo)
	emailSvc := service.NewEmailService(cfg)
	notifSvc := service.NewNotificationService(notificationRepo, userRepo, emailSvc)
	questionSvc := service.NewQuestionService(questionRepo, categoryRepo, notifSvc)
	answerSvc := service.NewAnswerService(answerRepo, questionRepo, notifSvc)
	commentSvc := service.NewCommentService(commentRepo, answerRepo, notifSvc)
	adminSvc := service.NewAdminService(userRepo, categoryRepo, tagRepo)
	attachmentStorage := service.NewMinioStorage(minioClient, cfg.MinioBucket)
	attachmentSvc := service.NewAttachmentService(attachmentRepo, attachmentStorage)
	userSvc := service.NewUserService(userRepo, questionRepo, answerRepo)
	searchSvc := service.NewSearchService(searchRepo)
	importSvc := service.NewImportService(questionRepo, answerRepo)

	handlers := &handler.Handlers{
		Auth:          handler.NewAuthHandler(authSvc),
		Questions:     handler.NewQuestionsHandler(questionSvc),
		Answers:       handler.NewAnswersHandler(answerSvc),
		Comments:      handler.NewCommentsHandler(commentSvc),
		Attachments:   handler.NewAttachmentsHandler(attachmentSvc),
		Users:         handler.NewUsersHandler(userSvc),
		Search:        handler.NewSearchHandler(searchSvc),
		Admin:         handler.NewAdminHandler(adminSvc),
		Import:        handler.NewImportHandler(importSvc),
		Notifications: handler.NewNotificationHandler(notifSvc),
	}

	r := chi.NewRouter()
	r.Use(middleware.Recovery)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.CORS(cfg.CORSOrigin))
	r.Mount("/", handler.NewRouter(handlers, authSvc))

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 300 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server started", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown: %v", err)
	}

	slog.Info("server stopped")
}
