package main

import (
	"agent-balam/handlers"
	"agent-balam/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PrepareRoutes registers all API routes with the Gin engine.
func PrepareRoutes(
	r *gin.Engine,
	authH *handlers.AuthHandler,
	agentH *handlers.AgentHandler,
	familyH *handlers.FamilyHandler,
	clientH *handlers.ClientHandler,
	planH *handlers.PlanHandler,
	policyH *handlers.PolicyHandler,
	fupH *handlers.FUPHandler,
	commH *handlers.CommissionHandler,
	gstH *handlers.GSTHandler,
	loanH *handlers.LoanHandler,
	sbH *handlers.SBHandler,
	leadH *handlers.LeadHandler,
	reportH *handlers.ReportHandler,
	adminH *handlers.AdminHandler,
	dashboardH *handlers.DashboardHandler,
	searchH *handlers.SearchHandler,
) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/v1")

	// Auth (public)
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authH.Login)
		auth.POST("/register", authH.Register)
		auth.POST("/refresh", middlewares.Auth(), authH.Refresh)
		auth.POST("/change-password", middlewares.Auth(), authH.ChangePassword)
	}

	// All routes below require authentication
	protected := v1.Group("", middlewares.Auth())

	// Agent
	agent := protected.Group("/agent")
	{
		agent.GET("/profile", agentH.GetProfile)
		agent.PUT("/profile", agentH.UpdateProfile)
	}

	// Families
	families := protected.Group("/families")
	{
		families.GET("", familyH.List)
		families.POST("", familyH.Create)
		families.GET("/:familyCode", familyH.Get)
		families.PUT("/:familyCode", familyH.Update)
	}

	// Clients
	clients := protected.Group("/clients")
	{
		clients.GET("", clientH.List)
		clients.POST("", clientH.Create)
		clients.GET("/search", clientH.Search)
		clients.GET("/:id", clientH.Get)
		clients.PUT("/:id", clientH.Update)
	}

	// Plans
	protected.GET("/plans", planH.List)

	// Policies
	policies := protected.Group("/policies")
	{
		policies.GET("", policyH.List)
		policies.POST("", policyH.Create)
		policies.GET("/:policyNo", policyH.Get)
		policies.PUT("/:policyNo", policyH.Update)
	}

	// FUP
	fup := protected.Group("/fup")
	{
		fup.GET("/due", fupH.Due)
		fup.POST("/update", fupH.Update)
		fup.GET("/history/:policyNo", fupH.History)
		fup.GET("/multipledue/:policyNo", fupH.MultipleDues)
	}

	// Commission
	commission := protected.Group("/commission")
	{
		commission.GET("", commH.List)
		commission.POST("", commH.Create)
		commission.GET("/summary", commH.Summary)
		commission.GET("/calculate", commH.Calculate)
	}

	// GST
	protected.GET("/gst/calculate", gstH.Calculate)

	// Loans
	loans := protected.Group("/loans")
	{
		loans.GET("", loanH.List)
		loans.POST("", loanH.Create)
	}

	// Survival Benefits
	sb := protected.Group("/sb")
	{
		sb.GET("", sbH.List)
		sb.POST("", sbH.Create)
		sb.PUT("/:id/mark-paid", sbH.MarkPaid)
	}

	// Leads & Activities
	leads := protected.Group("/leads")
	{
		leads.GET("", leadH.ListLeads)
		leads.POST("", leadH.CreateLead)
	}

	activities := protected.Group("/activities")
	{
		activities.GET("", leadH.ListActivities)
		activities.POST("", leadH.CreateActivity)
		activities.GET("/today", leadH.TodayActivities)
	}

	// Reports
	reports := protected.Group("/reports")
	{
		reports.GET("/cashflow", reportH.Cashflow)
		reports.GET("/cashinout", reportH.CashInOut)
		reports.GET("/status", reportH.CurrentStatus)
		reports.GET("/calendar", reportH.Calendar)
		reports.POST("/refresh", reportH.Refresh)
	}

	// Admin — protected, requires auth
	admin := protected.Group("/admin")
	{
		admin.POST("/import-sql", adminH.ImportSQL)
	}

	// Agent data import — clears all existing data then imports from INSERT-only SQL file
	protected.POST("/agent/import", adminH.ResetAndImport)

	// Dashboard — aggregated home-screen summary
	protected.GET("/dashboard", dashboardH.Summary)

	// Global search — families, clients, policies in one call
	protected.GET("/search", searchH.Search)
}
