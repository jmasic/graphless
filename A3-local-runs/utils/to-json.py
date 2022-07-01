import json
import os
import sys
import zlib


VERTEX_FILE_SUFFIX = '.v'
EDGE_FILE_SUFFIX = '.e'
DEFAULT_CHUNKS = 15
DEFAULT_VERTEX_FILE = "dota-league/dota-league.v"
DEFAULT_DIRECTED_ARG = "directed"

vertex_file = DEFAULT_VERTEX_FILE
chunks = DEFAULT_CHUNKS
directed_arg = DEFAULT_DIRECTED_ARG
if len(sys.argv) > 1:
    vertex_file = sys.argv[1]
if len(sys.argv) > 2:
    chunks = int(sys.argv[2])
if len(sys.argv) > 3: # NOTE: This should actually be a mandatory argument!
    directed_arg = sys.argv[3]

IS_DIRECTED = directed_arg == 'directed'
if vertex_file[-2:] != VERTEX_FILE_SUFFIX:
    print(f'Expecting \'{vertex_file}\' to finish with {VERTEX_FILE_SUFFIX}')
    sys.exit(-1)
graph_name = vertex_file.split('/')[-1][:-2]
graph_folder = '/'.join(vertex_file.split('/')[:-1]) + '/'
output_folder = graph_folder + 's3/'

print(f'Using \'{graph_folder}\' as folder, \'{vertex_file}\' as vertex file.')
print(f'The output for \'{graph_name}\', composed of {chunks} chunk(s), will be in \'{output_folder}\'.')

os.system(f'mkdir -p {output_folder}')

def compress_with_zlib(file_content):
    b = bytearray()
    b.extend(map(ord, file_content))
    return zlib.compress(b)


# [{ "i": 1, "e": [{"i": 2}] }, { "i": 2, "e": [{"i": 3}] }, { "i": 3, "e": [{"i": 1}] }]
# [{ "i": 1, "e": [{"i": 2, "v": 0.54}] }, { "i": 2, "e": [{"i": 3, "v": 0.54}] }, { "i": 3, "e": [{"i": 1, "v": 0.54}] }]

vertices = []
vertex_by_id = {}

with open(graph_folder + graph_name + VERTEX_FILE_SUFFIX, "r") as vertex_file:
    for vertex_line in vertex_file:
        vertex_id = int(vertex_line.strip())
        vertex = { 'i': vertex_id, 'e': [] }
        vertices.append(vertex)
        vertex_by_id[vertex_id] = vertex

edges_count = 0
with open(graph_folder + graph_name + EDGE_FILE_SUFFIX, "r") as edges_file:
    for edge_line in edges_file:
        if edges_count % 10000 == 0:
            print(f'Reached line {edges_count}')
        edge_tokens = edge_line.split(" ")
        id_1 = int(edge_tokens[0])
        id_2 = int(edge_tokens[1])
        edge = { 'i': id_2 }
        if len(edge_tokens) > 2:
            weight = float(edge_tokens[2])
            edge['v'] = weight
        vertex_by_id[id_1]['e'].append(edge)
        if not IS_DIRECTED:
            reverse_edge = { 'i': id_1 }
            if len(edge_tokens) > 2:
                weight = float(edge_tokens[2])
                reverse_edge['v'] = weight
            vertex_by_id[id_2]['e'].append(reverse_edge)
        edges_count += 1


with open(output_folder + 'graphFileKey-dota-league-properties', 'w') as properties_file:
    properties_obj = { "numberOfVertices": len(vertices), "numberOfEdges": edges_count, "numberOfBuckets": chunks }
    json_properties = json.dumps(properties_obj)
    properties_file.write(json_properties)


n = len(vertices)
chunk_length = (n // chunks) + 1

for i in range(0, chunks):
    print(f'Writing chunk {i}...')
    with open(output_folder + f'graphFileKey-dota-league-{i}', 'wb') as output_file:
        output_file.truncate(0)
        chunk_start = i * chunk_length
        chunk_end = chunk_start + chunk_length
        vertices_in_chunk = vertices[chunk_start:chunk_end]
        json_output = json.dumps(vertices_in_chunk)
        compressed_json = compress_with_zlib(json_output)
        output_file.write(compressed_json)


print(f'Finished preparing input files for \'{graph_name}\'.')
