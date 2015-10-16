#!/bin/sh
#create a new chain named SHADOWSOCKS
iptables -t nat -F
iptables -t nat -X SHADOWSOCKS
iptables -t mangle -F
iptables -t nat -N SHADOWSOCKS
# Ignore your shadowsocks server's addresses
# It's very IMPORTANT, just be careful.
iptables -t nat -A SHADOWSOCKS -m mark --mark 2 -j RETURN
# Ignore LANs IP address
# Anything else should be redirected to shadowsocks's local port
iptables -t nat -A SHADOWSOCKS -p tcp -j REDIRECT --to-ports 8024
# Apply the rules
iptables -t nat -A PREROUTING -p tcp -j SHADOWSOCKS

iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE

iptables -t mangle -A OUTPUT -m owner --gid-owner 8347 -j MARK --set-mark 2
iptables -t mangle -A OUTPUT -m owner --gid-owner 8347 -j CONNMARK --save-mark
