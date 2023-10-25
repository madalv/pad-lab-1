defmodule Rec.Client do
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

  def get_recs_for_user(user_id, recs_nr) do
    GenServer.call(__MODULE__, {:get_recs_for_user, user_id, recs_nr})
  end

  def get_recs_for_course(course_id, recs_nr) do
    GenServer.call(__MODULE__, {:get_recs_for_course, course_id, recs_nr})
  end

  #### GENSERVER METHODS ####

  def handle_call({:get_recs_for_user, user_id, recs_nr}, _from, state) do
    if state[:is_closed] do
      channel = state[:channel]
      request = %Proto.Rec.UserRecsRequest{user_id: user_id, recs_nr: recs_nr}

      resp =
        channel |> Proto.Rec.RecService.Stub.get_recs_for_user(request, timeout: state[:timeout])

      cnt = state[:error_cnt]

      case resp do
        {:error, _} ->
          {:reply, resp, %{state | error_cnt: cnt + 1}}

        _ ->
          {:reply, resp, state}
      end
    else
      resp = {:error, "circuit breaker for rec svc open, wait #{state[:recovery_delay]} ms"}
      {:reply, resp, state}
    end
  end

  def handle_call({:get_recs_for_course, course_id, recs_nr}, _from, state) do
    if state[:is_closed] do
      channel = state[:channel]
      request = %Proto.Rec.CourseRecsRequest{course_id: course_id, recs_nr: recs_nr}

      resp =
        channel
        |> Proto.Rec.RecService.Stub.get_recs_for_course(request, timeout: state[:timeout])

      cnt = state[:error_cnt]

      case resp do
        {:error, _} ->
          {:reply, resp, %{state | error_cnt: cnt + 1}}

        _ ->
          {:reply, resp, state}
      end
    else
      resp = {:error, "circuit breaker for rec svc open, wait #{state[:recovery_delay]} ms"}
      {:reply, resp, state}
    end
  end

  #### CIRCUIT BREAKER METHODS ####

  def handle_info(:check_errors, state) do
    if state[:error_cnt] > state[:error_threshold] do
      Logger.warning("Error threshold reached, circuit opened for rec client")
      :timer.send_after(state[:recovery_delay], self(), :close_circuit)
      {:noreply, %{state | is_closed: false}}
    else
      :timer.send_after(state[:timeout] * 3, self(), :check_errors)
      {:noreply, state}
    end
  end

  def handle_info(:close_circuit, state) do
    Logger.info("Circuit closed for rec client")
    :timer.send_after(state[:timeout] * 3, self(), :check_errors)
    {:noreply, %{state | is_closed: true, error_cnt: 0}}
  end
end
