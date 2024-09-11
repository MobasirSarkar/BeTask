#!/bin/sh

air -c .air.toml &

cd /app/client
pnpm run dev
