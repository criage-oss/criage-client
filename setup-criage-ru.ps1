# Quick setup script for criage.ru domain
# Run this script to deploy custom domain configuration

Write-Host "🌐 Setting up criage.ru domain for GitHub Pages..." -ForegroundColor Green

# Check if we're in the correct directory
if (-not (Test-Path "website")) {
    Write-Error "❌ Please run this script from the criage project root directory"
    exit 1
}

# Verify CNAME file exists
if (-not (Test-Path "website/CNAME")) {
    Write-Error "❌ CNAME file not found in website/ directory"
    exit 1
}

# Check CNAME content
$cnameContent = Get-Content "website/CNAME" -Raw
if ($cnameContent.Trim() -ne "criage.ru") {
    Write-Warning "⚠️ CNAME content: '$($cnameContent.Trim())' - should be 'criage.ru'"
}

Write-Host "📁 Files ready for deployment:" -ForegroundColor Yellow
Write-Host "  ✅ website/index.html (English homepage)"
Write-Host "  ✅ website/index_ru.html (Russian homepage)"  
Write-Host "  ✅ website/docs.html (English docs)"
Write-Host "  ✅ website/docs_ru.html (Russian docs)"
Write-Host "  ✅ website/logo.png (Logo)"
Write-Host "  ✅ website/CNAME (Domain: criage.ru)"

Write-Host ""
Write-Host "🚀 Deploying to GitHub..." -ForegroundColor Green

# Add and commit changes
git add website/CNAME
git add website/*.html
git add website/*.png
git commit -m "Add custom domain criage.ru and website files"

# Push to GitHub
git push origin main

Write-Host ""
Write-Host "✅ Deployment initiated!" -ForegroundColor Green
Write-Host ""
Write-Host "📋 Next steps in Cloudflare:" -ForegroundColor Yellow
Write-Host "1. Add A records for criage.ru:"
Write-Host "   - 185.199.108.153"
Write-Host "   - 185.199.109.153" 
Write-Host "   - 185.199.110.153"
Write-Host "   - 185.199.111.153"
Write-Host ""
Write-Host "2. Set Proxy status to ☁️ Proxied"
Write-Host ""
Write-Host "📋 Next steps in GitHub:" -ForegroundColor Yellow
Write-Host "1. Go to repository Settings → Pages"
Write-Host "2. Set Custom domain to: criage.ru"
Write-Host "3. Enable 'Enforce HTTPS' after domain verification"
Write-Host ""
Write-Host "🌐 Your website will be available at: https://criage.ru" -ForegroundColor Green
Write-Host "⏱️ DNS propagation may take 1-48 hours" -ForegroundColor Cyan
Write-Host ""
Write-Host "📚 For detailed instructions, see: CLOUDFLARE_SETUP.md" -ForegroundColor Blue 