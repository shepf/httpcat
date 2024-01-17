# 通用变量定义
# 安装日志
INSTALL_LOG=/tmp/httpcat_install.log
SERVER_IP="::"
is_ipv4="true"
is_need_uninstall="false"

function check_ipv4_address()
{
    IP=$1
    VALID_CHECK=$(echo $IP|awk -F. '$1<=255&&$2<=255&&$3<=255&&$4<=255{print "yes"}')
    if echo $IP|grep -E "^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$">/dev/null; then
    	if [ ${VALID_CHECK:-no} == "yes" ]; then
    		return 1
    	else
    		return 0
    	fi
    else
    	return 0
    fi
}

#处理安装脚本入参
while getopts i:f option
do
    case "$option" in
        i)
            #判断IPv4地址的合法性
            SERVER_IP=$OPTARG
            if echo $SERVER_IP | grep ":" > /dev/null; then
              is_ipv4="false"
            else
              is_ipv4="true"
              check_ipv4_address $SERVER_IP
              if [ $? -eq 0 ]; then
                echo "Input IP address was not correct."
                exit 1
              fi
            fi
            echo "Input server ip address is $SERVER_IP";;
        f)
            is_need_uninstall="true"
            echo "Need uninstall old data.";;
        \?)
            echo "Usage: ./httpcat_SERVER_V1.0.0.run [-i 1.1.1.1] [-f]"
            echo "-f means auto uninstall old httpcat"
            echo "-i means server ip address"
            exit 1;;
    esac
done


function httpcat::unintall(){
    uninstall_log=/tmp/httpcat_uninstall.log
    echo "uninstall httpcat" > $uninstall_log

    # 停止服务
    systemctl stop httpcat >> $uninstall_log  2>&1
    systemctl kill httpcat >> $uninstall_log  2>&1
    systemctl disable httpcat  >> $uninstall_log  2>&1

    echo -e "uninstall httpcat success\n" | tee -a $uninstall_log

    #exit 0
}


function httpcat::intall(){
  # 执行命令并提取版本号
  version=$(./httpcat -v | grep 'Version:' | awk '{print $2}')

  # 检查版本号是否为空
  if [ -n "$version" ]; then
      # 根据版本号给出不同的界面提示
      echo "正在安装 $version 版本"
  else
      echo "命令执行(./httpcat -v )失败，安装失败"
  fi


  execute_cmd "mkdir -p /home/web/website/upload/"
  execute_cmd "mkdir -p /home/web/website/httpcat_web/"
  execute_cmd "mkdir -p /etc/httpdcat/"

  execute_cmd "cp httpcat /usr/local/bin/httpcat  -rf"
  execute_cmd "cp conf/svr.yml /etc/httpdcat/svr.yml  -rf"


  cp httpcat.service /etc/systemd/system/httpcat.service  -rf
  sudo systemctl daemon-reload
  sudo systemctl restart httpcat

  # 安装前端
  rm /home/web/website/dist -rf
  rm /home/web/website/dist.zip -rf
  cp dist.zip /home/web/website/ -rf
  cd /home/web/website/
  unzip dist.zip
  rm httpcat_web -rf
  mv dist httpcat_web



}

#执行命令
execute_cmd() {
  cmd=$1
  eval $cmd
  ret=$?
  if [ $ret -ne 0 ]; then
    echo "Execute $cmd failed."
    # 注意这里：一旦有任何一个组件安装失败,则直接退出安装过程。
    # 这种更加保守的方式可以避免部分安装可能带来的隐患或不可预测的问题。
    exit 1
  fi
}



httpcat::intall $1 $2 && echo -e "\033[16C[ \033[32;49;1m OK \033[39;49;0m ]" ||\
echo -e "\033[16C[ \033[31;49;1m False \033[39;49;0m ]"
exit 0






