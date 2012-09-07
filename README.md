proxybfs
================================================================================
A TCP proxy utility

USAGE
--------------------------------------------------------------------------------
A straightforward TCP proxy from localhost's 8080 port, to google:
making HTTP requests to 127.0.0.1:8080 will return google content
(this is probably naive and overoptismistic, given HTTP/1.0's HOST)

    ./proxybfs -l 8080 -c google.com

A listening proxy connects two active connections: for instance,
buddies want to swap files from behind their respective NATs

    ./proxybfs -l 8080 -l 5050

A connecting proxy actively tries to connect to two clients: here, we
send all telnet traffic to sshd, and vice-versa

    ./proxybfs -c 22 -c 21

DESIGN PROCESS
--------------------------------------------------------------------------------
This was originally meant to deliver connections to a mobile vm,
allowing people to help out with debuging easily.

TODO
--------------------------------------------------------------------------------
- multiple listens -> 1 connect, funnel
- 1 listen -> n connects, "load balance"
- n listen -> n connects, some godawful scheme
