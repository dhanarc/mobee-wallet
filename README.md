# Mobee Wallet


### Introduction
Here is the backend apps to store user balance using ledger system.
This wallet apps store precision using big.Float and store in database as string to maintain the precisions.

### How to Run
1. If you want using docker environment run `docker-compose up`
2. Copy configs file using this config
   ```shell
   cp config.example.yml config.yml
   cp .env.example .env
   ```
3. Adjust your config with your environment (postgresql DSN)
4. Update your .env with your environment also (this is used for running database migrations)
5. Run db migrations
   ```shell
   goose up
   ```
6. Run services using this command
   ```shell
   go run main.go http
   ```
   if you see, logs like below:
   ```shell
   2025/01/30 18:13:45 INFO http server started at port: 8000
   ```
   The Apis ready to receive requests.
7. Postman collection is already available in root project, you can directly access this api using this api collections.
