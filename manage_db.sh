#!/usr/bin/env bash
#
# manage_db.sh
# a script that makes setonotes database management easier
# or at least less painful
# hopefully
#
# Usage:
#
# Dump the current database to output_db-schema.psql and output_db-all.psql
#
#     manage_db.sh --dump output_db
#
# Attempt to load input_db.psql
# (note that this may overwrite the current database, depending on the contents
# of input_db.psql)
#
#     manage_db.sh --load input_db.psql

# exit immediately if anything fails
set -e

case $1 in
	"--dump")
		echo "dumping schema to $2-schema.psql..."
		sudo -u postgres pg_dump setonotes --schema-only --clean \
			> $2-schema-only.psql
		echo "dumping data to $2-all.psql..."
		sudo -u postgres pg_dump setonotes --clean > $2-all.psql
		echo "dumped schema and data to $2-{schema,all}.psql"
		;;
	"--load")
		echo "loading from $2..."
		sudo -u postgres psql setonotes < $2
		echo "loaded $2 successfully"
		;;
	*)
		echo -n "usage: manage_db.sh <--load <input_psql> | "
		echo    "--dump <output_psql>>"
		exit 1
		;;
esac

exit 0
