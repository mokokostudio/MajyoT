import Global from "../Global";

export default class GameLayer{
    private layers:{[id:number]:Laya.Sprite} = {};

    constructor(){
        this.initLayer(GameLayerID.Bottom);
        this.initLayer(GameLayerID.BaseUi);
        this.initLayer(GameLayerID.Navigation); 
        this.initLayer(GameLayerID.PopUi);
        this.initLayer(GameLayerID.Top);
        this.initLayer(GameLayerID.Tips);
    }

    private initLayer(id:GameLayerID)
    {
        let sprite = new Laya.Sprite();
        sprite.y = Global.barHeight;
        this.layers[id] = sprite;
        Laya.stage.addChild(sprite);
    }

    public getLayer(id:GameLayerID)
    {
        return this.layers[id];
    }
    
    static _instance: GameLayer;
    public static get instance():GameLayer{
        return GameLayer._instance || (GameLayer._instance = new GameLayer());
    }

    public clearStage()
    {
        Laya.stage.removeChildren(0,Laya.stage.numChildren-1);
    }

    public resetStage()
    {
        this.AddLayer(GameLayerID.Bottom);
        this.AddLayer(GameLayerID.BaseUi);
        this.AddLayer(GameLayerID.Navigation);
        this.AddLayer(GameLayerID.PopUi);
        this.AddLayer(GameLayerID.Top);
        this.AddLayer(GameLayerID.Tips);
    }

    private AddLayer(id:GameLayerID)
    {
        Laya.stage.addChild(this.layers[id]);
    }
}

export const enum GameLayerID{
    Bottom = 8000,
    BaseUi = 8100,
    Navigation = 8200,
    PopUi = 8700,
    Top = 9000,
    Tips = 9100,
}