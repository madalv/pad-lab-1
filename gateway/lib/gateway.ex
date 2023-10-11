defmodule Gateway do
  use Application
  require Logger

  def start(_type, _args) do
    Logger.info("Starting Gateway...")
    children = [
      {Plug.Cowboy, scheme: :http, plug: Gateway.Router, options: [port: 8080]}
    ]

    opts = [strategy: :one_for_one, name: Gateway.Supervisor]
    Logger.info("Started Router...")
    Supervisor.start_link(children, opts)
  end
end
