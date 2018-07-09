# go-manga

Build & Upload

```
$ ./deploy_bin.sh
server                100%   14MB   1.3MB/s   00:11    
Deployed server
update_books          100%   13MB   1.4MB/s   00:09    
Deployed update_books
mailfile              100% 5389KB   1.1MB/s   00:05    
Deployed mailfile
```

Restart server process

```
$ ssh server
$ cd go/src/github.com/akiraak/go-manga/
$ ./newreplace.sh
$ ps aux | grep server
akiraak  21219  0.0  2.7  37876 16840 ?        Sl   Apr06  47:44 ./server
$ kill 21219
```

Check new server process

```
$ ps aux | grep server
akiraak   5915  0.0  1.3  12320  7980 ?        Sl   17:24   0:00 ./server
```
