# API Router

The API Router is a module in the Chirpy application responsible for handling various API endpoints. It uses the Chi router package to define and handle different routes for the Chirps, Users, Tokens, and Webhooks APIs.

## Endpoints

The following endpoints are defined in the API Router:

### Test API

- `/healthz` (GET): Endpoint for health checks and readiness probes.
- `/metrics` (GET): Endpoint for exposing application metrics.
- `/validate_chirp` (POST): Endpoint for validating a chirp.

### Chirps API

- `/chirps` (GET): Endpoint for fetching all chirps. Supports optional sorting by `sort` query parameter in ascending (`sort=asc`) or descending (`sort=desc`) order.
- `/chirps/{id}` (GET): Endpoint for fetching a single chirp by its ID.
- `/chirps` (POST): Endpoint for creating a new chirp.
- `/chirps/{id}` (DELETE): Endpoint for deleting a chirp by its ID.

### Users API

- `/users` (POST): Endpoint for creating a new user.
- `/login` (POST): Endpoint for user login. It issues a JWT access token for authentication.
- `/users` (PUT): Endpoint for updating a user's email and password.

### Token API

- `/refresh` (POST): Endpoint for refreshing an expired JWT access token.
- `/revoke` (POST): Endpoint for revoking a JWT access token.

### Webhooks API

- `/polka/webhooks` (POST): Endpoint for handling webhooks from the Polka service.

## Handlers

The API Router uses various handler functions to process incoming requests and generate responses. Each endpoint has its corresponding handler function responsible for handling the specific logic associated with that endpoint.

## Dependencies

The API Router relies on the following packages:

- `net/http`: For creating the HTTP server and handling HTTP requests.
- `github.com/go-chi/chi/v5`: A lightweight and fast HTTP router for Go.
- `database`: A custom package that provides access to the database.
- `jwt`: A package for handling JSON Web Tokens (JWT) for authentication.

## Usage

To use the API Router, you need to pass an `apiConfig` object to the `createAPIRouter` function. The `apiConfig` object should contain a database connection (`DB`) and a JWT secret (`jwtSecret`) for handling authentication.

Example:

```go
import (
	"net/http"
	"github.com/your-username/chirpy/apiRouter"
	"github.com/your-username/chirpy/database"
)

func main() {
	// Create a new database connection
	db := database.New()

	// Create an API configuration with the database and JWT secret
	apiCfg := &apiRouter.APIConfig{
		DB:        db,
		JWTSecret: "your-jwt-secret-key",
	}

	// Create the API router
	router := apiRouter.CreateAPIRouter(apiCfg)

	// Start the HTTP server
	http.ListenAndServe(":8080", router)
}
