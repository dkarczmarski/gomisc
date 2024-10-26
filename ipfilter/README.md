# ipfilter

ipfilter is a simple tool to manage access from selected IP addresses. 
It should not be considered a security feature for production environments. 
The intention is to easily open the firewall to selected hosts on development 
environments, where we often don't want to expose them to the public internet.

## proof of concept

run on:
- OS: Ubuntu 23.10
- user: myuser

```
> cat /etc/sudoers.d/myuser 
myuser  ALL=(ALL) NOPASSWD:/usr/sbin/ufw

```