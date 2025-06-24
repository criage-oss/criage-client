# 🌍 Criage Localization System

Criage features a dynamic localization system that supports multiple languages without requiring code modifications. The system automatically detects available languages by scanning translation files and allows users to add new localizations simply by creating JSON files.

## 🗂️ File Structure

```
criage/
├── locale/                           # Main application localization
│   ├── translations_en.json         # English
│   ├── translations_ru.json         # Russian
│   ├── translations_de.json         # German
│   └── translations_es.json         # Spanish
└── repository/
    └── locale/                       # Repository server localization
        ├── translations_en.json     # English
        ├── translations_ru.json     # Russian
        ├── translations_de.json     # German
        └── translations_es.json     # Spanish
```

## 🚀 Features

### Dynamic Language Detection

- **Automatic scanning** of `translations_*.json` files at startup
- **No hardcoded languages** in the source code
- **Hot reloading** - restart the application to pick up new languages
- **Thread-safe** implementation with mutex protection

### Supported Language Codes

The system supports standard **ISO 639-1** language codes:

#### Basic Format

- `en` - English
- `ru` - Russian  
- `de` - German
- `fr` - French
- `es` - Spanish
- `zh` - Chinese
- `ja` - Japanese
- `pt` - Portuguese
- `it` - Italian
- `ko` - Korean

#### Extended Format (with regions)

- `en-US` - English (United States)
- `en-GB` - English (United Kingdom)
- `zh-CN` - Chinese (Simplified)
- `zh-TW` - Chinese (Traditional)
- `fr-CA` - French (Canada)
- `es-MX` - Spanish (Mexico)

### Fallback System

1. **Primary**: Selected or detected language
2. **Secondary**: English (`en`) if available
3. **Tertiary**: First available language in the list
4. **Ultimate**: Hardcoded strings in code

## 📋 Adding New Languages

### Step 1: Create Translation Files

Create files named `translations_<language_code>.json` in the appropriate directories:

```bash
# For main application
locale/translations_fr.json        # French
locale/translations_zh.json        # Chinese
locale/translations_ja.json        # Japanese

# For repository server
repository/locale/translations_fr.json  # French
repository/locale/translations_zh.json  # Chinese
repository/locale/translations_ja.json  # Japanese
```

### Step 2: Add Translations

#### Main Application Template (`locale/translations_<lang>.json`)

```json
{
  "app_description": "High-performance package manager",
  "cmd_install": "Install package", 
  "cmd_uninstall": "Uninstall package",
  "cmd_update": "Update package",
  "cmd_search": "Search packages",
  "cmd_list": "List packages",
  "cmd_info": "Package information",
  "cmd_build": "Build package",
  "cmd_create": "Create package",
  "cmd_publish": "Publish package",
  "cmd_config": "Configuration settings",
  "installing_package": "Installing package %s...",
  "uninstalling_package": "Uninstalling package %s...",
  "updating_package": "Updating package %s...",
  "building_package": "Building package %s...",
  "searching_packages": "Searching packages...",
  "packages_found": "Found %d packages:",
  "no_packages_found": "No packages found matching your query.",
  "package_name": "Name",
  "package_version": "Version",
  "package_description": "Description",
  "package_author": "Author",
  "package_installed": "Package %s installed successfully!",
  "package_uninstalled": "Package %s uninstalled successfully!",
  "package_updated": "Package %s updated successfully!",
  "package_not_found": "Package %s not found.",
  "dependency_installing": "Installing dependency: %s",
  "dependency_failed": "Failed to install dependency: %s"
}
```

#### Repository Server Template (`repository/locale/translations_<lang>.json`)

```json
{
  "repository_name": "Criage Package Repository",
  "server_started": "Server started on http://localhost:%d",
  "server_stopping": "Stopping server...",
  "package_uploaded": "Package uploaded",
  "package_download": "Package download",
  "upload_successful": "Package uploaded successfully",
  "upload_failed": "Failed to upload package",
  "invalid_package": "Invalid package format",
  "package_exists": "Package already exists",
  "package_not_found": "Package not found",
  "api_list_packages": "Listing packages",
  "api_package_info": "Getting package information",
  "api_download_package": "Downloading package",
  "api_upload_package": "Uploading package",
  "storage_error": "Storage error",
  "validation_error": "Validation error",
  "authentication_required": "Authentication required",
  "permission_denied": "Permission denied"
}
```

### Step 3: Restart Application

After adding translation files, restart the application:

```bash
# Check available languages
go run test_localization.go
# Output: Supported languages: [de en es fr ru]
```

## 💡 Usage Examples

### Adding Portuguese Language

1. Create `locale/translations_pt.json`:

