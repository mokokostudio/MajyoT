import MainUIBinder from "../../bin/res/UI/MainUI/MainUIBinder";
import UI_MainUI from "../../bin/res/UI/MainUI/UI_MainUI";
import { GameLayerID } from "../game/GameLayer";
import Mediator from "../module/Mediator";
export default class MainUiPanelMediator extends Mediator{
    protected skin:UI_MainUI;

	constructor() {
        super("res/UI/MainUI",GameLayerID.BaseUi);
    }

	protected init()
	{
		this.uiPkgName = "MainUI";
		fairygui.UIPackage.addPackage("res/UI/MainUI");
		MainUIBinder.bindAll();
		this.skin = UI_MainUI.createInstance();
		this.skin.setSize(Laya.stage.width, Laya.stage.height);
		this.afterAllReady();
	}

    protected afterAllReady()
    {
        this.skin.m_ItemTip.visible = false;
		this.skin.m_BossTip.visible = false;
		this.skin.m_PropInfo.visible = false;
    }
}