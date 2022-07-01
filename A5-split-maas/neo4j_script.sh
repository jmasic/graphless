#!/bin/bash


# https://us-east-2.console.aws.amazon.com/cloudwatch/home?region=us-east-2#logsV2:log-groups/log-group/$252Faws$252Flambda$252FWorkerFunction/log-events/2022$252F05$252F24$252F$255B$2524LATEST$255D1f89d2a80f634a539704b6a6ab7d90a5

# 32 vCPUs: c6i.8xlarge // r5n.8xlarge
# 16 vCPUs: c5n.4xlarge
# NOTE: open ports 7687 and optionally 7474

sudo apt update

sudo apt install -y apt-transport-https ca-certificates curl software-properties-common net-tools \
                    openjdk-11-jre-headless

sudo wget -O - https://debian.neo4j.com/neotechnology.gpg.key | sudo apt-key add -
echo 'deb https://debian.neo4j.com stable 4.4' | sudo tee -a /etc/apt/sources.list.d/neo4j.list
sudo apt-get update -y
sudo apt-get install -y neo4j=1:4.4.6

# grant access to anyone
sudo sed -i "/^#dbms.default_listen_address=0.0.0.0/c\dbms.default_listen_address=0.0.0.0" /etc/neo4j/neo4j.conf
sudo sed -i "/^#dbms.security.auth_enabled=false/c\dbms.security.auth_enabled=false" /etc/neo4j/neo4j.conf
# disk IO tuning
sudo sed -i "/^dbms.tx_log.rotation.retention_policy=1 days/c\dbms.tx_log.rotation.retention_policy=keep_none" /etc/neo4j/neo4j.conf
echo "dbms.checkpoint.interval.time=1m" | sudo tee -a /etc/neo4j/neo4j.conf
echo "dbms.checkpoint.interval.tx=100000" | sudo tee -a /etc/neo4j/neo4j.conf
# network threadpool tuning
echo "dbms.connector.bolt.thread_pool_max_size=1000" | sudo tee -a /etc/neo4j/neo4j.conf
# memory sizing
sudo sed -i "/^#dbms.memory.heap.initial_size=512/c\dbms.memory.heap.initial_size=128g" /etc/neo4j/neo4j.conf
sudo sed -i "/^#dbms.memory.heap.max_size=512/c\dbms.memory.heap.max_size=128g" /etc/neo4j/neo4j.conf
sudo sed -i "/^#dbms.memory.pagecache.size=10g/c\dbms.memory.pagecache.size=64g" /etc/neo4j/neo4j.conf

sudo mkdir -p /etc/systemd/system/neo4j.service.d/
sudo echo "[Service]" | sudo tee /etc/systemd/system/neo4j.service.d/override.conf
sudo echo "LimitNOFILE=300000" | sudo tee -a /etc/systemd/system/neo4j.service.d/override.conf
sudo sed -i '$ a fs.file-max=200000' /etc/sysctl.conf
ulimit -n 40000

sudo systemctl enable neo4j.service
sudo systemctl restart neo4j.service


echo "CREATE INDEX agg_index IF NOT EXISTS FOR (a: Aggregator) ON (a.k, a.s);" | sudo tee init.neo
echo "CREATE INDEX rec_by_superstep_index IF NOT EXISTS FOR (r: Recipient) ON (r.s);" | sudo tee -a init.neo
echo "CREATE INDEX msg_index IF NOT EXISTS FOR (m: Message) ON (m.s, m.r);" | sudo tee -a init.neo
echo "CREATE CONSTRAINT rec_uniqueness IF NOT EXISTS FOR (r:Recipient) REQUIRE (r.s, r.i) IS UNIQUE;" | sudo tee -a init.neo
echo "CREATE CONSTRAINT vertex_id_unique IF NOT EXISTS FOR (v: Vertex) REQUIRE v.i IS UNIQUE;" | sudo tee -a init.neo
echo "CREATE CONSTRAINT active_workers_unique IF NOT EXISTS FOR (aw: ActiveWorkers) REQUIRE aw.c IS UNIQUE;" | sudo tee -a init.neo
echo "CREATE CONSTRAINT fw_t_unique IF NOT EXISTS FOR (fw: FinishedWorker) REQUIRE fw.t IS UNIQUE;" | sudo tee -a init.neo
cypher-shell -f init.neo


echo "ulimit -n 40000" | sudo tee wipe-db.sh
echo "sudo systemctl stop neo4j.service" | sudo tee -a wipe-db.sh
echo "sudo rm -rf /var/lib/neo4j/data/databases/*" | sudo tee -a wipe-db.sh
echo "sudo rm -rf /var/lib/neo4j/data/transactions/*" | sudo tee -a wipe-db.sh
echo "sudo systemctl restart neo4j.service" | sudo tee -a wipe-db.sh
echo "cypher-shell -f init.neo" | sudo tee -a wipe-db.sh
chmod +x wipe-db.sh
