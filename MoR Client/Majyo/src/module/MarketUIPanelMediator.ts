import MarketUIBinder from "../../bin/res/UI/MarketUI/MarketUIBinder";
import UI_MarketUI from "../../bin/res/UI/MarketUI/UI_MarketUI";
import { GameLayerID } from "../game/GameLayer";
import Mediator from "../module/Mediator";
export default class MarketUIPanelMediator extends Mediator{
    protected skin:UI_MarketUI;

	constructor() {
        super("res/UI/MarketUI",GameLayerID.BaseUi);
    }

	protected init()
	{
		this.uiPkgName = "MarketUI";
		fairygui.UIPackage.addPackage("res/UI/MarketUI");
		MarketUIBinder.bindAll();
		this.skin = UI_MarketUI.createInstance();
		this.skin.setSize(Laya.stage.width, Laya.stage.height);
		this.afterAllReady();
	}

    protected afterAllReady()
    {
        this.skin.m_Grounding.visible = false;
    }
}