defmodule Course.Client do
  use GenServer
  require Logger

  def start_link(init_arg) do
    GenServer.start_link(__MODULE__, init_arg, name: __MODULE__)
  end

  def init({channel, timeout}) do
    :timer.send_after(timeout * 3, self(), :check_errors)

    {:ok,
     %{
       channel: channel,
       timeout: timeout,
       error_threshold: 3,
       error_cnt: 0,
       is_closed: true,
       recovery_delay: 10_000
     }}
  end

  #### CLIENT FACING METHODS ####

  def get_course(id) do
    GenServer.call(__MODULE__, {:get_course, id})
  end

  def get_chapter(id) do
    GenServer.call(__MODULE__, {:get_chapter, id})
  end

  def create_course(create_course_dto) do
    GenServer.call(__MODULE__, {:create_course, create_course_dto})
  end

  def create_chapter(create_chapter_dto) do
    GenServer.call(__MODULE__, {:create_chapter, create_chapter_dto})
  end

  def get_courses(page, limit) do
    GenServer.call(__MODULE__, {:get_courses, page, limit})
  end

  def enroll_user(user_id, course_id) do
    GenServer.call(__MODULE__, {:enroll_user, user_id, course_id})
  end

  #### SERVER METHODS ####

  def handle_call({:enroll_user, user_id, course_id}, _from, state) do
    if state[:is_closed] do
      channel = state[:channel]

      request = %Proto.Course.EnrollRequest{
        user_id: user_id,
        course_id: course_id
      }

      resp =
        channel
        |> Proto.Course.CourseService.Stub.enroll_user(request, timeout: state[:timeout])

      cnt = state[:error_cnt]

      case resp do
        {:error, _} ->
          {:reply, resp, %{state | error_cnt: cnt + 1}}

        _ ->
          {:reply, resp, state}
      end
    else
      resp = {:error, "circuit breaker for course svc open, wait #{state[:recovery_delay]} ms"}
      {:reply, resp, state}
    end
  end

  def handle_call({:create_chapter, dto}, _from, state) do
    if state[:is_closed] do
      channel = state[:channel]

      request = %Proto.Course.CreateChapterRequest{
        course_id: dto.course_id,
        title: dto.title,
        body: dto.body
      }

      resp =
        channel
        |> Proto.Course.CourseService.Stub.create_chapter(request, timeout: state[:timeout])

      cnt = state[:error_cnt]

      case resp do
        {:error, _} ->
          {:reply, resp, %{state | error_cnt: cnt + 1}}

        _ ->
          {:reply, resp, state}
      end
    else
      resp = {:error, "circuit breaker for course svc open, wait #{state[:recovery_delay]} ms"}
      {:reply, resp, state}
    end
  end

  def handle_call({:create_course, dto}, _from, state) do
    if state[:is_closed] do
      channel = state[:channel]

      request = %Proto.Course.CreateCourseRequest{
        author_id: dto.author_id,
        title: dto.title,
        description: dto.description,
        category_ids: dto.category_ids
      }

      resp =
        channel
        |> Proto.Course.CourseService.Stub.create_course(request, timeout: state[:timeout])

      cnt = state[:error_cnt]

      case resp do
        {:error, _} ->
          {:reply, resp, %{state | error_cnt: cnt + 1}}

        _ ->
          {:reply, resp, state}
      end
    else
      resp = {:error, "circuit breaker for course svc open, wait #{state[:recovery_delay]} ms"}
      {:reply, resp, state}
    end
  end

  def handle_call({:get_courses, page, limit}, _from, state) do
    if state[:is_closed] do
      channel = state[:channel]

      request = %Proto.Course.PaginationQuery{
        page: page,
        limit: limit
      }

      resp =
        channel |> Proto.Course.CourseService.Stub.get_courses(request, timeout: state[:timeout])

      cnt = state[:error_cnt]

      case resp do
        {:error, _} ->
          {:reply, resp, %{state | error_cnt: cnt + 1}}

        _ ->
          {:reply, resp, state}
      end
    else
      resp = {:error, "circuit breaker for course svc open, wait #{state[:recovery_delay]} ms"}
      {:reply, resp, state}
    end
  end

  def handle_call({:get_course, id}, _from, state) do
    if state[:is_closed] do
      channel = state[:channel]
      request = %Proto.Course.CourseId{id: id}

      resp =
        channel |> Proto.Course.CourseService.Stub.get_course(request, timeout: state[:timeout])

      cnt = state[:error_cnt]

      case resp do
        {:error, _} ->
          {:reply, resp, %{state | error_cnt: cnt + 1}}

        _ ->
          {:reply, resp, state}
      end
    else
      resp = {:error, "circuit breaker for course svc open, wait #{state[:recovery_delay]} ms"}
      {:reply, resp, state}
    end
  end

  def handle_call({:get_chapter, id}, _from, state) do
    if state[:is_closed] do
      channel = state[:channel]
      request = %Proto.Course.ChapterId{id: id}

      resp =
        channel |> Proto.Course.CourseService.Stub.get_chapter(request, timeout: state[:timeout])

      cnt = state[:error_cnt]

      case resp do
        {:error, _} ->
          {:reply, resp, %{state | error_cnt: cnt + 1}}

        _ ->
          {:reply, resp, state}
      end
    else
      resp = {:error, "circuit breaker for course svc open, wait #{state[:recovery_delay]} ms"}
      {:reply, resp, state}
    end
  end

  #### CIRCUIT BREAKER METHODS ####

  def handle_info(:check_errors, state) do
    if state[:error_cnt] > state[:error_threshold] do
      Logger.warning("Error threshold reached, circuit opened for course client")
      :timer.send_after(state[:recovery_delay], self(), :close_circuit)
      {:noreply, %{state | is_closed: false}}
    else
      :timer.send_after(state[:timeout] * 3, self(), :check_errors)
      {:noreply, state}
    end
  end

  def handle_info(:close_circuit, state) do
    Logger.info("Circuit closed for course client")
    :timer.send_after(state[:timeout] * 3, self(), :check_errors)
    {:noreply, %{state | is_closed: true, error_cnt: 0}}
  end
end
