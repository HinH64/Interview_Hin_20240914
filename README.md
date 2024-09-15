Table of Contents
-----------------

- [Features](#features)
- [Env Variables](#env-variables)
- [Project Structure](#project-structure)
- [API Endpoints](#api-endpoints)
- [Player Management System](#player-management-system)
- [Challenge Game System](#challenge-game-system)
- [Endless Challenge System](#endless-challenge-system)
- [Game Log Collector](#game-log-collector)
- [Payment Processing System](#payment-processing-system)

Features
--------

- [Gin](https://github.com/gin-gonic/gin)
- DotEnv with [Viper](https://github.com/spf13/viper)
- Request Validation with [Ozzo Validation](https://github.com/go-ozzo/ozzo-validation)
- MongoDB ODM ([mgm](https://github.com/Kamva/mgm))
- Docker support
- Swagger with [gin-swagger](https://github.com/swaggo/gin-swagger)

Env Variables
--------------
The application uses the following environment variables:
SERVER_PORT: The port on which the server will run. Default is "8080".
SERVER_ADDR: The address on which the server will listen.
MONGO_URI: The MongoDB connection URI.
MONGO_DATABASE: The name of the MongoDB database to use.
MONGO_TEST_DATABASE: The name of the MongoDB database to use for testing.
MODE: The running mode of the application. Can be either "debug" or "release".

Project Structure
-----------------

```
├── controllers         # contains api functions and main business logic
├── docs                # swagger files 
├── logs
├── middlewares         # request/response middlewares
│   └── validators      # data/request validators
├── models              
│   └── db              # collection models
├── routes              # router initialization
└── services            # general service & database actions
```

API Endpoints
-----------------
Please refer to the Swagger documentation for the detailed API endpoints and their descriptions.

Player Management System
-----------------
GET /players - List all players
POST /players - Create a new player
GET /players/{id} - Get a specific player
PUT /players/{id} - Update a player
DELETE /players/{id} - Delete a player

GET /levels - List all levels
POST /levels - Create a new level

Challenge Game System
-----------------
GET /rooms - List all rooms
POST /rooms - Create a new room
GET /rooms/{id} - Get a specific room
PUT /rooms/{id} - Update a room
DELETE /rooms/{id} - Delete a room

GET /reservations - List all reservations
POST /reservations - Create a new reservation

Endless Challenge System
-----------------
GET /reservations - List all reservations
POST /reservations - Create a new reservation

GET /challengegames - List all challenge games
POST /challengegames - Create a new challenge game
GET /challengegames/{id} - Get a specific challenge game
PUT /challengegames/{id} - Update a challenge game
DELETE /challengegames/{id} - Delete a challenge game

Game Log Collector
-----------------
GET /logs - List all logs
POST /logs - Create a new log

Payment Processing System
-----------------
POST /payments - Create a new payment
GET /payments/{id} - Get a specific payment