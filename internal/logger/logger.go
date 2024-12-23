package logger

import (
  "log"
  "os"
)
var Logger = log.New(os.Stdout, "log_ak0_2: ", log.Ldate)
