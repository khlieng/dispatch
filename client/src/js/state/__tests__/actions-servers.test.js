import { setServerName } from '../servers';

describe('setServerName()', () => {
  it('passes valid names to the server', () => {
    const name = 'cake';
    const server = 'srv';

    expect(setServerName(name, server)).toMatchObject({
      socket: {
        type: 'set_server_name',
        data: { name, server }
      }
    });
  });

  it('does not pass invalid names to the server', () => {
    expect(setServerName('', 'srv').socket).toBeUndefined();
    expect(setServerName('   ', 'srv').socket).toBeUndefined();
  });
});
