
contract NftItem {
    collection_address: Address;
    item_index: Int; 
    owner: Address;
    individual_content: Cell;
    is_initialized: Bool;
    lock : Bool;

    init(collection_address: Address, item_index: Int, owner: Address , individual_content: Cell){
        self.collection_address = collection_address;
        self.item_index = item_index;
        self.owner = collection_address;
        self.individual_content = individual_content;
        self.is_initialized = false;
        self.lock = false;
    }

    receive(msg: Lock){
        let ctx: Context = context();
        require((ctx.sender == self.owner) || (ctx.sender == self.collection_address), "not owner");
        self.lock = true;
        let commit: StringBuilder = beginComment();
        commit.append("Lock:");
        commit.append(self.item_index.toString());
        emit(commit.toCell());
    }

    receive(msg: Unlock){
        let ctx: Context = context();
        require((ctx.sender == self.owner) || (ctx.sender == self.collection_address), "not owner");
        self.lock = false;
        let commit: StringBuilder = beginComment();
        commit.append("UnLock:");
        commit.append(self.item_index.toString());
        emit(commit.toCell());
    }

    receive(msg: Transfer){
        let ctx: Context = context();
        
        // Only Owner of the this NFT Item can transfer it.
        require((ctx.sender == self.owner) || (ctx.sender == self.collection_address), "not owner");

        require(!self.lock, "not unlock");

        let msgValue: Int = ctx.value; // Check the gasCost for storage
        let tonBalanceBeforeMsg: Int = myBalance() - msgValue;
        let storageFee: Int = minTonsForStorage - min(tonBalanceBeforeMsg, minTonsForStorage);
        msgValue = msgValue - (storageFee + gasConsumption);

        if (self.is_initialized == false) {  // Initial Transfer, aka the "Minting" of the NFT
            self.is_initialized = true;
            self.owner = msg.new_owner;
            self.lock = true;
            send(SendParameters{
                to: msg.response_destination,
                value: 0,
                mode:  SendIgnoreErrors + SendRemainingValue,
                body: Excesses { query_id: msg.query_id }.toCell()
            }); 
            let commit: StringBuilder = beginComment();
            commit.append("Mint:");
            commit.append(self.item_index.toString());
            emit(commit.toCell());
        } else {
            
            self.owner = msg.new_owner;  // change current owner to the new_owner
            if (msg.forward_amount > 0) {
                send(SendParameters{
                    to: msg.new_owner,
                    value: msg.forward_amount,
                    mode: SendIgnoreErrors, 
                    bounce: false,
                    body: OwnershipAssigned{
                        query_id: msg.query_id,
                        prev_owner: ctx.sender,
                        forward_payload: msg.forward_payload
                    }.toCell()
                }); 
            }

            msgValue = msgValue - ctx.readForwardFee(); 
            if (msg.response_destination != null) { 
                send(SendParameters{ 
                    to: msg.response_destination,
                    value: msgValue - msg.forward_amount,
                    mode: SendPayGasSeparately,
                    body: Excesses { query_id: msg.query_id }.toCell()
                });
            } 
        }
    }
    
    receive(msg: GetStaticData){ 
        let ctx: Context = context();
        send(SendParameters {
            to: ctx.sender,
            value: 0,
            mode: 64,  // (return msg amount except gas fees) 
            bounce: true,
            body: ReportStaticData{
                query_id: msg.query_id,
                index_id: self.item_index,
                collection: self.collection_address
            }.toCell()
        });
    }

    // --------- Get Function  --------- //
    get fun get_nft_data(): GetNftData {
        let b: StringBuilder = beginString();
        let collectionData: String = self.individual_content.asSlice().asString();
        b.append(collectionData);
        b.append(self.item_index.toString());
        b.append(".json");

        return GetNftData {
            is_initialized: self.is_initialized, 
            index: self.item_index, 
            collection_address: self.collection_address, 
            owner_address: self.owner,
            individual_content: b.toCell()
        };
    }

    get fun get_lock_status(): Bool {
        return self.lock;
    }
}