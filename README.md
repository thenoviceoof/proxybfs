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

LICENSE
--------------------------------------------------------------------------------
MIT LICENSE
Copyright (c) <2012> <thenoviceoof>

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
