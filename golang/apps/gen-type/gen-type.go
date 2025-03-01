package main

import (
	"de-net/libs/3_infrastructure/db_manager"
	"de-net/libs/4_common/env_vars"
	"de-net/libs/4_common/smart_context"
	"os"

	"gorm.io/gen"
)

type Querier interface {
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {
	env_vars.LoadEnvVars()
	os.Setenv("LOG_LEVEL", "info")
	logger := smart_context.NewSmartContext()

	dbm, err := db_manager.NewDbManager(logger)
	if err != nil {
		logger.Fatalf("NewDbManager failed: %v", err)
	}

	logger = logger.WithDB(dbm.GetGORM())

	// Конфигурация генератора
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./libs/2_generated_models/model", // Папка для генерации
		OutFile:           "model.gen.go",
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.UseDB(logger.GetDB())

	// Генерация всех таблиц
	g.GenerateAllTable()

	// Применение базовых настроек
	g.ApplyBasic()

	// Выполнение генерации
	g.Execute()

}
