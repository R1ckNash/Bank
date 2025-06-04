package config

type Config struct {
	Env   string `yaml:"env" env-default:"local" env-required:"true"`
	DBUrl string `yaml:"dbUrl" env-required:"true"`
}

func MustLoad() *Config {
	//configPath := os.Getenv("CONFIG_PATH")
	//if configPath == "" {
	//	log.Fatal("Config path could not be empty")
	//}
	//
	//if _, err := os.Stat(configPath); os.IsNotExist(err) {
	//	log.Fatalf("config file does not exists: %s", configPath)
	//}

	var cfg Config

	cfg.DBUrl = "postgres://postgres:postgres@postgres:5432/bank?sslmode=disable"
	cfg.Env = "dev"

	//if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
	//	log.Fatalf("could not read config: %s", err)
	//}

	return &cfg
}
