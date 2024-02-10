# **EWallet API**

The **EWallet API** is an implementation of an electronic wallet system that provides transaction processing in a payment environment. This application is designed as an HTTP server that provides a REST API for interacting with wallets.

The functionality of the project is based on the processing of the following basic operations:

1. Creating a wallet with an initial balance.
2. Transfer of funds between wallets.
3. Getting the wallet's transaction history.
4. Getting the current status of the wallet.

Each operation is represented by a corresponding API endpoint, providing a convenient and reliable way to manage wallets and transactions.

The application is implemented using the Go 1.21 programming language and PostgreSQL database in accordance with the stack requirements.

## Installation

1. Download Golang version 1.21+ from the official website: https://go.dev/dl/

2. Clone the repository:

```git clone https://github.com/MisterGnida/ewallet-rest.git``` 

3. Download all dependencies:

```go mod tidy```

4. Download and install Docker, docker-compose: https://www.docker.com/products/docker-desktop/

5. Download and install PostgreSQL: https://www.postgresql.org/download/

6. Create your database

7. Open ```config/server.toml``` and change this line:

```database_url = "user=your_username password=yourpass dbname=your_db host=db sslmode=disable"```

Enter your username, password, database name

8. Open ```docker-compose.yml``` and configure it

9. Open the terminal and navigate to the project directory

10. Build the server:
```docker-compose build```

11. Launch the server:
```docker-compose up```

**OPTIONALLY**. For tests change dbURL in ```server_test.go``` and ```wallet_repo_test.go```

You also can execute ```Makefile``` to launch the server.

## Usage examples

We need an HTTP Client to test the API, for example, Postman.

Download Postman: https://www.postman.com/downloads/

1. POST Request at ```/api/v1/wallet```

![изображение](https://github.com/MisterGnida/ewallet-rest/assets/102946972/d4137871-8bdf-476e-b1b2-d44d8557df6b)

2. GET Request at ```/api/v1/wallet/{walletId}```

![изображение](https://github.com/MisterGnida/ewallet-rest/assets/102946972/0e91781b-6bef-4371-8a6d-6b4708d0ca2f)

3. POST Request at ```/api/v1/wallet/{walletId}/send```

![изображение](https://github.com/MisterGnida/ewallet-rest/assets/102946972/75bb4827-5c5c-402c-825e-6deed846f6d2)

4. GET Request at ```/api/v1/wallet/{walletId}/history```

![изображение](https://github.com/MisterGnida/ewallet-rest/assets/102946972/c2d90419-a8ae-4729-bbda-37987bfc1bf2)

