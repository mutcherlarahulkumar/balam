package main

import (
	"agent-balam/db"
	"agent-balam/db/repository"
	"agent-balam/domain"
	"agent-balam/handlers"
	"agent-balam/middlewares"
	"log/slog"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env.local first, fall back to .env
	if err := godotenv.Load(".env.local"); err != nil {
		_ = godotenv.Load(".env")
	}

	database, err := db.New()
	if err != nil {
		slog.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	defer database.Close()

	if err := db.RunMigrations(database); err != nil {
		slog.Error("migrations failed", "err", err)
		os.Exit(1)
	}

	// Wire encryption for sensitive bank-detail fields (aadhar, pan, ckyc)
	repository.InitCrypto(domain.Encrypt, domain.Decrypt)

	// Repositories
	agentRepo := repository.NewAgentRepo(database)
	familyRepo := repository.NewFamilyRepo(database)
	clientRepo := repository.NewClientRepo(database)
	planRepo := repository.NewPlanRepo(database)
	policyRepo := repository.NewPolicyRepo(database)
	fupRepo := repository.NewFUPRepo(database)
	commRepo := repository.NewCommissionRepo(database)
	loanRepo := repository.NewLoanRepo(database)
	sbRepo := repository.NewSBRepo(database)
	leadRepo := repository.NewLeadRepo(database)
	reportRepo := repository.NewReportRepo(database)

	// Domain services
	authSvc := domain.NewAuthService(agentRepo)
	agentSvc := domain.NewAgentService(agentRepo)
	familySvc := domain.NewFamilyService(familyRepo, clientRepo, policyRepo)
	clientSvc := domain.NewClientService(clientRepo, familyRepo, policyRepo)
	planSvc := domain.NewPlanService(planRepo)
	policySvc := domain.NewPolicyService(policyRepo, planRepo, fupRepo, loanRepo, sbRepo, clientRepo)
	fupSvc := domain.NewFUPService(fupRepo, policyRepo)
	commSvc := domain.NewCommissionService(commRepo, policyRepo)
	gstSvc := domain.NewGSTService(policyRepo, planRepo)
	loanSvc := domain.NewLoanService(loanRepo, policyRepo)
	sbSvc := domain.NewSBService(sbRepo, policyRepo)
	leadSvc := domain.NewLeadService(leadRepo, clientRepo)
	reportSvc := domain.NewReportService(reportRepo, familyRepo)

	// Handlers
	authH := handlers.NewAuthHandler(authSvc, agentSvc)
	agentH := handlers.NewAgentHandler(agentSvc)
	familyH := handlers.NewFamilyHandler(familySvc)
	clientH := handlers.NewClientHandler(clientSvc)
	planH := handlers.NewPlanHandler(planSvc)
	policyH := handlers.NewPolicyHandler(policySvc)
	fupH := handlers.NewFUPHandler(fupSvc, agentSvc)
	commH := handlers.NewCommissionHandler(commSvc)
	gstH := handlers.NewGSTHandler(gstSvc)
	loanH := handlers.NewLoanHandler(loanSvc)
	sbH := handlers.NewSBHandler(sbSvc)
	leadH := handlers.NewLeadHandler(leadSvc)
	reportH := handlers.NewReportHandler(reportSvc)
	adminH := handlers.NewAdminHandler(database)

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(middlewares.CORS())

	PrepareRoutes(router, authH, agentH, familyH, clientH, planH, policyH, fupH, commH, gstH, loanH, sbH, leadH, reportH, adminH)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	slog.Info("Balam API starting", "port", port)
	if err := router.Run(":" + port); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
