FROM ubuntu:24.04

# Install dependencies
RUN apt-get update && apt-get install -y \
    vim \
    curl \
    ca-certificates \
    sqlite3 \
    cron \
    && rm -rf /var/lib/apt/lists/*

# Copy local squix binary (v0.4.0)
COPY squix /usr/local/bin/squix
RUN chmod +x /usr/local/bin/squix

# Download presenterm v0.16.1
RUN mkdir -p /tmp/presenterm && \
    curl -L https://github.com/mfontanini/presenterm/releases/download/v0.16.1/presenterm-0.16.1-x86_64-unknown-linux-gnu.tar.gz | tar xz -C /tmp/presenterm && \
    find /tmp/presenterm -type f -name presenterm -exec mv {} /usr/local/bin/ \; && \
    chmod +x /usr/local/bin/presenterm && \
    rm -rf /tmp/presenterm

# Download ttyd v1.7.4
RUN curl -L https://github.com/tsl0922/ttyd/releases/download/1.7.4/ttyd.x86_64 -o /usr/local/bin/ttyd && \
    chmod +x /usr/local/bin/ttyd

# Create restricted demo user with rbash
RUN useradd -m -s /bin/rbash scrat && \
    usermod -aG sudo scrat

# Create directory for allowed commands
RUN mkdir -p /home/scrat/bin

# Create symlinks to ONLY allowed commands
RUN ln -s /usr/local/bin/squix /home/scrat/bin/squix && \
    ln -s /usr/local/bin/presenterm /home/scrat/bin/presenterm && \
    ln -s /usr/local/bin/presenterm /home/scrat/bin/demo && \
    ln -s /usr/bin/sqlite3 /home/scrat/bin/sqlite3 && \
    ln -s /usr/bin/ls /home/scrat/bin/ls && \
    ln -s /usr/bin/vim /home/scrat/bin/vim && \
    ln -s /usr/bin/cat /home/scrat/bin/cat && \
    ln -s /bin/bash /home/scrat/bin/bash && \
    chown -R scrat:scrat /home/scrat/bin

# Create cron job for hourly database reset (using sqlite3 backup to avoid cp command issues)
RUN echo '0 * * * * root sqlite3 /root/tree-templates/squirrels.db ".backup /home/scrat/tree/squirrels.db" && chmod 666 /home/scrat/tree/squirrels.db && cp /root/.config/squix/config.yaml /home/scrat/.config/squix/config.yaml && echo "Database and config reset at $(date)" >> /var/log/db-reset.log' > /etc/cron.d/squix-reset && \
    chmod 0644 /etc/cron.d/squix-reset && \
    crontab /etc/cron.d/squix-reset

# Create reset-demo script using sqlite3 (available in restricted PATH)
RUN echo '#!/bin/bash\n\
echo "Resetting demo database..."\n\
sqlite3 /root/tree-templates/squirrels.db ".backup /home/scrat/tree/squirrels.db"\n\
cat /root/.config/squix/config.yaml > /home/scrat/.config/squix/config.yaml\n\
echo "Database and config reset complete! Starting fresh."' > /home/scrat/bin/reset-demo && \
    chmod +x /home/scrat/bin/reset-demo && \
    chown scrat:scrat /home/scrat/bin/reset-demo

WORKDIR /root/tree
ENV EDITOR=vim
ENV TERM=xterm-256color

# Copy bashrc with welcome message and demo alias (for both root and scrat)
COPY .bashrc /root/.bashrc
COPY .bashrc /home/scrat/.bashrc
RUN chown scrat:scrat /home/scrat/.bashrc && \
    echo 'export PATH=/home/scrat/bin' >> /home/scrat/.bashrc

# Copy workspace templates for restoration
COPY workspace /root/tree-templates

# Copy pre-populated squix config as root (read-only)
COPY squix-config /root/.config/squix

# Make configs read-only for non-root users
RUN chmod 755 /root && \
    chmod 644 /root/.bashrc && \
    chmod 755 /root/.config/squix && \
    find /root/.config/squix -type f -exec chmod 644 {} \; && \
    chown -R root:root /root/.config/squix /root/.bashrc

# Create tree directory for scrat user
RUN mkdir -p /home/scrat/tree && \
    chown -R scrat:scrat /home/scrat/tree

# Copy entrypoint script
COPY entrypoint.sh /root/entrypoint.sh
RUN chmod +x /root/entrypoint.sh

EXPOSE 8080

# Use entrypoint script to initialize database
ENTRYPOINT ["/root/entrypoint.sh"]

# Default to interactive bash shell via ttyd with login shell for proper initialization
CMD ["ttyd", "--port", "8080", "--writable", "-t", "fontSize=14", "su", "-", "scrat"]
