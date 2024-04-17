import child_process from "child_process";
import path from "path";

import * as cli from "./cli";

jest.mock("child_process");
jest.mock("path");

describe("cli", () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  it("should run the binary with the correct arguments", () => {
    const execFileSyncSpy = jest.spyOn(child_process, "execFileSync").mockReturnValue("");
    const resolveSpy = jest.spyOn(path, "resolve").mockReturnValue("dummy/path");

    cli.run();

    expect(resolveSpy).toHaveBeenCalled();
    expect(execFileSyncSpy).toHaveBeenCalledWith("dummy/path", [], {
      stdio: "inherit",
    });
  });
});
