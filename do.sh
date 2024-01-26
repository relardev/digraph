#!/bin/bash
set -e

echo "---------------- flattening ----------------"

PROJECT=$1
ENTRY=$2
BUILD=$3
URL=$4

go run cmd/flat/main.go \
	--config=_projects/$PROJECT/config.json \
	--entrypoint=$ENTRY \
	--basepath=/home/user/workspace/go-di-graph/_projects/$PROJECT/ \
	--buildpackage=$BUILD \
	--flatpackage=flattened

echo "---------------- depcalc ----------------"

go run cmd/depcalc/main.go \
	--entrypoint=$ENTRY \
	--path=/home/user/workspace/go-di-graph/_projects/$PROJECT/flattened > _projects/$PROJECT/deptree.json

echo "---------------- enhancer ----------------"

go run cmd/enhancer/main.go \
	--config_path=_projects/$PROJECT/config.json \
	--tree_path=_projects/$PROJECT/deptree.json \
	--project_name=$PROJECT \
	--base_url=$URL > _projects/$PROJECT/enhanced.json 

echo "---------------- render ----------------"

# go run cmd/render-d2/main.go --graph_path=_projects/$PROJECT/enhanced.json > render.d2
go run cmd/render-html/main.go --graph_path=_projects/$PROJECT/enhanced.json > render.html
