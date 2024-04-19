/**
 * Inspired by https://github.com/mizdra/eslint-interactive/blob/a5ab787c4ccc780a2999b88d59d719cd6c1e651d/e2e-test/global-installation/index.test.ts
 */
"use strict";

const { spawn } = require("child_process");
const { rimraf } = require("rimraf");
const { mkdirp } = require("mkdirp");

const FILE_NAME = "for-benchmark";
const MAX = 10;

const LF = String.fromCharCode(0x0a); // \n
const DOWN = String.fromCharCode(0x1b, 0x5b, 0x42); // ↓
const ENTER = String.fromCharCode(0x0d); // enter

async function wait(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function clear() {
  await rimraf("./src");
  await mkdirp("./src");
}

async function bench(type) {
  if (type === "plop") {
    const plop = spawn("./node_modules/.bin/plop");
    await wait(1000);
    plop.stdin.write(FILE_NAME);
    await wait(1000);
    const plopMeasureStart = performance.now();
    plop.stdin.write(LF);
    const plopMeasureEnd = performance.now();
    return plopMeasureEnd - plopMeasureStart;
  }

  if (type === "scaffdog") {
    const scaffdog = spawn("./node_modules/.bin/scaffdog", ["generate"]);
    await wait(1000);
    scaffdog.stdin.write(LF);
    await wait(1000);
    scaffdog.stdin.write(DOWN);
    await wait(1000);
    scaffdog.stdin.write(ENTER);
    await wait(1000);
    scaffdog.stdin.write(FILE_NAME);
    await wait(1000);
    const scaffdogMeasureStart = performance.now();
    scaffdog.stdin.write(LF);
    const scaffdogMeasureEnd = performance.now();
    return scaffdogMeasureEnd - scaffdogMeasureStart;
  }

  if (type === "moldable") {
    const moldable = spawn("./node_modules/.bin/moldable");
    await wait(1000);
    moldable.stdin.write(LF);
    await wait(1000);
    moldable.stdin.write(FILE_NAME);
    await wait(1000);
    const moldableMeasureStart = performance.now();
    moldable.stdin.write(LF);
    const moldableMeasureEnd = performance.now();
    return moldableMeasureEnd - moldableMeasureStart;
  }

  return 0;
}

(async () => {
  const type = process.argv[2];

  const results = [];
  for (let i = 0; i < MAX; i++) {
    const result = await bench(type);
    console.log(`[${type}] ${i + 1} time: ${result}ms`);
    results.push(result);
    await clear();
  }

  console.log(
    `[${type}] average time of 10 times: ${results.reduce((acc, cur) => acc + cur, 0) / MAX}ms`
  );

  process.exit(0);
})();
