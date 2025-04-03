package main

const configTemplate = `
app:
  name: "%s"

config_center:
  nacos:
    host: 127.0.0.1
    port: 8848
    namespace: b652ebeb-fc2c-4f5f-8398-316914b0be27
    data_id: goboot
    group: DEFAULT_GROUP
    log_dir: ./logs/nacos
    cache_dir: ./cache/nacos

http:
  port: 8080
  addr: 0.0.0.0
  gin_mode: release
`
