import MainUIBinder from "../../bin/res/UI/MainUI/MainUIBinder";
import UI_TopUI from "../../bin/res/UI/MainUI/UI_TopUI";
import { GameLayerID } from "../game/GameLayer";
import Mediator from "./Mediator";
export default class TopUIPanelMediator extends Mediator{
    protected skin:UI_TopUI;

	constructor() {
        super("res/UI/MainUI",GameLayerID.Navigation);
    }

	protected init()
	{
		this.uiPkgName = "MainUI";
		fairygui.UIPackage.addPackage("res/UI/MainUI");
		MainUIBinder.bindAll();
		this.skin = UI_TopUI.createInstance();
		this.skin.setSize(Laya.stage.width, Laya.stage.height);
		this.afterAllReady();
	}

    protected afterAllReady()
    {
        
    }
}