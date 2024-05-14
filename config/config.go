package config

import "os"

func MongoURI() string {
	return os.Getenv("MONGO_URI")
}

func DbName() string {
	return os.Getenv("DB_NAME")
}

func CollectionName() string {
	return os.Getenv("COLLECTION_NAME")
}

func JwtSecretKey() string {
	return os.Getenv("JWT_SECRET_KEY")
}

func Port() string {
	return os.Getenv("PORT")
}
