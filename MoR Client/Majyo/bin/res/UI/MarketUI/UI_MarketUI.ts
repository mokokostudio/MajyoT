/** This is an automatically generated class by FairyGUI. Please do not modify it. **/

export default class UI_MarketUI extends fgui.GComponent {

	public m_List:fgui.GList;
	public m_Grounding:fgui.GGroup;
	public static URL:string = "ui://n9tq8la7ploc0";

	public static createInstance():UI_MarketUI {
		return <UI_MarketUI>(fgui.UIPackage.createObject("MarketUI", "MarketUI"));
	}

	protected onConstruct():void {
		this.m_List = <fgui.GList>(this.getChildAt(2));
		this.m_Grounding = <fgui.GGroup>(this.getChildAt(19));
	}
}