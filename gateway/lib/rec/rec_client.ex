defmodule Rec.Client do
  use GenServer
  require Logger

  def start_link(channel) do
    GenServer.start_link(__MODULE__, channel, name: __MODULE__)
  end

  def init(channel) do
    {:ok, channel}
  end

  #### CLIENT FACING METHODS ####

  def get_recs_for_user(user_id, recs_nr) do
    GenServer.call(__MODULE__, {:get_recs_for_user, user_id, recs_nr})
  end

  def get_recs_for_course(course_id, recs_nr) do
    GenServer.call(__MODULE__, {:get_recs_for_course, course_id, recs_nr})
  end

  #### GENSERVER METHODS ####

  def handle_call({:get_recs_for_user, user_id, recs_nr}, _from, channel) do
    request = %Proto.Rec.UserRecsRequest{user_id: user_id, recs_nr: recs_nr}
    resp = channel |> Proto.Rec.RecService.Stub.get_recs_for_user(request, timeout: 2000)
    {:reply, resp, channel}
  end

  def handle_call({:get_recs_for_course, course_id, recs_nr}, _from, channel) do
    request = %Proto.Rec.CourseRecsRequest{course_id: course_id, recs_nr: recs_nr}
    resp = channel |> Proto.Rec.RecService.Stub.get_recs_for_course(request, timeout: 2000)
    {:reply, resp, channel}
  end
end
