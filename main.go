package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"stock-control-back/config"
	"stock-control-back/usecase"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env load error", err)
	}

	mongoClient, err = mongo.Connect(
		context.Background(), options.Client().ApplyURI(config.MongoURI()))
	if err != nil {
		log.Fatal("Connection mongodb error", err)
	}

	err = mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Ping failed", err)
	}

	log.Print("Mongo connected")
}

func main() {
	defer mongoClient.Disconnect(context.Background())

	userColl := mongoClient.Database(config.DbName()).Collection(config.CollectionName())
	productService := usecase.ProductService{MongoCollection: userColl}
	userService := usecase.UserService{MongoCollection: userColl}
	signinService := usecase.SigninService{MongoCollection: userColl}

	r := mux.NewRouter()

	log.Println("Server is running on port 4444")
	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)
	r.HandleFunc("/auth/signup", userService.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/auth/signin", signinService.Signin).Methods(http.MethodPost)

	protectedRoutes := r.PathPrefix("/").Subrouter()
	protectedRoutes.Use(authMiddleware)
	protectedRoutes.HandleFunc("/product", productService.CreateProduct).Methods(http.MethodPost)
	protectedRoutes.HandleFunc("/product", productService.FindAllProducts).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/product/{id}", productService.DeleteProductByID).Methods(http.MethodDelete)

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"*",
		},
	})

	log.Fatal(http.ListenAndServe(":"+config.Port(), corsOpts.Handler(r)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Running"))
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromHeader(r)
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(config.JwtSecretKey()), nil
		})
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, ok := claims["userID"].(string)
		if !ok || userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}

	splitToken := strings.Split(bearerToken, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}
