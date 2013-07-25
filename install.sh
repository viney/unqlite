#! /bin/sh

# 下载源码
wget -c http://unqlite.org/db/unqlite-db-116.zip

# 解压源码
unzip unqlite-db-116.zip

# 编译源码
gcc -Wall -fPIC -c *.c
gcc -shared -Wl,-soname,libunqlite.so.1 -o libunqlite.so.1.0 *.o

# 建立软链接
sudo cp `pwd`/libunqlite.so.1.0 /usr/local/lib/
sudo cp `pwd`/unqlite.h /usr/local/include/
sudo ln -sf /usr/local/lib/libunqlite.so.1.0 /usr/local/lib/libunqlite.so.1
sudo ln -sf /usr/local/lib/libunqlite.so.1 /usr/local/lib/libunqlite.so

# 建立共享
ldconfig /usr/local/lib/libunqlite.so

# 下载golang unqlite驱动包
git clone git@github.com:ceh/gounqlite.git
