import { Address, beginCell, toNano } from 'ton-core';
import { NftCollection } from '../wrappers/NftCollection';
import { NetworkProvider } from '@ton-community/blueprint';

export async function run(provider: NetworkProvider) {
    const nftCollection = provider.open(
        await NftCollection.fromInit(
            provider.sender().address!,
            toNano('1'),
            beginCell().storeInt(0x01, 8).storeStringTail('https://weapondata.majyo.vip/weapons/').endCell(),
            {
                $$type: 'CollectionInit',
                total_items: 20001n,
                wallet_address: Address.parse('0QD4oc41Pp3PdPcd74BA9WR9CaimEt14Ve0I8OovZ4hQOKb_'),
                price: toNano('1'),
                today: 0n,
                today_mints: 0n,
                number_of_per_day_mints: 500n,
                white_mint_start_time: 0n,
                public_mint_start_time: 0n,
                promotion_start_time: 0n,
                promotion_end_time: 0n,
                promotion_price: 0n,
            },
            {
                $$type: 'RoyaltyParams',
                numerator: 8n,
                denominator: 100n,
                destination: provider.sender().address!,
            }
        )
    );

    await nftCollection.send(
        provider.sender(),
        {
            value: toNano('0.05'),
        },
        {
            $$type: 'Deploy',
            queryId: 1n,
        }
    );

    await provider.waitForDeploy(nftCollection.address);

    console.log('Collection Data', await nftCollection.getGetCollectionData());
}
