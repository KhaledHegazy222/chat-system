# Backend Challenge - Chat System

## Overview

This project is a chat system that allows for the creation of applications, chats, and messages. Each application has a unique token, and each chat and message has an incrementing number within its application or chat. The project also supports full-text search on message content using Elasticsearch.

The system is built using **Ruby on Rails** for the API and workers. It uses **MySQL** as the main datastore, **Redis** for queuing, and **Elasticsearch** for message search. The API is RESTful, providing endpoints to create, read, and update applications, chats, and messages. This project is containerized using **Docker**.

## Table of Contents

1. [Prerequisites](#prerequisites)
1. [Getting Started](#getting-started)
1. [API Documentation](#api-documentation)

## Prerequisites

Before you begin, ensure you have the following installed:

- Docker
- Docker Compose

## Getting Started

### Clone the Repository:

```bash
git clone https://github.com/KhaledHegazy222/chat-system.git
cd chat-system
```

### Set Up Environment Variables:

Create a .env file in the root directory and add the following variables, you can also use `.env.example` (rename it `.env`)

```
REDIS_HOST=redis
REDIS_PORT=6379

DB_HOST=db
DB_PORT=3306
DB_DATABASE= chat_system
DB_ROOT_PASSWORD= my-secret-pw
DB_USER=admin
DB_PASSWORD=admin

RAILS_ENV=development

ELASTICSEARCH_URL=http://elastic_search:9200
ELASTICSEARCH_DISCOVERY_TYPE=single-node
```

### Start the Application:

To run the application and all its services (MySQL, Redis, Elasticsearch, Sidekiq, Rails), simply run:

```
docker compose up --build
```

### Run Tests:

```bash
# Run GO API Tests
docker compose exec go_app go test --cover ./internal/api
# Run RAILS API Tests
docker compose exec rails_app rails test
```

### Seed the Database (Optional):

```bash
docker compose exec rails_app rake chat:seed
```

## API Documentation

### Overview

This API consists of two services: **Go** and **Ruby**. Each service has its own set of endpoints, which are described below.

---

### Go Service Endpoints

#### Create Chat

- URL: `/applications/:application_token/chats`
- Method: POST
- Request Body:

```json
{
  "title": "string"
}
```

- Response:
  - 201 Created: Returns the created chat.

```json
{
  "number": "number",
  "application_token": "string",
  "title": "string"
}
```

#### Create Message

- URL: `/applications/:application_token/chats/:chat_number/messages`
- Method: POST
- Request Body:

```json
{
  "content": "string"
}
```

- Response:
  - 201 Created: Returns the created resource.

```json
{
  "number": "number",
  "chat_number": "number",
  "application_token": "string",
  "content": "string"
}
```

## Rails Service Endpoints

#### Create Application

- URL: `/applications`
- Method: POST
- Request Body:

```json
{
  "name": "string"
}
```

- Response:
  - 201 Created: Returns the created resource.

```json
{
  "name": "string",
  "token": "string"
}
```

#### List All Applications

- URL: `/applications`
- Method: GET

- Response:
  - 200 OK: Returns the created resource.

```json
[
  {
    "name": "string",
    "token": "string"
  },
  {
    "name": "string",
    "token": "string"
  }
]
```

#### Fetch Application by token

- URL: `/applications/:token`
- Method: GET

- Response:
  - 200 OK: Returns the created resource.

```json
{
  "name": "string",
  "token": "string"
}
```

#### List All Chats

- URL: `/chats`
- Method: GET

- Response:
  - 200 OK: Returns the created resource.

```json
[
  {
    "number": "number",
    "title": "string"
  },
  {
    "number": "number",
    "title": "string"
  }
]
```

#### Fetch Chat by app token and chat number

- URL: `/applications/:application_token/chats/:chat_number`
- Method: GET

- Response:
  - 201 Created: Returns the created resource.

```json
{
  "number": "number",
  "title": "string"
}
```

#### Edit Chat title

- URL: `/applications/:application_token/chats/:chat_number`
- Method: POST
- Request Body:

```json
{
  "title": "string"
}
```

- Response:
  - 200 OK: Returns the created resource.

```json
{
  "status": "success"
}
```

#### List All Messages

- URL: `/messages/`
- Method: GET

- Response:
  - 200 OK

```json
[
  {
    "number": "number",
    "content": "string"
  },
  {
    "number": "number",
    "content": "string"
  }
]
```

#### List All Message in a specific chat

- URL: `/applications/:application_token/chats/:chat_number/messages?q=`
- Method: GET
- Query Parameters:
  - key: `q` 
  - usage: search string to use elastic search in messages body
- Response:
  - 200 OK

```json
[
  {
    "number": "number",
    "content": "string"
  },
  {
    "number": 1,
    "content": "string"
  }
]
```

#### Fetch Message by app token and chat number and chat number

- URL: `/applications/:application_token/chats/:chat_number/messages/:messages_number`
- Method: GET

- Response:
  - 200 OK

```json
{
  "number": "number",
  "content": "string"
}
```

#### Edit Message content

- URL: `/applications/:application_token/chats/:chat_number/messages/:messages_number`
- Method: POST
- Request Body:

```json
{
  "content": "string"
}
```

- Response:
  - 200 OK

```json
{
  "status": "success"
}
```
