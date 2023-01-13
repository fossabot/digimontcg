#!/bin/sh
set -e

# Create digimontcg group (if it doesn't exist)
if ! getent group digimontcg >/dev/null; then
    groupadd --system digimontcg
fi

# Create digimontcg user (if it doesn't exist)
if ! getent passwd digimontcg >/dev/null; then
    useradd                        \
        --system                   \
        --gid digimontcg           \
        --shell /usr/sbin/nologin  \
        digimontcg
fi

# Reload systemd to pickup services
systemctl daemon-reload

# Restart components (if enabled)
components="digimontcg"
for component in $components; do
    if systemctl is-enabled $component >/dev/null
    then
        systemctl restart $component >/dev/null
    fi
done
