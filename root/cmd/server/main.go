package main

import (
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"

    "github.com/FUADIKAMIL/taskify/internal/handler"
    "github.com/FUADIKAMIL/taskify/internal/repository"
    "github.com/FUADIKAMIL/taskify/internal/router"
    "github.com/FUADIKAMIL/taskify/internal/service"
    "github.com/FUADIKAMIL/taskify/pkg/db"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println(".env gagal dimuat!")
    }

    
    cfg := db.LoadConfigFromEnv()

    dsn := db.BuildDSN(cfg)

    connect := db.Connect("postgres")

    dbConn, err := connect(dsn)
    if err != nil {
        log.Fatalf("db connect error: %v", err)
    }
    defer dbConn.Close()

    
    userRepo := repository.NewUserRepo(dbConn)
    taskRepo := repository.NewTaskRepo(dbConn)

    
    authSvc := service.NewAuthService(userRepo)
    taskSvc := service.NewTaskService(taskRepo)

    
    authHandlers := router.AuthHandlers{
        Register: handler.RegisterHandler(authSvc),
        Login:    handler.LoginHandler(authSvc),
    }

    taskHandlers := router.TaskHandlers{
        GetAll:  handler.ListTaskHandler(taskSvc),
        Create:  handler.CreateTaskHandler(taskSvc),
        GetOne:  handler.GetOneTaskHandler(taskSvc),
        Update:  handler.UpdateTaskHandler(taskSvc),
        Delete:  handler.DeleteTaskHandler(taskSvc),
        Toggle:  handler.ToggleTaskHandler(taskSvc),
    }

    
    r := router.NewRouter(authHandlers, taskHandlers)

    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("server running on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
