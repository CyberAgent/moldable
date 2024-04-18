module.exports = function (plop) {
  plop.setGenerator("benchmark", {
    description: "This is a benchmark generator",
    prompts: [
      {
        type: "input",
        name: "name",
        message: "Please enter any text.",
      },
    ],
    actions: [
      {
        type: "add",
        path: "src/plop/{{name}}.tsx",
        templateFile: "plop-templates/benchmark.hbs",
      },
    ],
  });
};
