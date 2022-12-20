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

# Create directory for extracted shiv binaries
mkdir -p /var/lib/digimontcg
chown digimontcg:digimontcg /var/lib/digimontcg
chmod 0750 /var/lib/digimontcg

# Reload systemd to pickup services
systemctl daemon-reload

# Start or restart components
components="digimontcg"
for component in $components; do
    if systemctl is-enabled $component >/dev/null
    then
        systemctl restart $component >/dev/null
    fi
done
