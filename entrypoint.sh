#!/bin/bash

set -e

echo "Starting Docker GC Cron with schedule: ${CRON}"

# Generate crontab file
cat > /etc/crontabs/root << EOF
${CRON} /app/docker-gc-cron >> /var/log/cron.log 2>&1
EOF

# Start cron daemon
crond -f -l 2