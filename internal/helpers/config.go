package helpers

import (
  "encoding/json"
  "errors"
  "log/slog"
  "os"
)

type Config struct{
  ServerAddress string `json:"serverAddress"`
  Port string `json:"port"`
  LogPath string `json:"logPath"`
  EnableOtel bool `json:"enableOtel"`
}

func LoadConfig(path, env string) (Config, error){
  var c Config
  if env!="dev" && env!="prod" {
    return c, errors.New("environment must be either 'dev' or 'prod'")
  }
  data, err := os.ReadFile(path)
  if err != nil {
    return c, err
  }

  var cfg map[string]json.RawMessage
  if err = json.Unmarshal(data, &cfg); err != nil {
      slog.Error("Fatal error parsing json", "error", err)
      os.Exit(1)
  }
  err = json.Unmarshal(cfg[env], &c)
  if err != nil {
    return c, err
  }
  return c, nil
}

