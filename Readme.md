# Gobank

Gobank is a simple bank application that allows users to create accounts, deposit, withdraw and transfer money between accounts.
It is a simple RESTful API that uses a PostgreSQL database to store user information and account balances.
Uses JWT for authentication and authorization.

This project has been created following the tutorial [Complete JSON API project in Golang](https://www.youtube.com/playlist?list=PL0xRBLFXXsP6nudFDqMXzrvQCZrxSOm-2) on his [Youtube Channel](https://www.youtube.com/@anthonygg_).

## Installation

1. Clone the repository
2. Make sure docker and docker-compose are installed on your machine
   3 Run the following command to start the application

```bash
docker-compose up -d && docker-compose logs -f web
```
