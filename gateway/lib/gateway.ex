defmodule Gateway do
  use Application
  require Logger
  import Dotenvy

  def start(_type, _args) do
    Logger.info("Starting Gateway...")

    case System.get_env("MIX_ENV") do
      "prod" ->
        source(".prod.env")
        Logger.info(env!("REC_SVC_ADDRESS"))
        Logger.info(env!("COURSE_SVC_ADDRESS"))

      "dev" ->
        source(".local.env")
        Logger.info(env!("REC_SVC_ADDRESS"))
        Logger.info(env!("COURSE_SVC_ADDRESS"))
        Logger.info(env!("REDIS_URL"))

      _ ->
        Logger.error("Unkown env!")
    end

    {:ok, course_channel} = GRPC.Stub.connect(env!("COURSE_SVC_ADDRESS"))
    Logger.info("Successfully connected to Course Svc")

    {:ok, rec_channel} = GRPC.Stub.connect(env!("REC_SVC_ADDRESS"))
    Logger.info("Successfully connected to Rec Svc")

    {:ok, redis_conn} = Redix.start_link(env!("REDIS_URL"), name: :redis)
    Logger.info("Successfully connected to Redis")

    Redix.command(:redis, ["SET", 50, "sviernerewvre_d", "ddfdf"])

    children = [
      {Plug.Cowboy, scheme: :http, plug: Gateway.Router, options: [port: 8080]},
      %{
        id: Rec.Client,
        start: {Rec.Client, :start_link, [{rec_channel, 2000}]}
      },
      %{
        id: Course.Client,
        start: {Course.Client, :start_link, [{course_channel, 2000}]}
      },
      %{
        id: Cache,
        start: {Cache, :start_link, [redis_conn]}
      }
    ]

    opts = [strategy: :one_for_one, name: Gateway.Supervisor]
    Logger.info("Started Router...")
    Supervisor.start_link(children, opts)
  end
end
