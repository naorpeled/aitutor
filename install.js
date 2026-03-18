#!/usr/bin/env node

const { execFileSync } = require("child_process");
const { createWriteStream, readFileSync, chmodSync, unlinkSync, mkdtempSync, renameSync, rmSync } = require("fs");
const crypto = require("crypto");
const os = require("os");
const https = require("https");
const path = require("path");
const { pipeline } = require("stream/promises");

const VERSION = require("./package.json").version;
const REPO = "naorpeled/aitutor";
const MAX_REDIRECTS = 10;

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
  if (depth > MAX_REDIRECTS) {
    return Promise.reject(new Error(`Too many redirects (max ${MAX_REDIRECTS})`));
  }
  if (!url.startsWith("https://")) {
    return Promise.reject(new Error(`Refusing non-HTTPS URL: ${url}`));
  }
  return new Promise((resolve, reject) => {
    https
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

async function fetchChecksums() {
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/checksums.txt`;
  const res = await follow(url);
  const chunks = [];
  for await (const chunk of res) chunks.push(chunk);
  const body = Buffer.concat(chunks).toString("utf8");
  const checksums = {};
  for (const line of body.trim().split("\n")) {
    const [hash, filename] = line.trim().split(/\s+/);
    if (hash && filename) checksums[filename] = hash;
  }
  return checksums;
}

function verifyChecksum(filePath, expectedHash) {
  const data = readFileSync(filePath);
  const actual = crypto.createHash("sha256").update(data).digest("hex");
  if (actual !== expectedHash) {
    throw new Error(
      `Checksum mismatch for ${path.basename(filePath)}\n` +
        `  Expected: ${expectedHash}\n` +
        `  Actual:   ${actual}\n` +
        `The downloaded file may have been tampered with.`
    );
  }
}

async function install() {
  const url = getDownloadURL();
  const archiveName = path.basename(url);
  const binName = getBinaryName();
  const binPath = path.join(__dirname, binName);

  console.log(`Downloading aitutor v${VERSION}...`);
  console.log(`  ${url}`);

  const [res, checksums] = await Promise.all([follow(url), fetchChecksums()]);
  const expectedHash = checksums[archiveName];
  if (!expectedHash) {
    throw new Error(`No checksum found for ${archiveName} in checksums.txt`);
  }

  if (process.platform === "win32") {
    const zipPath = path.join(__dirname, "aitutor.zip");
    const tmpDir = mkdtempSync(path.join(os.tmpdir(), "aitutor-"));
    await pipeline(res, createWriteStream(zipPath));
    verifyChecksum(zipPath, expectedHash);
    const psEscape = (s) => s.replace(/'/g, "''");
    execFileSync("powershell.exe", [
      "-NoProfile",
      "-NonInteractive",
      "-Command",
      `Expand-Archive -LiteralPath '${psEscape(zipPath)}' -DestinationPath '${psEscape(tmpDir)}' -Force -ErrorAction Stop`,
    ], { stdio: "ignore" });
    renameSync(path.join(tmpDir, binName), binPath);
    rmSync(tmpDir, { recursive: true, force: true });
    unlinkSync(zipPath);
  } else {
    const tarPath = path.join(__dirname, "aitutor.tar.gz");
    await pipeline(res, createWriteStream(tarPath));
    verifyChecksum(tarPath, expectedHash);
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
