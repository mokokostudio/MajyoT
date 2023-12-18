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
            value: toNano('0.08'),
        },
        {
            $$type: 'WithdrawJettons',
            query_id: 1n,
            amount: toNano('1'),
            wallet_address: Address.parseFriendly('kQBiv9Xjkbs0fuhnJEAisnpe-7QU5CIjozs78Unadgy1aQlV').address,
            response_destination: Address.parseFriendly('0QD4oc41Pp3PdPcd74BA9WR9CaimEt14Ve0I8OovZ4hQOKb_').address,
            destination: Address.parseFriendly('0QD4oc41Pp3PdPcd74BA9WR9CaimEt14Ve0I8OovZ4hQOKb_').address,
        }
    );

    ui.clearActionPrompt();
    ui.write('Withdraw Jettons');
}
