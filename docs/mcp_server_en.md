# 🤖 Criage MCP Server

Documentation for MCP server that integrates Criage with AI tools through the Model Context Protocol.

## Overview

Criage MCP Server provides full package manager functionality through the MCP (Model Context Protocol) for integration with AI tools.

### What is MCP Server?

MCP Server allows AI assistants to directly interact with Criage for:

- 🔍 **Searching and installing packages** - AI can find suitable packages and install them
- 📦 **Managing dependencies** - automatic conflict resolution and updates
- 🆕 **Creating new packages** - generating structure and configuring manifests
- 🔨 **Building and publishing packages** - automating the release process
- 📊 **Getting information** - analyzing package and repository status

## Installation

### Requirements

- Go 1.21 or higher
- Installed Criage
- Claude Desktop or other MCP client

### Building the server

```bash
# Navigate to mcp-server directory
cd mcp-server

# Install dependencies
go mod tidy

# Build server
go build -o criage-mcp-server .
```

## Claude Desktop Integration

### Windows

Edit the file `%APPDATA%\Claude\config.json`:

```json
{
  "mcpServers": {
    "criage": {
      "command": "C:\\path\\to\\criage-mcp-server.exe",
      "args": [],
      "env": {}
    }
  }
}
```

### Linux/macOS

Edit the file `~/.config/claude-desktop/config.json`:

```json
{
  "mcpServers": {
    "criage": {
      "command": "/path/to/criage-mcp-server",
      "args": [],
      "env": {}
    }
  }
}
```

## Available Tools

### 📦 Package Management

- `install_package` - Install package from repository
- `uninstall_package` - Remove installed package
- `update_package` - Update package to latest version
- `list_packages` - List all installed packages

### 🔍 Search and Exploration

- `search_packages` - Search packages by keywords
- `package_info` - Detailed package information
- `repository_info` - Repository information

### 🛠️ Development

- `create_package` - Create new package
- `build_package` - Build package
- `publish_package` - Publish package to repository

## Usage Examples

### Search and install package

**User request:** "Find a JSON package and install it"

**AI actions:**
1. Uses `search_packages` with "JSON" query
2. Analyzes results and selects suitable package
3. Uses `install_package` for installation
4. Confirms successful installation

### Create new project

**User request:** "Create a new package for REST API client"

**AI actions:**
1. Uses `create_package` with appropriate parameters
2. Configures manifest with description and dependencies
3. Suggests installing necessary development packages
4. Explains the created project structure

## Configuration

MCP Server uses standard Criage configuration from `~/.criage/config.json`:

```json
{
  "repositories": [
    {
      "name": "default",
      "url": "http://localhost:8080", 
      "priority": 1,
      "enabled": true
    }
  ],
  "global_path": "~/.criage/packages",
  "local_path": "./criage_modules",
  "cache_path": "~/.criage/cache",
  "timeout": 30
}
```

### Debugging

```bash
# Enable verbose logging
export CRIAGE_DEBUG=1

# Run server
./criage-mcp-server
```

## Troubleshooting

### Server won't start

**Solutions:**
- Check Go installation: `go version`
- Ensure correct path in config.json
- Check file permissions for server executable
- Restart Claude Desktop

### Tools not showing

**Solutions:**
- Check JSON syntax in configuration
- Ensure server path is correct
- Restart Claude Desktop
- Check logs for errors

## Security

⚠️ **Warning:** MCP Server runs with full user privileges. Make sure you trust the AI system before providing access to package management.

## Compatibility

✅ MCP Server is fully compatible with the main Criage client and supports all archive formats and operations.

---

[← Back to documentation](docs.html) | [GitHub](https://github.com/Zu-Krein/criage) 