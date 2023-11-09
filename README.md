# Warehouse API Project

This project is an API for managing warehouse data.

## Table of Contents

- [Build and Run](#build-and-run)
- [Usage of the API](#usage-of-the-api)
- [Testing](#testing)
- [Generate Documentation](#generate-documentation)

## Build and Run

To build and run the project, follow these steps:

1. Clone the repository:

    ```bash
    git clone https://github.com/DmitriiKumancev/lamoda-test.git
    cd lamoda-test
    ```

2. Start the containers using the `Makefile`:

    ```bash
    make postgres
    ```

3. Create the database:

    ```bash
    make createdb
    ```

4. Create migrations:

    ```bash
    make migratecreate
    ```

5. Apply migrations:

    ```bash
    make migrateup
    ```

This will start PostgreSQL, apply migrations, and run the application container.

## Usage of the API

1. Get started:

    ```bash
    cd app/cmd
    go run main.go
    ```

The API is accessible at `http://localhost:8080`. You can use Swagger UI for documentation and example requests.

Swagger UI: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## Testing

To run unit tests, execute the following command:

```bash
make test
```

## Generate Documentation

To generate documentation using Swag, run the following command:

```bash
make swagger
```

---

**Author:** [Dmitrii Kumancev](https://github.com/DmitriiKumancev)
