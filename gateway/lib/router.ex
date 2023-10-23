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


  # TODO create chapter

  # TODO create course
  post "/courses" do

  end

  # TODO get courses, paginated

  # TODO get chapter


  get "/courses/:id" do
    # assume request has passed auth middleware and it has decoded this user id from a token
    uid = Plug.Conn.get_req_header(conn, "user-id")

    reply = Course.Client.get_course(id)

    case reply do
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

      {:error, %GRPC.RPCError{status: 4, message: _msg}} ->
        conn
        |> put_resp_header("Connection", "close")
        |> send_resp(408, "Request timeout")

      error ->
        Logger.error(inspect(error))
        send_resp(conn, 500, "failed to get recs")
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

      {:error, %GRPC.RPCError{status: 4, message: _msg}} ->
          conn
          |> put_resp_header("Connection", "close")
          |> send_resp(408, "Request timeout")

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
          {:error, %GRPC.RPCError{message: "timeout when waiting for server", status: _}} ->
           conn
            |> put_resp_header("Connection", "close")
            |> send_resp(408, "Request timeout")
          {:error, _} ->
            send_resp(conn, 500, "failed to encode response to json")
        end

      {:error, %GRPC.RPCError{status: 4, message: _msg}} ->
        conn
        |> put_resp_header("Connection", "close")
        |> send_resp(408, "Request timeout")

      error ->
        Logger.error(inspect(error))
        send_resp(conn, 500, "failed to get recs")
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
