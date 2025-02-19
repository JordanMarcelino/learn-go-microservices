package feign

import (
	"net/http"
	"time"

	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/jwtutils"
)

var (
	feignClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	jwtUtil = jwtutils.NewJwtUtil()
)

const (
	APPLICATION_JSON = "application/json"
)
