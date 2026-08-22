import type { FC } from 'react';
import { useMemo } from 'react';

export function compute(x: number): number {
  return x * 2;
}

export class Counter {
  value = 0;
  inc(): void {
    this.value += 1;
  }
}
