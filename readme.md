## Quick start
### Payment service
#### Requirement
- A functioning postgres database with the credentials populated in the `.env.payment` file in the payment folder.
- A tunneling software e.g ngrok to be used to receive callbacks from daraja api
- MPESA daraja api credentials to be used to make requests to the api.

#### Steps
1. Create a database and replace the database credentials in the .env.payment file.
2. Create a publicly available endpoint using the tunneling software or use this command in any unix terminal `ssh -R 80:localhost:5002 nokey@localhost.run` to generate a random public endpoint.
3. Update the `CALLBACK_BASEURL` var in the .env.payment file in the payment directory to the base url provided by the tunneling software.
4. Run `make start_payment_service` in the terminal (in the payment directory the default port is `:5001` for the grpc endpoints and `:5002` for the REST endpoints)
    - The rest endpoint is used to receive callbacks from daraja api , it is mapped to the public endpoint.
5. Some functionality in this service communicate with the `orders service`. Ensure that the payment service is up and healthy to test all the functionality of this api

### Order service
#### Requirements
- A functioning postgres database with the credentials populated in the `.env.orders` file in the orders folder.

#### Steps
1. Create a database and replace the database credentials in the .env.orders file.
2. Run `make start_order_service` in the terminal (in the orders directory the default port is `:5000`)
3. Some functionality in this service communicate with the `payment service`. Ensure that the payment service is up and healthy to test all the functionality of this api
