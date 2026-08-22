import { Foo } from './mod';

test('Foo', () => {
  Foo();
});

describe('suite', () => {
  it('nested', () => {
    Foo();
  });
});
