import os
import sys
import time

os.system('docker kill neo4j')
os.system('docker run --rm -p 7687:7687 -p 7474:7474 --env NEO4J_AUTH=neo4j/n --name neo4j -d neo4j:4.4.6')

res = os.system('docker exec -it neo4j cypher-shell -uneo4j -pn "MATCH (n) RETURN COUNT(n)"')
while res != 0:
    print('neo4j not ready yet')
    time.sleep(0.5)
    res = os.system("docker exec -it neo4j cypher-shell -uneo4j -pn \"MATCH (n) RETURN COUNT(n)\"")

print('neo4j started!')

os.system('docker exec -it neo4j cypher-shell -uneo4j -pn "CREATE INDEX agg_index IF NOT EXISTS FOR (a: Aggregator) ON (a.k, a.s);"')
os.system('docker exec -it neo4j cypher-shell -uneo4j -pn "CREATE INDEX rec_by_superstep_index IF NOT EXISTS FOR (r: Recipient) ON (r.s);"')
os.system('docker exec -it neo4j cypher-shell -uneo4j -pn "CREATE INDEX msg_index IF NOT EXISTS FOR (m: Message) ON (m.s, m.r);"')
os.system('docker exec -it neo4j cypher-shell -uneo4j -pn "CREATE CONSTRAINT rec_uniqueness FOR (r:Recipient) REQUIRE (r.s, r.i) IS UNIQUE;"')
os.system('docker exec -it neo4j cypher-shell -uneo4j -pn "CREATE CONSTRAINT vertex_id_unique IF NOT EXISTS FOR (v: Vertex) REQUIRE v.i IS UNIQUE;"')

print("Indexes correctly created")

os.system('docker logs -f neo4j')
