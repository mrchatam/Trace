import { formatName } from "./format";

/** Greet a caller by name (fixture symbol for P0-X indexing). */
export function greet(name: string): string {
  return `Hello, ${formatName(name)}!`;
}
