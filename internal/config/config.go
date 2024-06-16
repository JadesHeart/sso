package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPC          `yaml:"grpc"`
}

type GRPC struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

/*
MustLoad Функция возвращает значение конфига
если не получается по какой-то причине нормально его создать то дропает панику
*/
func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("путь до конфига пустой")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) { // проверяем существует ли файл с конфигом
		panic("конфигурационный файл не существует: " + configPath)
	}

	var config Config

	// с помощью библиотеки "github.com/ilyakaznacheev/cleanenv" читаем конфиг в переменную config
	// если не выходит, то выдаём панику
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic("не удалось прочитать конфиг: " + err.Error())
	}

	return &config
}

// Функция парсит путь до конфига сначала из флага и если не выходит из переменных окружения
func fetchConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "config", "config/local.yaml", "путь до конфига")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}
