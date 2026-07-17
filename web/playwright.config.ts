import { defineConfig } from "@playwright/test";
import { mkdtempSync } from "node:fs";
import { tmpdir } from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";

const PORT = 8099;
const dbDir = mkdtempSync(path.join(tmpdir(), "palimpsest-e2e-"));
const rootDir = path.dirname(fileURLToPath(import.meta.url));

// Drives the real built binary (embedded frontend + chi API + sqlite), not the
// vite dev server — this exercises the same code path a user actually hits.
// Tests run serially against one server/DB since pages created in one spec
// are read by others (e.g. the `[[` page-link picker).
export default defineConfig({
  testDir: "./e2e",
  timeout: 30_000,
  fullyParallel: false,
  workers: 1,
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI ? "line" : "list",
  use: {
    baseURL: `http://localhost:${PORT}`,
    trace: "retain-on-failure",
  },
  webServer: {
    command: "./palimpsest",
    cwd: path.resolve(rootDir, ".."),
    env: {
      PORT: String(PORT),
      DB_PATH: path.join(dbDir, "palimpsest-e2e.db"),
    },
    url: `http://localhost:${PORT}/health`,
    reuseExistingServer: false,
    timeout: 10_000,
  },
});
