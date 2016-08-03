




# 创建用户

```shell
rabbitmqctl add_user moooofly moooofly
```

# 设置用户角色

```shell
rabbitmqctl set_user_tags moooofly administrator
```

# 设置用户权限

```shell
rabbitmqctl set_permissions -p / moooofly ".\*" ".\*" ".\*"
```

# 单机集群构建

```shell
RABBITMQ_NODE_PORT=5672 RABBITMQ_NODENAME=rabbit_1 RABBITMQ_SERVER_START_ARGS="-rabbitmq_management listener [{port,15672}]" rabbitmq-server -detached

RABBITMQ_NODE_PORT=5673 RABBITMQ_NODENAME=rabbit_2 RABBITMQ_SERVER_START_ARGS="-rabbitmq_management listener [{port,15673}]" rabbitmq-server -detached

RABBITMQ_NODE_PORT=5674 RABBITMQ_NODENAME=rabbit_3 RABBITMQ_SERVER_START_ARGS="-rabbitmq_management listener [{port,15674}]" rabbitmq-server -detached

rabbitmqctl -n rabbit_2 stop_app
rabbitmqctl -n rabbit_2 join_cluster rabbit_1@`hostname -s`
rabbitmqctl -n rabbit_2 start_app

rabbitmqctl -n rabbit_3 stop_app
rabbitmqctl -n rabbit_3 join_cluster rabbit_1@`hostname -s`
rabbitmqctl -n rabbit_3 start_app
```

