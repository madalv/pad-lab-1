defmodule Gateway.Router do
  require Logger
  use Plug.Router

  plug(Plug.Parsers,
    parsers: [:urlencoded, :multipart, :json],
    pass: ["text/*"],
    json_decoder: Jason
  )

  plug(:match)
  plug(:dispatch)

  get "/courses" do
    uid = conn |> extract_user_id()

    case Hammer.check_rate("#{uid}", 2000, 2) do
      {:deny, limit} ->
        Logger.info("DENY request for user #{uid}; limit: #{limit}")
        conn |> send_resp(429, "too many requests")

      {:allow, _count} ->
        {page, _} = extract_query_param("page", conn) |> Integer.parse()
        {limit, _} = extract_query_param("limit", conn) |> Integer.parse()
        reply = Course.Client.get_courses(page, limit)

        case reply do
          {:ok, recs} ->
            send_encoded_json(recs, conn)

          {:error, %GRPC.RPCError{status: 4, message: _msg}} ->
            conn
            |> put_resp_header("Connection", "close")
            |> send_resp(408, "request timeout")

          error ->
            Logger.error(inspect(error))
            send_resp(conn, 500, "failed to get courses list")
        end
    end
  end

  get "/chapters/:id" do
    # assume request has passed auth middleware and it has decoded this user id from a token
    uid = conn |> extract_user_id()

    # rate limit request
    case Hammer.check_rate("#{uid}", 2000, 2) do
      {:deny, limit} ->
        Logger.info("DENY request for user #{uid}; limit: #{limit}")
        conn |> send_resp(429, "too many requests")

      {:allow, _count} ->
        reply = Course.Client.get_chapter(id)

        # check reply
        case reply do
          {:ok, course} ->
            send_encoded_json(course, conn)

          {:error, %GRPC.RPCError{status: 4, message: _msg}} ->
            conn
            |> put_resp_header("Connection", "close")
            |> send_resp(408, "request timeout")

          error ->
            Logger.error(inspect(error))
            send_resp(conn, 500, "failed to get chapter")
        end
    end
  end

  post "/chapters" do
    # assume request has passed auth middleware and it has decoded this user id from a token
    uid = conn |> extract_user_id()

    case Hammer.check_rate("#{uid}", 2000, 2) do
      {:deny, limit} ->
        Logger.info("DENY request for user #{uid}; limit: #{limit}")
        conn |> send_resp(429, "too many requests")

      {:allow, _count} ->
        res =
          conn.params
          |> Nestru.decode_from_map(Chapter.Dto.CreateChapter)

        case res do
          {:ok, dto} ->
            reply = Course.Client.create_chapter(dto)

            case reply do
              {:ok, id} ->
                send_encoded_json(id, conn)

              {:error, %GRPC.RPCError{status: 4, message: _msg}} ->
                conn
                |> put_resp_header("Connection", "close")
                |> send_resp(408, "request timeout")

              error ->
                Logger.error(inspect(error))
                send_resp(conn, 500, "failed to create chapter")
            end

          {:error, %{message: msg}} ->
            conn
            |> put_resp_content_type("text/plain")
            |> send_resp(400, msg)
        end
    end
  end

  post "/courses/:id/enroll" do
    # assume request has passed auth middleware and it has decoded this user id from a token
    uid = conn |> extract_user_id()

    case Hammer.check_rate("#{uid}", 2000, 2) do
      {:deny, limit} ->
        Logger.info("DENY request for user #{uid}; limit: #{limit}")
        conn |> send_resp(429, "too many requests")

      {:allow, _count} ->
        reply = Course.Client.enroll_user(uid, id)

        case reply do
          {:ok, _} ->
            send_resp(conn, 200, "user enrolled")

          {:error, %GRPC.RPCError{status: 4, message: _msg}} ->
            conn
            |> put_resp_header("Connection", "close")
            |> send_resp(408, "request timeout")

          error ->
            Logger.error(inspect(error))
            send_resp(conn, 500, "failed to enroll user")
        end

      {:error, %{message: msg}} ->
        conn
        |> put_resp_content_type("text/plain")
        |> send_resp(400, msg)
    end
  end

  post "/courses" do
    # assume request has passed auth middleware and it has decoded this user id from a token
    uid = conn |> extract_user_id()

    case Hammer.check_rate("#{uid}", 2000, 2) do
      {:deny, limit} ->
        Logger.info("DENY request for user #{uid}; limit: #{limit}")
        conn |> send_resp(429, "too many requests")

      {:allow, _count} ->
        res =
          conn.params
          |> Nestru.decode_from_map(Course.Dto.CreateCourse)

        case res do
          {:ok, dto} ->
            reply = Course.Client.create_course(dto)

            case reply do
              {:ok, id} ->
                send_encoded_json(id, conn)

              {:error, %GRPC.RPCError{status: 4, message: _msg}} ->
                conn
                |> put_resp_header("Connection", "close")
                |> send_resp(408, "request timeout")

              error ->
                Logger.error(inspect(error))
                send_resp(conn, 500, "failed to create course")
            end

          {:error, %{message: msg}} ->
            conn
            |> put_resp_content_type("text/plain")
            |> send_resp(400, msg)
        end
    end
  end

  get "/courses/:id" do
    # assume request has passed auth middleware and it has decoded this user id from a token
    uid = conn |> extract_user_id()

    # rate limit request
    case Hammer.check_rate("#{uid}", 2000, 2) do
      {:deny, limit} ->
        Logger.info("DENY request for user #{uid}; limit: #{limit}")
        conn |> send_resp(429, "too many requests")

      {:allow, _count} ->
        # check cache first
        case Cache.get_value(id) do
          {:ok, nil} ->
            reply = Course.Client.get_course(id)

            # check reply
            case reply do
              {:ok, course} ->
                case Protobuf.JSON.encode(course) do
                  {:ok, json} ->
                    Cache.set_value(id, json)

                    conn
                    |> put_resp_content_type("application/json")
                    |> send_resp(200, json)
                end

              {:error, %GRPC.RPCError{status: 4, message: _msg}} ->
                conn
                |> put_resp_header("Connection", "close")
                |> send_resp(408, "request timeout")

              error ->
                Logger.error(inspect(error))
                send_resp(conn, 500, "failed to get course")
            end

          {:ok, json} ->
            conn
            |> put_resp_content_type("application/json")
            |> send_resp(200, json)
        end
    end
  end

  get "/courses/:id/recommendations" do
    uid = conn |> extract_user_id()

    case Hammer.check_rate("#{uid}", 2000, 2) do
      {:deny, limit} ->
        Logger.info("DENY request for user #{uid}; limit: #{limit}")
        conn |> send_resp(429, "too many requests")

      {:allow, _count} ->
        {recs_nr, _} = extract_query_param("recs_nr", conn) |> Integer.parse()

        # check cache first
        case Cache.get_value("#{recs_nr}_#{id}") do
          {:ok, nil} ->
            reply = Rec.Client.get_recs_for_course(id, recs_nr)

            case reply do
              {:ok, recs} ->
                case Protobuf.JSON.encode(recs) do
                  {:ok, json} ->
                    Cache.set_value("#{recs_nr}_#{id}", json)

                    conn
                    |> put_resp_content_type("application/json")
                    |> send_resp(200, json)

                  {:error, _} ->
                    send_resp(conn, 500, "failed to encode response to json")
                end

              {:error, %GRPC.RPCError{status: 4, message: _msg}} ->
                conn
                |> put_resp_header("Connection", "close")
                |> send_resp(408, "request timeout")

              error ->
                Logger.error(inspect(error))
                send_resp(conn, 500, "failed to get recs")
            end

          {:ok, json} ->
            conn
            |> put_resp_content_type("application/json")
            |> send_resp(200, json)
        end
    end
  end

  get "/users/:id/recommendations" do
    # extract recs nr from query
    {recs_nr, _} = extract_query_param("recs_nr", conn) |> Integer.parse()

    reply = Rec.Client.get_recs_for_user(id, recs_nr)

    case reply do
      {:ok, recs} ->
        send_encoded_json(recs, conn)

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

  ### PRIVATE METHODS ###

  defp extract_user_id(conn) do
    case Plug.Conn.get_req_header(conn, "user-id") do
      [] -> "bruh"
      uid -> Enum.at(uid, 0)
    end
  end

  defp send_encoded_json(pb, conn) do
    encoded = Protobuf.JSON.encode(pb)

    case encoded do
      {:ok, json} ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(200, json)

      {:error, _} ->
        send_resp(conn, 500, "failed to encode response to json")
    end
  end

  defp extract_query_param(param, conn) do
    case Map.get(conn.query_params, param) do
      nil ->
        conn
        |> put_resp_content_type("text/plain")
        |> send_resp(400, "missing #{param} query parameter")

      val ->
        val
    end
  end
end
