# GoLedger Challenge - Besu Edition

This REST API connects to a Besu Network with a Smart Contrat deployed, sets a value to it, reads and syncronize the value to a PostgreSQL database deployed with Docker

## **Documentation**

Built with Go’s standard net/http package
- Folder Structure
```
	go-besu-api/
	├── besu/      			 # Besu network setup and smart contract deploy

	├── api/   
 	  └── abi/                 	 # Contains compiled ABI JSON for the smart contract
	  	└── SimpleStorage.json
	  └── cmd/
	  	└── api/
	  		└── main.go     # Entry point — initializes app, loads env, starts server
	  └── database/            	# PostgreSQL connection setup
		└── db.go
	  └── handler/             	# REST API route definitions (/set, /get, /sync, /check)
		└── handler.go
	  └── contract/            	# Smart contract interaction logic (get/set via Besu)
		└── contract.go
	  └── go.mod / go.sum      	# Go module dependencies
	  └── .env                 	# Environment configuration (RPC URL, DB credentials, contract address)
```
-  REST API
Built with Go’s standard net/http package.
Exposes 4 endpoints:
	- /get
	- /set
	- /sync
	- /check
  
- Besu Network (via Hardhat)
	- The app connects to a local Besu node using the Ethereum JSON-RPC interface (ethclient).
	- Contract interactions use the contract ABI (abi/SimpleStorage.json) for encoding/decoding method calls.
	- Uses the go-ethereum libraries to call contract methods
 
- PostgreSQL Database
	- Stores blockchain state locally for comparison and audit purposes.
	- Has a single table storage `(id SERIAL, address TEXT, value INTEGER)` to store values synced from the smart contract.
	- db.go sets up the DB connection using credentials from .env.
- Environment Setup via .env
```
PRIVATE_KEY=8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63
# Update CONTRACT_ADDRESS variable with the address printed and the end of the deploy script including 0x
CONTRACT_ADDRESS=0x42699A7612A82f1d9C36148af9C77354759b210b
NODE_URL=http://127.0.0.1:8545
# Update PG_CONN with your database connection string
PG_CONN=postgres://admin:admin@localhost:5432/simplestorage?sslmode=disable
CHAIN_ID=1337
```
 
## Prerequisites
- Linux 
- [Go](https://go.dev/)
- [Besu](https://besu.hyperledger.org/private-networks)
- [Hardhat](https://hardhat.org/hardhat-runner/docs/getting-started#overview)
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
- [jq ](https://jqlang.org/download/)
 

### Install
1. Clone this repository
```
git clone https://github.com/bergsantana/goledger-challenge-besu.git
```
2. Navigate to besu folder, then deploy the smart contract
```
cd goledger-challenge-besu/besu
./startDev.sh
```
3. Navigate to the api folder and start the database
Run 
```
docker compose up -d
```
4. 

# The Endpoints
After installing and running your API you should see the available endpoints on the console:
<img width="811" height="237" alt="image" src="https://github.com/user-attachments/assets/f1d31c61-f3cd-4d9c-ab56-81f150ba149d" />


## [GET] /set?value=123
Sets the value for the smart contract deployed on the Besu network

 
## [GET] /get
Retrieve the current value of the smart contract variable from the blockchain.

## [GET] /sync
Synchronize the value of the smart contract variable from the blockchain to the SQL database.

## [GET] /check
- Compare the value stored in the database with the current value of the smart contract variable.
- Return `true` if they are the same, otherwise return `false`.

 
 
## Postman Demo
![ezgif-855738f1561370](https://github.com/user-attachments/assets/8855802c-52d0-4bce-ae51-409d74f1b71f)

