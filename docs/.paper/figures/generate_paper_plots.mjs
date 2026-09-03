// docs/.paper/figures/generate_paper_plots.mjs
import puppeteer from 'puppeteer';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const executablePath = process.env.PUPPETEER_EXECUTABLE_PATH || '/run/current-system/sw/bin/chromium';

async function renderHTMLFileToPNG(htmlFilePath, outPngPath, width = 1240, height = 600) {
  const browser = await puppeteer.launch({
    headless: true,
    executablePath: fs.existsSync(executablePath) ? executablePath : undefined,
    args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-dev-shm-usage', '--disable-gpu']
  });
  const page = await browser.newPage();
  await page.setViewport({ width, height, deviceScaleFactor: 2 });
  
  await page.goto(`file://${path.resolve(htmlFilePath)}`, { waitUntil: 'networkidle0' });

  const element = await page.$('svg');
  if (element) {
    await element.screenshot({ path: outPngPath });
  } else {
    await page.screenshot({ path: outPngPath });
  }
  await browser.close();
  console.log(`Generated: ${outPngPath}`);
}

async function main() {
  const figures = [
    { html: 'fig1_workflow.html', png: 'fig1_workflow.png', w: 1240, h: 600 },
    { html: 'fig2_empirical_results.html', png: 'fig2_empirical_results.png', w: 1260, h: 620 },
    { html: 'fig3_dataflow.html', png: 'fig3_dataflow.png', w: 1240, h: 500 },
    { html: 'fig4_lifecycle.html', png: 'fig4_lifecycle.png', w: 1240, h: 510 },
  ];

  for (const fig of figures) {
    const htmlPath = path.join(__dirname, fig.html);
    const pngPath = path.join(__dirname, fig.png);
    if (fs.existsSync(htmlPath)) {
      await renderHTMLFileToPNG(htmlPath, pngPath, fig.w, fig.h);
    } else {
      console.warn(`Missing: ${htmlPath}`);
    }
  }
  console.log('All publication figures successfully generated with high-DPI serif typography.');
}

main().catch(console.error);
