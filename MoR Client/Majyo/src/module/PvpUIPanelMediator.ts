import PVPUIBinder from "../../bin/res/UI/PVPUI/PVPUIBinder";
import UI_PVPUI from "../../bin/res/UI/PVPUI/UI_PVPUI";
import { GameLayerID } from "../game/GameLayer";
import Mediator from "../module/Mediator";
export default class PvpUIPanelMediator extends Mediator{
    protected skin:UI_PVPUI;

	constructor() {
        super("res/UI/PVPUI",GameLayerID.BaseUi);
    }

	protected init()
	{
		this.uiPkgName = "PVPUI";
		fairygui.UIPackage.addPackage("res/UI/PVPUI");
		PVPUIBinder.bindAll();
		this.skin = UI_PVPUI.createInstance();
		this.skin.setSize(Laya.stage.width, Laya.stage.height);
		this.afterAllReady();
	}

    protected afterAllReady()
    {
        
    }
}