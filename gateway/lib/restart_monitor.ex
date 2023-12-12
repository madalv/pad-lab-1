defmodule RestartMonitor do
  use GenServer
  require Logger

  def start_link(init_arg) do
    GenServer.start_link(__MODULE__, init_arg, name: __MODULE__)
  end

  def init(timeout) do
    :timer.send_after(timeout, self(), :check_restarts)

    {:ok, conn} = K8s.Conn.from_service_account(insecure_skip_tls_verify: true)
    Logger.info("Successfully Init k8s Client")

    # get pod list from kubernetes api
    operation = K8s.Client.list("v1", "Pod", namespace: "default")
    {:ok, deployments} = K8s.Client.run(conn, operation)

    # compute pod restart map
    restarts_per_pod = Enum.reduce(deployments["items"], %{}, fn pod, acc ->
      count = Enum.at(pod["status"]["containerStatuses"], 0) |> Map.get("restartCount", 0)
      Map.put(acc, pod["metadata"]["name"], count)
    end)

    {:ok,
     %{
       conn: conn,
       restarts_per_pod: restarts_per_pod,
       timeout: timeout,
     }}
  end

  def handle_info(:check_restarts, state) do
    operation = K8s.Client.list("v1", "Pod", namespace: "default")
    {:ok, deployments} = K8s.Client.run(state[:conn], operation)

    new_cnt =
    Enum.reduce(deployments["items"], state[:restarts_per_pod], fn pod, acc ->
      name = pod["metadata"]["name"]
      count = Enum.at(pod["status"]["containerStatuses"], 0) |> Map.get("restartCount", 0)
      old_count = Map.get(state[:restarts_per_pod], name)


    if count > old_count do
      Logger.warning("Pod #{name} has restarted #{count - old_count} times in the last #{state[:timeout]} ms.")
      Map.put(acc, name, count)
      else
        acc
      end
    end)

    :timer.send_after(state[:timeout], self(), :check_restarts)
    {:noreply, %{state | restarts_per_pod: new_cnt}}
  end

end
