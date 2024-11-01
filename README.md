# Post Service

This repository contains the Post Service for a social media platform. This service is responsible for handling post-related functionalities, such as creating, retrieving, updating, and deleting posts. It is designed to work with an Oracle database and is integrated with Eureka for service discovery.

To view all services for this social media system, lets visit: `https://github.com/goddie9x?tab=repositories&q=social`

## Features

- **CRUD Operations**: Supports create, read, update, and delete operations for posts.
- **Oracle Database Integration**: Connects to an Oracle database for persistent data storage.
- **Service Discovery**: Registers with Eureka for dynamic service registration and discovery.
- **Environment Configuration**: Configurable using environment variables specified in a `.env` file.

## Environment Variables

The following environment variables are expected in the `.env` file for configuring the Post Service:

| Variable                        | Description                                 |
|---------------------------------|---------------------------------------------|
| `DB_HOST`                       | Hostname of the Oracle database.           |
| `DB_PORT`                       | Port for connecting to the Oracle database.|
| `DB_SERVICE`                    | Service name of the Oracle database.       |
| `DB_USER`                       | Username for accessing the Oracle database.|
| `DB_PASSWORD`                   | Password for the database user.            |
| `DB_CONNECTION_TIMEOUT`         | Timeout for database connections in seconds.|
| `DB_SSL`                        | Whether to use SSL for the database connection (true/false).|
| `API_PORT`                      | Port on which the Post Service will listen for incoming requests.|
| `APP_IP_ADDRESS`                | IP address of the Post Service.            |
| `EUREKA_APP_NAME`              | Name of the application for Eureka registration.|
| `HOST_NAME`                     | Hostname for the Post Service.             |
| `EUREKA_DISCOVERY_SERVER_URL`   | URL for the Eureka discovery server.       |

### Example `.env` File

Here is an example of how your `.env` file should look:

```env
DB_HOST=oracle
DB_PORT=1521
DB_SERVICE=XE
DB_USER=system
DB_PASSWORD=thisIsJustTheTestPassword123
DB_CONNECTION_TIMEOUT=90
DB_SSL=false
API_PORT=3005
APP_IP_ADDRESS=post-service
EUREKA_APP_NAME=post-service
HOST_NAME=post-service
EUREKA_DISCOVERY_SERVER_URL=http://discovery-server:8761/eureka
```

## Docker Configuration

The Post Service can be easily deployed using Docker. Below is the `Dockerfile` used to build the service image.

### Dockerfile

```dockerfile
FROM alpine:latest

# Install necessary packages
RUN apk add --no-cache file

# Set the working directory
WORKDIR /app

# Copy the application binary and environment file
COPY ./main ./.env .

# Expose the API port
EXPOSE 3005

# Command to run the service
CMD ["./main"]
```

### Building the Docker Image

To build the Docker image for the Post Service, run the following command in the root directory of the repository:

```bash
docker build -t post-service .
```

### Running the Docker Container

Once the image is built, you can run the Post Service container without specifying environment variables, as they are included in the `.env` file:

```bash
docker run -d \
  --name post-service \
  -p 3005:3005 \
  post-service
```

## Running the Post Service Locally

If you prefer to run the service locally without Docker, ensure that you have the required environment variables set up in your local environment. You can then run the application using your preferred method (e.g., using a Java IDE or command line).

## Service Registration

The Post Service automatically registers with the Eureka discovery server at `http://discovery-server:8761/eureka`. Ensure that the discovery server is running before starting the Post Service.

## Contributing

Contributions are welcome! Please fork this repository and submit a pull request with your changes. Ensure that your changes are tested and documented.

## License

This project is licensed under the MIT License. See `LICENSE` for more details.