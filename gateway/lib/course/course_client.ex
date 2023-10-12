defmodule Course.Client do
  use GenServer
  require Logger

  def start_link(channel) do
    GenServer.start_link(__MODULE__, channel, name: __MODULE__)
  end

  def init(channel) do
    {:ok, channel}
  end

  #### CLIENT FACING METHODS ####

  def get_course(id) do
    GenServer.call(__MODULE__, {:get_course, id})
  end


  #### SERVER METHODS ####

  def handle_call({:get_course, id}, _from, channel) do
    request = %Proto.Course.CourseId{id: id}
    resp = channel |> Proto.Course.CourseService.Stub.get_course(request)
    {:reply, resp, channel}
  end

end
