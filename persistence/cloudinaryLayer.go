package persistence 

import (
	"fmt"
	"log"
	"os"
	"github.com/rootsongjc/cloudinary-go"
)

type CloudinaryAuthentication struct {
	CloudName 	string
	APIKey	  	string
	APISecret	string
}

type Options struct {
	Headers map[string]string
}

func CloudinaryConnection()  *cloudinary.Service {
	auth := &CloudinaryAuthentication{
		CloudName: os.Getenv("CLOUDINARY_NAME"),
		APIKey: os.Getenv("CLOUDINARY_KEY"),
		APISecret: os.Getenv("CLOUDINARY_SECRET"),
	}
	s, err := cloudinary.Dial(fmt.Sprintf("cloudinary://%s:%s@%s", auth.APIKey, auth.APISecret, auth.CloudName))
	if err != nil {
		log.Fatal(err)
	}
	return s
}