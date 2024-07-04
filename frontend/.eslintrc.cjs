module.exports = {
    root: true,
    parser: "@typescript-eslint/parser",
    plugins: [
        "unused-imports",
        "eslint-plugin-react",
        "@typescript-eslint",
        "prettier",
    ],
    extends: [
        "eslint:recommended",
        "plugin:react/recommended",
        "plugin:react-hooks/recommended",
        "plugin:@typescript-eslint/eslint-recommended",
        "plugin:@typescript-eslint/recommended",
        "prettier",
    ],
    rules: {
        "react/react-in-jsx-scope": 0,
        "@typescript-eslint/no-var-requires": 0,
        "@typescript-eslint/no-empty-interface": 0,
        "no-empty-pattern": 0,
        "react/prop-types": 0,
        "react/no-unknown-property": [0, {}],
        "no-undef": 0,
        "react-hooks/exhaustive-deps": 0,
        "@typescript-eslint/no-explicit-any": 0,
        "prettier/prettier": 2,
        "no-unused-vars": "off", // or "@typescript-eslint/no-unused-vars": "off",
        "unused-imports/no-unused-imports": "error",
        "unused-imports/no-unused-vars": [
            "warn",
            {
                vars: "all",
                varsIgnorePattern: "^_",
                args: "after-used",
                argsIgnorePattern: "^_",
            },
        ],
    },
    settings: {
        react: {
            version: "detect",
        },
    },
};