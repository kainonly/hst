# etcd

历史版本

- 3.4.0(unrelease)
- 3.3.15
- 3.2.34


拉取 etcd:tag 镜像

```shell
docker pull kainonly/etcd:tag
```

重命名镜像

```shell
docker tag kainonly/etcd:tag k8s.gcr.io/etcd:tag
```

删除镜像

```shell
docker rmi kainonly/etcd:tag
```

设置 docker-compose

```shell
version: '3.7'
services:
  etcd:
    image: kainonly/etcd:tag
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
