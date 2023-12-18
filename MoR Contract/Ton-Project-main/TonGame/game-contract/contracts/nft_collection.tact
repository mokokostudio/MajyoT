import "@stdlib/stoppable";
import "@stdlib/ownable";
import "@stdlib/deploy";

import "./message.tact";
import "./nft_item.tact";

const minTonsForStorage: Int = ton("0.03");
const gasConsumption: Int = ton("0.03");

contract NftCollection with Deployable, Resumable, Stoppable,OwnableTransferable{
    next_item_index: Int as uint32 = 1;
    owner: Address;
    stopped: Bool;
    royalty_params: RoyaltyParams;                      
    collection_content: Cell;
    collection_init: CollectionInit;
    

    init(owner: Address, price: Int, collection_content: Cell, collection_init: CollectionInit,royalty_params: RoyaltyParams){
        self.owner = owner;
        self.stopped = false;
        self.collection_content = collection_content; 
        self.royalty_params = royalty_params;
        self.collection_init = collection_init;
    }

    receive(msg: SetCollectionInitWalletAddress ){
        require(sender() == self.owner, "Only deployer is allowed to set");
        self.collection_init.wallet_address = msg.wallet_address;
    }

    receive(msg: SetCollectionInitPrice){
        require(sender() == self.owner, "Only deployer is allowed to set");
        self.collection_init.price = msg.price;
    }

    receive(msg: SetCollectionInitNumberOfPerDayMints){
        require(sender() == self.owner, "Only deployer is allowed to set");
        self.collection_init.number_of_per_day_mints = msg.number_of_per_day_mints;
    }

    receive(msg: SetCollectionInitWhiteMintStartTime){
        require(sender() == self.owner, "Only deployer is allowed to set");
        self.collection_init.white_mint_start_time = msg.white_mint_start_time;
    }

    receive(msg: SetCollectionInitPublicMintStartTime){
        require(sender() == self.owner, "Only deployer is allowed to set");
        self.collection_init.public_mint_start_time = msg.public_mint_start_time;
    }

    receive(msg: SetCollectionInitPromotionStartTime){
        require(sender() == self.owner, "Only deployer is allowed to set");
        self.collection_init.promotion_start_time = msg.promotion_start_time;
    }

    receive(msg: SetCollectionInitPromotionEndTime){
        require(sender() == self.owner, "Only deployer is allowed to set");
        self.collection_init.promotion_end_time = msg.promotion_end_time;
    }

    receive(msg: SetCollectionInitPromotionPrice){
        require(sender() == self.owner, "Only deployer is allowed to set");
        self.collection_init.promotion_price = msg.promotion_price;
    }

    receive(msg: AdminUnlock){
        require(sender() == self.owner, "Only deployer is allowed to set");
        let nft_init: StateInit = self.getNftItemInit(msg.tokenId);
        let ctx: Context = context(); 
        let msgValue: Int = ctx.value;
        send(SendParameters{
            to: contractAddress(nft_init),
            bounce: false,
            value: msgValue - gasConsumption,
            mode: SendRemainingValue + SendIgnoreErrors,
            body: Unlock{
                query_id: msg.query_id
            }.toCell()
        });
    }

    receive("Mint"){
        // require(sender() == self.owner, "Only deployer is allowed to set");
        let ctx: Context = context(); 
        let msgValue: Int = ctx.value;
        let tonBalanceBeforeMsg: Int = myBalance() - msgValue;
        let storageFee: Int = minTonsForStorage - min(tonBalanceBeforeMsg, minTonsForStorage);
        msgValue = msgValue - (storageFee + gasConsumption);
        require(msgValue > 0, "Low Price");
        self.mint(ctx.sender, msgValue, false);
    }

    receive("withdraw safe") {
        require(sender() == self.owner, "Only deployer is allowed to withdraw");
        send(SendParameters{
            to: sender(),
            bounce: true,
            value: myBalance() - context().value - minTonsForStorage,
            mode: SendRemainingValue + SendIgnoreErrors
        });
    }

    receive(msg: WithdrawJettons){
        require(sender() == self.owner, "Only deployer is allowed to withdraw");
        send(
            SendParameters{
                to: msg.wallet_address,
                value: ton("0.05"),
                mode: SendPayGasSeparately + SendIgnoreErrors,
                bounce: false,
                body: TokenTransfer{
                    query_id:  msg.query_id,
                    amount: msg.amount,
                    destination: msg.destination,
                    response_destination: msg.response_destination,
                    custom_payload: null,
                    forward_ton_amount: 0,
                    forward_payload: beginCell().storeUint(1,8).asSlice()
                }.toCell()
            }
        );
    }

    receive(msg: TokenNotification){
        let ctx: Context = context(); 
        let msgValue: Int = ctx.value;
        let tonBalanceBeforeMsg: Int = myBalance() - msgValue;
        let storageFee: Int = minTonsForStorage - min(tonBalanceBeforeMsg, minTonsForStorage);
        msgValue = msgValue - (storageFee + gasConsumption);
        require(msgValue > 0, "Low Price");
        let payload: Int = msg.forward_payload.loadUint(8);

        if(sender() == self.collection_init.wallet_address){
            
            if(!self.check_mintable(payload)){
                self.sendTokens(
                    sender(),
                    msg.query_id,
                    msg.amount,
                    msg.from
                );
                return;
            }
            let return_amount: Int = 0;
               
            if( payload == 3 ){
                return_amount = msg.amount - self.collection_init.promotion_price;
            }else if( payload == 2 ){
                
            }else if( payload == 1 ){
                return_amount = msg.amount - self.collection_init.price;
            }
            if( return_amount > 0 ){
                self.sendTokens(
                    sender(),
                    msg.query_id,
                    return_amount,
                    msg.from
                );
            }else if( return_amount < 0) {
                self.sendTokens(
                    sender(),
                    msg.query_id,
                    msg.amount,
                    msg.from
                );
                return;
            }else if( return_amount == 0 ){

            }
            self.mint(ctx.sender, msgValue, payload == 2);
        }else{
            self.sendTokens(
                    sender(),
                    msg.query_id,
                    msg.amount,
                    msg.from
                );
            return;
        }
    }
    
    fun sendTokens(
        to: Address,
        query_id: Int,
        amount: Int,
        destination: Address
    ) {
        send(
            SendParameters{
                to: to,
                value: ton("0.02"),
                mode: SendPayGasSeparately + SendIgnoreErrors,
                bounce: false,
                body: TokenTransfer{
                    query_id: query_id,
                    amount: amount,
                    destination: destination,
                    response_destination: destination,
                    custom_payload: null,
                    forward_ton_amount: 0,
                    forward_payload: beginCell().storeUint(1,8).asSlice()
                }.toCell()
            }
        );
    }

    // ===== Private Methods ===== //
    fun mint(sender: Address, msgValue: Int, whiteMint: Bool) {
        require(self.next_item_index >= 0, "non-sequential NFTs");
        require(self.next_item_index < self.collection_init.total_items, "no more NFTs");
        require(self.collection_init.today_mints < self.collection_init.number_of_per_day_mints, "today no more NFTs");
        
        let nft_init: StateInit = self.getNftItemInit(self.next_item_index);
        send(SendParameters{
                to: contractAddress(nft_init), 
                value: msgValue, 
                bounce: false,
                mode: SendIgnoreErrors,
                body: Transfer {
                    query_id: 0,
                    new_owner: sender,
                    response_destination: self.owner,
                    custom_payload: emptyCell(),
                    forward_amount: 0,
                    forward_payload: emptySlice()
                }.toCell(),
                code: nft_init.code,
                data: nft_init.data
            });
        self.next_item_index = self.next_item_index + 1;
        if(!whiteMint){
            self.collection_init.today_mints = self.collection_init.today_mints + 1;
        }
    }

    receive(msg: GetRoyaltyParams) {   
        let ctx: Context = context(); // get sender Info
        send(SendParameters{
            to: ctx.sender,
            value: 0,
            mode: 64, 
            bounce: false,
            body: ReportRoyaltyParams {
                query_id: msg.query_id,
                numerator:  self.royalty_params.numerator,
                denominator: self.royalty_params.denominator,
                destination: self.owner
            }.toCell()
        });        
    }

    fun check_mintable(type : Int): Bool {
        let today: Int = now() / 86400;
        if (today != self.collection_init.today) {
            self.collection_init.today = today;
            self.collection_init.today_mints = 0;
        }

        let isPublicMint: Bool = self.collection_init.public_mint_start_time <= now();  
        let isWhiteMint: Bool = self.collection_init.white_mint_start_time <= now();
        let isPromotionMint: Bool = self.collection_init.promotion_start_time <= now() && now() <= self.collection_init.promotion_end_time;

        let mintable: Bool = false;

        // PublicMint
        if (type == 1){ 
            mintable = isPublicMint && self.next_item_index < self.collection_init.total_items && self.collection_init.today_mints < self.collection_init.number_of_per_day_mints;
        // WhiteMint
        } else if (type == 2) {
            mintable = isWhiteMint && self.next_item_index < self.collection_init.total_items && self.collection_init.today_mints < self.collection_init.number_of_per_day_mints;
        // PromotionMint
        } else if (type == 3) {
            mintable = isPromotionMint && self.next_item_index < self.collection_init.total_items && self.collection_init.today_mints < self.collection_init.number_of_per_day_mints;
        }

        return mintable;
    }

    // ------------------ Get Function  ------------------ //
    get fun get_collection_data(): CollectionData {     
        let b: StringBuilder = beginString();
        let collectionDataString: String = self.collection_content.asSlice().asString();
        b.append(collectionDataString);
        b.append("meta.json"); // You can changed this your self.
        return CollectionData{
            next_item_index: self.next_item_index, 
            collection_content: b.toCell(), 
            owner_address: self.owner
        };
    }

    get fun get_nft_address_by_index(item_index: Int): Address?{      
        let initCode: StateInit = self.getNftItemInit(item_index);
        return contractAddress(initCode);
    }

    get fun getNftItemInit(item_index: Int): StateInit {
        return initOf NftItem(myAddress(), item_index, self.owner, self.collection_content);
    }
    
    get fun get_nft_content(index: Int, individual_content: Cell): Cell { 
        let b: StringBuilder = beginString();
        let ic: String = individual_content.asSlice().asString();
        b.append(ic);
        return b.toCell();
    }

    get fun get_collection_init(): CollectionInit {
        return self.collection_init;
    }

    get fun royalty_params(): RoyaltyParams {
        return self.royalty_params;
    }
}


