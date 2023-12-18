message(0x693d3950) GetRoyaltyParams {
    query_id: Int as uint64;
}
message(0xa8cb00ad) ReportRoyaltyParams {
    query_id: Int as uint64;
    numerator: Int as uint16;
    denominator: Int as uint16;
    destination: Address;
}

message(0x50) SetCollectionInitWalletAddress {
    query_id: Int as uint64;
    wallet_address: Address;
}

message(0x51) SetCollectionInitPrice {
    query_id: Int as uint64;
    price: Int;
}

message(0x52) SetCollectionInitNumberOfPerDayMints {
    query_id: Int as uint64;
    number_of_per_day_mints: Int;
}

message(0x53) SetCollectionInitWhiteMintStartTime {
    query_id: Int as uint64;
    white_mint_start_time: Int;
}

message(0x54) SetCollectionInitPublicMintStartTime {
    query_id: Int as uint64;
    public_mint_start_time: Int;
}

message(0x55) SetCollectionInitPromotionStartTime {
    query_id: Int as uint64;
    promotion_start_time: Int;
}

message(0x56) SetCollectionInitPromotionEndTime {
    query_id: Int as uint64;
    promotion_end_time: Int;
}

message(0x57) SetCollectionInitPromotionPrice {
    query_id: Int as uint64;
    promotion_price: Int;
}


struct CollectionInit {
    total_items: Int as uint32;
    wallet_address: Address;
    price: Int;

    today: Int;
    today_mints: Int;
    number_of_per_day_mints: Int;

    white_mint_start_time: Int;

    public_mint_start_time: Int;

    promotion_start_time: Int;
    promotion_end_time: Int;
    promotion_price: Int;
}



struct CollectionData {
    next_item_index: Int;
    collection_content: Cell;
    owner_address: Address;
}

struct RoyaltyParams {
    numerator: Int;
    denominator: Int;
    destination: Address;
}

message(0x11) Lock {
    query_id: Int as uint64;
}

message(0x12) AdminLock {
    query_id: Int as uint64;
    tokenId: Int as uint64;
}

message(0x13) Unlock {
    query_id: Int as uint64;
}

message(0x14) AdminUnlock {
    query_id: Int as uint64;
    tokenId: Int as uint64;
}

message(0x20) Burn {
    query_id: Int as uint64;
}

message(0x30) WithdrawJettons{
    query_id: Int as uint64;
    amount: Int as coins;
    wallet_address: Address;
    response_destination: Address;
    destination: Address;
}

message(0x7362d09c) TokenNotification {
    query_id: Int as uint64;
    amount: Int as coins;
    from: Address;
    forward_payload: Slice as remaining; // Comment Text message when Transfer the jetton 
}

message(0xf8a7ea5) TokenTransfer {
    query_id: Int as uint64;
    amount: Int as coins;
    destination: Address;
    response_destination: Address;
    custom_payload: Cell?;
    forward_ton_amount: Int as coins;
    forward_payload: Slice as remaining;
}

message(0x5fcc3d14) Transfer {
    query_id: Int as uint64;
    new_owner: Address;
    response_destination: Address;
    custom_payload: Cell?;
    forward_amount: Int as coins;
    forward_payload: Slice as remaining;
}
message(0x05138d91) OwnershipAssigned {
    query_id: Int as uint64;
    prev_owner: Address;
    forward_payload: Slice as remaining;
}
message(0xd53276db) Excesses {
    query_id: Int as uint64;
}
message(0x2fcb26a2) GetStaticData {
    query_id: Int as uint64;
}
message(0x8b771735) ReportStaticData {
    query_id: Int as uint64;
    index_id: Int;
    collection: Address;
}
struct GetNftData {
    is_initialized: Bool;
    index: Int;
    collection_address: Address;
    owner_address: Address;
    individual_content: Cell;
}