# etcd

拉取 k8s.gcr.io 镜像

```shell
docker push kainonly/etcd:3.2.24
// or
docker pull ccr.ccs.tencentyun.com/kainonly/etcd:3.2.24
```

重命名镜像

```shell
docker tag kainonly/etcd:3.2.24 k8s.gcr.io/etcd:3.2.24
// or
docker tag ccr.ccs.tencentyun.com/kainonly/etcd:3.2.24 k8s.gcr.io/etcd:3.2.24
```

删除镜像

```shell
docker rmi kainonly/etcd:3.2.24
// or
docker rmi ccr.ccs.tencentyun.com/kainonly/etcd:3.2.24
```