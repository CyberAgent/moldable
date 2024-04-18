"use strict";

const { spawn } = require("child_process");
const { rimraf } = require("rimraf");
const { mkdirp } = require("mkdirp");

const FILE_NAME = "for-benchmark";

const LF = String.fromCharCode(0x0a); // \n
const DOWN = String.fromCharCode(0x1b, 0x5b, 0x42); // ↓
const ENTER = String.fromCharCode(0x0d); // enter

async function wait(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function readStream(stream) {
  let result = "";
  for await (const line of stream) {
    result += line;
  }
  return result;
}

async function clear() {
  await rimraf("./src");
  await mkdirp("./src");
}

(async () => {
  const type = process.argv[2];

  if (type === "plop") {
    const plop = spawn("./node_modules/.bin/plop");
    await wait(1000);
    plop.stdin.write(FILE_NAME);
    await wait(1000);
    const plopMeasureStart = performance.now();
    plop.stdin.write(LF);
    const plopMeasureEnd = performance.now();
    console.log(`plop: ${plopMeasureEnd - plopMeasureStart}ms`);
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
    console.log(`scaffdog: ${scaffdogMeasureEnd - scaffdogMeasureStart}ms`);
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
    console.log(`moldable: ${moldableMeasureEnd - moldableMeasureStart}ms`);
  }

  await clear();
})();
