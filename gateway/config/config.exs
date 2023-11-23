use Mix.Config

config :hammer,
  backend: {Hammer.Backend.ETS, [expiry_ms: 60_000 * 60 * 4, cleanup_interval_ms: 60_000 * 10]}


config :gateway, RedisCache,
  mode: :redis_cluster,
  redis_cluster: [
    configuration_endpoints: [
      endpoint1_conn_opts: [
        host: "redis-cluster",
        port: 6379,
      ]
    ]
]
