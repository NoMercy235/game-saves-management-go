package api

import (
	"time"
	"math/rand"
	"os"
	"sync"
)

var cwd, _ = os.Getwd()
var FILES_PATH = cwd + "/tmp/"
var ERR_TIME = 0 * time.Second
var DELAY_TIME = 500 * time.Millisecond
var PING_TIME = 10 * time.Second
var MESSAGE_TIME = time.Duration(rand.Int31n(5000)) * time.Millisecond
var MUTEX = &sync.Mutex{}
var EXECUTE_COMMAND_DELAY = 300 * time.Millisecond