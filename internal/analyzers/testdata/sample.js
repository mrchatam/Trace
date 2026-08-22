import { useState } from 'react';
import './styles.css';

export function greet(name) {
  return 'hello ' + name;
}

export class Greeter {
  constructor() {}
  sayHi() {
    return greet('world');
  }
}
