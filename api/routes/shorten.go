package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/aayush-makhija/url-shortener/database"
	"github.com/aayush-makhija/url-shortener/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type request struct {
	URL         string        `json: "url"`
	customShort string        `json: "short"`
	expiry      time.Duration `json: "expiry"`
}

type response struct {
	URL             string        `json: "url"`
	customShort     string        `json: "short"`
	expiry          time.Duration `json: "expiry"`
	xRateRemaining  int           `json: "rate_limit"`
	xRateLimitReset time.Duration `json: "rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error:": "cannot parse JSON"})
	}

	//implement rate limiting
	r2 := database.CreateClient(1)
	defer r2.Close()
	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		val, _ = r2.Get(database.Ctx, c.IP()).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	//check if the input is an actual URL

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	//check for domain error

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "you can't hack the system (: "})
	}

	//enforce https, SSL

	body.URL = helpers.EnforceHTTP(body.URL)

	var id string

	if body.customShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.customShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL custom short is already in use",
		})
	}
	if body.expiry == 0 {
		body.expiry = 24
	}

	err = r.Set(database.Ctx, id, body.URL, body.expiry*3600*time.Second).Err()

	if err != nil {
		return C.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server.",
		})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.expiry,
		xRateRemaining:  10,
		xRateLimitReset: 30,
	}

	r2.Decr(database.Ctx, c.IP())

	cal, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.xRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.xRateLimitReset = ttl / time.Nanosecond / timme.Minute

	resp.customShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)
}
