var data = {
  "id": "b600161",
  "specification": "1.0.0",
  "description": "",
  "system": {
    "platform": {
      "name": "Graphmat",
      "acronym": "gmat",
      "version": "?",
      "link": "?"
    },
    "environment": {
      "name": "",
      "acronym": "",
      "version": "",
      "link": "",
      "machines": [
        {
          "quantity": "",
          "cpu": "",
          "memory": "",
          "network": "",
          "storage": ""
        }
      ]
    },
    "tool": {
      "graphalytics-core": {
        "name": "graphalytics-core",
        "version": "1.0.0",
        "link": "https://github.com/ldbc/ldbc_graphalytics"
      },
      "graphalytics-platforms-graphmat": {
        "name": "graphalytics-platforms-graphmat",
        "version": "0.2-SNAPSHOT",
        "link": "https://github.com/atlarge-research/graphalytics-platforms-graphmat"
      }
    },
    "pricing": ""
  },
  "benchmark": {
    "type": "custom",
    "name": "Custom Benchmark",
    "duration": "6482835",
    "timeout": "3600",
    "resources": {},
    "output": {
      "required": "-",
      "directory": "./output/"
    },
    "validation": {
      "required": "-",
      "directory": "/home/hduser/graphalytics-datasets"
    },
    "configurations": {
      "graph.datagen-8_9-fb.edge-properties.types": "real",
      "environment.machine.memory": "",
      "graph.datagen-8_3-zf.cdlp.max-iterations": "10",
      "graph.datagen-9_4-fb.sssp.weight-property": "weight",
      "graph.datagen-8_9-fb.sssp.source-vertex": "6",
      "graph.graph500-24.edge-file": "graph500-24.e",
      "graph.datagen-8_4-fb.vertex-file": "datagen-8_4-fb.v",
      "graph.datagen-9_1-fb.meta.vertices": "16087483",
      "graph.example-undirected.meta.vertices": "9",
      "graph.com-friendster.cdlp.max-iterations": "10",
      "graph.datagen-7_9-fb.vertex-file": "datagen-7_9-fb.v",
      "environment.machine.storage": "",
      "graph.datagen-8_7-zf.bfs.source-vertex": "6",
      "benchmark.runner.max-memory": "",
      "graph.datagen-8_3-zf.meta.vertices": "53525014",
      "graph.datagen-8_1-fb.algorithms": "bfs",
      "graph.datagen-7_7-zf.meta.edges": "32791267",
      "graph.datagen-8_7-zf.vertex-file": "datagen-8_7-zf.v",
      "graph.datagen-8_8-zf.meta.vertices": "168308893",
      "graph.datagen-9_2-zf.vertex-file": "datagen-9_2-zf.v",
      "graph.dota-league.cdlp.max-iterations": "10",
      "graph.datagen-8_5-fb.sssp.source-vertex": "6",
      "graph.datagen-8_8-zf.edge-properties.types": "real",
      "graph.graph500-22.cdlp.max-iterations": "10",
      "graph.graph500-23.algorithms": "bfs",
      "graph.graph500-22.edge-file": "graph500-22.e",
      "graph.datagen-8_0-fb.edge-properties.names": "weight",
      "graph.datagen-8_0-fb.meta.vertices": "1706561",
      "graph.example-undirected.edge-properties.names": "weight",
      "graph.datagen-8_0-fb.algorithms": "bfs",
      "graph.datagen-9_4-fb.edge-properties.types": "real",
      "graph.datagen-7_7-zf.algorithms": "bfs",
      "graph.datagen-8_0-fb.bfs.source-vertex": "6",
      "graph.twitter_mpi.cdlp.max-iterations": "10",
      "graph.graph500-22.bfs.source-vertex": "248533",
      "graph.example-directed.bfs.source-vertex": "1",
      "graph.datagen-9_1-fb.edge-properties.types": "real",
      "graph.datagen-8_5-fb.vertex-file": "datagen-8_5-fb.v",
      "graph.datagen-8_9-fb.cdlp.max-iterations": "10",
      "graph.com-friendster.meta.edges": "1806067135",
      "graph.graph500-26.meta.vertices": "32804978",
      "graph.datagen-7_8-zf.algorithms": "bfs",
      "graph.datagen-9_2-zf.sssp.weight-property": "weight",
      "graph.datagen-9_2-zf.sssp.source-vertex": "6",
      "graph.datagen-7_5-fb.algorithms": "bfs",
      "graph.graph500-24.bfs.source-vertex": "2592222",
      "graph.datagen-7_6-fb.edge-file": "datagen-7_6-fb.e",
      "graph.datagen-8_8-zf.bfs.source-vertex": "6",
      "graph.dota-league.meta.vertices": "61170",
      "graph.datagen-7_7-zf.meta.vertices": "13180508",
      "graph.datagen-7_5-fb.meta.edges": "34185747",
      "graph.example-directed.edge-properties.types": "real",
      "graph.datagen-7_8-zf.edge-file": "datagen-7_8-zf.e",
      "graph.datagen-9_2-zf.pr.damping-factor": "0.85",
      "graph.datagen-9_3-zf.pr.num-iterations": "10",
      "graph.datagen-9_3-zf.edge-properties.types": "real",
      "graph.datagen-7_5-fb.edge-properties.names": "weight",
      "graph.datagen-8_2-zf.bfs.source-vertex": "6",
      "graph.twitter_mpi.meta.edges": "1963263508",
      "graph.datagen-7_7-zf.edge-properties.names": "weight",
      "graph.dota-league.edge-properties.types": "real",
      "graph.datagen-9_4-fb.edge-properties.names": "weight",
      "graph.datagen-9_3-zf.vertex-file": "datagen-9_3-zf.v",
      "graph.datagen-7_9-fb.edge-properties.names": "weight",
      "graph.datagen-8_9-fb.sssp.weight-property": "weight",
      "graph.datagen-8_2-zf.sssp.source-vertex": "6",
      "graph.graph500-26.cdlp.max-iterations": "10",
      "graph.graph500-26.bfs.source-vertex": "62455266",
      "graph.datagen-9_4-fb.meta.vertices": "29310565",
      "platform.graphmat.hosts": "10.164.0.2,10.164.0.3,10.164.0.4,10.164.0.5,10.164.0.6,10.164.0.7,10.164.0.8,10.164.0.9,10.164.0.10,10.164.0.11,10.164.0.12,10.164.0.13,10.164.0.14,10.164.0.15,10.164.0.16,10.164.0.17",
      "graph.datagen-8_4-fb.sssp.weight-property": "weight",
      "graph.datagen-8_1-fb.meta.edges": "134267822",
      "graph.example-undirected.sssp.source-vertex": "2",
      "graph.twitter_mpi.algorithms": "bfs",
      "graph.datagen-7_5-fb.directed": "false",
      "graph.graph500-25.directed": "false",
      "graph.datagen-9_1-fb.pr.num-iterations": "10",
      "graph.datagen-7_9-fb.edge-file": "datagen-7_9-fb.e",
      "graph.datagen-8_4-fb.bfs.source-vertex": "6",
      "graphs.cache-directory": "/home/hduser/graphalytics-platforms-graphmat/graphalytics-1.0.0-graphmat-0.2-SNAPSHOT/cache",
      "graph.dota-league.pr.damping-factor": "0.85",
      "graph.datagen-9_3-zf.edge-properties.names": "weight",
      "graph.twitter_mpi.edge-file": "twitter_mpi.e",
      "graph.datagen-7_5-fb.edge-properties.types": "real",
      "graph.datagen-8_5-fb.directed": "false",
      "graph.datagen-8_9-fb.directed": "false",
      "graph.datagen-8_5-fb.cdlp.max-iterations": "10",
      "graph.datagen-7_6-fb.meta.edges": "42162988",
      "graph.datagen-8_2-zf.pr.num-iterations": "10",
      "graph.datagen-8_0-fb.meta.edges": "107507376",
      "graph.datagen-9_0-fb.vertex-file": "datagen-9_0-fb.v",
      "graph.graph500-23.bfs.source-vertex": "7348998",
      "graph.datagen-9_1-fb.sssp.source-vertex": "6",
      "graph.datagen-8_4-fb.edge-file": "datagen-8_4-fb.e",
      "graph.datagen-9_0-fb.pr.damping-factor": "0.85",
      "graph.datagen-8_6-fb.vertex-file": "datagen-8_6-fb.v",
      "graph.datagen-8_1-fb.directed": "false",
      "graph.example-undirected.directed": "false",
      "graph.datagen-7_9-fb.directed": "false",
      "platform.acronym": "gmat",
      "graph.datagen-7_5-fb.edge-file": "datagen-7_5-fb.e",
      "platform.graphmat.intermediate-dir": "intermediate",
      "platform.graphmat.home": "/home/hduser/GraphMat",
      "graph.com-friendster.edge-file": "com-friendster.e",
      "graph.dota-league.sssp.weight-property": "weight",
      "graph.datagen-7_9-fb.meta.vertices": "1387587",
      "graph.datagen-8_1-fb.edge-properties.names": "weight",
      "graph.graph500-25.vertex-file": "graph500-25.v",
      "platform.graphmat.command.prefix": "env I_MPI_DEBUG\u003d2 I_MPI_FABRICS_LIST\u003dtmi,dapl,tcp I_MPI_TMI_PROVIDER\u003dpsm2 KMP_AFFINITY\u003dscatter OMP_NUM_THREADS\u003d16 salloc -N 1 --ntasks-per-node\u003d1 mpiexec.hydra numactl -i all",
      "graph.datagen-8_2-zf.algorithms": "bfs",
      "graph.com-friendster.pr.damping-factor": "0.85",
      "graph.datagen-9_1-fb.directed": "false",
      "graph.datagen-9_3-zf.meta.edges": "1309998551",
      "graph.example-undirected.bfs.source-vertex": "2",
      "graph.datagen-7_6-fb.pr.damping-factor": "0.85",
      "platform.graphmat.command.convert": "env I_MPI_DEBUG\u003d2 I_MPI_FABRICS_LIST\u003dtmi,dapl,tcp I_MPI_TMI_PROVIDER\u003dpsm2 KMP_AFFINITY\u003dscatter OMP_NUM_THREADS\u003d16 salloc -N 1 --ntasks-per-node\u003d1 mpiexec.hydra numactl -i all %s %s",
      "environment.machine.cpu": "",
      "graph.datagen-7_7-zf.sssp.weight-property": "weight",
      "environment.acronym": "",
      "graph.dota-league.meta.edges": "50870313",
      "graph.datagen-8_1-fb.vertex-file": "datagen-8_1-fb.v",
      "graph.datagen-8_7-zf.pr.num-iterations": "10",
      "graph.datagen-9_0-fb.meta.edges": "1049527225",
      "graph.datagen-7_8-zf.meta.vertices": "16521886",
      "graph.example-undirected.edge-file": "example-undirected.e",
      "graph.datagen-8_8-zf.edge-file": "datagen-8_8-zf.e",
      "graph.graph500-24.vertex-file": "graph500-24.v",
      "graph.datagen-7_6-fb.algorithms": "bfs",
      "graph.graph500-25.meta.vertices": "17062472",
      "graph.datagen-7_8-zf.directed": "false",
      "graph.datagen-9_3-zf.sssp.weight-property": "weight",
      "platform.graphmat.command.run": "env I_MPI_DEBUG\u003d2 I_MPI_FABRICS_LIST\u003dtmi,dapl,tcp I_MPI_TMI_PROVIDER\u003dpsm2 KMP_AFFINITY\u003dscatter OMP_NUM_THREADS\u003d16 salloc -N 1 --ntasks-per-node\u003d1 mpiexec.hydra numactl -i all %s %s",
      "graph.datagen-8_5-fb.bfs.source-vertex": "6",
      "graph.example-directed.vertex-file": "example-directed.v",
      "graph.datagen-8_3-zf.directed": "false",
      "graph.example-directed.algorithms": "bfs",
      "graph.datagen-7_6-fb.edge-properties.names": "weight",
      "graph.datagen-8_1-fb.cdlp.max-iterations": "10",
      "graph.datagen-8_9-fb.edge-properties.names": "weight",
      "graph.example-undirected.meta.edges": "12",
      "graph.datagen-8_8-zf.directed": "false",
      "graph.datagen-8_9-fb.algorithms": "bfs",
      "graph.graph500-23.pr.damping-factor": "0.85",
      "graph.datagen-8_7-zf.directed": "false",
      "graph.datagen-8_0-fb.vertex-file": "datagen-8_0-fb.v",
      "graph.graph500-24.meta.vertices": "8870942",
      "platform.graphmap.num_local_processes": "16",
      "graph.datagen-8_6-fb.edge-properties.names": "weight",
      "benchmark.runner.port": "8012",
      "graph.datagen-7_9-fb.pr.damping-factor": "0.85",
      "graph.datagen-8_8-zf.sssp.weight-property": "weight",
      "graph.datagen-9_0-fb.edge-properties.types": "real",
      "graph.datagen-8_3-zf.algorithms": "bfs",
      "graph.datagen-7_5-fb.sssp.source-vertex": "6",
      "graph.com-friendster.directed": "false",
      "graph.datagen-9_2-zf.meta.edges": "1042340732",
      "graph.datagen-7_7-zf.directed": "false",
      "graph.datagen-8_0-fb.sssp.weight-property": "weight",
      "benchmark.description": "",
      "graph.datagen-8_3-zf.vertex-file": "datagen-8_3-zf.v",
      "graph.graph500-22.algorithms": "bfs",
      "graph.graph500-23.pr.num-iterations": "10",
      "benchmark.custom.validation-required": "false",
      "graph.datagen-8_1-fb.sssp.source-vertex": "6",
      "graph.datagen-9_0-fb.meta.vertices": "12857671",
      "graph.datagen-8_2-zf.sssp.weight-property": "weight",
      "graph.datagen-7_7-zf.vertex-file": "datagen-7_7-zf.v",
      "graph.datagen-8_5-fb.meta.vertices": "4599739",
      "graph.graph500-24.pr.damping-factor": "0.85",
      "graph.datagen-8_0-fb.sssp.source-vertex": "6",
      "graph.datagen-9_1-fb.edge-properties.names": "weight",
      "graph.datagen-8_8-zf.algorithms": "bfs",
      "benchmark.custom.graphs": "example-directed",
      "graph.datagen-8_9-fb.edge-file": "datagen-8_9-fb.e",
      "graph.datagen-8_7-zf.cdlp.max-iterations": "10",
      "graph.example-directed.sssp.source-vertex": "1",
      "graph.datagen-9_4-fb.pr.damping-factor": "0.85",
      "graph.datagen-9_2-zf.edge-properties.names": "weight",
      "environment.machine.network": "",
      "graph.twitter_mpi.vertex-file": "twitter_mpi.v",
      "platform.graphmat.num-machines": "1",
      "graph.datagen-8_1-fb.bfs.source-vertex": "6",
      "graph.datagen-9_3-zf.pr.damping-factor": "0.85",
      "graph.datagen-8_1-fb.pr.num-iterations": "10",
      "graph.datagen-8_2-zf.pr.damping-factor": "0.85",
      "graph.graph500-26.edge-file": "graph500-26.e",
      "graph.datagen-9_1-fb.meta.edges": "1342158397",
      "graph.datagen-9_3-zf.directed": "false",
      "graph.datagen-8_3-zf.pr.num-iterations": "10",
      "graph.datagen-7_8-zf.pr.damping-factor": "0.85",
      "graph.graph500-25.algorithms": "bfs",
      "graph.datagen-8_3-zf.bfs.source-vertex": "6",
      "graph.datagen-9_1-fb.sssp.weight-property": "weight",
      "graph.graph500-23.meta.vertices": "4610222",
      "graph.graph500-26.directed": "false",
      "graph.datagen-8_0-fb.edge-properties.types": "real",
      "environment.link": "",
      "graph.dota-league.pr.num-iterations": "10",
      "graph.datagen-7_8-zf.meta.edges": "41025255",
      "graph.datagen-9_4-fb.meta.edges": "2588948669",
      "graph.datagen-8_4-fb.cdlp.max-iterations": "10",
      "graph.datagen-7_6-fb.directed": "false",
      "platform.link": "?",
      "graph.datagen-9_1-fb.pr.damping-factor": "0.85",
      "benchmark.executor.port": "8011",
      "platform.graphmat.enable-slurm": "False",
      "graph.datagen-8_0-fb.pr.damping-factor": "0.85",
      "graph.graph500-22.pr.num-iterations": "10",
      "graph.datagen-8_2-zf.edge-file": "datagen-8_2-zf.e",
      "graph.datagen-8_4-fb.edge-properties.types": "real",
      "graph.datagen-8_6-fb.directed": "false",
      "graph.example-directed.cdlp.max-iterations": "2",
      "graph.datagen-9_1-fb.edge-file": "datagen-9_1-fb.e",
      "benchmark.custom.output-required": "true",
      "graph.graph500-23.meta.edges": "129333677",
      "graph.datagen-8_2-zf.edge-properties.names": "weight",
      "graph.com-friendster.algorithms": "bfs",
      "graph.graph500-23.directed": "false",
      "graph.datagen-8_7-zf.algorithms": "bfs",
      "graph.datagen-7_6-fb.sssp.source-vertex": "6",
      "graph.datagen-8_8-zf.sssp.source-vertex": "6",
      "graph.datagen-8_4-fb.meta.vertices": "3809084",
      "graph.datagen-7_6-fb.bfs.source-vertex": "6",
      "graph.datagen-8_3-zf.edge-file": "datagen-8_3-zf.e",
      "graph.graph500-22.meta.vertices": "2396657",
      "graph.com-friendster.bfs.source-vertex": "101",
      "graph.datagen-7_9-fb.meta.edges": "85670523",
      "environment.machine.quantity": "",
      "graph.datagen-7_5-fb.sssp.weight-property": "weight",
      "graph.graph500-24.algorithms": "bfs",
      "graph.example-undirected.pr.num-iterations": "2",
      "graphs.validation-directory": "/home/hduser/graphalytics-datasets",
      "graph.example-directed.sssp.weight-property": "weight",
      "graph.datagen-8_1-fb.sssp.weight-property": "weight",
      "graph.datagen-8_6-fb.algorithms": "bfs",
      "graph.datagen-8_7-zf.edge-properties.names": "weight",
      "graph.datagen-8_9-fb.pr.num-iterations": "10",
      "graph.datagen-8_0-fb.cdlp.max-iterations": "10",
      "graph.datagen-8_1-fb.edge-properties.types": "real",
      "graph.graph500-23.edge-file": "graph500-23.e",
      "graph.datagen-8_7-zf.meta.vertices": "145050709",
      "graph.datagen-8_6-fb.meta.edges": "421988619",
      "graph.datagen-8_9-fb.meta.vertices": "10572901",
      "platform.graphmat.command.prefix-no-slurm": "env I_MPI_DEBUG\u003d2 I_MPI_FABRICS_LIST\u003dtmi,dapl,tcp I_MPI_TMI_PROVIDER\u003dpsm2,KMP_AFFINITY\u003dscatter,OMP_NUM_THREADS\u003d16 mpirun -np 16 --host 10.164.0.2,10.164.0.3,10.164.0.4,10.164.0.5,10.164.0.6,10.164.0.7,10.164.0.8,10.164.0.9,10.164.0.10,10.164.0.11,10.164.0.12,10.164.0.13,10.164.0.14,10.164.0.15,10.164.0.16,10.164.0.17 -x LD_LIBRARY_PATH",
      "graph.datagen-8_9-fb.pr.damping-factor": "0.85",
      "graph.dota-league.bfs.source-vertex": "287770",
      "graph.datagen-7_6-fb.meta.vertices": "754147",
      "graph.datagen-7_9-fb.cdlp.max-iterations": "10",
      "graph.datagen-8_5-fb.meta.edges": "332026902",
      "graph.datagen-8_6-fb.sssp.weight-property": "weight",
      "benchmark.name": "Custom Benchmark",
      "graph.datagen-7_7-zf.sssp.source-vertex": "6",
      "graph.datagen-8_4-fb.pr.num-iterations": "10",
      "graph.datagen-8_5-fb.edge-properties.names": "weight",
      "benchmark.custom.timeout": "3600",
      "graph.datagen-7_5-fb.bfs.source-vertex": "6",
      "graph.datagen-9_0-fb.edge-file": "datagen-9_0-fb.e",
      "graph.datagen-8_1-fb.edge-file": "datagen-8_1-fb.e",
      "graph.datagen-8_8-zf.meta.edges": "413354288",
      "graph.datagen-7_8-zf.sssp.weight-property": "weight",
      "graph.datagen-9_1-fb.cdlp.max-iterations": "10",
      "graph.graph500-24.cdlp.max-iterations": "10",
      "graph.graph500-25.pr.damping-factor": "0.85",
      "graph.example-undirected.pr.damping-factor": "0.85",
      "platform.name": "Graphmat",
      "graph.datagen-9_2-zf.bfs.source-vertex": "6",
      "graph.graph500-25.edge-file": "graph500-25.e",
      "graph.datagen-8_3-zf.edge-properties.types": "real",
      "graph.datagen-9_4-fb.sssp.source-vertex": "6",
      "graph.datagen-8_2-zf.meta.vertices": "43734497",
      "graph.datagen-8_2-zf.vertex-file": "datagen-8_2-zf.v",
      "graph.graph500-26.meta.edges": "1051922853",
      "graph.graph500-25.cdlp.max-iterations": "10",
      "graph.datagen-9_4-fb.cdlp.max-iterations": "10",
      "graph.datagen-9_1-fb.algorithms": "bfs",
      "graph.datagen-8_9-fb.vertex-file": "datagen-8_9-fb.v",
      "benchmark.custom.algorithms": "BFS",
      "graph.datagen-9_3-zf.bfs.source-vertex": "6",
      "graph.datagen-9_4-fb.vertex-file": "datagen-9_4-fb.v",
      "graph.datagen-9_2-zf.directed": "false",
      "graph.datagen-9_2-zf.edge-file": "datagen-9_2-zf.e",
      "graph.datagen-9_3-zf.edge-file": "datagen-9_3-zf.e",
      "graph.example-directed.directed": "true",
      "graph.datagen-8_1-fb.pr.damping-factor": "0.85",
      "graph.datagen-9_2-zf.cdlp.max-iterations": "10",
      "graph.graph500-26.algorithms": "bfs",
      "graph.datagen-7_9-fb.edge-properties.types": "real",
      "graph.datagen-8_5-fb.algorithms": "bfs",
      "environment.name": "",
      "graph.datagen-8_3-zf.sssp.source-vertex": "6",
      "graphs.output-directory": "./output/",
      "graph.datagen-8_4-fb.pr.damping-factor": "0.85",
      "graph.example-directed.meta.edges": "17",
      "graph.datagen-8_5-fb.pr.num-iterations": "10",
      "graph.example-undirected.sssp.weight-property": "weight",
      "graph.example-directed.pr.num-iterations": "2",
      "graph.datagen-7_6-fb.cdlp.max-iterations": "10",
      "graph.graph500-22.pr.damping-factor": "0.85",
      "graph.datagen-8_0-fb.pr.num-iterations": "10",
      "graph.graph500-22.meta.edges": "64155735",
      "graph.datagen-9_2-zf.meta.vertices": "434943376",
      "graph.datagen-7_7-zf.pr.damping-factor": "0.85",
      "graph.datagen-7_8-zf.edge-properties.types": "real",
      "graph.datagen-8_4-fb.algorithms": "bfs",
      "graph.datagen-8_8-zf.pr.damping-factor": "0.85",
      "graph.datagen-7_8-zf.sssp.source-vertex": "6",
      "graph.datagen-7_5-fb.vertex-file": "datagen-7_5-fb.v",
      "graph.datagen-8_1-fb.meta.vertices": "2072117",
      "graph.datagen-7_9-fb.sssp.weight-property": "weight",
      "graph.dota-league.directed": "false",
      "graph.twitter_mpi.bfs.source-vertex": "420",
      "graph.datagen-8_6-fb.sssp.source-vertex": "6",
      "graph.datagen-7_7-zf.cdlp.max-iterations": "10",
      "graph.datagen-8_9-fb.bfs.source-vertex": "6",
      "graph.datagen-9_0-fb.edge-properties.names": "weight",
      "graph.graph500-24.pr.num-iterations": "10",
      "graph.graph500-24.meta.edges": "260379520",
      "graph.graph500-26.pr.damping-factor": "0.85",
      "graph.datagen-8_3-zf.edge-properties.names": "weight",
      "graph.datagen-7_8-zf.pr.num-iterations": "10",
      "graph.datagen-7_9-fb.algorithms": "bfs",
      "graph.datagen-8_6-fb.pr.damping-factor": "0.85",
      "graph.datagen-8_8-zf.cdlp.max-iterations": "10",
      "graph.datagen-7_6-fb.pr.num-iterations": "10",
      "graph.datagen-7_8-zf.vertex-file": "datagen-7_8-zf.v",
      "graph.datagen-7_5-fb.pr.damping-factor": "0.85",
      "graph.datagen-8_0-fb.edge-file": "datagen-8_0-fb.e",
      "graph.datagen-8_7-zf.edge-properties.types": "real",
      "graph.datagen-8_6-fb.cdlp.max-iterations": "10",
      "graph.dota-league.edge-file": "dota-league.e",
      "graph.graph500-26.vertex-file": "graph500-26.v",
      "benchmark.type": "custom",
      "graph.datagen-9_0-fb.algorithms": "bfs",
      "graph.datagen-8_7-zf.sssp.source-vertex": "6",
      "graph.graph500-24.directed": "false",
      "graph.datagen-8_6-fb.meta.vertices": "5667674",
      "graph.datagen-8_5-fb.edge-properties.types": "real",
      "graph.datagen-9_4-fb.directed": "false",
      "graph.example-directed.edge-file": "example-directed.e",
      "graph.datagen-7_5-fb.pr.num-iterations": "10",
      "graph.datagen-7_6-fb.vertex-file": "datagen-7_6-fb.v",
      "graph.datagen-8_2-zf.cdlp.max-iterations": "10",
      "graph.datagen-9_0-fb.cdlp.max-iterations": "10",
      "graph.datagen-8_7-zf.edge-file": "datagen-8_7-zf.e",
      "graph.dota-league.sssp.source-vertex": "287770",
      "graph.datagen-9_2-zf.edge-properties.types": "real",
      "graph.graph500-23.vertex-file": "graph500-23.v",
      "graph.graph500-22.vertex-file": "graph500-22.v",
      "graph.datagen-7_9-fb.sssp.source-vertex": "6",
      "graph.datagen-7_7-zf.bfs.source-vertex": "6",
      "graph.datagen-9_0-fb.sssp.source-vertex": "6",
      "graph.datagen-9_1-fb.bfs.source-vertex": "6",
      "platform.graphmat.num-threads": "16",
      "graph.datagen-8_3-zf.pr.damping-factor": "0.85",
      "graphs.names": "[dota-league, com-friendster, twitter_mpi, graph500-22, graph500-23, graph500-24, graph500-25, graph500-26, datagen-7_5-fb, datagen-7_6-fb, datagen-7_7-zf, datagen-7_8-zf, datagen-7_9-fb, datagen-8_0-fb, datagen-8_1-fb, datagen-8_2-zf, datagen-8_3-zf, datagen-8_4-fb, datagen-8_5-fb, datagen-8_6-fb, datagen-8_7-zf, datagen-8_8-zf, datagen-8_9-fb, datagen-9_0-fb, datagen-9_1-fb, datagen-9_2-zf, datagen-9_3-zf, datagen-9_4-fb, example-directed, example-undirected]",
      "graph.twitter_mpi.pr.damping-factor": "0.85",
      "graph.datagen-8_2-zf.directed": "false",
      "graph.datagen-9_0-fb.sssp.weight-property": "weight",
      "graph.datagen-7_9-fb.pr.num-iterations": "10",
      "graph.datagen-9_3-zf.cdlp.max-iterations": "10",
      "graph.datagen-7_6-fb.sssp.weight-property": "weight",
      "graph.datagen-8_3-zf.sssp.weight-property": "weight",
      "graph.dota-league.edge-properties.names": "weight",
      "graph.datagen-8_7-zf.meta.edges": "340157363",
      "graph.datagen-7_8-zf.bfs.source-vertex": "6",
      "graph.com-friendster.meta.vertices": "65608366",
      "graph.graph500-22.directed": "false",
      "graph.graph500-23.cdlp.max-iterations": "10",
      "graph.datagen-8_6-fb.bfs.source-vertex": "6",
      "graph.datagen-9_4-fb.edge-file": "datagen-9_4-fb.e",
      "graph.com-friendster.pr.num-iterations": "10",
      "graph.twitter_mpi.pr.num-iterations": "10",
      "graph.datagen-7_8-zf.edge-properties.names": "weight",
      "graph.datagen-8_2-zf.edge-properties.types": "real",
      "graph.datagen-9_0-fb.bfs.source-vertex": "6",
      "graph.datagen-8_4-fb.edge-properties.names": "weight",
      "graph.datagen-8_5-fb.edge-file": "datagen-8_5-fb.e",
      "graph.datagen-8_5-fb.sssp.weight-property": "weight",
      "graph.example-directed.edge-properties.names": "weight",
      "graph.datagen-9_1-fb.vertex-file": "datagen-9_1-fb.v",
      "graph.datagen-8_4-fb.meta.edges": "269479177",
      "graph.datagen-9_4-fb.pr.num-iterations": "10",
      "platform.graphmat.command.convert-no-slurm": "env I_MPI_DEBUG\u003d2 I_MPI_FABRICS_LIST\u003dtmi,dapl,tcp I_MPI_TMI_PROVIDER\u003dpsm2,KMP_AFFINITY\u003dscatter,OMP_NUM_THREADS\u003d16 mpirun -np 16 --host 10.164.0.2,10.164.0.3,10.164.0.4,10.164.0.5,10.164.0.6,10.164.0.7,10.164.0.8,10.164.0.9,10.164.0.10,10.164.0.11,10.164.0.12,10.164.0.13,10.164.0.14,10.164.0.15,10.164.0.16,10.164.0.17 -x LD_LIBRARY_PATH %s %s",
      "graph.graph500-25.bfs.source-vertex": "24460635",
      "graph.example-undirected.cdlp.max-iterations": "2",
      "graph.datagen-7_7-zf.edge-properties.types": "real",
      "graph.datagen-9_2-zf.pr.num-iterations": "10",
      "graph.datagen-8_7-zf.sssp.weight-property": "weight",
      "graph.graph500-25.pr.num-iterations": "10",
      "graph.datagen-9_4-fb.algorithms": "bfs",
      "platform.version": "?",
      "graph.datagen-8_8-zf.vertex-file": "datagen-8_8-zf.v",
      "graph.example-undirected.vertex-file": "example-undirected.v",
      "graph.datagen-8_4-fb.sssp.source-vertex": "6",
      "graph.datagen-8_8-zf.edge-properties.names": "weight",
      "graph.example-undirected.edge-properties.types": "real",
      "graph.graph500-26.pr.num-iterations": "10",
      "graph.dota-league.vertex-file": "dota-league.v",
      "graph.datagen-7_7-zf.edge-file": "datagen-7_7-zf.e",
      "graph.datagen-9_3-zf.meta.vertices": "555270053",
      "graph.datagen-9_3-zf.algorithms": "bfs",
      "graph.twitter_mpi.directed": "true",
      "system.pricing": "",
      "platform.graphmat.command.run-no-slurm": "env I_MPI_DEBUG\u003d2 I_MPI_FABRICS_LIST\u003dtmi,dapl,tcp I_MPI_TMI_PROVIDER\u003dpsm2,KMP_AFFINITY\u003dscatter,OMP_NUM_THREADS\u003d16 mpirun -np 16 --host 10.164.0.2,10.164.0.3,10.164.0.4,10.164.0.5,10.164.0.6,10.164.0.7,10.164.0.8,10.164.0.9,10.164.0.10,10.164.0.11,10.164.0.12,10.164.0.13,10.164.0.14,10.164.0.15,10.164.0.16,10.164.0.17 -x LD_LIBRARY_PATH %s %s",
      "graph.example-undirected.algorithms": "bfs",
      "graph.datagen-8_2-zf.meta.edges": "106440188",
      "graph.datagen-8_5-fb.pr.damping-factor": "0.85",
      "graph.dota-league.algorithms": "bfs",
      "graph.datagen-8_6-fb.edge-file": "datagen-8_6-fb.e",
      "graph.datagen-8_6-fb.edge-properties.types": "real",
      "graphs.root-directory": "/home/hduser/graphalytics-datasets",
      "graph.datagen-7_7-zf.pr.num-iterations": "10",
      "graph.com-friendster.vertex-file": "com-friendster.v",
      "graph.graph500-25.meta.edges": "523602831",
      "graph.datagen-7_5-fb.meta.vertices": "633432",
      "graph.twitter_mpi.meta.vertices": "52579678",
      "graph.datagen-8_6-fb.pr.num-iterations": "10",
      "graph.datagen-7_5-fb.cdlp.max-iterations": "10",
      "graph.datagen-9_2-zf.algorithms": "bfs",
      "graph.datagen-7_6-fb.edge-properties.types": "real",
      "graph.datagen-9_0-fb.pr.num-iterations": "10",
      "graph.datagen-8_3-zf.meta.edges": "130579909",
      "graph.datagen-8_7-zf.pr.damping-factor": "0.85",
      "graph.datagen-8_8-zf.pr.num-iterations": "10",
      "graph.example-directed.pr.damping-factor": "0.85",
      "graph.datagen-8_0-fb.directed": "false",
      "graph.datagen-8_9-fb.meta.edges": "848681908",
      "graph.datagen-9_0-fb.directed": "false",
      "graph.datagen-7_9-fb.bfs.source-vertex": "6",
      "graph.datagen-7_8-zf.cdlp.max-iterations": "10",
      "environment.version": "",
      "graph.datagen-9_4-fb.bfs.source-vertex": "6",
      "graph.datagen-9_3-zf.sssp.source-vertex": "6",
      "graph.example-directed.meta.vertices": "10",
      "graph.datagen-8_4-fb.directed": "false"
    }
  },
  "result": {
    "experiments": {
      "e699636": {
        "id": "e699636",
        "type": "custom:exp",
        "jobs": [
          "j491887",
          "j662286",
          "j651465",
          "j632848",
          "j759682",
          "j596864",
          "j555830",
          "j607365",
          "j511404",
          "j711040",
          "j562422",
          "j636523",
          "j622183",
          "j520674",
          "j587384",
          "j874925",
          "j693360",
          "j832855",
          "j677624",
          "j731524",
          "j550676",
          "j608341",
          "j881613",
          "j636793",
          "j486792",
          "j854196",
          "j780744",
          "j734618",
          "j872326",
          "j840029",
          "j514359",
          "j824190",
          "j507770",
          "j737752",
          "j910992",
          "j524439"
        ]
      }
    },
    "jobs": {
      "j832855": {
        "id": "j832855",
        "algorithm": "WCC",
        "dataset": "example-undirected",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r711444"
        ]
      },
      "j555830": {
        "id": "j555830",
        "algorithm": "CDLP",
        "dataset": "example-directed",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r552385"
        ]
      },
      "j824190": {
        "id": "j824190",
        "algorithm": "LCC",
        "dataset": "datagen-8_4-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r647311"
        ]
      },
      "j636523": {
        "id": "j636523",
        "algorithm": "CDLP",
        "dataset": "example-undirected",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r772774"
        ]
      },
      "j910992": {
        "id": "j910992",
        "algorithm": "LCC",
        "dataset": "datagen-7_9-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r604861"
        ]
      },
      "j607365": {
        "id": "j607365",
        "algorithm": "CDLP",
        "dataset": "datagen-8_4-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r898666"
        ]
      },
      "j780744": {
        "id": "j780744",
        "algorithm": "SSSP",
        "dataset": "dota-league",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r488604"
        ]
      },
      "j596864": {
        "id": "j596864",
        "algorithm": "BFS",
        "dataset": "example-undirected",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r700889"
        ]
      },
      "j632848": {
        "id": "j632848",
        "algorithm": "BFS",
        "dataset": "datagen-8_9-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r771152"
        ]
      },
      "j587384": {
        "id": "j587384",
        "algorithm": "WCC",
        "dataset": "dota-league",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r658444"
        ]
      },
      "j737752": {
        "id": "j737752",
        "algorithm": "LCC",
        "dataset": "datagen-8_9-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r575649"
        ]
      },
      "j677624": {
        "id": "j677624",
        "algorithm": "PR",
        "dataset": "example-directed",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r585397"
        ]
      },
      "j562422": {
        "id": "j562422",
        "algorithm": "CDLP",
        "dataset": "datagen-7_9-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r773002"
        ]
      },
      "j491887": {
        "id": "j491887",
        "algorithm": "BFS",
        "dataset": "example-directed",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r709841"
        ]
      },
      "j636793": {
        "id": "j636793",
        "algorithm": "PR",
        "dataset": "example-undirected",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r470371"
        ]
      },
      "j881613": {
        "id": "j881613",
        "algorithm": "PR",
        "dataset": "datagen-7_9-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r500976"
        ]
      },
      "j662286": {
        "id": "j662286",
        "algorithm": "BFS",
        "dataset": "datagen-8_4-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r901913"
        ]
      },
      "j511404": {
        "id": "j511404",
        "algorithm": "CDLP",
        "dataset": "dota-league",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r744733"
        ]
      },
      "j608341": {
        "id": "j608341",
        "algorithm": "PR",
        "dataset": "datagen-8_9-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r794507"
        ]
      },
      "j854196": {
        "id": "j854196",
        "algorithm": "SSSP",
        "dataset": "datagen-8_4-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r485181"
        ]
      },
      "j711040": {
        "id": "j711040",
        "algorithm": "CDLP",
        "dataset": "datagen-8_9-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r813951"
        ]
      },
      "j550676": {
        "id": "j550676",
        "algorithm": "PR",
        "dataset": "dota-league",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r630811"
        ]
      },
      "j759682": {
        "id": "j759682",
        "algorithm": "BFS",
        "dataset": "datagen-7_9-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r786790"
        ]
      },
      "j872326": {
        "id": "j872326",
        "algorithm": "SSSP",
        "dataset": "datagen-7_9-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r568729"
        ]
      },
      "j622183": {
        "id": "j622183",
        "algorithm": "WCC",
        "dataset": "example-directed",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r761644"
        ]
      },
      "j874925": {
        "id": "j874925",
        "algorithm": "WCC",
        "dataset": "datagen-8_9-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r489083"
        ]
      },
      "j693360": {
        "id": "j693360",
        "algorithm": "WCC",
        "dataset": "datagen-7_9-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r640938"
        ]
      },
      "j524439": {
        "id": "j524439",
        "algorithm": "LCC",
        "dataset": "example-undirected",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r603080"
        ]
      },
      "j520674": {
        "id": "j520674",
        "algorithm": "WCC",
        "dataset": "datagen-8_4-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r808888"
        ]
      },
      "j651465": {
        "id": "j651465",
        "algorithm": "BFS",
        "dataset": "dota-league",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r738819"
        ]
      },
      "j486792": {
        "id": "j486792",
        "algorithm": "SSSP",
        "dataset": "example-directed",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r755973"
        ]
      },
      "j731524": {
        "id": "j731524",
        "algorithm": "PR",
        "dataset": "datagen-8_4-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r613261"
        ]
      },
      "j734618": {
        "id": "j734618",
        "algorithm": "SSSP",
        "dataset": "datagen-8_9-fb",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r919579"
        ]
      },
      "j514359": {
        "id": "j514359",
        "algorithm": "LCC",
        "dataset": "example-directed",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r517723"
        ]
      },
      "j840029": {
        "id": "j840029",
        "algorithm": "SSSP",
        "dataset": "example-undirected",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r703209"
        ]
      },
      "j507770": {
        "id": "j507770",
        "algorithm": "LCC",
        "dataset": "dota-league",
        "scale": "1",
        "repetition": "1",
        "runs": [
          "r551627"
        ]
      }
    },
    "runs": {
      "r794507": {
        "id": "r794507",
        "timestamp": "1534514254245",
        "success": "true",
        "load_time": "1168.597",
        "makespan": "59.560",
        "processing_time": "2.595"
      },
      "r640938": {
        "id": "r640938",
        "timestamp": "1534516829017",
        "success": "true",
        "load_time": "113.881",
        "makespan": "7.532",
        "processing_time": "0.367"
      },
      "r604861": {
        "id": "r604861",
        "timestamp": "1534516758956",
        "success": "true",
        "load_time": "113.881",
        "makespan": "12.566",
        "processing_time": "4.297"
      },
      "r551627": {
        "id": "r551627",
        "timestamp": "1534512267725",
        "success": "true",
        "load_time": "51.944",
        "makespan": "28.825",
        "processing_time": "26.517"
      },
      "r711444": {
        "id": "r711444",
        "timestamp": "1534517112473",
        "success": "true",
        "load_time": "1.377",
        "makespan": "1.426",
        "processing_time": "0.052"
      },
      "r761644": {
        "id": "r761644",
        "timestamp": "1534510680694",
        "success": "true",
        "load_time": "1.402",
        "makespan": "1.508",
        "processing_time": "0.038"
      },
      "r630811": {
        "id": "r630811",
        "timestamp": "1534512340788",
        "success": "true",
        "load_time": "51.944",
        "makespan": "2.733",
        "processing_time": "0.324"
      },
      "r709841": {
        "id": "r709841",
        "timestamp": "1534510736787",
        "success": "true",
        "load_time": "1.402",
        "makespan": "1.502",
        "processing_time": "0.033"
      },
      "r500976": {
        "id": "r500976",
        "timestamp": "1534516849042",
        "success": "true",
        "load_time": "113.881",
        "makespan": "8.730",
        "processing_time": "0.443"
      },
      "r700889": {
        "id": "r700889",
        "timestamp": "1534517070413",
        "success": "true",
        "load_time": "1.377",
        "makespan": "1.441",
        "processing_time": "0.038"
      },
      "r488604": {
        "id": "r488604",
        "timestamp": "1534512164216",
        "success": "true",
        "load_time": "81.838",
        "makespan": "2.465",
        "processing_time": "0.155"
      },
      "r772774": {
        "id": "r772774",
        "timestamp": "1534517084429",
        "success": "true",
        "load_time": "1.377",
        "makespan": "1.424",
        "processing_time": "0.044"
      },
      "r813951": {
        "id": "r813951",
        "timestamp": "1534514508467",
        "success": "true",
        "load_time": "1168.597",
        "makespan": "85.198",
        "processing_time": "33.759"
      },
      "r647311": {
        "id": "r647311",
        "timestamp": "1534511375459",
        "success": "true",
        "load_time": "360.069",
        "makespan": "35.543",
        "processing_time": "14.228"
      },
      "r773002": {
        "id": "r773002",
        "timestamp": "1534516804008",
        "success": "true",
        "load_time": "113.881",
        "makespan": "11.940",
        "processing_time": "4.281"
      },
      "r568729": {
        "id": "r568729",
        "timestamp": "1534517033895",
        "success": "true",
        "load_time": "163.785",
        "makespan": "8.019",
        "processing_time": "0.090"
      },
      "r901913": {
        "id": "r901913",
        "timestamp": "1534511497544",
        "success": "true",
        "load_time": "360.069",
        "makespan": "17.739",
        "processing_time": "0.208"
      },
      "r658444": {
        "id": "r658444",
        "timestamp": "1534512252703",
        "success": "true",
        "load_time": "51.944",
        "makespan": "2.419",
        "processing_time": "0.187"
      },
      "r919579": {
        "id": "r919579",
        "timestamp": "1534516503816",
        "success": "true",
        "load_time": "1819.985",
        "makespan": "58.841",
        "processing_time": "0.639"
      },
      "r898666": {
        "id": "r898666",
        "timestamp": "1534511454511",
        "success": "true",
        "load_time": "360.069",
        "makespan": "29.936",
        "processing_time": "10.766"
      },
      "r744733": {
        "id": "r744733",
        "timestamp": "1534512323772",
        "success": "true",
        "load_time": "51.944",
        "makespan": "4.822",
        "processing_time": "2.493"
      },
      "r786790": {
        "id": "r786790",
        "timestamp": "1534516783979",
        "success": "true",
        "load_time": "113.881",
        "makespan": "7.334",
        "processing_time": "0.094"
      },
      "r552385": {
        "id": "r552385",
        "timestamp": "1534510722763",
        "success": "true",
        "load_time": "1.402",
        "makespan": "1.484",
        "processing_time": "0.060"
      },
      "r603080": {
        "id": "r603080",
        "timestamp": "1534517056380",
        "success": "true",
        "load_time": "1.377",
        "makespan": "1.437",
        "processing_time": "0.046"
      },
      "r517723": {
        "id": "r517723",
        "timestamp": "1534510694720",
        "success": "true",
        "load_time": "1.402",
        "makespan": "1.525",
        "processing_time": "0.081"
      },
      "r575649": {
        "id": "r575649",
        "timestamp": "1534514326409",
        "success": "true",
        "load_time": "1168.597",
        "makespan": "107.404",
        "processing_time": "49.511"
      },
      "r703209": {
        "id": "r703209",
        "timestamp": "1534517127884",
        "success": "true",
        "load_time": "1.386",
        "makespan": "1.435",
        "processing_time": "0.039"
      },
      "r613261": {
        "id": "r613261",
        "timestamp": "1534511341431",
        "success": "true",
        "load_time": "360.069",
        "makespan": "21.302",
        "processing_time": "0.945"
      },
      "r771152": {
        "id": "r771152",
        "timestamp": "1534514446443",
        "success": "true",
        "load_time": "1168.597",
        "makespan": "49.644",
        "processing_time": "0.533"
      },
      "r470371": {
        "id": "r470371",
        "timestamp": "1534517098450",
        "success": "true",
        "load_time": "1.377",
        "makespan": "1.404",
        "processing_time": "0.037"
      },
      "r755973": {
        "id": "r755973",
        "timestamp": "1534510665286",
        "success": "true",
        "load_time": "6.159",
        "makespan": "1.461",
        "processing_time": "0.020"
      },
      "r585397": {
        "id": "r585397",
        "timestamp": "1534510708743",
        "success": "true",
        "load_time": "1.402",
        "makespan": "1.477",
        "processing_time": "0.039"
      },
      "r489083": {
        "id": "r489083",
        "timestamp": "1534514606499",
        "success": "true",
        "load_time": "1168.597",
        "makespan": "51.907",
        "processing_time": "2.252"
      },
      "r485181": {
        "id": "r485181",
        "timestamp": "1534512049197",
        "success": "true",
        "load_time": "521.506",
        "makespan": "20.347",
        "processing_time": "0.231"
      },
      "r808888": {
        "id": "r808888",
        "timestamp": "1534511423491",
        "success": "true",
        "load_time": "360.069",
        "makespan": "18.450",
        "processing_time": "0.744"
      },
      "r738819": {
        "id": "r738819",
        "timestamp": "1534512308750",
        "success": "true",
        "load_time": "51.944",
        "makespan": "2.303",
        "processing_time": "0.049"
      }
    }
  }
}