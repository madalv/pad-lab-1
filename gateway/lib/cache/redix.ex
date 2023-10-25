defmodule Cache do
  use GenServer
  require Logger

  def start_link(init_arg) do
    GenServer.start_link(__MODULE__, init_arg, name: __MODULE__)
  end

  def init(conn) do
    {:ok, conn}
  end

  def set_value(key, val) do
    GenServer.call(__MODULE__, {:set, key, val})
  end

  def get_value(key) do
    GenServer.call(__MODULE__, {:get, key})
  end

  def handle_call({:get, key}, _from, conn) do
    repl = Redix.command(:redis, ["GET", key])
    Logger.debug(inspect(repl))
    {:reply, repl, conn}
  end

  def handle_call({:set, key, val}, _from, conn) do
    repl = Redix.command(:redis, ["SETEX", key, 60, val])
    # Logger.debug(inspect(repl))
    {:reply, repl, conn}
  end
end
