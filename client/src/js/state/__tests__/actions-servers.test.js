import { connect, setServerName } from '../servers';

describe('connect()', () => {
  it('sets host and port correctly', () => {
    expect(connect('cake.com:1881', '', {})).toMatchObject({
      socket: {
        data: {
          host: 'cake.com',
          port: '1881'
        }
      }
    });

    expect(connect('cake.com', '', {})).toMatchObject({
      socket: {
        data: {
          host: 'cake.com',
          port: undefined
        }
      }
    });
  });
});

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
