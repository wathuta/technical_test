orders_pb_gen:
	protoc ./protos/orders/orders.proto  --go_out=./protos_gen/orders --proto_path=.  --go-grpc_out=./protos_gen/orders
customers_pb_gen:
	protoc ./protos/orders/customers.proto  --go_out=./protos_gen/customers --proto_path=.  --go-grpc_out=./protos_gen/customers
products_pb_gen:
	protoc ./protos/orders/products.proto  --go_out=./protos_gen/products --proto_path=.  --go-grpc_out=./protos_gen/products
payment_pb_gen:
	protoc ./protos/payment/payment.proto  --go_out=./protos_gen/payment  --proto_path=. --go-grpc_out=./protos_gen/payment
start_order_service:
	cd orders && make run-api
start_payment_service:
	cd payment && make run-api