```json
{
  "app_description": "Gerenciador de pacotes de alto desempenho",
  "cmd_install": "Instalar pacote",
  "cmd_search": "Pesquisar pacotes",
  "installing_package": "Instalando pacote %s...",
  "packages_found": "%d pacotes encontrados:"
}
```

2. Create `repository/locale/translations_pt.json`:

```json
{
  "repository_name": "Repositório de Pacotes Criage",
  "server_started": "Servidor iniciado em http://localhost:%d",
  "package_uploaded": "Pacote enviado",
  "upload_successful": "Pacote enviado com sucesso"
}
```

3. Restart the application:

```bash
go run test_locale_structure.go
# Output: Supported languages: [de en es pt ru]
```

### Setting Language Programmatically

```go
package main

import (
    "criage/pkg"
    "fmt"
)

func main() {
    loc := pkg.GetLocalization()
    
    // Get current language
    fmt.Printf("Current language: %s\n", loc.GetLanguage())
    
    // Get supported languages
    fmt.Printf("Supported: %v\n", loc.GetSupportedLanguages())
    
    // Switch language
    if err := loc.SetLanguage("de"); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    // Use translation
    fmt.Println(pkg.T("app_description"))
    // Output: Hochleistungs-Paketmanager
}
```

## 🔧 Technical Details

### Architecture

- **Singleton pattern** for global access
- **Mutex protection** for thread safety
- **Memory-efficient** lazy loading of translations
- **Regex-based** file scanning: `^translations_([a-z]{2}(?:-[A-Z]{2})?)\.json$`

### System Language Detection

The system automatically detects the user's language based on:

1. `LANGUAGE` environment variable
2. `LC_ALL` environment variable
3. `LANG` environment variable
4. Fallback to English if none match available translations

### Performance

- **Fast startup** - only scans files once during initialization
- **Low memory footprint** - only loads active language
- **Efficient lookups** - uses Go's native map for O(1) access

## 🌟 Best Practices

### Translation Quality

- **Keep keys descriptive** but concise
- **Use placeholders** for dynamic content: `"installing %s"`
- **Maintain consistency** across similar contexts
- **Test all translations** before committing

### File Organization

- ✅ Use `locale/` subdirectories for clean structure
- ✅ Follow naming convention: `translations_<code>.json`
- ✅ Keep translations synchronized between main app and repository
- ✅ Use UTF-8 encoding for all files

### Development Workflow

```bash
# 1. Add translation keys to English file first
vim locale/translations_en.json

# 2. Copy structure to new language file
cp locale/translations_en.json locale/translations_fr.json

# 3. Translate all values (keep keys in English)
vim locale/translations_fr.json

# 4. Test the new language
go run test_localization.go

# 5. Repeat for repository if needed
cp repository/locale/translations_en.json repository/locale/translations_fr.json
vim repository/locale/translations_fr.json
```

## 🤝 Contributing

We welcome contributions for new languages! To add support for your language:

1. **Fork** the repository
2. **Create** translation files for your language
3. **Test** the translations work correctly
4. **Submit** a pull request with:
   - Translation files for both main app and repository
   - Brief description of language coverage
   - Note about any special considerations

### Language Coverage Status

| Language | Code | Main App | Repository | Status |
|----------|------|----------|------------|---------|
| English | `en` | ✅ | ✅ | Complete |
| Russian | `ru` | ✅ | ✅ | Complete |
| German | `de` | ✅ | ✅ | Complete |
| Spanish | `es` | ⚠️ | ✅ | Partial |
| French | `fr` | ❌ | ❌ | Needed |
| Chinese | `zh` | ❌ | ❌ | Needed |
| Japanese | `ja` | ❌ | ❌ | Needed |

## 📚 API Reference

### Localization Object Methods

```go
// Get singleton instance
loc := pkg.GetLocalization()

// Get current language code
lang := loc.GetLanguage()

// Get list of supported languages
languages := loc.GetSupportedLanguages()

// Set language (returns error if not supported)
err := loc.SetLanguage("de")

// Translate key with current language
text := pkg.T("app_description")

// Translate with placeholders
text := pkg.T("installing_package", "example-package")
```

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `LANGUAGE` | Primary language preference | `LANGUAGE=de_DE.UTF-8` |
| `LC_ALL` | System locale override | `LC_ALL=fr_FR.UTF-8` |
| `LANG` | System language setting | `LANG=es_ES.UTF-8` |

## 🚀 Future Enhancements

- **Pluralization support** for complex grammar rules
- **Context-aware translations** for different UI contexts
- **Real-time language switching** without restart
- **Translation validation** tools for developers
- **Automatic translation updates** from translation services
- **RTL language support** for Arabic, Hebrew, etc.

---

The Criage localization system is designed to be simple, efficient, and developer-friendly. Adding support for your language takes just minutes, making Criage accessible to users worldwide! 🌍
