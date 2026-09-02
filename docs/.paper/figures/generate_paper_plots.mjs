// docs/.paper/figures/generate_paper_plots.mjs
import puppeteer from 'puppeteer';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

async function renderSVGToPNG(svgHtml, outPath, width = 1200, height = 700) {
  const browser = await puppeteer.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });
  const page = await browser.newPage();
  await page.setViewport({ width, height, deviceScaleFactor: 2 });
  await page.setContent(`
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="utf-8">
      <style>
        @import url('https://fonts.googleapis.com/css2?family=Crimson+Pro:ital,wght@0,400;0,600;0,700;1,400&family=JetBrains+Mono:wght@400;500&display=swap');
        body {
          margin: 0;
          padding: 20px;
          background: #ffffff;
          font-family: 'Crimson Pro', 'Times New Roman', Times, serif;
          display: flex;
          justify-content: center;
          align-items: center;
        }
        svg {
          background: #ffffff;
        }
        text {
          font-family: 'Crimson Pro', 'Times New Roman', Times, serif;
        }
        .mono {
          font-family: 'JetBrains Mono', monospace;
        }
      </style>
    </head>
    <body>
      ${svgHtml}
    </body>
    </html>
  `, { waitUntil: 'networkidle0' });

  const element = await page.$('svg');
  if (element) {
    await element.screenshot({ path: outPath, omitBackground: false });
  } else {
    await page.screenshot({ path: outPath, fullPage: true });
  }
  await browser.close();
  console.log(`Generated: ${outPath}`);
}

