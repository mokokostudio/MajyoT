(function () {
    'use strict';

    class GameConfig {
        constructor() {
        }
        static init() {
            var reg = Laya.ClassUtils.regClass;
        }
    }
    GameConfig.width = 560;
    GameConfig.height = 1120;
    GameConfig.scaleMode = "fixedwidth";
    GameConfig.screenMode = "none";
    GameConfig.alignV = "top";
    GameConfig.alignH = "left";
    GameConfig.startScene = "";
    GameConfig.sceneRoot = "";
    GameConfig.debug = false;
    GameConfig.stat = false;
    GameConfig.physicsDebug = false;
    GameConfig.exportSceneToJson = true;
    GameConfig.init();

    class Global {
    }
    Global.barHeight = 0;

    class GameLayer {
        constructor() {
            this.layers = {};
            this.initLayer(8000);
            this.initLayer(8100);
            this.initLayer(8200);
            this.initLayer(8700);
            this.initLayer(9000);
            this.initLayer(9100);
        }
        initLayer(id) {
            let sprite = new Laya.Sprite();
            sprite.y = Global.barHeight;
            this.layers[id] = sprite;
            Laya.stage.addChild(sprite);
        }
        getLayer(id) {
            return this.layers[id];
        }
        static get instance() {
            return GameLayer._instance || (GameLayer._instance = new GameLayer());
        }
        clearStage() {
            Laya.stage.removeChildren(0, Laya.stage.numChildren - 1);
        }
        resetStage() {
            this.AddLayer(8000);
            this.AddLayer(8100);
            this.AddLayer(8200);
            this.AddLayer(8700);
            this.AddLayer(9000);
            this.AddLayer(9100);
        }
        AddLayer(id) {
            Laya.stage.addChild(this.layers[id]);
        }
    }

    class Facade {
        static toggle(clazz, toggleState = 0) {
            let mediator;
            if (clazz["name"] in Facade.mediators) {
                mediator = Facade.mediators[clazz["name"]];
                switch (toggleState) {
                    case 1:
                        if (mediator.toggleState == -1) {
                            mediator.show();
                        }
                        break;
                    case -1:
                        if (mediator.toggleState == 1) {
                            mediator.hide();
                        }
                        break;
                    case 0:
                        if (mediator.toggleState == -1) {
                            mediator.show();
                        }
                        else {
                            mediator.hide();
                        }
                        break;
                }
            }
            else {
                if (toggleState != -1) {
                    mediator = new clazz();
                    Facade.mediators[clazz["name"]] = mediator;
                }
            }
            return mediator;
        }
        static mediatorExec(clazz, handlerName, ...args) {
            if (clazz["name"] in Facade.mediators) {
                let mediator = Facade.mediators[clazz["name"]];
                if (mediator.toggleState == 1 && mediator["skin"]) {
                    mediator[handlerName].apply(mediator, args);
                }
                else {
                    mediator.handlerName = handlerName;
                    mediator.functionArgs = args;
                }
            }
        }
        static getMediatorActivityState(clazz) {
            if (clazz["name"] in Facade.mediators) {
                let mediator = Facade.mediators[clazz["name"]];
                return mediator.toggleState == 1;
            }
            return false;
        }
        static ChangePanel(clazz) {
            if (!Facade.curPanelMediator || Facade.curPanelMediator != Facade.mediators[clazz.name] || Facade.curPanelMediator.toggleState == -1) {
                this.CloseCurPanel();
                Facade.curPanelMediator = Facade.toggle(clazz, 1);
            }
        }
        static CloseCurPanel() {
            if (Facade.curPanelMediator && (Facade.curPanelMediator.toggleState == 1 || !Facade.curPanelMediator.isInit)) {
                Facade.curPanelMediator.hide();
            }
        }
        static closeAllPanel() {
            for (let className in Facade.mediators) {
                let mediator = Facade.mediators[className];
                if (mediator.uiPkgName != "alertpanel" && mediator.toggleState == 1) {
                    mediator.hide();
                }
            }
        }
    }
    Facade.mediators = {};

    class UI_TopUI extends fgui.GComponent {
        static createInstance() {
            return (fgui.UIPackage.createObject("MainUI", "TopUI"));
        }
        onConstruct() {
            this.m_nameTxt = (this.getChildAt(6));
        }
    }
    UI_TopUI.URL = "ui://89m4e0wjddsn2g";

    class UI_BottomUI extends fgui.GComponent {
        static createInstance() {
            return (fgui.UIPackage.createObject("MainUI", "BottomUI"));
        }
        onConstruct() {
            this.m_c1 = this.getControllerAt(0);
            this.m_HomeBtn = (this.getChildAt(1));
            this.m_NftBtn = (this.getChildAt(2));
            this.m_MarketBtn = (this.getChildAt(3));
            this.m_PvpBtn = (this.getChildAt(4));
        }
    }
    UI_BottomUI.URL = "ui://89m4e0wjddsn2h";

    class UI_MainUI extends fgui.GComponent {
        static createInstance() {
            return (fgui.UIPackage.createObject("MainUI", "MainUI"));
        }
        onConstruct() {
            this.m_nameTxt = (this.getChildAt(56));
            this.m_nameTxt_2 = (this.getChildAt(57));
            this.m_nameTxt_3 = (this.getChildAt(59));
            this.m_nameTxt_4 = (this.getChildAt(60));
            this.m_nameTxt_5 = (this.getChildAt(62));
            this.m_nameTxt_6 = (this.getChildAt(63));
            this.m_nameTxt_7 = (this.getChildAt(64));
            this.m_nameTxt_8 = (this.getChildAt(65));
            this.m_PropInfo = (this.getChildAt(66));
            this.m_nameTxt_9 = (this.getChildAt(69));
            this.m_nameTxt_10 = (this.getChildAt(72));
            this.m_BossTip = (this.getChildAt(73));
            this.m_ItemTip = (this.getChildAt(87));
        }
    }
    UI_MainUI.URL = "ui://89m4e0wjun3u0";

    class MainUIBinder {
        static bindAll() {
            fgui.UIObjectFactory.setExtension(UI_TopUI.URL, UI_TopUI);
            fgui.UIObjectFactory.setExtension(UI_BottomUI.URL, UI_BottomUI);
            fgui.UIObjectFactory.setExtension(UI_MainUI.URL, UI_MainUI);
        }
    }

    class Mediator {
        constructor(resUrl, layerID = 8700) {
            this.layerID = layerID;
            this.resUrl = resUrl;
            this.loadUI();
            Laya.stage.on(Laya.Event.RESIZE, this, this.autoScreen);
        }
        loadUI() {
            fgui.UIPackage.loadPackage(this.resUrl, Laya.Handler.create(this, this.onLoaded));
        }
        onLoaded() {
            if (!this.isInit) {
                this.isInit = true;
                this.init();
            }
            if (this.toggleState == -1) {
                return;
            }
            GameLayer.instance.getLayer(this.layerID).addChild(fairygui.GRoot.inst.displayObject);
            this.show();
            if (this.handlerName) {
                this[this.handlerName].apply(this, this.functionArgs);
                this.handlerName = this.functionArgs = undefined;
            }
        }
        init() {
        }
        afterAllReady() {
        }
        awake() {
        }
        sleep() {
        }
        autoScreen() {
            if (this["skin"]) {
                this["skin"].setSize(Laya.stage.width, Laya.stage.height - Global.barHeight);
            }
        }
        show() {
            this.toggleState = 1;
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
        hide() {
            this.toggleState = -1;
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
        isShow() {
            if (!this["skin"]) {
                return false;
            }
            let displayObject = this["skin"].displayObject;
            if (displayObject.parent) {
                return true;
            }
            else {
                return false;
            }
        }
    }

    class MainUiPanelMediator extends Mediator {
        constructor() {
            super("res/UI/MainUI", 8100);
        }
        init() {
            this.uiPkgName = "MainUI";
            fairygui.UIPackage.addPackage("res/UI/MainUI");
            MainUIBinder.bindAll();
            this.skin = UI_MainUI.createInstance();
            this.skin.setSize(Laya.stage.width, Laya.stage.height);
            this.afterAllReady();
        }
        afterAllReady() {
            this.skin.m_ItemTip.visible = false;
            this.skin.m_BossTip.visible = false;
            this.skin.m_PropInfo.visible = false;
        }
    }

    class TopUIPanelMediator extends Mediator {
        constructor() {
            super("res/UI/MainUI", 8200);
        }
        init() {
            this.uiPkgName = "MainUI";
            fairygui.UIPackage.addPackage("res/UI/MainUI");
            MainUIBinder.bindAll();
            this.skin = UI_TopUI.createInstance();
            this.skin.setSize(Laya.stage.width, Laya.stage.height);
            this.afterAllReady();
        }
        afterAllReady() {
        }
    }

    class UI_MarketUI extends fgui.GComponent {
        static createInstance() {
            return (fgui.UIPackage.createObject("MarketUI", "MarketUI"));
        }
        onConstruct() {
            this.m_List = (this.getChildAt(2));
            this.m_Grounding = (this.getChildAt(19));
        }
    }
    UI_MarketUI.URL = "ui://n9tq8la7ploc0";

    class UI_Market_Info extends fgui.GComponent {
        static createInstance() {
            return (fgui.UIPackage.createObject("MarketUI", "Market Info"));
        }
        onConstruct() {
            this.m_nameTxt = (this.getChildAt(2));
            this.m_nameTxt_2 = (this.getChildAt(3));
            this.m_nameTxt_3 = (this.getChildAt(4));
            this.m_nameTxt_4 = (this.getChildAt(6));
            this.m_nameTxt_5 = (this.getChildAt(8));
        }
    }
    UI_Market_Info.URL = "ui://n9tq8la7q9prt";

    class MarketUIBinder {
        static bindAll() {
            fgui.UIObjectFactory.setExtension(UI_MarketUI.URL, UI_MarketUI);
            fgui.UIObjectFactory.setExtension(UI_Market_Info.URL, UI_Market_Info);
        }
    }

    class MarketUIPanelMediator extends Mediator {
        constructor() {
            super("res/UI/MarketUI", 8100);
        }
        init() {
            this.uiPkgName = "MarketUI";
            fairygui.UIPackage.addPackage("res/UI/MarketUI");
            MarketUIBinder.bindAll();
            this.skin = UI_MarketUI.createInstance();
            this.skin.setSize(Laya.stage.width, Laya.stage.height);
            this.afterAllReady();
        }
        afterAllReady() {
            this.skin.m_Grounding.visible = false;
        }
    }

    class UI_NFTUI extends fgui.GComponent {
        static createInstance() {
            return (fgui.UIPackage.createObject("NFTUI", "NFTUI"));
        }
        onConstruct() {
            this.m_c1 = this.getControllerAt(0);
            this.m_MintBtn = (this.getChildAt(4));
        }
    }
    UI_NFTUI.URL = "ui://o2qwdak911fs21";

    class NFTUIBinder {
        static bindAll() {
            fgui.UIObjectFactory.setExtension(UI_NFTUI.URL, UI_NFTUI);
        }
    }

    class NftUIPanelMediator extends Mediator {
        constructor() {
            super("res/UI/NFTUI", 8100);
        }
        init() {
            this.uiPkgName = "NFTUI";
            fairygui.UIPackage.addPackage("res/UI/NFTUI");
            NFTUIBinder.bindAll();
            this.skin = UI_NFTUI.createInstance();
            this.skin.setSize(Laya.stage.width, Laya.stage.height);
            this.afterAllReady();
        }
        afterAllReady() {
        }
    }

    class UI_PVPUI extends fgui.GComponent {
        static createInstance() {
            return (fgui.UIPackage.createObject("PVPUI", "PVPUI"));
        }
        onConstruct() {
            this.m_HistoryBtn = (this.getChildAt(2));
            this.m_nameTxt = (this.getChildAt(3));
            this.m_nameTxt_2 = (this.getChildAt(8));
            this.m_nameTxt_3 = (this.getChildAt(9));
            this.m_nameTxt_4 = (this.getChildAt(11));
            this.m_nameTxt_5 = (this.getChildAt(12));
            this.m_nameTxt_6 = (this.getChildAt(14));
            this.m_nameTxt_7 = (this.getChildAt(15));
            this.m_nameTxt_8 = (this.getChildAt(16));
            this.m_nameTxt_9 = (this.getChildAt(17));
            this.m_nameTxt_10 = (this.getChildAt(18));
            this.m_nameTxt_11 = (this.getChildAt(19));
        }
    }
    UI_PVPUI.URL = "ui://n6cyu7a6ja100";

    class UI_PVP_Info extends fgui.GComponent {
        static createInstance() {
            return (fgui.UIPackage.createObject("PVPUI", "PVP Info"));
        }
        onConstruct() {
            this.m_nameTxt = (this.getChildAt(1));
            this.m_nameTxt_2 = (this.getChildAt(2));
            this.m_nameTxt_3 = (this.getChildAt(3));
        }
    }
    UI_PVP_Info.URL = "ui://n6cyu7a6ja10n";

    class PVPUIBinder {
        static bindAll() {
            fgui.UIObjectFactory.setExtension(UI_PVPUI.URL, UI_PVPUI);
            fgui.UIObjectFactory.setExtension(UI_PVP_Info.URL, UI_PVP_Info);
        }
    }

    class PvpUIPanelMediator extends Mediator {
        constructor() {
            super("res/UI/PVPUI", 8100);
        }
        init() {
            this.uiPkgName = "PVPUI";
            fairygui.UIPackage.addPackage("res/UI/PVPUI");
            PVPUIBinder.bindAll();
            this.skin = UI_PVPUI.createInstance();
            this.skin.setSize(Laya.stage.width, Laya.stage.height);
            this.afterAllReady();
        }
        afterAllReady() {
        }
    }

    class BottomUiPanelMediator extends Mediator {
        constructor() {
            super("res/UI/MainUI", 8200);
        }
        init() {
            this.uiPkgName = "BottomUi";
            fairygui.UIPackage.addPackage("res/UI/MainUI");
            MainUIBinder.bindAll();
            this.skin = UI_BottomUI.createInstance();
            this.skin.setSize(Laya.stage.width, Laya.stage.height);
            this.afterAllReady();
        }
        afterAllReady() {
            this.skin.m_HomeBtn.onClick(this, this.homeBtnClick);
            this.skin.m_NftBtn.onClick(this, this.NftBtnClick);
            this.skin.m_MarketBtn.onClick(this, this.marketBtnClick);
            this.skin.m_PvpBtn.onClick(this, this.pvpBtnClick);
        }
        homeBtnClick() {
            Facade.toggle(MainUiPanelMediator, 1);
            Facade.toggle(NftUIPanelMediator, -1);
            Facade.toggle(MarketUIPanelMediator, -1);
            Facade.toggle(PvpUIPanelMediator, -1);
        }
        NftBtnClick() {
            Facade.toggle(MainUiPanelMediator, -1);
            Facade.toggle(NftUIPanelMediator, 1);
            Facade.toggle(MarketUIPanelMediator, -1);
            Facade.toggle(PvpUIPanelMediator, -1);
        }
        marketBtnClick() {
            Facade.toggle(MainUiPanelMediator, -1);
            Facade.toggle(NftUIPanelMediator, -1);
            Facade.toggle(MarketUIPanelMediator, 1);
            Facade.toggle(PvpUIPanelMediator, -1);
        }
        pvpBtnClick() {
            Facade.toggle(MainUiPanelMediator, -1);
            Facade.toggle(NftUIPanelMediator, -1);
            Facade.toggle(MarketUIPanelMediator, -1);
            Facade.toggle(PvpUIPanelMediator, 1);
        }
    }

    var Stage = Laya.Stage;
    class Main {
        constructor() {
            if (window["Laya3D"])
                Laya3D.init(GameConfig.width, GameConfig.height);
            else
                Laya.init(GameConfig.width, GameConfig.height, Laya["WebGL"]);
            Laya["Physics"] && Laya["Physics"].enable();
            Laya["DebugPanel"] && Laya["DebugPanel"].enable();
            if (Laya.Browser.onPC) {
                Laya.stage.scaleMode = Stage.SCALE_SHOWALL;
                Laya.stage.alignH = Stage.ALIGN_CENTER;
                Laya.stage.alignV = Stage.ALIGN_MIDDLE;
                Laya.stage.screenMode = Stage.SCREEN_NONE;
            }
            else {
                Laya.stage.scaleMode = Stage.SCALE_FIXED_WIDTH;
                Laya.stage.alignH = Stage.ALIGN_LEFT;
                Laya.stage.alignV = Stage.ALIGN_TOP;
                Laya.stage.screenMode = Stage.SCREEN_VERTICAL;
            }
            Laya.URL.exportSceneToJson = GameConfig.exportSceneToJson;
            if (GameConfig.debug || Laya.Utils.getQueryString("debug") == "true")
                Laya.enableDebugPanel();
            if (GameConfig.physicsDebug && Laya["PhysicsDebugDraw"])
                Laya["PhysicsDebugDraw"].enable();
            if (GameConfig.stat)
                Laya.Stat.show();
            Laya.alertGlobalError(true);
            Laya.ResourceVersion.enable("version.json", Laya.Handler.create(this, this.onVersionLoaded), Laya.ResourceVersion.FILENAME_VERSION);
        }
        onVersionLoaded() {
            Laya.AtlasInfoManager.enable("fileconfig.json", Laya.Handler.create(this, this.onConfigLoaded));
        }
        onConfigLoaded() {
            GameConfig.startScene && Laya.Scene.open(GameConfig.startScene);
            GameLayer.instance;
            Laya.stage.addChild(fgui.GRoot.inst.displayObject);
            Facade.toggle(TopUIPanelMediator, 1);
            Facade.toggle(BottomUiPanelMediator, 1);
            Facade.toggle(MainUiPanelMediator, 1);
        }
    }
    new Main();

}());
//# sourceMappingURL=bundle.js.map
