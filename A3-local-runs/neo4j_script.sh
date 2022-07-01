#!/bin/bash


# https://us-east-2.console.aws.amazon.com/cloudwatch/home?region=us-east-2#logsV2:log-groups/log-group/$252Faws$252Flambda$252FWorkerFunction/log-events/2022$252F05$252F24$252F$255B$2524LATEST$255D1f89d2a80f634a539704b6a6ab7d90a5

# c6i.8xlarge
# NOTE: open ports 7687 and optionally 7474

sudo apt update

sudo apt install -y apt-transport-https ca-certificates curl software-properties-common net-tools

sudo curl -fsSL https://debian.neo4j.com/neotechnology.gpg.key | sudo apt-key add -

sudo add-apt-repository -y "deb https://debian.neo4j.com stable 4.0"

sudo apt install -y neo4j

sudo sed -i "/^#dbms.default_listen_address=0.0.0.0/c\dbms.default_listen_address=0.0.0.0" /etc/neo4j/neo4j.conf
sudo sed -i "/^#dbms.security.auth_enabled=false/c\dbms.security.auth_enabled=false" /etc/neo4j/neo4j.conf
sudo sed -i "/^dbms.tx_log.rotation.retention_policy=1 days/c\dbms.tx_log.rotation.retention_policy=keep_none" /etc/neo4j/neo4j.conf
echo "dbms.connector.bolt.thread_pool_max_size=1000" | sudo tee -a /etc/neo4j/neo4j.conf

echo "[Service]" | sudo tee /etc/systemd/system/neo4j.service.d/override.conf
echo "LimitNOFILE=60000" | sudo tee -a /etc/systemd/system/neo4j.service.d/override.conf

sudo systemctl enable neo4j.service
sudo systemctl restart neo4j.service
