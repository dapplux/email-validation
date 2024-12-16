# Zerobounce Email Validation API

A simple Go-based client for validating email addresses using the [ZeroBounce Email Validation API](https://www.zerobounce.net/docs/email-validation-api-quickstart/).

## Features

- Validate email addresses via the ZeroBounce API.
- Optional IP address validation.
- Rate limiting for high-volume requests (up to 5000 requests/second).
- Easy configuration with environment variables.
- Exposes a local HTTP API for convenient use.

---

## Setup

### Prerequisites

- **Go**: Version 1.18 or later.
- **ZeroBounce API Key**: Obtain an API key by signing up at [ZeroBounce](https://www.zerobounce.net).

### Steps to Install and Run

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/zerobounce-api-client.git
   cd zerobounce-api-client
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Create a `.env` file in the project root with your API key:

   ```env
   ZEROBOUNCE_API_KEY=your-api-key-here
   ```

4. Run the application:

   ```bash
   go run main.go
   ```

---

## API Usage

### Endpoint

- **URL**: `http://localhost:8080/validate-email`
- **Method**: `POST`
- **Request Body**:
  - `email` (required): The email address to validate.
  - `ip_address` (optional): The IP address associated with the email.

### Request and Response Examples

#### Validate Email via cURL

Validate an email without an IP address:

```bash
curl -X POST http://localhost:8080/validate-email \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com"}'
```

Validate an email with an IP address:

```bash
curl -X POST http://localhost:8080/validate-email \
     -H "Content-Type: application/json" \
     -d '{"email":"valid@example.com", "ip_address":"156.124.12.145"}'
```

#### Response Example (Successful Validation)

```json
{
  "status": "valid",
  "error": ""
}
```

#### Response Example (Error)

```json
{
  "status": "invalid",
}
```

---

## Configuration

### Environment Variables

The application uses a `.env` file for configuration:

- `ZEROBOUNCE_API_KEY`: Your ZeroBounce API key.

### Rate Limiting

The application enforces a rate limit of 5000 requests per second by default. You can adjust this in the `NewZerobounceProvider` function:

```go
RateLimiter: rate.NewLimiter(rate.Every(time.Second/5000), 5000),
```

---

## Project Structure

- **`ZerobounceProvider`**: Encapsulates all logic for interacting with the ZeroBounce API.
- **HTTP API**: Exposes an endpoint at `http://localhost:8080/validate-email` for validating emails via the provider.

---

## Requirements

- **Go**: Version 1.18 or later.
- **ZeroBounce API Key**: Required to access the validation API.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

---

## Contributions

Feel free to submit issues or pull requests to improve this project. Contributions are always welcome!
