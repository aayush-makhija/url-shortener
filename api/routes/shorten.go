package routes

import "time"

type request struct {
	URL string  `json: "url"`
	customShort string `json: "short"`
	expiry time.Duration `json: "expiry"`

}

type response struct {
	URL string `json: "url"`
	customShort string `json: "short"`
	expiry time.Duration `json: "expiry"`
	xRateRemaining int `json: "rate_limit"`
	xRateLimitReset time.Duration `json: "rate_limit_reset"`
}