#! /usr/bin/env node

import child_process from "child_process";
import path from "path";

export function getBinaryPath() {
  const availablePlatforms = ["darwin", "linux", "windows"];
  const availableArchs = ["x64", "arm64"];

  const { platform: pf, arch } = process;
  let platform: string = pf;
  if (pf === "win32") {
    platform = "windows";
  }
  const ext = platform === "windows" ? ".exe" : "";

  if (!availablePlatforms.includes(platform)) {
    console.error(`Moldable does not presently support ${platform}.`);
    return "";
  }
  if (!availableArchs.includes(arch)) {
    console.error(`Moldable does not presently support ${arch}.`);
    return "";
  }

  const binaryFile = `moldable-${platform}-${arch}${ext}`;

  return path.resolve(__dirname, binaryFile);
}

export function run() {
  child_process.execFileSync(getBinaryPath(), process.argv.slice(2), { stdio: "inherit" });
}

run();
