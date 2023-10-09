# import grpc
# from grpc import RpcError
# import your_grpc_pb2  # Import your generated gRPC protobuf file
# import your_grpc_pb2_grpc  # Import the gRPC stubs

# # Create a gRPC channel and stub
# channel = grpc.insecure_channel('localhost:50051')  # Replace with your server address
# stub = your_grpc_pb2_grpc.YourServiceStub(channel)

# # Set the timeout in seconds (e.g., 5 seconds)
# timeout_seconds = 5

# try:
#     # Make a gRPC request with the specified timeout
#     response = stub.YourRpcMethod(your_request, timeout=timeout_seconds)
#     print("Response:", response)
# except RpcError as e:
#     # Handle gRPC errors, including timeout
#     if e.code() == grpc.StatusCode.DEADLINE_EXCEEDED:
#         print("Request timed out after {} seconds.".format(timeout_seconds))
#     else:
#         print("gRPC error:", e.details())

# # Don't forget to close the channel when done
# channel.close()