#启动本地tracker
###
 # @Description: 
 # @Author: neozhang
 # @Date: 2022-05-15 00:10:40
 # @LastEditors: neozhang
 # @LastEditTime: 2022-05-15 00:10:42
### 
fdfs_trackerd ./conf/tracker.conf restart
#启动本地storage
fdfs_storaged ./conf/storage.conf restart