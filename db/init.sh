#!/bin/bash

scriptDir="$(dirname $0)"

runPSQL() {
  psql -U $DB_SUPER_USER -f $1 >&2
}

for directory in $scriptDir/schemas/*; do
  for filename in $directory/*.sql; do
    [ -e "$filename" ] || continue
    echo Running PSQL for $filename
    runPSQL $filename
  done
done
