




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
