# Criage Website

This is a minisite for the Criage package manager, containing project information, documentation, and usage instructions.

## 🌐 Language Support

The website supports both English and Russian languages:

- **English** - Default language for international users (index.html, docs.html)
- **Russian** - Localization for Russian-speaking users (index_ru.html, docs_ru.html)

Navigation between languages is available through language switcher buttons in the header.

## 📁 Structure

### English (Default)

- `index.html` - Main page with Criage capabilities overview
- `docs.html` - Comprehensive usage documentation

### Russian

- `index_ru.html` - Главная страница с обзором возможностей Criage
- `docs_ru.html` - Подробная документация по использованию

### Common

- `README.md` - This documentation file

## ✨ Features

### Main Page (index.html / index_ru.html)

- 🎨 Modern design with gradients and animations
- 📱 Responsive layout for all devices
- 🚀 Key features section
- 💻 Installation command examples
- 📋 Interactive code blocks with copy functionality
- 📊 Project statistics
- 🌐 Language switcher

### Documentation Page (docs.html / docs_ru.html)

- 📖 Comprehensive documentation for all features
- 🔍 Sidebar navigation menu
- 📝 Step-by-step instructions
- ⚠️ Important notes and warnings
- 🛠️ Troubleshooting section
- 🌐 Language switcher

## 🎨 Design Features

- **Color Scheme**: Blue (#2563eb) as primary color
- **Typography**: Inter font for modern appearance
- **Components**: Cards with shadows, buttons with hover effects
- **Animations**: Smooth transitions and element appearance animations
- **Interactivity**: Code copy buttons, smooth scrolling, language switching

## 🚀 Running Locally

### Quick Start Options

```bash
# Using Python HTTP server
cd website
python -m http.server 8000

# Using Node.js
npx http-server . -p 8000

# Using Go
go install github.com/shurcooL/goexec@latest
goexec 'http.ListenAndServe(":8000", http.FileServer(http.Dir(".")))'
```

Then open <http://localhost:8000> in your browser.

### Deployment

The site consists of static HTML files and can be deployed on any web server:

- GitHub Pages
- Netlify
- Vercel
- Any static hosting provider

## 🔧 Technical Details

- **Languages**: HTML5, CSS3, JavaScript (ES6+)
- **Libraries**: No external dependencies
- **Compatibility**: Modern browsers (Chrome, Firefox, Safari, Edge)
- **Size**: Minimal file size
- **Performance**: Optimized styles and scripts

## 🎨 Customization

### Changing Color Scheme

Main colors are defined in CSS variables at the beginning of each file:

```css
:root {
    --primary-color: #2563eb;
    --primary-dark: #1d4ed8;
    --secondary-color: #64748b;
    /* ... */
}
```

### Adding New Sections

1. Add HTML markup to the appropriate file
2. Add styles to the `<style>` section
3. Add JavaScript functionality if needed

### Updating Content

- Command information: update sections in both language versions
- Statistics: change numbers in the `.stats` section
- Links: update GitHub and other resource links
- Translations: keep both language versions synchronized

## 🌐 Language Management

When updating content:

1. **English First**: Update English files (index.html, docs.html)
2. **Russian Translation**: Apply corresponding changes to Russian files
3. **Navigation**: Ensure language switcher links are correct
4. **Consistency**: Keep structure and functionality identical between languages

## 📝 Content Guidelines

- Keep both language versions functionally identical
- Ensure language switcher buttons work correctly
- Maintain consistent styling across all versions
- Test all interactive features in both languages

## 🤝 Support

If you have suggestions for improving the website, please create an issue in the main Criage repository.
