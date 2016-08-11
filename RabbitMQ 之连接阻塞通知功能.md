


Blocked Connection Notifications

It is sometimes desirable for clients to receive a notification when their connection gets blocked due to the broker running low on resources (memory or disk).

We have introduced an AMQP protocol extension in which the broker sends to the client a connection.blocked method when the connection gets blocked, and connection.unblocked when it is unblocked.

To receive these notifications, the client must present a capabilities table in its client-properties in which there is a key connection.blocked and a boolean value true. See the capabilities section for further details on this. Our supported clients indicate this capability by default and provide a way to register handlers for the connection.blocked and connection.unblocked methods.

# Using Blocked Connection Notifications with Java Client

lve

# Using Blocked Connection Notifications with .NET Client

lve

# When Notifications are Sent

A connection.blocked notification is sent to publishing connections the first time RabbitMQ is low on a resource. For example, when a RabbitMQ node detects that it is low on RAM, it sends connection.blocked to all connected publishing clients supporting this feature. If before the connections are unblocked the node also starts running low on disk space, another connection.blocked will not be sent.

A connection.unblocked is sent when all resource alarms have cleared and the connection is fully unblocked


----------

官网原文：[这里](http://www.rabbitmq.com/connection-blocked.html)

