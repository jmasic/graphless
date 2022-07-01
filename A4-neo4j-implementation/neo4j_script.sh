#!/bin/bash


# https://us-east-2.console.aws.amazon.com/cloudwatch/home?region=us-east-2#logsV2:log-groups/log-group/$252Faws$252Flambda$252FWorkerFunction/log-events/2022$252F05$252F24$252F$255B$2524LATEST$255D1f89d2a80f634a539704b6a6ab7d90a5

# c6i.8xlarge // r5n.8xlarge
# NOTE: open ports 7687 and optionally 7474

sudo apt update

sudo apt install -y apt-transport-https ca-certificates curl software-properties-common net-tools \
                    openjdk-11-jre-headless

sudo wget -O - https://debian.neo4j.com/neotechnology.gpg.key | sudo apt-key add -
echo 'deb https://debian.neo4j.com stable 4.4' | sudo tee -a /etc/apt/sources.list.d/neo4j.list
sudo apt-get update -y
sudo apt-get install -y neo4j=1:4.4.6

sudo sed -i "/^#dbms.default_listen_address=0.0.0.0/c\dbms.default_listen_address=0.0.0.0" /etc/neo4j/neo4j.conf
sudo sed -i "/^#dbms.security.auth_enabled=false/c\dbms.security.auth_enabled=false" /etc/neo4j/neo4j.conf
sudo sed -i "/^dbms.tx_log.rotation.retention_policy=1 days/c\dbms.tx_log.rotation.retention_policy=keep_none" /etc/neo4j/neo4j.conf
echo "dbms.connector.bolt.thread_pool_max_size=1000" | sudo tee -a /etc/neo4j/neo4j.conf
echo "dbms.checkpoint.interval.time=1m" | sudo tee -a /etc/neo4j/neo4j.conf
echo "dbms.checkpoint.interval.tx=100000" | sudo tee -a /etc/neo4j/neo4j.conf

echo "[Service]" | sudo tee /etc/systemd/system/neo4j.service.d/override.conf
echo "LimitNOFILE=60000" | sudo tee -a /etc/systemd/system/neo4j.service.d/override.conf

sudo systemctl enable neo4j.service
sudo systemctl restart neo4j.service

echo "CREATE INDEX agg_index IF NOT EXISTS FOR (a: Aggregator) ON (a.k, a.s);" | sudo tee -a init.neo
echo "CREATE INDEX rec_by_superstep_index IF NOT EXISTS FOR (r: Recipient) ON (r.s);" | sudo tee -a init.neo
echo "CREATE INDEX msg_index IF NOT EXISTS FOR (m: Message) ON (m.s, m.r);" | sudo tee -a init.neo
echo "CREATE CONSTRAINT rec_uniqueness FOR (r:Recipient) REQUIRE (r.s, r.i) IS UNIQUE;" | sudo tee -a init.neo
echo "CREATE CONSTRAINT vertex_id_unique IF NOT EXISTS FOR (v: Vertex) REQUIRE v.i IS UNIQUE;" | sudo tee -a init.neo
cypher-shell -f init.neo

# TODO: Execute init.neo once neo4j starts
