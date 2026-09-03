// docs/.paper/figures/generate_paper_plots.mjs
import puppeteer from 'puppeteer';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const executablePath = process.env.PUPPETEER_EXECUTABLE_PATH || '/run/current-system/sw/bin/chromium';

async function renderHTMLFileToPNG(htmlFilePath, outPngPath, width = 1200, height = 540) {
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

  // Force light theme and clean white background on viewer DOM
  await page.evaluate(() => {
    document.documentElement.setAttribute('data-theme', 'light');
    document.documentElement.classList.remove('dark');
    document.documentElement.classList.add('light');
    document.body.style.backgroundColor = '#ffffff';
    document.body.style.color = '#0f172a';
    
    // If there is an archify theme button, click light if needed
    const lightBtn = document.querySelector('[data-theme-value="light"], [aria-label*="Light"], .theme-toggle');
    if (lightBtn && document.documentElement.getAttribute('data-theme') !== 'light') {
      lightBtn.click();
    }
  });

  // Wait for layout settlement
  await new Promise(r => setTimeout(r, 200));

  // Find main svg or root container
  const svgEl = await page.$('svg');
  if (svgEl) {
    await svgEl.screenshot({ path: outPngPath, omitBackground: false });
  } else {
    await page.screenshot({ path: outPngPath, omitBackground: false });
  }
  
  await browser.close();
  console.log(`Generated light-theme publication PNG: ${outPngPath}`);
}

async function main() {
  const figures = [
    { html: 'fig1_workflow.html', png: 'fig1_workflow.png', w: 1200, h: 560 },
    { html: 'fig2_empirical_results.html', png: 'fig2_empirical_results.png', w: 1200, h: 560 },
    { html: 'fig3_dataflow.html', png: 'fig3_dataflow.png', w: 1200, h: 480 },
    { html: 'fig4_lifecycle.html', png: 'fig4_lifecycle.png', w: 1200, h: 500 },
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
  console.log('All publication figures rendered with clean white background and sharp typography.');
}

main().catch(console.error);
