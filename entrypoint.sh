#!/usr/bin/env sh
set -e

echo "Running $1"
if [ "$1" = 'smd-init' ]; then
	# This directory has to exist either way, but hopefully a persistent storage is mounted here.
	mkdir -p /persistent_migrations

	# Make sure the migrations make their way to the persistent mounted storage.
	cp /migrations/*.sql /persistent_migrations/

	echo "Migrations copied to persistent location."
fi
exec "$@"