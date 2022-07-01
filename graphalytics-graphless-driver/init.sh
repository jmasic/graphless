#/bin/sh

set -eu
set -o pipefail

GRAPHS_DIR=${1:-~/graphalytics/graphs}
GRAPHLESS_PLATFORM=${2:-Aws}
SKIP_COMPILE=${3:-compile}

GRAPHALYTICS_VERSION=1.4.0
PROJECT_VERSION=0.1-SNAPSHOT
PROJECT=graphalytics-$GRAPHALYTICS_VERSION-graphless-$PROJECT_VERSION
GRAPHLESS_DIR="~/graphalytics/graphless"

GRAPHLESS_PLATFORM=`echo "$GRAPHLESS_PLATFORM" | awk '{print tolower($0)}'`
GRAPHLESS_PLATFORM="$(tr '[:lower:]' '[:upper:]' <<< ${GRAPHLESS_PLATFORM:0:1})${GRAPHLESS_PLATFORM:1}"

if [ "$SKIP_COMPILE" == "compile" ]; then
  rm -rf $PROJECT
  mvn clean package -Dgraphless.environment=$GRAPHLESS_PLATFORM -Dgraphless.directory=$GRAPHLESS_DIR
  tar xf $PROJECT-bin.tar.gz
fi
cd $PROJECT/

cp -r config-template config
# set directories
sed -i.bkp "s|^.*graphs.root-directory =.*$|graphs.root-directory = $GRAPHS_DIR|g" config/benchmark.properties
sed -i.bkp "s|^.*graphs.validation-directory =.*$|graphs.validation-directory = $GRAPHS_DIR|g" config/benchmark.properties

bin/sh/run-benchmark.sh
