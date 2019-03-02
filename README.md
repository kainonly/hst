# etcd

拉取 etcd:3.2.24 镜像

```shell
docker push kainonly/etcd:3.2.24
```

重命名镜像

```shell
docker tag kainonly/etcd:3.2.24 k8s.gcr.io/etcd:3.2.24
```

删除镜像

```shell
docker rmi kainonly/etcd:3.2.24
```

设置 docker-compose

```shell
version: '3.7'
services:
  etcd:
    image: kainonly/etcd:3.2.24
    restart: always
    command: etcd --config-file=/etc/config.yml
    volumes:
      - ./etcd/config.yml:/etc/config.yml
      - /var/lib/etcd:/var/data
    ports:
      - 2379:2379
      - 2380:2380
```

- `/etc/config.yml` etcd 配置文件
- `/var/data` 存储目录
