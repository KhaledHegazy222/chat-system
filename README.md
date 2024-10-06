# Backend Challenge - Chat System

## Overview

This project is a chat system that allows for the creation of applications, chats, and messages. Each application has a unique token, and each chat and message has an incrementing number within its application or chat. The project also supports full-text search on message content using Elasticsearch.

The system is built using **Ruby on Rails** for the API and workers. It uses **MySQL** as the main datastore, **Redis** for queuing, and **Elasticsearch** for message search. The API is RESTful, providing endpoints to create, read, and update applications, chats, and messages. This project is containerized using **Docker**.

## Table of Contents

1. [Prerequisites](#prerequisites)
1. [Getting Started](#getting-started)

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

### Seed the Database (Optional):

```bash
docker compose exec rails_app rake chat:seed
```
