# Go, Gin, Redis, PG and Gorm - Basic API utilising caching

![Gopher with Umbrella - Project Title Image](https://raw.githubusercontent.com/egonelbre/gophers/63b1f5a9f334f9e23735c6e09ac003479ffe5df5/vector/dandy/umbrella.svg)


A simple CRUD API using Gin, Docker, Gorm, and Redis. Initial project was highly-coupled and had a handful. The project has been refactored and utilises layered architecture and interfaces to decouple controllers (handlers), caching, and storage for handling data with **PostgreSQL and Redis caching**. Implements a **decoupled caching layer**, where Redis optimizes reads while PostgreSQL remains the source of truth. Uses **Gin for API routing**, **GORM for database access**, and uses a **repository abstraction**.  

## Features  
- **Layered architecture** (Handlers → Cached Repository → Postgres Repository)  
- **Cache-first reads**: Queries Redis first, falls back to Postgres on a cache miss  
- **Write-through caching**: Updates Postgres and invalidates Redis cache on writes  
- **Dependency inversion**: The caching layer holds a reference to the Postgres repository, not the other way around  
- **Extensible design**: Easily swap Redis/Postgres without breaking API logic  


## Getting Started

### Prerequisites
- Docker
- Docker Compose
- Go (1.18+ recommended)

## Installation

1. Clone this repository:
   ```sh
   git clone https://github.com/travboz/gin-redis-gorm-pg.git
   cd gin-redis-gorm-pg
   ```
2. Run docker compose for Postgres and Redis containers:
    ```sh
    make up
    ```
3. Run server:
    ```sh
    make run
    ```
4. Navigate to `http://localhost:8080` and call an endpoint


## API endpoints

| Method   | Endpoint                     | Description                         |
|----------|-----------------------------|-------------------------------------|
| `POST`   | `/v1/products/`              | Create a new product               |
| `GET`    | `/v1/products/{id}`          | Get product by ID                  |
| `PUT`    | `/v1/products/{id}`          | Update a product                   |
| `DELETE` | `/v1/products/{id}`          | Delete a product                   |
| `POST`   | `/v1/products/invalidate/{id}` | Invalidate product in cache        |
| `GET`    | `/v1/products/recent`        | Get recent products                |

## Example usage

### JSON payload structures

#### Create product payload

```json
{
  "id": "1",
  "name": "bbq shapes",
  "price": 5
}
```

#### Update product payload

```json
{
  "name": "bbq shapes",
  "price": 15
}
```

### Endpoint example usage
#### Create a user
```sh
curl -X POST "http://localhost:8080/v1/products" \
     -H "Content-Type: application/json" \
     -d '{
       "id": "2",
       "name": "pizza shapes",
       "price": 6
     }'
```

#### Update a user
```sh
curl -X PUT "http://localhost:8080/v1/products/3" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "pizza shapes",
       "price": 12
     }'
```

#### Get a product by ID
```sh
curl -X GET "http://localhost:8080/v1/products/3"
```

#### Delete a product
```sh
curl -X DELETE "http://localhost:8080/v1/products/3"
```

#### Invalidate product in cache
```sh
curl -X POST "http://localhost:8080/v1/products/invalidate/3"
```

#### Get recent products
```sh
curl -X GET "http://localhost:8080/v1/products/recent"
```


## Contributing
Feel free to fork and submit PRs!

## License:
`MIT`

This should work for GitHub! Let me know if you need any tweaks. 


## Image
Image by [Egon Elbre](https://github.com/egonelbre), used under CC0-1.0 license.
