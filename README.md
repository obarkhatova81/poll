# LivePoll

## Overview

LivePoll is a real-time polling system that allows users to create polls, vote in real-time, and view live results.

## Features

- **Poll Creation**: Users can create polls with multiple choice answers.
- **Voting**: Real-time voting for active polls.
- **Live Results**: Real-time update of poll results using WebSockets.

## Requirements

- Go >= 1.20
- Redis
- Docker


## Build the Docker Image

To build the Docker image for the LivePoll application, run the following command:

```bash
docker build -t livepoll .
```

## Run the Docker Container

To run the Docker container, use the following command:

```bash
docker run -p 8080:8080 -p 8081:8081 --name livepoll-container livepoll
```

## Endpoints

### HTTP Endpoints

- **POST /polls**
  Create a new poll with a question and multiple choice answers.

- **GET /polls**
  Retrieve a list of all polls.

- **GET /polls/{id}**
  Retrieve a specific poll by its unique ID.

- **PUT /polls/{id}**
  Update a specific poll by its unique ID.

- **DELETE /polls/{id}**
  Delete a specific poll by its unique ID.

### WebSocket Endpoint

- **/ws**
  Establish a WebSocket connection to receive real-time updates on poll results.

## Local Development

To run the application locally, follow these steps:

**Clone the Repository**
   ```sh
   git clone https://github.com/obarkhatova81/poll.git
   cd poll
   go mod tidy
```
**Run Redis**
   ```sh
   docker run --name redis -p 6379:6379 -d redis
  ```
    
**Run the Application**
```sh
   go run cmd/main.go
  ```

**Access the Application**

HTTP server: http://localhost:8080
WebSocket server: ws://localhost:8080/ws

