#!/bin/bash

# NOTE: This actually works on Amazon AMI's on EC2
# 32 vCPUs: r4.8xlarge

# $ renice +10 <pid>

NO_OF_INSTANCES=$1
REDIS=redis
START_PORT=$2
END_PORT=$(($START_PORT + NO_OF_INSTANCES - 1))
CURRENT_USER=$3
REDIS_FULL_NAME=redis-4.0.11

if [[ $EUID -ne 0 ]]; then
  echo "This script must be run as root"
  exit 1
fi

if [[ $3 == "start" ]]; then
  for port in `seq $START_PORT $END_PORT`;
  do
    systemctl start $REDIS"_"$port
  done
	exit 0
fi

if [[ $3 == "stop" ]];then
  for port in `seq $START_PORT $END_PORT`;
  do
    systemctl stop $REDIS"_"$port
  done
	exit 0
fi

if [ -z $NO_OF_INSTANCES ] || [ -z $CURRENT_USER ];
then
  echo "Usage: redis-script.sh <number-of-instances> <fist-port> <user or start/stop>"
  exit 1
fi

# Create redis group if it does not exist
/bin/egrep  -iq "^$REDIS:" /etc/group
if [ $? -eq 0 ]; then
  echo "Group $REDIS exists in /etc/group"
else 
  echo "Group $REDIS does not exist in /etc/group. Creating it now"
  groupadd redis
fi


# Create redis user if it does not exist
/bin/egrep  -iq "^$REDIS:" /etc/passwd
if [ $? -eq 0 ]; then
  echo "User $REDIS exists in /etc/passwd"
else 
  echo "User $REDIS does not exist in /etc/passwd. Creating it now"
  useradd -M $REDIS -g $REDIS
fi

# Download and compile redis
if [ -f /usr/local/bin/redis-server ]; then
  echo "Redis is already installed"
else 
  echo "Redis is not installed. Installing it now"
  sudo yum groupinstall "Development Tools"
  cd $HOME
  wget http://download.redis.io/releases/$REDIS_FULL_NAME.tar.gz
  tar xzf $REDIS_FULL_NAME.tar.gz
  cd $REDIS_FULL_NAME
  make distclean
  make
  cp src/redis-server /usr/local/bin/
  cp src/redis-cli /usr/local/bin/
  chown redis:redis /usr/local/bin/redis-server
  chown $CURRENT_USER:redis /usr/local/bin/redis-cli
  mkdir -p /etc/redis
  mkdir -p /var/log/redis
  mkdir -p /var/lib/redis
  mkdir -p /var/run/redis
  chown redis:root /var/lib/redis
  chown redis:root /var/run/redis
  chown redis:root /var/log/redis

  # Create systemd service files
  echo "Creating template systemd service file"
  echo "#Redis service file" > $REDIS.service
  echo "[Unit]" >> $REDIS.service
  echo "Description=Redis In-Memory Data Store" >> $REDIS.service
  echo "After=network.target" >> $REDIS.service
  echo "PartOf=redis.target" >> $REDIS.service
  echo "" >> $REDIS.service
  echo "[Service]" >> $REDIS.service
  echo "User=redis" >> $REDIS.service
  echo "Group=redis" >> $REDIS.service
  echo "Type=simple" >> $REDIS.service
  echo "LimitNOFILE=65536" >> $REDIS.service
  echo "PIDFile=/var/run/redis/redis_7379.pid" >> $REDIS.service
  echo "ExecStart=\"/usr/local/bin/redis-server\" \"/etc/redis/redis_7379.conf\"" >> $REDIS.service
  echo "ExecStop=\"/usr/local/bin/redis-cli\" shutdown" >> $REDIS.service
  echo "" >> $REDIS.service
  echo "[Install]" >> $REDIS.service
  echo "WantedBy=multi-user.target redis.target" >> $REDIS.service
fi

# Create redis instances
for i in `seq $START_PORT $END_PORT`;
do
	echo "Creating redis instance $i"
	cp $HOME/$REDIS_FULL_NAME/redis.conf /etc/redis/redis_$i.conf
	sed -i "/^port/c\port $i" /etc/redis/redis_$i.conf
	sed -i "/^bind 127.0.0/c\bind 0.0.0.0" /etc/redis/redis_$i.conf
	sed -i "/^supervised no/c\supervised systemd" /etc/redis/redis_$i.conf
	sed -i "/^pidfile/c\pidfile /var/run/redis/redis_$i.pid" /etc/redis/redis_$i.conf
	sed -i "/^logfile/c\logfile /var/log/redis/redis_$i.log" /etc/redis/redis_$i.conf
	sed -i "/^save/c\#save" /etc/redis/redis_$i.conf
	sed -i "/^dir/c\dir /var/lib/redis/redis_$i" /etc/redis/redis_$i.conf
	sed -i "/^tcp-backlog/c\tcp-backlog 8192" /etc/redis/redis_$i.conf
	sed -i "/^# requirepass foobared/c\requirepass gr4phlezz__2018" /etc/redis/redis_$i.conf
	if [ ! -d /var/lib/redis/redis_$i ];then
	  mkdir /var/lib/redis/redis_$i
	  chown redis:root /var/lib/redis/redis_$i
 	fi

	if [ ! -f /etc/systemd/system/redis_$i.service ];then
	  cp $HOME/$REDIS_FULL_NAME/redis.service /etc/systemd/system/redis_$i.service
	  sed -i "s/7379/$i/g" /etc/systemd/system/redis_$i.service
	fi

	mkdir -p /etc/systemd/system/redis_$i.service.d
  echo "[Service]" | sudo tee /etc/systemd/system/redis_$i.service.d/cpu-affinity.conf
  affinity="$((i-START_PORT))"
  echo "CPUAffinity=$affinity" | sudo tee -a /etc/systemd/system/redis_$i.service.d/cpu-affinity.conf

  for cpunum in $(cat /sys/devices/system/cpu/cpu*/topology/thread_siblings_list | cut -s -d, -f2- | tr ',' '\n' | sort -un)
  do
    echo 0 > /sys/devices/system/cpu/cpu$cpunum/online
  done
done

# Create and apply OS config for redis
if grep -Fxq "#REDIS_SETTINGS" /etc/rc.local
then
  echo "OS settings for Redis are already applied"
  exit 0
else
  sed -i '$ a #REDIS_SETTINGS' /etc/rc.local
  sed -i '$ a echo never > /sys/kernel/mm/transparent_hugepage/enabled' /etc/rc.local
  sed -i '$ a echo 1 > /proc/sys/net/ipv4/tcp_tw_reuse' /etc/rc.local
  sed -i '$ a #REDIS_SETTINGS' /etc/sysctl.conf
  sed -i '$ a vm.overcommit_memory=1' /etc/sysctl.conf
  sed -i '$ a net.ipv4.tcp_max_syn_backlog=65536' /etc/sysctl.conf
  sed -i '$ a net.core.somaxconn=65536' /etc/sysctl.conf
  sed -i '$ a fs.file-max=200000' /etc/sysctl.conf
  sysctl --system
fi
