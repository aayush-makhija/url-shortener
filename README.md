
---

# URL Shortener

## Overview

This project is a URL shortening service built with Go, Redis, and the Fiber web framework. It allows users to shorten long URLs and access them via a short alias. The service also includes rate limiting to prevent abuse and ensures the shortened URLs are unique.

## Features

- Shorten long URLs to a custom alias
- Redirect short URLs to the original URL
- Rate limiting to prevent abuse
- Validation of URLs to ensure they are valid
- Avoid domain errors
- Enforce HTTP/HTTPS

## Prerequisites

- Docker
- Docker Compose
- Go 1.16 or higher

## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/aayush-makhija/url-shortener.git
   cd url-shortener
   ```

2. **Set up environment variables:**
   Create a `.env` file in the root directory of the project and add the following environment variables:
   ```
   DB_ADDR=localhost:6379
   DB_PASS=your_redis_password
   DOMAIN=your_domain
   API_QUOTA=10
   APP_PORT=:3000
   ```

## Usage

1. **Run the application:**
   Using Docker Compose:
   ```sh
   docker-compose up --build
   ```

   Without Docker:
   ```sh
   go mod tidy
   go run main.go
   ```

2. **API Endpoints:**

   - **Shorten URL:**
     ```
     POST /api/v1
     {
       "url": "https://example.com",
       "short": "customAlias",  // Optional
       "expiry": 24  // Optional, in hours
     }
     ```

   - **Resolve URL:**
     ```
     GET /:url
     ```

## Acknowledgments

- [Fiber](https://gofiber.io/) - Fast, Express-inspired web framework for Go
- [Redis](https://redis.io/) - In-memory data structure store, used as a database, cache, and message broker
- [GoRedis](https://github.com/go-redis/redis) - Type-safe Redis client for Golang

---

