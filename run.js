#!/usr/bin/env node

const { execFileSync } = require("child_process");
const path = require("path");

const binName = process.platform === "win32" ? "aitutor.exe" : "aitutor";
const binPath = path.join(__dirname, binName);

try {
  execFileSync(binPath, process.argv.slice(2), { stdio: "inherit" });
} catch (err) {
  if (err.status !== null) {
    process.exit(err.status);
  }
  console.error(`Failed to run aitutor: ${err.message}`);
  console.error(`Try reinstalling: npm install -g aitutor`);
  process.exit(1);
}
