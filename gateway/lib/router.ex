defmodule Gateway.Router do
  require Logger
  use Plug.Router

  plug Plug.Parsers,
    parsers: [:urlencoded, :multipart, :json],
    pass: ["text/*"],
    json_decoder: Jason
  plug(:match)
  plug(:dispatch)

  # TODO add separate handler for course
  # TODO validate dtos before sending to services
  # TODO add rate limiter

  post "/courses" do
    
  end

  get "/courses/:id" do
    reply = Course.Client.get_course(id)

    case reply do
      {:error, error} ->
        send_resp(conn, 500, error)

      {:ok, course} ->
        encoded = Protobuf.JSON.encode(course)

        case encoded do
          {:ok, json} ->
            conn
            |> put_resp_content_type("application/json")
            |> send_resp(200, json)
          {:error, _} ->
            send_resp(conn, 500, "failed to encode response to json")
        end
    end
  end

  get "/courses/:id/recommendations" do
    # extract recs nr from query
    recs_nr = case Map.get(conn.query_params, "recs_nr") do
      nil ->
        conn
        |> put_resp_content_type("text/plain")
        |> send_resp(400, "missing recs_nr query parameter")
      nr ->
        {val, _} = Integer.parse(nr)
        val
      end

    reply = Rec.Client.get_recs_for_course(id, recs_nr)

    case reply do
      {:ok, recs} ->
        encoded = Protobuf.JSON.encode(recs)

        case encoded do
          {:ok, json} ->
            conn
            |> put_resp_content_type("application/json")
            |> send_resp(200, json)
          {:error, _} ->
            send_resp(conn, 500, "failed to encode response to json")
        end

       error ->
        Logger.error(inspect(error))
        send_resp(conn, 500, "failed to get recs")
    end
  end

  get "/users/:id/recommendations" do
    # extract recs nr from query
    recs_nr = case Map.get(conn.query_params, "recs_nr") do
      nil ->
        conn
        |> put_resp_content_type("text/plain")
        |> send_resp(400, "missing recs_nr query parameter")
      nr ->
        {val, _} = Integer.parse(nr)
        val
      end


    reply = Rec.Client.get_recs_for_user(id, recs_nr)

    case reply do
      {:ok, recs} ->
        encoded = Protobuf.JSON.encode(recs)

        case encoded do
          {:ok, json} ->
            conn
            |> put_resp_content_type("application/json")
            |> send_resp(200, json)
          {:error, _} ->
            send_resp(conn, 500, "failed to encode response to json")
        end
      {:error, error} ->
        encoded = Protobuf.JSON.encode(error)
        send_resp(conn, 500, encoded)
    end

  end


  get "/status" do
    send_resp(conn, 200, "STATUS: SERVING")
  end

  get _ do
    Logger.debug("Required path /#{conn.path_info} not found!")
    send_resp(conn, 404, "Oops... Nothing here...")
  end
end
