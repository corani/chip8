package web

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/corani/chip-8/internal/chip8"
	fiber "github.com/gofiber/fiber/v2"
)

type App struct {
	log  *log.Logger
	roms []string
	root string
}

func New(log *log.Logger, chip8 *chip8.Chip8) *App {
	app := &App{
		log: log,
	}

	app.findRoms()

	return app
}

func (app *App) Run() error {
	web := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
		StreamRequestBody: true,
		ErrorHandler:      app.HandleError,
	})

	web.Static("/", "./static").
		Name("Static")

	web.Get("/api/roms", app.HandleRoms).
		Name("List ROMs")
	web.Get("/api/rom/*", app.HandleRom).
		Name("Load ROM")

	if err := web.Listen(":8384"); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (app *App) HandleError(c *fiber.Ctx, err error) error {
	log.Errorf("Error: %v", err)

	return c.SendStatus(fiber.StatusInternalServerError)
}

func (app *App) HandleRoms(c *fiber.Ctx) error {
	return c.JSON(app.roms)
}

func (app *App) HandleRom(c *fiber.Ctx) error {
	name := c.Params("*")
	if name == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	name = app.root + name

	// urldecode name
	name, err := url.PathUnescape(name)
	if err != nil {
		app.log.Errorf("Error decoding URL: %v", err)

		return c.SendStatus(fiber.StatusBadRequest)
	}

	romData, err := os.ReadFile(name)
	if err != nil {
		app.log.Errorf("Error reading ROM file: %v", err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Send(romData)
}

func (app *App) findRoms() {
	// find the root folder of the project (the one with `go.mod`), then recursively
	// search for `.ch8` files and collect them into `app.roms`.
	root, err := os.Getwd()
	if err != nil {
		app.log.Errorf("Error getting working directory: %v", err)

		return
	}

	for {
		info, err := os.Stat(root + "/go.mod")
		if err == nil && !info.IsDir() {
			break
		}

		root = filepath.Dir(root)
		if root == "/" {
			app.log.Errorf("Error: could not find project root")

			return
		}
	}

	app.log.Infof("Project root: %s", root)
	app.root = root + "/"

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			app.log.Errorf("Error walking path %s: %v", path, err)

			return nil
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".ch8" {
			app.roms = append(app.roms, strings.TrimPrefix(path, app.root))
		}

		return nil
	})
	if err != nil {
		app.log.Errorf("Error walking project root: %v", err)
	}
}
