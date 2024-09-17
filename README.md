# Laptop-Lounge-Microservice

Welcome to the **Laptop-Lounge-Microservice** project! This repository contains a microservices-based e-commerce platform specializing in laptop sales. The project is built using Go (Golang) and leverages gRPC for inter-service communication.

## Table of Contents

- [Introduction](#introduction)
- [Architecture Overview](#architecture-overview)
- [Directory Structure](#directory-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Running the Services](#running-the-services)
- [API Documentation](#api-documentation)
- [Contributing](#contributing)
- [License](#license)

## Introduction

Laptop-Lounge-Microservice is an e-commerce platform designed with a microservices architecture to ensure scalability, maintainability, and ease of deployment. Each service is responsible for a specific domain within the application, promoting a separation of concerns.

## Architecture Overview

The application consists of the following microservices:

- **API Gateway**: Acts as the entry point for all client requests, routing them to the appropriate services.
- **Admin Service**: Handles administrative tasks such as product management and order oversight.
- **Customer Service**: Manages user authentication, profiles, and customer-related data.
- **Product Service**: Oversees the product catalog, inventory management, and product details.
- **Cart Service**: Manages shopping cart functionalities for users.
- **Order Service**: Processes orders, payments, and maintains order histories.

These services communicate via gRPC and are built using clean architecture principles. The project also includes Protocol Buffer definitions (`.proto` files) for defining service interfaces.

## Directory Structure

```
Laptop-Lounge-Microservice/
├── Admin-Service/
├── API-Gateway/
├── Cart-Service/
├── Customer-Service/
├── Order-Service/
└── Product-Service/
```

Each service directory follows a similar structure:

- **cmd/**: Entry point of the service (`main.go`).
- **pkg/**: Contains the core packages, including API handlers, configurations, database connections, dependency injection, domain models, and more.
- **pb/**: Protocol Buffers generated code and `.proto` files.
- **.env**: Environment variables specific to the service.
- **go.mod & go.sum**: Go module files for dependency management.

## Prerequisites

Before running the application, ensure you have the following installed:

- **Go**: Version 1.16 or higher.
- **Protocol Buffers Compiler (protoc)**: For compiling `.proto` files.
- **Docker**: For containerization (optional but recommended).
- **Make**: For using the provided `Makefile` scripts (if available).

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/ShahabazSulthan/Laptop_Lounge_Microservice.git
   cd Laptop-Lounge-Microservice
   ```

2. **Set up environment variables**:

   - Copy the `.env` file in each service directory and configure them as needed.
   - Ensure database URLs, ports, and other configurations are correctly set.

3. **Install dependencies**:

   Navigate to each service directory and run:

   ```bash
   go mod download
   ```

## Running the Services

### Using Docker Compose (Recommended)

If you prefer containerization, you can use Docker Compose to run all services simultaneously.

1. **Build and run the containers**:

   ```bash
   docker-compose up --build
   ```

2. **Verify that all services are running**:

   Access the logs to ensure each service has started without errors.

### Running Locally

To run the services locally without Docker:

1. **Compile Protocol Buffers**:

   Navigate to the `pb/` directory of each service and run:

   ```bash
   protoc --go_out=. --go-grpc_out=. *.proto
   ```

2. **Run each service**:

   Open separate terminal windows for each service and run:

   ```bash
   go run cmd/main.go
   ```

   Ensure that ports used by each service do not conflict.

## API Documentation

API endpoints and their usage are documented using Swagger (or any other API documentation tool you prefer). To access the API documentation:

1. **Run the API Gateway**:

   Ensure the API Gateway service is running.

## Contributing

We welcome contributions! If you'd like to contribute:

1. **Fork the repository**.
2. **Create a new branch** for your feature or bug fix.
3. **Submit a pull request** with a detailed description of your changes.

Please ensure your code adheres to the existing code style and include appropriate tests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Feel free to explore the codebase, open issues, and contribute to making Laptop-Lounge-Microservice a robust e-commerce platform!

For any questions or support, please contact [shahabazsulthan4@gmail.com]
