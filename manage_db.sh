#!/usr/bin/env bash
#
# manage_db.sh
# a script that makes setonotes database management easier
# or at least less painful
# hopefully
#
# Usage:
#
# Dump the current database to YYYYmmdd-HHMMSS-schema-only.psql and
# YYYYmmdd-HHMMSS-all.psql
#
#     manage_db.sh --dump
#
# Attempt to load input_db.psql
# (note that this may overwrite the current database, depending on the contents
# of input_db.psql)
#
#     manage_db.sh --load input_db.psql

# exit immediately if anything fails
set -e

case $1 in
	# if the first argument is --dump
	"--dump")
		# then proceed with dumping
		FILE_PREFIX=$(date +%Y%m%d-%H%M%S)
		echo "dumping schema to $FILE_PREFIX-schema.psql..."
		sudo -u postgres pg_dump setonotes --schema-only --clean \
			> $FILE_PREFIX-schema-only.psql
		echo "dumping data to $FILE_PREFIX-all.psql..."
		sudo -u postgres pg_dump setonotes --clean \
			> $FILE_PREFIX-all.psql
		echo "dumped schema and data to $FILE_PREFIX-{schema,all}.psql"
		;; # no fallthrough to other cases
	# if the first argument is --load
	"--load")
		# then attempt to load whatever's in the second argument
		echo "loading from $2..."
		sudo -u postgres psql setonotes < $2
		echo "loaded $2 successfully"
		;;
	# if the first argument isn't --dump or --load
	*)
		# then show usage and exit
		echo "usage: manage_db.sh <--dump | --load <input_psql>>"
		exit 1
		;;
esac

exit 0
