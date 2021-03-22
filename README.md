# Simple WebDAV Server

This is a tiny implementation of a user-scoped WebDAV server.

1. Authentication through .htpasswd file
2. All users have their own root directories
3. No databases required


## Installation

To install the server you simply download the binary.

### Linux

```bash
# Download the binary
wget -O webdav https://github.com/ChappIO/simple-webdav/releases/latest/download/webdav_linux_amd64
# Move it into a more appropriate place
sudo mv webdav /usr/local/bin/webdav
# Test
webdav --version
```
