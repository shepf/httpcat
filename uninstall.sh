#!/bin/bash

uninstall_log=/tmp/httpcat_uninstall.log
echo "uninstall httpcat" > $uninstall_log

# 停止服务
systemctl stop httpcat >> $uninstall_log 2>&1
systemctl kill httpcat >> $uninstall_log 2>&1
systemctl disable httpcat >> $uninstall_log 2>&1

echo -e "uninstall httpcat success\n" | tee -a $uninstall_log

exit 0