import MainUIBinder from "../../bin/res/UI/MainUI/MainUIBinder";
import UI_BottomUI from "../../bin/res/UI/MainUI/UI_BottomUI";
import { GameLayerID } from "../game/GameLayer";
import Mediator from "../module/Mediator";
import Facade from "./Facade";
import MainUiPanelMediator from "./MainUiPanelMediator";
import MarketUIPanelMediator from "./MarketUIPanelMediator";
import NftUIPanelMediator from "./NftUIPanelMediator";
import PvpUIPanelMediator from "./PvpUIPanelMediator";
import { ToggleState } from "./ToggleState";
export default class BottomUiPanelMediator extends Mediator{
    protected skin:UI_BottomUI;

	constructor() {
        super("res/UI/MainUI",GameLayerID.Navigation);
    }

	protected init()
	{
		this.uiPkgName = "BottomUi";
		fairygui.UIPackage.addPackage("res/UI/MainUI");
		MainUIBinder.bindAll();
		this.skin = UI_BottomUI.createInstance();
		this.skin.setSize(Laya.stage.width, Laya.stage.height);
		this.afterAllReady();
	}

    protected afterAllReady()
    {
		this.skin.m_HomeBtn.onClick(this,this.homeBtnClick);
		this.skin.m_NftBtn.onClick(this,this.NftBtnClick);
		this.skin.m_MarketBtn.onClick(this,this.marketBtnClick);
		this.skin.m_PvpBtn.onClick(this,this.pvpBtnClick);
    }

	private homeBtnClick()
	{
		Facade.toggle(MainUiPanelMediator,ToggleState.SHOW);
		Facade.toggle(NftUIPanelMediator,ToggleState.HIDE);
		Facade.toggle(MarketUIPanelMediator,ToggleState.HIDE);
		Facade.toggle(PvpUIPanelMediator,ToggleState.HIDE); 
	}

	private NftBtnClick()
	{
		Facade.toggle(MainUiPanelMediator,ToggleState.HIDE);
		Facade.toggle(NftUIPanelMediator,ToggleState.SHOW);
		Facade.toggle(MarketUIPanelMediator,ToggleState.HIDE);
		Facade.toggle(PvpUIPanelMediator,ToggleState.HIDE); 
	}

	private marketBtnClick()
	{
		Facade.toggle(MainUiPanelMediator,ToggleState.HIDE);
		Facade.toggle(NftUIPanelMediator,ToggleState.HIDE);
		Facade.toggle(MarketUIPanelMediator,ToggleState.SHOW);
		Facade.toggle(PvpUIPanelMediator,ToggleState.HIDE); 
	}

	private pvpBtnClick()
	{
		Facade.toggle(MainUiPanelMediator,ToggleState.HIDE);
		Facade.toggle(NftUIPanelMediator,ToggleState.HIDE);
		Facade.toggle(MarketUIPanelMediator,ToggleState.HIDE);
		Facade.toggle(PvpUIPanelMediator,ToggleState.SHOW); 
	}
}