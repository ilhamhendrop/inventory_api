package main

import (
	"inventory-app/internal/api"
	"inventory-app/internal/config"
	"inventory-app/internal/database"
	"inventory-app/internal/repository"
	"inventory-app/internal/service"
	"inventory-app/internal/util"
	"time"

	jwt "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cnf := config.Get()
	dbConnect := database.GetMySQLDB(cnf.MysqlDB)
	cacheConnect := database.GetRedisCache(cnf.RedisDB)

	app := fiber.New()
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))
	app.Use(helmet.New())

	jwt := jwt.New(jwt.Config{
		SigningKey: jwt.SigningKey{
			Key: []byte(cnf.Jwt.Key),
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return util.Unauthorized(ctx, err)
		},
	})

	authRepo := repository.NewAuth(dbConnect, cacheConnect, cnf.Server.Host)
	userRepo := repository.NewUser(dbConnect, cacheConnect, cnf.Server.Host)
	categorieRepo := repository.NewCategorie(dbConnect, cacheConnect, cnf.Server.Host)
	productRepo := repository.NewProduct(dbConnect, cacheConnect, cnf.Server.Host)
	warehouseRepo := repository.NewWarehouse(dbConnect, cacheConnect, cnf.Server.Host)
	maintenanceRepo := repository.NewMaintenance(dbConnect, cacheConnect, cnf.Server.Host)
	financeRepo := repository.NewFinance(dbConnect, cacheConnect, cnf.Server.Host)

	authService := service.NewAuth(cnf, authRepo)
	userService := service.NewUser(userRepo)
	categorieService := service.NewCategorie(categorieRepo)
	productService := service.NewProduct(productRepo, categorieRepo)
	warehouseService := service.NewWarehouse(warehouseRepo, productRepo, categorieRepo)
	maintenanceService := service.NewMaintenance(maintenanceRepo, productRepo, categorieRepo, warehouseRepo, userRepo)
	financeService := service.NewFinance(financeRepo, maintenanceRepo, userRepo, productRepo, categorieRepo)

	api.NewAuth(app, authService)
	api.NewUser(app, userService, jwt)
	api.NewCategorie(app, categorieService, jwt)
	api.NewProduct(app, productService, jwt)
	api.NewWarehouse(app, warehouseService, jwt)
	api.NewMaintenance(app, maintenanceService, jwt)
	api.NewFinance(app, financeService, jwt)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
