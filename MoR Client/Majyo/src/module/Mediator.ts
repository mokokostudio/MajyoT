import { ToggleState } from "./ToggleState";
import Global from "../Global";
import GameLayer, { GameLayerID } from "../game/GameLayer";

export default class Mediator {
    public uiPkgName:string;
    protected resUrl: string;
    public toggleState: ToggleState;
    public isInit: boolean;
    public handlerName: string;
    public functionArgs: any[];
    private layerID: GameLayerID;

    constructor(resUrl: string, layerID: GameLayerID = GameLayerID.PopUi) {
        this.layerID = layerID;
        this.resUrl = resUrl;
        this.loadUI();
        Laya.stage.on(Laya.Event.RESIZE, this, this.autoScreen)
    }

    private loadUI() {
        fgui.UIPackage.loadPackage(this.resUrl, Laya.Handler.create(this, this.onLoaded));
    }

    protected onLoaded(): void {
        if (!this.isInit) {
            this.isInit = true;
            this.init();
        }
        if(this.toggleState == ToggleState.HIDE){
            return;
        }
        GameLayer.instance.getLayer(this.layerID).addChild(fairygui.GRoot.inst.displayObject);
        this.show();
        if (this.handlerName) {
            this[this.handlerName].apply(this, this.functionArgs);
            this.handlerName = this.functionArgs = undefined;
        }
    }

    protected init() {

    }

    protected afterAllReady()
    {
        
    }

    protected awake() {

    }

    protected sleep() {

    }

    protected autoScreen() {
        if (this["skin"]) {
            this["skin"].setSize(Laya.stage.width, Laya.stage.height - Global.barHeight);
        }
    }

    public show() {
        this.toggleState = ToggleState.SHOW;
        if (!this["skin"]) {
            return;
        }
        let displayObject = this["skin"].displayObject;
        if (displayObject) {
            GameLayer.instance.getLayer(this.layerID).addChild(displayObject);
        }
        if (this["skin"] && this["skin"]["m_btnClose"]) {
            this["skin"]["m_btnClose"].onClick(this, this.hide);
        }
        this.autoScreen();
        this.awake();
    }

    public hide() {
        this.toggleState = ToggleState.HIDE;
        if (!this["skin"]) {
            return;
        }
        let displayObject = this["skin"].displayObject;
        if (displayObject.parent) {
            displayObject.parent.removeChild(displayObject);
        }
        if (this["skin"] && this["skin"]["m_btnClose"]) {
            this["skin"]["m_btnClose"].offClick(this, this.hide);
        }
        this.sleep();
    }

    public isShow(): boolean {
        if (!this["skin"]) {
            return false;
        }
        let displayObject = this["skin"].displayObject;
        if (displayObject.parent) {
            return true;
        } else {
            return false;
        }
    }
}