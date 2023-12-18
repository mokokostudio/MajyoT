import { Address, toNano } from 'ton-core';
import { NftCollection } from '../wrappers/NftCollection';
import { NetworkProvider, sleep } from '@ton-community/blueprint';

export async function run(provider: NetworkProvider, args: string[]) {
    const ui = provider.ui();

    const address = Address.parse(args.length > 0 ? args[0] : await ui.input('NftCollection address'));

    if (!(await provider.isContractDeployed(address))) {
        ui.write(`Error: Contract at address ${address} is not deployed!`);
        return;
    }

    const nftCollection = provider.open(NftCollection.fromAddress(address));

    await nftCollection.send(
        provider.sender(),
        {
            value: toNano('0.1'),
        },
        'Mint'
    );

    ui.clearActionPrompt();
    ui.write('Mint successfully!');
}
