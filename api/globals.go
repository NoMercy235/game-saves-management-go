package api

import (
	"time"
	"math/rand"
	"os"
	"sync"
)

var cwd, _ = os.Getwd()
var FILES_PATH = cwd + "/"
var ERR_TIME = 0 * time.Second
var DELAY_TIME = 500 * time.Millisecond
var PING_TIME = 10 * time.Second
var CLOCK_SYNC_TIME = 10 * time.Second
var MESSAGE_TIME = time.Duration(rand.Int31n(5000)) * time.Millisecond
var MUTEX = &sync.Mutex{}
var EXECUTE_COMMAND_DELAY = DELAY_TIME + 500 * time.Millisecond
var COMMAND_QUEUE_LIMIT = 100