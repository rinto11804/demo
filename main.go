package main

import (
	"fmt"
	"log"
	"os"

	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func init(){
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main(){
	engine := html.New("./views",".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./public")


	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
	log.Printf("error: %v", err)
	return
}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)


	app.Get("/",func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	
	app.Post("/",func(c *fiber.Ctx) error {
		//AKIAZFHROKMHCNXEEDHR
		//VpVKFp3nfDfqFbiEZJyAsvVuZP5VL3/ik7cMot59
		file,err := c.FormFile("upload")

		if err != nil {
			return err
		}

		// Open file
		f,err := file.Open()
		
		if err !=nil {
			return err
		}

		//upload to s3

		result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("activitypoint"),
			Key:    aws.String(file.Filename),
			Body:   f,
			ACL: "public-read",
		})
		
		fmt.Println(result)

		return c.Render("index", fiber.Map{})
	})

	app.Listen(":" + os.Getenv("PORT"))
}