export function isString(value: unknown): value is string {
    return typeof value === "string";
};

export function isStringArray(value: unknown): value is string[] {
  return Array.isArray(value) && value.every((v) => typeof v === "string");
}