// 1. Workflow Diagram (Figure 1)
const fig1SVG = `
<svg width="1100" height="520" viewBox="0 0 1100 520" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <marker id="arrow" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#1e293b"/>
    </marker>
    <marker id="arrow-blue" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#1d4ed8"/>
    </marker>
    <marker id="arrow-amber" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#b45309"/>
    </marker>
    <marker id="arrow-rose" viewBox="0 0 10 10" refX="6" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
      <path d="M 0 1.5 L 8 5 L 0 8.5 z" fill="#be123c"/>
    </marker>
    <filter id="shadow" x="-5%" y="-5%" width="110%" height="110%">
      <feDropShadow dx="0" dy="2" stdDeviation="3" flood-opacity="0.08"/>
    </filter>
  </defs>

  <!-- Background Panel -->
  <rect x="10" y="10" width="1080" height="500" rx="10" fill="#f8fafc" stroke="#cbd5e1" stroke-width="1.5" />
  
  <!-- Title / Header -->
  <text x="550" y="42" font-size="20" font-weight="700" text-anchor="middle" fill="#0f172a">Figure 1: MAgHARCM Multi-Agent Repository Translation Architecture &amp; Orchestration Workflow</text>
  <text x="550" y="64" font-size="13" font-style="italic" text-anchor="middle" fill="#475569">Four-agent topological decomposition with bounded context assembly, multi-stage validation, and compiler repair loops</text>

  <!-- Stage 1: Input Codebase -->
  <g transform="translate(35, 110)">
    <rect width="160" height="150" rx="6" fill="#ffffff" stroke="#94a3b8" stroke-width="1.5" filter="url(#shadow)" />
    <rect width="160" height="28" rx="6" fill="#f1f5f9" stroke="#94a3b8" stroke-width="1" />
    <text x="80" y="19" font-size="13" font-weight="700" text-anchor="middle" fill="#1e293b">Source Repository</text>
    <text x="14" y="52" font-size="12" fill="#334155">• Source AST ($P_s$)</text>
    <text x="14" y="74" font-size="12" fill="#334155">• Build Manifests</text>
    <text x="14" y="96" font-size="12" fill="#334155">• Test Suite ($T_s$)</text>
    <text x="14" y="118" font-size="12" fill="#334155">• Library Deps ($D_s$)</text>
    <text x="80" y="140" font-size="11" font-weight="600" text-anchor="middle" fill="#64748b" class="mono">C / Go / Java</text>
  </g>

  <!-- Arrow 1 -> Analyzer -->
  <path d="M 195 185 L 240 185" stroke="#1e293b" stroke-width="1.8" marker-end="url(#arrow)" fill="none" />

  <!-- Stage 2: Analyzer Agent -->
  <g transform="translate(245, 100)">
    <rect width="180" height="170" rx="6" fill="#ffffff" stroke="#2563eb" stroke-width="1.8" filter="url(#shadow)" />
    <rect width="180" height="30" rx="6" fill="#eff6ff" stroke="#2563eb" stroke-width="1" />
    <text x="90" y="20" font-size="14" font-weight="700" text-anchor="middle" fill="#1d4ed8">1. Analyzer Agent</text>
    <text x="14" y="55" font-size="12" fill="#1e293b">Tree-Sitter Directory Walk</text>
    <text x="14" y="77" font-size="12" fill="#1e293b">Müller Strategy Selection:</text>
    <text x="24" y="95" font-size="11" fill="#475569" class="mono">DIRECT | INCR | PILOT</text>
    <text x="14" y="118" font-size="12" fill="#1e293b">Dependency Resolution</text>
    <text x="14" y="140" font-size="12" fill="#1e293b">Target Architecture Design</text>
    <rect x="14" y="148" width="152" height="16" rx="3" fill="#dbeafe" />
    <text x="90" y="160" font-size="10" font-weight="600" text-anchor="middle" fill="#1e40af">Liquid MoE 8B (Reasoning)</text>
  </g>

  <!-- Arrow Analyzer -> Planning -->
  <path d="M 425 185 L 470 185" stroke="#1d4ed8" stroke-width="1.8" marker-end="url(#arrow-blue)" fill="none" />

  <!-- Stage 3: Planning Agent -->
  <g transform="translate(475, 100)">
    <rect width="180" height="170" rx="6" fill="#ffffff" stroke="#7c3aed" stroke-width="1.8" filter="url(#shadow)" />
    <rect width="180" height="30" rx="6" fill="#f5f3ff" stroke="#7c3aed" stroke-width="1" />
    <text x="90" y="20" font-size="14" font-weight="700" text-anchor="middle" fill="#6d28d9">2. Planning Agent</text>
    <text x="14" y="55" font-size="12" fill="#1e293b">AST Fragment Extraction</text>
    <text x="14" y="77" font-size="12" fill="#1e293b">Reverse-Topo Ordering</text>
    <text x="24" y="95" font-size="11" fill="#475569">Back-Edge Cycle Break</text>
    <text x="14" y="118" font-size="12" fill="#1e293b">Manifest Rewriter ($D_s \to D_t$)</text>
    <text x="14" y="140" font-size="12" fill="#1e293b">Skeleton Generation</text>
    <rect x="14" y="148" width="152" height="16" rx="3" fill="#ede9fe" />
    <text x="90" y="160" font-size="10" font-weight="600" text-anchor="middle" fill="#5b21b6">Part A (Code) &amp; Part B (Tests)</text>
  </g>

  <!-- Arrow Planning -> Translator -->
  <path d="M 655 185 L 700 185" stroke="#1e293b" stroke-width="1.8" marker-end="url(#arrow)" fill="none" />

  <!-- Stage 4: Translator Agent -->
  <g transform="translate(705, 100)">
    <rect width="180" height="170" rx="6" fill="#ffffff" stroke="#059669" stroke-width="1.8" filter="url(#shadow)" />
    <rect width="180" height="30" rx="6" fill="#ecfdf5" stroke="#059669" stroke-width="1" />
    <text x="90" y="20" font-size="14" font-weight="700" text-anchor="middle" fill="#047857">3. Translator Agent</text>
    <text x="14" y="55" font-size="12" fill="#1e293b">Chunked Code Generation</text>
    <text x="14" y="77" font-size="12" fill="#1e293b">Symbol Navigator ($\le 4$KB)</text>
    <text x="14" y="99" font-size="12" fill="#1e293b">Prior Memory ($\le 4$KB)</text>
    <text x="14" y="121" font-size="12" fill="#1e293b">Part B Test Synthesis</text>
    <text x="14" y="140" font-size="12" fill="#1e293b">Disk Sync &amp; Checkpointing</text>
    <rect x="14" y="148" width="152" height="16" rx="3" fill="#d1fae5" />
    <text x="90" y="160" font-size="10" font-weight="600" text-anchor="middle" fill="#065f46">Qwen3-4B (Coding SLM)</text>
  </g>

  <!-- Arrow Translator -> Validator -->
  <path d="M 885 185 L 930 185" stroke="#059669" stroke-width="1.8" marker-end="url(#arrow)" fill="none" />

  <!-- Stage 5: Validator Agent -->
  <g transform="translate(935, 100)">
    <rect width="130" height="170" rx="6" fill="#ffffff" stroke="#e11d48" stroke-width="1.8" filter="url(#shadow)" />
    <rect width="130" height="30" rx="6" fill="#fff1f2" stroke="#e11d48" stroke-width="1" />
    <text x="65" y="20" font-size="14" font-weight="700" text-anchor="middle" fill="#be123c">4. Validator</text>
    <text x="10" y="55" font-size="11" fill="#1e293b">1. Pre-AST Syntax</text>
    <text x="10" y="77" font-size="11" fill="#1e293b">2. Cargo Build</text>
    <text x="10" y="99" font-size="11" fill="#1e293b">3. Test Execution</text>
    <text x="10" y="121" font-size="11" fill="#1e293b">4. Weakening Guard</text>
    <text x="10" y="143" font-size="11" fill="#1e293b">5. Plateau Detector</text>
    <rect x="10" y="148" width="110" height="16" rx="3" fill="#ffe4e6" />
    <text x="65" y="160" font-size="10" font-weight="600" text-anchor="middle" fill="#9f1239">Cascade Gate</text>
  </g>

  <!-- Repair Loop (Feedback Arc) -->
  <path d="M 1000 270 L 1000 370 L 795 370 L 795 275" stroke="#be123c" stroke-width="1.8" stroke-dasharray="5,4" marker-end="url(#arrow-rose)" fill="none" />
  <rect x="830" y="355" width="140" height="24" rx="4" fill="#ffffff" stroke="#be123c" stroke-width="1" />
  <text x="900" y="371" font-size="11" font-weight="700" text-anchor="middle" fill="#be123c">Repair Mode (Iter $\le K$)</text>

  <!-- Checkpoints & Persistence Box -->
  <g transform="translate(245, 330)">
    <rect width="410" height="140" rx="6" fill="#ffffff" stroke="#d97706" stroke-width="1.5" stroke-dasharray="4,3" filter="url(#shadow)" />
    <rect width="410" height="26" rx="6" fill="#fef3c7" stroke="#d97706" stroke-width="1" />
    <text x="205" y="18" font-size="13" font-weight="700" text-anchor="middle" fill="#92400e">State Persistence, Checkpoints &amp; Defense Mechanisms</text>
    
    <text x="14" y="48" font-size="12" fill="#1e293b">• <tspan font-weight="700">Atomic Checkpoint Snapshotting:</tspan> Serialized under <tspan class="mono">.artifacts/&lt;run&gt;/checkpoints/</tspan></text>
    <text x="14" y="70" font-size="12" fill="#1e293b">• <tspan font-weight="700">CodaMOSA Coverage Plateau Detector:</tspan> Halts unproductive repair loops (2 stag. pairs)</text>
    <text x="14" y="92" font-size="12" fill="#1e293b">• <tspan font-weight="700">Adversarial Weakening Guard:</tspan> Rejects assertion drops or deletion of tests</text>
    <text x="14" y="114" font-size="12" fill="#1e293b">• <tspan font-weight="700">Context Bounding:</tspan> Caps raw prompt slices to 32KB to prevent local SLM saturation</text>
  </g>

  <!-- Output Box (Success Exit) -->
  <g transform="translate(935, 340)">
    <rect width="130" height="110" rx="6" fill="#f0fdf4" stroke="#16a34a" stroke-width="2" filter="url(#shadow)" />
    <text x="65" y="28" font-size="13" font-weight="700" text-anchor="middle" fill="#15803d">Target Crate</text>
    <text x="65" y="50" font-size="11" text-anchor="middle" fill="#166534">• Clean Compile</text>
    <text x="65" y="70" font-size="11" text-anchor="middle" fill="#166534">• 100% Tests Pass</text>
    <text x="65" y="90" font-size="11" text-anchor="middle" fill="#166534">• Preserved AST</text>
  </g>

  <!-- Arrow Validator -> Success Exit -->
  <path d="M 1000 270 L 1000 335" stroke="#16a34a" stroke-width="2" marker-end="url(#arrow)" fill="none" />
  <text x="1015" y="305" font-size="11" font-weight="700" fill="#16a34a">Success</text>
</svg>
`;

