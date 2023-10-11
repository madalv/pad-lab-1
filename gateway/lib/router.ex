defmodule Gateway.Router do
  require Logger
  use Plug.Router

  plug(Plug.Parsers,
  parsers: [:urlencoded, :json],
  pass: ["text/*"],
  json_decoder: Jason
)

  plug(:match)
  plug(:dispatch)

  get "/users/:id/recommendations" do
    send_resp(conn, 200, "ok man")
  end

  get "/courses/:id/recommendations" do
    send_resp(conn, 200, "ok man")
  end

  get "/" do
    send_resp(conn, 200, "STATUS: SERVING")
  end

  get _ do
    Logger.debug("Required path /#{conn.path_info} not found!")
    send_resp(conn, 404, "Oops... Nothing here...")
  end
end
