package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"tanken/backend/data-fetcher/rpc"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var (
	geo_postid_rdb *redis.Client
	post_cache_rdb *redis.Client
	db             *sql.DB
	ctx            = context.Background()
)

func initGeoPostIdRedis() {
	geo_postid_rdb = redis.NewClient(&redis.Options{
		Addr:     "geo-postid-redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := geo_postid_rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to geo-postid Redis: %v", err)
	} else {
		log.Println("Connected successfully to geo-postid Redis")
	}
}

func initPostCache() {
	post_cache_rdb = redis.NewClient(&redis.Options{
		Addr:     "post-cache:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := post_cache_rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to post-cache Redis: %v", err)
	} else {
		log.Println("Connected successfully to post-cache Redis")
	}
}

func initPostgres() {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DBNAME")
	port, portStringToIntErr := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if portStringToIntErr != nil {
		log.Fatalf("Error converting POSTGRES_PORT to int: %v", portStringToIntErr)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", //CHECK: sslmode=disable
		host, port, user, password, dbname)

	log.Default().Printf("psqlInfo: %v", psqlInfo)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
}

func initAWS() (*session.Session, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, fmt.Errorf("AWS_REGION is not set")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}
	return sess, nil
}

func initS3(sess *session.Session) *s3manager.Uploader {
	uploader := s3manager.NewUploader(sess)
	return uploader
}

func main() {
	initGeoPostIdRedis()
	initPostCache()

	initPostgres()

	// sess, _ := initAWS()
	// uploader := initS3(sess)

	// rpc.StartServer(geo_postid_rdb, post_cache_rdb, db, uploader)
	rpc.StartServer(geo_postid_rdb, post_cache_rdb, db)
}
