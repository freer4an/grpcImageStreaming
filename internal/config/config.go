package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App   app      `yaml:"app"`
	Paths paths    `yaml:"paths"`
	Db    database `yaml:"db"`
}

func (c *Config) GetDbUrl() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.Db.User,
		c.Db.Pass,
		c.Db.Host,
		c.Db.Port,
		c.Db.Name)
}

func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.App.host, c.App.port)
}

type app struct {
	host         string   `yaml:"host" env:"HOST"`
	port         string   `yaml:"port" env:"PORT"`
	ImageFormats []string `yaml:"image_formats"`
}

type paths struct {
	Images            string `yaml:"images"`
	OImagesStorage    string `yaml:"o_images_storage"`
	ThumbnailsStorage string `yaml:"thumbnails_storage"`
}

type database struct {
	Name    string `yaml:"name" env:"DB_NAME"`
	Host    string `yaml:"host" env:"DB_HOST"`
	Port    string `yaml:"port" env:"DB_PORT"`
	User    string `yaml:"user" env:"DB_USER"`
	Pass    string `yaml:"password" env:"DB_PASS"`
	SSLMode string `yaml:"ssl_mode" env-default:"disable"`
}

func New(path string) *Config {
	var config Config

	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		panic(err)
	}

	err = cleanenv.ReadEnv(&config)
	if err != nil {
		panic(err)
	}

	return &config
}
