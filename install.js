#!/usr/bin/env node

const { execFileSync } = require("child_process");
const { createWriteStream, copyFileSync, chmodSync, unlinkSync, mkdtempSync, rmSync } = require("fs");
const os = require("os");
const https = require("https");
const http = require("http");
const path = require("path");
const { pipeline } = require("stream/promises");

const VERSION = require("./package.json").version;
const REPO = "naorpeled/aitutor";

const PLATFORM_MAP = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const ARCH_MAP = {
  x64: "amd64",
  arm64: "arm64",
};

function getDownloadURL() {
  const platform = PLATFORM_MAP[process.platform];
  const arch = ARCH_MAP[process.arch];

  if (!platform || !arch) {
    throw new Error(
      `Unsupported platform: ${process.platform} ${process.arch}\n` +
        `Supported: darwin (amd64, arm64), linux (amd64, arm64), windows (amd64, arm64)`
    );
  }

  const ext = process.platform === "win32" ? "zip" : "tar.gz";
  return `https://github.com/${REPO}/releases/download/v${VERSION}/aitutor_${platform}_${arch}.${ext}`;
}

function getBinaryName() {
  return process.platform === "win32" ? "aitutor.exe" : "aitutor";
}

function follow(url, depth = 0) {
  if (depth > 10) {
    return Promise.reject(new Error("Too many redirects"));
  }
  const mod = url.startsWith("https") ? https : http;
  return new Promise((resolve, reject) => {
    mod
      .get(url, { headers: { "User-Agent": "aitutor-npm" } }, (res) => {
        if (
          res.statusCode >= 300 &&
          res.statusCode < 400 &&
          res.headers.location
        ) {
          return follow(res.headers.location, depth + 1).then(resolve, reject);
        }
        if (res.statusCode !== 200) {
          return reject(
            new Error(`Download failed: HTTP ${res.statusCode} from ${url}`)
          );
        }
        resolve(res);
      })
      .on("error", reject);
  });
}

async function install() {
  const url = getDownloadURL();
  const binName = getBinaryName();
  const binPath = path.join(__dirname, binName);

  console.log(`Downloading aitutor v${VERSION}...`);
  console.log(`  ${url}`);

  const res = await follow(url);

  if (process.platform === "win32") {
    const zipPath = path.join(__dirname, "aitutor.zip");
    const tmpDir = mkdtempSync(path.join(os.tmpdir(), "aitutor-"));
    try {
      await pipeline(res, createWriteStream(zipPath));
      const psEscape = (s) => s.replace(/'/g, "''");
      execFileSync("powershell.exe", [
        "-NoProfile",
        "-NonInteractive",
        "-Command",
        `Expand-Archive -LiteralPath '${psEscape(zipPath)}' -DestinationPath '${psEscape(tmpDir)}' -Force -ErrorAction Stop`,
      ]);
      copyFileSync(path.join(tmpDir, binName), binPath);
    } finally {
      rmSync(tmpDir, { recursive: true, force: true });
      try { unlinkSync(zipPath); } catch {}
    }
  } else {
    const tarPath = path.join(__dirname, "aitutor.tar.gz");
    await pipeline(res, createWriteStream(tarPath));
    execFileSync("tar", ["-xzf", tarPath, "-C", __dirname, binName], {
      stdio: "ignore",
    });
    unlinkSync(tarPath);
  }

  chmodSync(binPath, 0o755);
  console.log(`Installed aitutor to ${binPath}`);
}

install().catch((err) => {
  console.error(`Failed to install aitutor: ${err.message}`);
  console.error(
    `\nYou can install manually with: go install github.com/naorpeled/aitutor@latest`
  );
  process.exit(1);
});
