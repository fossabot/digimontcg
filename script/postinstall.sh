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

# Update config file permissions (idempotent)
chown root:digimontcg /etc/digimontcg.conf
chmod 0640 /etc/digimontcg.conf

# Reload systemd to pickup services
systemctl daemon-reload
