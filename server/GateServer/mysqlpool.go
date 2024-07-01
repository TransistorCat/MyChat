package main

// import (
// 	"fmt"
// 	"sync"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/jmoiron/sqlx"
// )

// type MySqlPool struct {
// 	db          *sqlx.DB
// 	maxPoolSize int
// 	pool        chan *sqlx.DB
// 	mu          sync.Mutex
// }

// func NewMySqlPool(url, user, pass, schema string, poolSize int) (*MySqlPool, error) {

// 	if err != nil {
// 		db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, url, schema))
// 		return nil, err
// 	}
// 	db.SetMaxOpenConns(poolSize)
// 	db.SetMaxIdleConns(poolSize)
// 	pool := make(chan *sqlx.DB, poolSize)
// 	for i := 0; i < poolSize; i++ {
// 		pool <- db
// 	}
// 	return &MySqlPool{
// 		db:          db,
// 		maxPoolSize: poolSize,
// 		pool:        pool,
// 	}, nil
// }

// func (p *MySqlPool) getConnection() (*sqlx.DB, error) {
// 	select {
// 	case db := <-p.pool:
// 		return db, nil
// 	default:
// 		return nil, fmt.Errorf("connection pool is full")
// 	}
// }

// func (p *MySqlPool) returnConnection(db *sqlx.DB) {
// 	p.pool <- db
// }

// func (p *MySqlPool) Close() {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	close(p.pool)
// 	for db := range p.pool {
// 		db.Close()
// 	}
// }

// func main() {
// 	// Initialize Gin
// 	r := gin.Default()

// 	// Initialize MySQL pool
// 	pool, err := NewMySqlPool("localhost:3306", "username", "password", "database_name", 10)
// 	if err != nil {
// 		log.Fatalf("Failed to create MySQL pool: %v", err)
// 	}
// 	defer pool.Close()

// 	// Endpoint to get a database connection
// 	r.GET("/getConnection", func(c *gin.Context) {
// 		db, err := pool.getConnection()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"message": "Got a database connection", "db": db})
// 	})

// 	// Endpoint to return a database connection
// 	r.GET("/returnConnection", func(c *gin.Context) {
// 		// Simulate some work with the connection
// 		time.Sleep(5 * time.Second)

// 		// Return the connection to the pool
// 		pool.returnConnection(nil)
// 		c.JSON(http.StatusOK, gin.H{"message": "Returned database connection"})
// 	})

// 	// Run the server
// 	if err := r.Run(":8080"); err != nil {
// 		log.Fatalf("Failed to run server: %v", err)
// 	}
// }
