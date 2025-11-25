package main

import (
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"

    "github.com/yourname/taskify/internal/handler"
    "github.com/yourname/taskify/internal/repository"
    "github.com/yourname/taskify/internal/router"
    "github.com/yourname/taskify/internal/service"
    "github.com/yourname/taskify/pkg/db"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println(".env not loaded, rely on environment variables")
    }

    dsn := db.GetDSNFromEnv()
    dbConn, err := db.Connect(dsn)
    if err != nil {
        log.Fatalf("db connect: %v", err)
    }
    defer dbConn.Close()

    // repositories
    userRepo := repository.NewUserRepo(dbConn)
    taskRepo := repository.NewTaskRepo(dbConn)

    // services
    authSvc := service.NewAuthService(userRepo)
    taskSvc := service.NewTaskService(taskRepo)

    // handlers
    authHandler := handler.NewAuthHandler(authSvc)
    taskHandler := handler.NewTaskHandler(taskSvc)

    // router
    r := router.NewRouter(authHandler, taskHandler)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("server running on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}