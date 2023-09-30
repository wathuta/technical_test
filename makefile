orders_pb_gen:
	protoc ./protos/orders/orders.proto  --go_out=./protos_gen/orders --proto_path=.  --go-grpc_out=./protos_gen/orders
customers_pb_gen:
	protoc ./protos/orders/customers.proto  --go_out=./protos_gen/customers --proto_path=.  --go-grpc_out=./protos_gen/customers
products_pb_gen:
	