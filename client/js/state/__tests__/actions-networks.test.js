import { connect, setNetworkName } from '../networks';

describe('setNetworkName()', () => {
  it('passes valid names to the network', () => {
    const name = 'cake';
    const network = 'srv';

    expect(setNetworkName(name, network)).toMatchObject({
      socket: {
        type: 'set_network_name',
        data: { name, network }
      }
    });
  });

  it('does not pass invalid names to the network', () => {
    expect(setNetworkName('', 'srv').socket).toBeUndefined();
    expect(setNetworkName('   ', 'srv').socket).toBeUndefined();
  });
});
