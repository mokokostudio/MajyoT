import Mediator from "./Mediator";
import { ToggleState } from "./ToggleState";

export default class Facade{
    
    public static mediators:{[index: string]: Mediator} = {};
    private static curPanelMediator:Mediator;

    public static toggle(clazz: { new (): Mediator; _instance?: Mediator }, toggleState: ToggleState = ToggleState.AUTO):Mediator
    {
        let mediator:Mediator;
        if(clazz["name"] in Facade.mediators){
            mediator = Facade.mediators[clazz["name"]];
            switch(toggleState){
                case ToggleState.SHOW:
                    if(mediator.toggleState == ToggleState.HIDE){
                        mediator.show();
                    }
                    break;
                case ToggleState.HIDE:
                    if(mediator.toggleState == ToggleState.SHOW){
                        mediator.hide();
                    }
                    break;   
                case ToggleState.AUTO:
                    if(mediator.toggleState == ToggleState.HIDE){
                        mediator.show();
                    }else{
                        mediator.hide();
                    }
                    break;     
            }
        }else{
            if(toggleState != ToggleState.HIDE){
                mediator = new clazz();
                Facade.mediators[clazz["name"]] = mediator;
            }
        }
        return mediator;
    }

    public static mediatorExec(clazz: { new (): Mediator; _instance?: Mediator }, handlerName: string, ...args)
    {
        if(clazz["name"] in Facade.mediators){
            let mediator:Mediator = Facade.mediators[clazz["name"]];
            if(mediator.toggleState == ToggleState.SHOW && mediator["skin"]){
                mediator[handlerName].apply(mediator, args);
            }else{
                mediator.handlerName = handlerName;
                mediator.functionArgs = args;
            }
        }
    }

    public static getMediatorActivityState(clazz: { new (): Mediator; _instance?: Mediator })
    {
        if(clazz["name"] in Facade.mediators){
            let mediator:Mediator = Facade.mediators[clazz["name"]];
            return mediator.toggleState == ToggleState.SHOW;
        }
        return false;
    }

    public static ChangePanel(clazz: { new (): Mediator; _instance?: Mediator })
    {
        if(!Facade.curPanelMediator || Facade.curPanelMediator != Facade.mediators[clazz.name] || Facade.curPanelMediator.toggleState == ToggleState.HIDE){
            this.CloseCurPanel();
            Facade.curPanelMediator = Facade.toggle(clazz , ToggleState.SHOW);
        }
    }

    public static CloseCurPanel()
    {
        if(Facade.curPanelMediator && (Facade.curPanelMediator.toggleState == ToggleState.SHOW || !Facade.curPanelMediator.isInit)){
            Facade.curPanelMediator.hide();
        }
    }

    public static closeAllPanel()
    {
        for(let className in Facade.mediators){
            let mediator = Facade.mediators[className];
            if(mediator.uiPkgName != "alertpanel" && mediator.toggleState == ToggleState.SHOW){
                mediator.hide();
            }
        }
    }

}