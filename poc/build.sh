#!/usr/bin/env bash

build_dir=$1
tar_file=${build_dir}.tar.gz

cd poc
pwd

echo "cleaning up previously generated content"

if [[ -d ${build_dir} ]]; then
    echo "removing existing build dir"
    rm -rf ${build_dir};
fi

if [[ -d ${tar_file} ]]; then
    rm -rf ${tar_file};
fi

mkdir ${build_dir}
mkdir ${build_dir}/bin

cp -r bin ${build_dir}
cp Dockerfile.recipe ${build_dir}
cp Dockerfile.stub ${build_dir}
cp docker-compose.yml ${build_dir}

echo "tarring poc docker-compose artifacts"
tar -cvzf ../${tar_file} ${build_dir}

echo "cleaning up"
rm -rf ${build_dir}

echo "successfully build poc artifacts tar: ${tar_file}"