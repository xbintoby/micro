# micro

# ETCD
nohup etcd --listen-peer-urls='http://192.168.3.117:2379' > /dev/null &


docker run -d \
  -p 2379:2379 \
  --name etcd \
  quay.io/coreos/etcd:v3.5.0 \
  etcd \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster-state new
