package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // blank import for file source driver
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq" // blank import for PostgreSQL driver
	"github.com/pkg/errors"
)

type Config struct {
	// Generic
	Environment  string `envconfig:"ENV" default:"staging"`
	Port         string `envconfig:"APP_PORT" default:"8080"`
	AllowedHosts string `envconfig:"ALLOWED_HOSTS" default:"*"`

	// DB
	DbDriver   string `envconfig:"DB_DRIVER" required:"false" default:"postgres"`
	DbUser     string `envconfig:"DB_USER" required:"false" default:"postgres"`
	DbPassword string `envconfig:"DB_PASSWORD" required:"false" default:"password"`
	DbHost     string `envconfig:"DB_HOST" required:"false" default:"localhost"`
	DbPort     string `envconfig:"DB_PORT" required:"false" default:"5432"`
	DbName     string `envconfig:"DB_NAME" required:"false" default:"local_todos"`
	DbSslMode  string `envconfig:"DB_SSL_MODE" required:"false" default:"disable"`

	// Migrations
	MigrationsDir string `envconfig:"MIGRATIONS_DIR" required:"false" default:"file://migrations"`
}

type Todo struct {
	ID             int    `json:"id"`
	TodoID         int    `json:"todo_id"`
	Description    string `json:"description"`
	ExpirationDate string `json:"expiration_date"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

const (
	apiBaseURL     = "/api/v1"
	todosAPIPrefix = "/todos"
)

func main() {
	// get config file path
	var cfgFilePath string
	flag.StringVar(
		&cfgFilePath,
		"config",
		"",
		"Path to config file. If not provided, config will be parsed from the environment.",
	)
	flag.Parse()

	// load config
	cfg, err := loadConfig(cfgFilePath)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// create postgresql database connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbSslMode,
	)

	fmt.Printf("Connecting to database: %s\n", connStr)

	// connect to the database
	db, err := sql.Open(cfg.DbDriver, connStr)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	// create the PostgreSQL driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create PostgreSQL driver instance: %v", err)
	}

	// create a new migration instance
	m, err := migrate.NewWithDatabaseInstance(cfg.MigrationsDir, cfg.DbName, driver)
	if err != nil {
		log.Fatalf("failed to create migration instance: %v", err)
	}

	// run migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	fmt.Println("migrations applied successfully!")

	// init Gin router
	router := gin.Default()

	// base url
	baseUrl := fmt.Sprintf("%s%s", apiBaseURL, todosAPIPrefix)

	// define the route for retrieving all todo items
	router.GET(baseUrl, func(c *gin.Context) {
		var todos []Todo

		rows, err := db.Query("SELECT * FROM todos")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve todo items"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var todo Todo
			err := rows.Scan(&todo.ID, &todo.TodoID, &todo.Description, &todo.ExpirationDate, &todo.CreatedAt, &todo.UpdatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve todo items"})
				return
			}
			todos = append(todos, todo)
		}

		c.JSON(http.StatusOK, todos)
	})

	// define the route for inserting a todo item
	router.POST(baseUrl, func(c *gin.Context) {
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert the todo item into the database
		insertQuery := `INSERT INTO todos (todo_id, description, expiration_date, created_at, updated_at)
						VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`
		err := db.QueryRow(insertQuery, todo.TodoID, todo.Description, todo.ExpirationDate).Scan(&todo.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert todo item"})
			return
		}

		c.JSON(http.StatusCreated, todo)
	})

	// define the route for updating a todo item
	router.PUT(fmt.Sprintf("%s/:id", baseUrl), func(c *gin.Context) {
		id := c.Param("id")

		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// update the todo item in the database
		updateQuery := `UPDATE todos SET todo_id=$1, description=$2, expiration_date=$3, updated_at=CURRENT_TIMESTAMP WHERE id=$4`
		_, err := db.Exec(updateQuery, todo.TodoID, todo.Description, todo.ExpirationDate, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Todo item updated successfully"})
	})

	// define the route for deleting a todo item
	router.DELETE(fmt.Sprintf("%s/:id", baseUrl), func(c *gin.Context) {
		id := c.Param("id")

		// delete the todo item from the database
		deleteQuery := "DELETE FROM todos WHERE id=$1"
		_, err := db.Exec(deleteQuery, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Todo item deleted successfully"})
	})

	// run server
	log.Fatal(router.Run(":8080"))
}

func loadConfig(cfgFile string) (cfg *Config, err error) {
	if cfgFile != "" {
		fmt.Printf("Loading config from file: %s", cfgFile)
		if cfg, err = readCfgFromFile(cfgFile); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	if cfg, err = readCfgFromEnv(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func readCfgFromFile(cfgFile string) (*Config, error) {
	if err := godotenv.Load(cfgFile); err != nil {
		return nil, errors.WithStack(err)
	}
	return readCfgFromEnv()
}

func readCfgFromEnv() (*Config, error) {
	cfg := Config{}
	fmt.Println("Loading config from env variables")
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.WithStack(err)
	}
	return &cfg, nil
}
