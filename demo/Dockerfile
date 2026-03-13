FROM ubuntu:24.04

# Install dependencies
RUN apt-get update && apt-get install -y \
    vim \
    curl \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# Copy local squix binary with table alignment fix
COPY squix /usr/local/bin/squix
RUN chmod +x /usr/local/bin/squix

# Download presenterm v0.16.1
RUN mkdir -p /tmp/presenterm && \
    curl -L https://github.com/mfontanini/presenterm/releases/download/v0.16.1/presenterm-0.16.1-x86_64-unknown-linux-gnu.tar.gz | tar xz -C /tmp/presenterm && \
    find /tmp/presenterm -type f -name presenterm -exec mv {} /usr/local/bin/ \; && \
    chmod +x /usr/local/bin/presenterm && \
    rm -rf /tmp/presenterm

# Download gotty v1.6.0
RUN curl -L https://github.com/sorenisanerd/gotty/releases/download/v1.6.0/gotty_v1.6.0_linux_amd64.tar.gz | tar xz && \
    mv gotty /usr/local/bin/ && \
    chmod +x /usr/local/bin/gotty

WORKDIR /root/tree
ENV EDITOR=vim

# Copy custom gotty assets with Fira Code font
COPY gotty-assets /root/.gotty-assets

# Copy bashrc with welcome message and demo alias
COPY .bashrc /root/.bashrc

# Copy pre-populated squix config
COPY squix-config /root/.config/squix

# Copy entrypoint script
COPY entrypoint.sh /root/entrypoint.sh
RUN chmod +x /root/entrypoint.sh

EXPOSE 8080

# Use entrypoint script to initialize database
ENTRYPOINT ["/root/entrypoint.sh"]

# Default to interactive bash shell via gotty with custom index
CMD ["--permit-write", "-p", "8080", "--index", "/root/.gotty-assets/index.html", "gotty", "bash"]
