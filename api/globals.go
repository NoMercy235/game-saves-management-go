package api

import (
	"time"
	"math/rand"
)

var ERR_TIME = 0 * time.Second
var DELAY_TIME = 500 * time.Millisecond
var PING_TIME = 10 * time.Second
var MESSAGE_TIME = time.Duration(rand.Int31n(5000)) * time.Millisecond