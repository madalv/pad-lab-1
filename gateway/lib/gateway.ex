defmodule Gateway do
  use Application
  require Logger

  def start(_type, _args) do
    Logger.info("Starting Gateway...")

    # TODO get addresses out of config
    {:ok, course_channel} = GRPC.Stub.connect("localhost:50052")
    Logger.info("Successfully connected to Course Svc")

    {:ok, rec_channel} = GRPC.Stub.connect("localhost:50051")
    Logger.info("Successfully connected to Rec Svc")

    # TODO consider using poolboy
    children = [
      {Plug.Cowboy, scheme: :http, plug: Gateway.Router, options: [port: 8080]},
      %{
        id: Course.Client,
        start: {Course.Client, :start_link, [course_channel]}
        },
      %{
        id: Rec.Client,
        start: {Rec.Client, :start_link, [rec_channel]}
      }
    ]

    opts = [strategy: :one_for_one, name: Gateway.Supervisor]
    Logger.info("Started Router...")
    Supervisor.start_link(children, opts)
  end
end
