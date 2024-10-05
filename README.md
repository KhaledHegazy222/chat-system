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

1- Clone the Repository:

```bash
git clone https://github.com/KhaledHegazy222/chat-system.git
cd chat-system
```

2.Set Up Environment Variables:

Create a .env file in the root directory and add the following variables:

```
MYSQL_USER=root
MYSQL_PASSWORD=secret
MYSQL_HOST=mysql
MYSQL_DATABASE=chat_system
REDIS_URL=redis://redis:6379/0
ELASTICSEARCH_URL=http://elasticsearch:9200
```

3. Start the Application:

To run the application and all its services (MySQL, Redis, Elasticsearch, Sidekiq, Rails), simply run:

```
docker-compose up --build
```

4. Seed the Database (Optional):

```bash
docker-compose exec web rails db:seed
```