// 2. Failure Taxonomy & Stability Comparison (Figure 2)
const fig2SVG = `
<svg width="1100" height="420" viewBox="0 0 1100 420" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <filter id="shadow2" x="-5%" y="-5%" width="110%" height="110%">
      <feDropShadow dx="0" dy="2" stdDeviation="3" flood-opacity="0.08"/>
    </filter>
  </defs>
  
  <rect x="10" y="10" width="1080" height="400" rx="10" fill="#f8fafc" stroke="#cbd5e1" stroke-width="1.5" />
  <text x="550" y="40" font-size="19" font-weight="700" text-anchor="middle" fill="#0f172a">Figure 2: Empirical Performance, Latency Trajectory &amp; Failure Root-Cause Taxonomy</text>
  <text x="550" y="62" font-size="13" font-style="italic" text-anchor="middle" fill="#475569">Multi-run empirical evaluation across benchmark suites with detailed failure distribution analysis</text>

  <!-- Panel A: Performance Overview Table / Chart -->
  <g transform="translate(35, 85)">
    <rect width="490" height="300" rx="6" fill="#ffffff" stroke="#94a3b8" stroke-width="1" filter="url(#shadow2)" />
    <rect width="490" height="28" rx="6" fill="#f1f5f9" stroke="#94a3b8" stroke-width="1" />
    <text x="245" y="19" font-size="13" font-weight="700" text-anchor="middle" fill="#1e293b">Table (a): Multi-Run Stability Results (K=3 Runs per Benchmark)</text>
    
    <!-- Table Header -->
    <rect x="10" y="38" width="470" height="24" fill="#e2e8f0" />
    <text x="20" y="54" font-size="11" font-weight="700" fill="#0f172a">Benchmark</text>
    <text x="140" y="54" font-size="11" font-weight="700" fill="#0f172a">Pair</text>
    <text x="210" y="54" font-size="11" font-weight="700" fill="#0f172a">Compile</text>
    <text x="280" y="54" font-size="11" font-weight="700" fill="#0f172a">Test Pass</text>
    <text x="360" y="54" font-size="11" font-weight="700" fill="#0f172a">Mean Wall (s)</text>
    <text x="445" y="54" font-size="11" font-weight="700" fill="#0f172a">Status</text>

    <!-- Row 1: GildedRose -->
    <rect x="10" y="66" width="470" height="48" fill="#f8fafc" />
    <text x="20" y="86" font-size="12" font-weight="600" fill="#1e293b">GildedRose</text>
    <text x="20" y="102" font-size="10" fill="#64748b">C Kata (7 files)</text>
    <text x="140" y="94" font-size="11" fill="#334155">C $\to$ Rust</text>
    <text x="210" y="94" font-size="11" font-weight="700" fill="#15803d">3 / 3 (100%)</text>
    <text x="280" y="94" font-size="11" font-weight="700" fill="#15803d">100%</text>
    <text x="360" y="94" font-size="11" fill="#334155">450.0 $\pm$ 45s</text>
    <text x="445" y="94" font-size="11" font-weight="700" fill="#15803d">CONV</text>

    <!-- Row 2: Gohistogram -->
    <rect x="10" y="118" width="470" height="48" fill="#ffffff" />
    <text x="20" y="138" font-size="12" font-weight="600" fill="#1e293b">Gohistogram</text>
    <text x="20" y="154" font-size="10" fill="#64748b">Go Lib (11 files)</text>
    <text x="140" y="146" font-size="11" fill="#334155">Go $\to$ Rust</text>
    <text x="210" y="146" font-size="11" fill="#334155">1 / 3 (33%)</text>
    <text x="280" y="146" font-size="11" fill="#334155">50%</text>
    <text x="360" y="146" font-size="11" fill="#334155">750.8 $\pm$ 82s</text>
    <text x="445" y="146" font-size="11" font-weight="600" fill="#b45309">PART</text>

    <!-- Row 3: Stats -->
    <rect x="10" y="170" width="470" height="48" fill="#f8fafc" />
    <text x="20" y="190" font-size="12" font-weight="600" fill="#1e293b">Stats</text>
    <text x="20" y="206" font-size="10" fill="#64748b">Go Math (73 files)</text>
    <text x="140" y="198" font-size="11" fill="#334155">Go $\to$ Rust</text>
    <text x="210" y="198" font-size="11" fill="#334155">0 / 3 (0%)</text>
    <text x="280" y="198" font-size="11" fill="#334155">0%</text>
    <text x="360" y="198" font-size="11" fill="#334155">900.0 $\pm$ 0s</text>
    <text x="445" y="198" font-size="11" font-weight="600" fill="#be123c">TIMEOUT</text>

    <!-- Row 4: Commons-Validator -->
    <rect x="10" y="222" width="470" height="48" fill="#ffffff" />
    <text x="20" y="242" font-size="12" font-weight="600" fill="#1e293b">Commons-Val.</text>
    <text x="20" y="258" font-size="10" fill="#64748b">Java Enterprise (100+ files)</text>
    <text x="140" y="250" font-size="11" fill="#334155">Java $\to$ Rust</text>
    <text x="210" y="250" font-size="11" fill="#334155">0 / 3 (0%)</text>
    <text x="280" y="250" font-size="11" fill="#334155">0%</text>
    <text x="360" y="250" font-size="11" fill="#334155">900.0 $\pm$ 0s</text>
    <text x="445" y="250" font-size="11" font-weight="600" fill="#be123c">TIMEOUT</text>
    
    <!-- Legend / Note -->
    <text x="20" y="288" font-size="10" font-style="italic" fill="#64748b">* Execution with local SLM topology (Liquid MoE 8B + Qwen3-4B) under bounded inference budget.</text>
  </g>

  <!-- Panel B: Failure Root-Cause Breakdown -->
  <g transform="translate(555, 85)">
    <rect width="500" height="300" rx="6" fill="#ffffff" stroke="#94a3b8" stroke-width="1" filter="url(#shadow2)" />
    <rect width="500" height="28" rx="6" fill="#f1f5f9" stroke="#94a3b8" stroke-width="1" />
    <text x="250" y="19" font-size="13" font-weight="700" text-anchor="middle" fill="#1e293b">Chart (b): Failure Mode &amp; Root-Cause Distribution</text>

    <!-- Category 1: Rust Borrow Checker & Lifetime -->
    <g transform="translate(20, 48)">
      <text x="0" y="14" font-size="12" font-weight="600" fill="#1e293b">Borrow Checker &amp; Lifetimes (E0382, E0502, E0597)</text>
      <rect x="0" y="22" width="460" height="18" rx="3" fill="#f1f5f9" />
      <rect x="0" y="22" width="193" height="18" rx="3" fill="#ef4444" />
      <text x="200" y="36" font-size="11" font-weight="700" fill="#1e293b">42% (Mutable Aliasing &amp; Ownership across modules)</text>
    </g>

    <!-- Category 2: Cyclic Dependencies & Module Visibility -->
    <g transform="translate(20, 106)">
      <text x="0" y="14" font-size="12" font-weight="600" fill="#1e293b">Cyclic Dependencies &amp; Module Visibility (E0432, E0603)</text>
      <rect x="0" y="22" width="460" height="18" rx="3" fill="#f1f5f9" />
      <rect x="0" y="22" width="124" height="18" rx="3" fill="#f97316" />
      <text x="132" y="36" font-size="11" font-weight="700" fill="#1e293b">27% (Circular imports across source files)</text>
    </g>

    <!-- Category 3: Standard Library & Type Mappings -->
    <g transform="translate(20, 164)">
      <text x="0" y="14" font-size="12" font-weight="600" fill="#1e293b">Foreign Runtime Semantics (Java/Go GC vs RAII)</text>
      <rect x="0" y="22" width="460" height="18" rx="3" fill="#f1f5f9" />
      <rect x="0" y="22" width="83" height="18" rx="3" fill="#eab308" />
      <text x="92" y="36" font-size="11" font-weight="700" fill="#1e293b">18% (Reflection, Nullability &amp; Class inheritance)</text>
    </g>

    <!-- Category 4: Local SLM Context Saturation -->
    <g transform="translate(20, 222)">
      <text x="0" y="14" font-size="12" font-weight="600" fill="#1e293b">Inference Latency &amp; Context Budget Saturation</text>
      <rect x="0" y="22" width="460" height="18" rx="3" fill="#f1f5f9" />
      <rect x="0" y="22" width="60" height="18" rx="3" fill="#3b82f6" />
      <text x="68" y="36" font-size="11" font-weight="700" fill="#1e293b">13% (Prompt throughput &amp; timeout ceilings)</text>
    </g>
  </g>
</svg>
`;

async function main() {
  const fig1Path = path.join(__dirname, 'fig1_workflow.png');
  const fig2Path = path.join(__dirname, 'fig2_empirical_results.png');
  
  await renderSVGToPNG(fig1SVG, fig1Path, 1140, 560);
  await renderSVGToPNG(fig2SVG, fig2Path, 1140, 460);
  console.log('All publication figures successfully generated with serif typography.');
}

main().catch(console.error);
