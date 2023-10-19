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

  def get_chapter(id) do
    GenServer.call(__MODULE__, {:get_chapter, id})
  end


  #### SERVER METHODS ####

  def handle_call({:get_course, id}, _from, channel) do
    request = %Proto.Course.CourseId{id: id}
    resp = channel |> Proto.Course.CourseService.Stub.get_course(request, timeout: 2000)
    {:reply, resp, channel}
  end

  def handle_call({:get_chapter, id}, _from, channel) do
    request = %Proto.Course.ChapterId{id: id}
    resp = channel |> Proto.Course.CourseService.Stub.get_chapter(request, timeout: 2000)
    {:reply, resp, channel}
  end

end
