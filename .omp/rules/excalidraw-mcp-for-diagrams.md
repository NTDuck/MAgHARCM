---
name: excalidraw-mcp-for-diagrams
description: "Force every graph, diagram, flowchart, sequence diagram, data-flow, state, lifecycle, ER, or architecture visualization in this project to be generated through the Excalidraw MCP server at https://mcp.excalidraw.com/mcp. No Mermaid, ASCII, or HTML diagrams; archify is uninstalled. Grandfathered: pre-existing archify-derived PNGs in docs/.paper/figures/ referenced by LaTeX \\includegraphics."
scope: ["text", "tool"]
alwaysApply: true
---

# Excalidraw MCP Is the Sole Diagram Renderer

Any diagram, graph, chart, flowchart, sequence, state, lifecycle, ER, mind-map, architecture sketch, data-flow, or topology visualization in this project MUST be produced through the **Excalidraw MCP server** configured at `.omp/mcp.json` (`excalidraw` → `https://mcp.excalidraw.com/mcp`).

## 1. Hard Rules

1. When the task asks for a diagram of any kind, the model MUST call the Excalidraw MCP server tools (e.g. `create_excalidraw_diagram`, `create_element`, `create_rectangle`, `create_ellipse`, `create_arrow`, `create_text`, etc.) instead of returning Mermaid, PlantUML, Graphviz DOT, ASCII art, HTML, SVG, PNG, or `archify` HTML.
2. NEVER emit a `mermaid` fenced code block for a final deliverable. If Mermaid is the only way the user expressed the source, translate the topology and call Excalidraw MCP tools to author a fresh Excalidraw scene; do not paste Mermaid through Excalidraw.
3. NEVER invoke the `archify` skill, the `archify.mjs` CLI, or any node-based renderer to produce a final diagram artifact. `archify` is uninstalled from this project and is no longer a permitted path.
4. NEVER commit `*.html`, `*.svg`, `*.png`, `*.jpg`, `*.webp`, or `*.pdf` files as the diagram source-of-truth for a NEW deliverable. Excalidraw owns the visual layer; persistence is the canvas URL or the Excalidraw scene returned by the MCP, optionally embedded in vault notes as `![Excalidraw](<url>)`.
5. For local snapshots that must live in this repository, save the JSON scene returned by the Excalidraw MCP under `.obsidian/MAgHARCM/diagrams/<Slug>.excalidraw` or under `docs/diagrams/<Slug>.excalidraw`, then link to it. The scene file is the canonical artifact; do not re-render to a raster.
6. EXCEPTION — LaTeX paper figures: `docs/.paper/sec_method.tex` includes `figures/fig1_workflow.png`, `figures/fig3_dataflow.png`, and `figures/fig4_lifecycle.png` via `\includegraphics`. These rasters and the existing archify-derived PNGs in `docs/.paper/figures/` are GRANDMOTHERED and MAY remain in the tree. Re-deriving them through Excalidraw MCP is a follow-up migration tracked elsewhere; this rule does NOT author NEW PNG figures for any purpose.
7. EXCEPTION — Excalidraw scene export: when the LaTeX toolchain genuinely requires a raster and the figure has been freshly authored through the Excalidraw MCP, exporting the scene to PNG/SVG for `\includegraphics` is permitted (the source remains Excalidraw).
8. When the user asks for a "diagram", "graph", "chart", "flowchart", "sequence diagram", "state diagram", "ER diagram", "data flow", or "architecture diagram" without naming a tool, default to Excalidraw MCP without asking.
## 2. Tool Selection Order

When authoring a diagram through Excalidraw MCP:

1. Prefer the high-level `create_excalidraw_diagram` (or equivalent scene-builder tool exposed by the server) when the user wants a complete scene in one call.
2. Use the lower-level element tools (`create_rectangle`, `create_ellipse`, `create_diamond`, `create_arrow`, `create_line`, `create_text`, `create_group`, etc.) when the diagram needs explicit layout control, custom styling, or iterative refinement.
3. For sequence diagrams, prefer the MCP's dedicated sequence tool if available; otherwise compose rectangles, lines, arrows, and text labels in the same order as the call sites.
4. For state/lifecycle diagrams, prefer the MCP's dedicated state tool if available; otherwise compose rounded rectangles for states, arrows for transitions, and text labels for guards.

## 3. Vault and Doc Integration

1. Paper notes in `.obsidian/MAgHARCM/research/papers/`, primitive specs, and ADR files embed diagrams as `![Excalidraw](<scene-url>)` links pointing at the Excalidraw canvas URL returned by the MCP.
2. `.obsidian/MAgHARCM/adhoc/Architecture-And-Dataflow.md` and similar reports must source their visualizations from Excalidraw MCP; do not paste Mermaid blocks there.
3. `docs/.paper/` LaTeX figures: the existing PNGs in `docs/.paper/figures/` referenced by `\includegraphics` (`fig1_workflow.png`, `fig3_dataflow.png`, `fig4_lifecycle.png`) and the visual-check PNGs in the same directory are GRANDMOTHERED archify raster outputs and MAY remain. They MUST NOT be re-rasterized via archify or any other non-Excalidraw renderer. The migration path is to regenerate each figure end-to-end through the Excalidraw MCP and replace the PNG via a future sprint commit; until that migration lands, the PNGs stay as-is. Archiving or restoring the deleted `*.html` source files or the `*.workflow.json` / `*.dataflow.json` / `*.lifecycle.json` / `*.arch.json` / `*.visual-check.json` authoring specs is forbidden; those were intentionally purged so that no future regenerator can reach for them.

## 4. Failure Modes

1. If the Excalidraw MCP server is unreachable, surface the connection error in the response and stop. Do not silently fall back to Mermaid, ASCII, or any other renderer.
2. If the server returns an error or a partial scene, retry once with a corrected payload; on a second failure, report the diagnostic to the user verbatim and ask whether to abort or proceed without a diagram.
3. Do not invent a fake Excalidraw scene; if no tool call succeeded, the diagram does not exist and the response must say so.

## 5. Quick Check

Before yielding any task that produced a diagram, confirm:

- The diagram came from an `mcp__excalidraw__*` tool call.
- The output is the Excalidraw scene/URL, not a Mermaid block or HTML page.
- No NEW `archify`, `mermaid`, `plantuml`, `dot`, `png`, `svg`, or `html` diagram artifact was committed as the diagram source-of-truth. Pre-existing archify-derived PNGs in `docs/.paper/figures/` referenced by `\includegraphics` are grandfathered and are not violations.
