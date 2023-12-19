import NFTUIBinder from "../../bin/res/UI/NFTUI/NFTUIBinder";
import UI_NFTUI from "../../bin/res/UI/NFTUI/UI_NFTUI";
import { GameLayerID } from "../game/GameLayer";
import Mediator from "../module/Mediator";
export default class NftUIPanelMediator extends Mediator{
    protected skin:UI_NFTUI;

	constructor() {
        super("res/UI/NFTUI",GameLayerID.BaseUi);
    }

	protected init()
	{
		this.uiPkgName = "NFTUI";
		fairygui.UIPackage.addPackage("res/UI/NFTUI");
		NFTUIBinder.bindAll();
		this.skin = UI_NFTUI.createInstance();
		this.skin.setSize(Laya.stage.width, Laya.stage.height);
		this.afterAllReady();
	}

    protected afterAllReady()
    {
        
    }
}