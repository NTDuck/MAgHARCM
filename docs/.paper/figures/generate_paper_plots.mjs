// docs/.paper/figures/generate_paper_plots.mjs
import puppeteer from 'puppeteer';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const executablePath = process.env.PUPPETEER_EXECUTABLE_PATH || '/run/current-system/sw/bin/chromium';

async function renderHTMLFileToPNG(htmlFilePath, outPngPath, width = 1280, height = 620) {
  const browser = await puppeteer.launch({
    headless: true,
    executablePath: fs.existsSync(executablePath) ? executablePath : undefined,
    args: [
      '--no-sandbox',
      '--disable-setuid-sandbox',
      '--disable-dev-shm-usage',
      '--disable-gpu'
    ]
  });
  const page = await browser.newPage();
  await page.setViewport({ width, height, deviceScaleFactor: 2.5 });
  
  // Emulate light color scheme for paper publication
  await page.emulateMediaFeatures([{ name: 'prefers-color-scheme', value: 'light' }]);

  await page.goto(`file://${path.resolve(htmlFilePath)}`, { waitUntil: 'networkidle0' });

  // Remove all viewer chrome, toolbars, buttons, headers, export menus, HUDs, lens, zoom controls
  await page.evaluate(() => {
    document.documentElement.setAttribute('data-theme', 'light');
    document.documentElement.classList.remove('dark');
    document.documentElement.classList.add('light');
    
    document.body.style.backgroundColor = '#ffffff';
    document.body.style.color = '#0f172a';
    document.body.style.margin = '0';
    document.body.style.padding = '0';

    // Remove all HUDs, toolbars, lens indicators, zoom controls, and menus
    const removeSelectors = [
      'header',
      'nav',
      'footer',
      '.viewer-toolbar',
      '.viewer-chrome',
      '.theme-toggle',
      '.nav-controls',
      '.chapter-rail',
      '.chapter-delta-preview',
      '[data-viewer-toolbar]',
      '[data-viewer-hud]',
      '.viewer-hud',
      '.nav-hud',
      '.viewport-hud',
      '.zoom-pill',
      '.lens-pill',
      '.export-btn',
      '.present-btn',
      '.viewer-header',
      '.toolbar',
      '.controls-bar',
      '.hud',
      '[class*="hud"]',
      '[class*="lens"]',
      '[class*="pill"]',
      '[class*="zoom"]',
      '[class*="toolbar"]'
    ];
    
    removeSelectors.forEach(sel => {
      document.querySelectorAll(sel).forEach(el => el.remove());
    });

    const svg = document.querySelector('svg');
    if (svg) {
      svg.style.backgroundColor = '#ffffff';
    }
  });

  // Wait for layout settlement
  await new Promise(r => setTimeout(r, 250));

  // Find the primary svg element and capture its exact bounding box
  const svgEl = await page.$('svg.diagram-svg') || await page.$('svg');
  if (svgEl) {
    await svgEl.screenshot({ path: outPngPath, omitBackground: false });
  } else {
    await page.screenshot({ path: outPngPath, omitBackground: false });
  }
  
  await browser.close();
  console.log(`Generated publication PNG (clean, no HUD): ${outPngPath}`);
}

async function main() {
  const figures = [
    { html: 'fig1_workflow.html', png: 'fig1_workflow.png', w: 1280, h: 580 },
    { html: 'fig2_empirical_results.html', png: 'fig2_empirical_results.png', w: 1280, h: 580 },
    { html: 'fig3_dataflow.html', png: 'fig3_dataflow.png', w: 1280, h: 500 },
    { html: 'fig4_lifecycle.html', png: 'fig4_lifecycle.png', w: 1280, h: 520 },
  ];

  for (const fig of figures) {
    const htmlPath = path.join(__dirname, fig.html);
    const pngPath = path.join(__dirname, fig.png);
    if (fs.existsSync(htmlPath)) {
      await renderHTMLFileToPNG(htmlPath, pngPath, fig.w, fig.h);
    } else {
      console.warn(`Missing HTML file: ${htmlPath}`);
    }
  }
  console.log('All publication figures rendered cleanly with zero HUDs/menus.');
}

main().catch(console.error);
