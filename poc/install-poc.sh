#!/usr/bin/env bash

work_dir=docdb-docker-compose
tar_file=${work_dir}.tar.gz
containers=(doc_db_api_stub doc_db_recipe_api)

run() {
  cd ${work_dir}

  echo "starting up docker-compose"
  docker-compose up -d --build
}

unpack_binaries() {
  if [[ -d ${work_dir} ]]; then
    echo "removing existing install"
    rm -rf ${work_dir};
  fi

  echo "untarring artifact"
  tar -xvf ${tar_file}

  echo "copying cert to bin dir"
  cp rds-combined-ca-bundle.pem ${work_dir}/bin

  echo "deleting far file ${tar_file}"
  rm ${tar_file}
}

stop_containers() {
  echo "stopping running instance"
  docker-compose stop

  for c in ${containers[@]}; do
    echo "stopping ${c} container"
    docker stop ${c} || true

    echo "removing ${c} container"
    docker rm ${c} || true
  done
}

stop_containers
unpack_binaries
run