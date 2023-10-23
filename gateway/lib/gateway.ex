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
      _ ->
        Logger.error("Unkown env!")
    end

    # TODO get addresses out of config
    {:ok, course_channel} = GRPC.Stub.connect(env!("COURSE_SVC_ADDRESS"))
    Logger.info("Successfully connected to Course Svc")

    {:ok, rec_channel} = GRPC.Stub.connect(env!("REC_SVC_ADDRESS"))
    Logger.info("Successfully connected to Rec Svc")

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
