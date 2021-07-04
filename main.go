package main

import (
	"fmt"
	"github.com/exam105-UPD/backend/logging"
	"os"
	"context"
	"log"
	"time"

	_searchHandler "github.com/exam105-UPD/backend/search/delivery/http"
	_searchUseCase "github.com/exam105-UPD/backend/search/usecase"
	_searchRepo "github.com/exam105-UPD/backend/search/repository"

	_loginHandler "github.com/exam105-UPD/backend/login/delivery/http"
	_loginUseCase "github.com/exam105-UPD/backend/login/usecase"
	_loginRepo "github.com/exam105-UPD/backend/login/repository"

	_questionHandler "github.com/exam105-UPD/backend/question/delivery/http"
	_questionUseCase "github.com/exam105-UPD/backend/question/usecase"
	_questionRepo "github.com/exam105-UPD/backend/question/repository"

	_middleware "github.com/exam105-UPD/backend/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoConnClient *mongo.Client

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	/* 	dbHost := viper.GetString(`database.host`)
	   	dbPort := viper.GetString(`database.port`)
	   	dbUser := viper.GetString(`database.user`)
	   	dbPass := viper.GetString(`database.pass`)
	   	dbName := viper.GetString(`database.name`)
	   	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	   	val := url.Values{}
	   	val.Add("parseTime", "1")
	   	val.Add("loc", "Asia/Jakarta")
	   	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
		dbConn, err := sql.Open(`mysql`, dsn) */

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dbConn := initializeMongoDatabase(ctx)
	logging.InitializeMessages()
	
	defer func() {
		err := dbConn.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _middleware.InitMiddleware()
	e.Use(middL.CORS)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, remote_ip:${remote_ip}, host:${host} \n",
	}))

	// **** Article wiring ****

	/* 	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	   	ar := _articleRepo.NewMysqlArticleRepository(dbConn)

	   	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	   	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)
	   	_articleHttpDelivery.NewArticleHandler(e, au)
	*/
	// *********************************

	// ****** Question Wiring ******
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	qsRepo := _questionRepo.NewQuestionRepository(dbConn)
	qsUC := _questionUseCase.NewQuestionUsecase(qsRepo, timeoutContext)
	_questionHandler.NewQuestionHandler(e, qsUC)

	//**********Login Wiring**************

	//timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	loginRepo := _loginRepo.NewLoginRepository(dbConn)
	loginUC := _loginUseCase.NewLoginUsecase(loginRepo, timeoutContext)
	_loginHandler.NewLoginHandler(e, loginUC)

	//**********Search Wiring**************

	//timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	searchRepo := _searchRepo.NewSearchRepository(dbConn) 
	searchUC := _searchUseCase.NewSearchUsecase(searchRepo, timeoutContext)
	_searchHandler.NewSearchHandler(e, searchUC)

	log.Println(e.Start(viper.GetString("server.address")))
}

func initializeMongoDatabase(ctx context.Context) *mongo.Client {

	// Set client options
	var mongoURL string
	env := os.Getenv("ENV_EXAM105")
	
	if env == "LOCAL"  {
		mongoURL = os.ExpandEnv("mongodb://${ENV_MONGO_USER}:${ENV_MONGO_PASS}@54.255.95.50:27017/?authSource=${ENV_MONGO_AUTH_DB}") // Local
	} else if env == "DEV" {
		mongoURL = os.ExpandEnv("mongodb://${ENV_MONGO_USER}:${ENV_MONGO_PASS}@mongodb:27017/?authSource=${ENV_MONGO_AUTH_DB}") // DEV
	} else if env == "PROD" {
		mongoURL = os.ExpandEnv("mongodb://${ENV_REPLICA_USER}:${ENV_REPLICA_PASS}@${ENV_REPLICA_HOST_1}:27017,${ENV_REPLICA_HOST_2}:27017,${ENV_REPLICA_HOST_3}:27017/${ENV_REPLICA_DB}?replicaSet=${ENV_REPLICA_SET_NAME}&authSource=admin")
		//mongoURL = os.ExpandEnv("mongodb://${ENV_REPLICA_USER}:${ENV_REPLICA_PASS}@${ENV_REPLICA_PUBLIC_HOST_1}:27017,${ENV_REPLICA_PUBLIC_HOST_2}:27017,${ENV_REPLICA_PUBLIC_HOST_3}:27017/${ENV_REPLICA_DB}?replicaSet=${ENV_REPLICA_SET_NAME}&authSource=admin")
	}

	log.Println("Environment: " + os.Getenv("ENV_EXAM105"))
	log.Println("Env User: " + os.Getenv("ENV_MONGO_USER"))
	log.Println("S3 User: " + os.Getenv("ENV_S3_USERNAME"))
	log.Println("Replica Set Name: " + os.Getenv("ENV_REPLICA_SET_NAME"))
	log.Println("Replica Public IP: " + os.Getenv("ENV_REPLICA_PUBLIC_HOST_1"))
	log.Println("Replica Prive IP - 1: " + os.Getenv("ENV_REPLICA_HOST_1"))
	log.Println("Replica Prive IP - 2: " + os.Getenv("ENV_REPLICA_HOST_2"))
	log.Println("Replica Prive IP - 3: " + os.Getenv("ENV_REPLICA_HOST_3"))

	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions = clientOptions.SetMaxPoolSize(100) //100 is default driver setting
	log.Println("Connection String: " + clientOptions.GetURI())

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		//log.Fatal(err.Error())
		panic("Couldn't Connect to ReplicaSet")
	} else {
		fmt.Println("Connected to MongoDB Replica Set. ")
	}

	//Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		//log.Fatal("Couldn't PING to the database \n", err.Error())		
		panic("Database Replication PING Issue *** "+ err.Error())	
		
	} else {
		fmt.Println(" New MongoDB Replica Set connection created ! ")
	}

	MongoConnClient = client
	return MongoConnClient

}

