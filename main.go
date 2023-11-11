package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Définir une route pour la page d'accueil
	app.Get("/", func(c *fiber.Ctx) error {
		files, _ := ListUploadedFiles("uploads")
		return c.Render("views/index.html", fiber.Map{
			"Files": files,
		})
	})

	// Définir un point de terminaison pour l'envoi de fichiers
	app.Post("/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).SendString("Mauvaise demande")
		}
		filename := file.Filename

		// Vérifier si un fichier avec le même nom existe déjà
		if _, err := os.Stat("uploads/" + filename); !os.IsNotExist(err) {
			return c.Status(409).SendString("Un fichier avec le même nom existe déjà")
		}

		uploadedFile, err := file.Open()
		if err != nil {
			return c.Status(500).SendString("Impossible d'ouvrir le fichier")
		}
		defer uploadedFile.Close()

		out, err := os.Create("uploads/" + filename)
		if err != nil {
			return c.Status(500).SendString("Impossible de créer le fichier")
		}
		defer out.Close()

		_, err = io.Copy(out, uploadedFile)
		if err != nil {
			return c.Status(500).SendString("Impossible de copier le fichier")
		}

		return c.SendString(fmt.Sprintf("Fichier %s téléchargé avec succès", filename))
	})

	// Définir une route pour visualiser les fichiers téléchargés
	app.Get("/view/:filename", func(c *fiber.Ctx) error {
		filename := c.Params("filename")
		return c.SendFile("uploads/" + filename)
	})

	app.Get("/files", func(c *fiber.Ctx) error {
		files, _ := ListUploadedFiles("uploads")
		return c.JSON(fiber.Map{"Files": files})
	})

	// Définir une route pour supprimer un fichier
	app.Delete("/delete/:filename", func(c *fiber.Ctx) error {
		filename := c.Params("filename")
		filePath := "uploads/" + filename
		if err := os.Remove(filePath); err != nil {
			return c.Status(500).SendString("Impossible de supprimer le fichier")
		}
		return c.SendString("Fichier supprimé avec succès")
	})

	// Lancer le serveur sur le port 8080 en réseau local
	app.Listen(":8080")
}

func ListUploadedFiles(directory string) ([]string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	var fileList []string
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}
	return fileList, nil
}
