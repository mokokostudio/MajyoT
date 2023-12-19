/** This is an automatically generated class by FairyGUI. Please do not modify it. **/

export default class UI_Market_Info extends fgui.GComponent {

	public m_nameTxt:fgui.GTextField;
	public m_nameTxt_2:fgui.GTextField;
	public m_nameTxt_3:fgui.GTextField;
	public m_nameTxt_4:fgui.GTextField;
	public m_nameTxt_5:fgui.GTextField;
	public static URL:string = "ui://n9tq8la7q9prt";

	public static createInstance():UI_Market_Info {
		return <UI_Market_Info>(fgui.UIPackage.createObject("MarketUI", "Market Info"));
	}

	protected onConstruct():void {
		this.m_nameTxt = <fgui.GTextField>(this.getChildAt(2));
		this.m_nameTxt_2 = <fgui.GTextField>(this.getChildAt(3));
		this.m_nameTxt_3 = <fgui.GTextField>(this.getChildAt(4));
		this.m_nameTxt_4 = <fgui.GTextField>(this.getChildAt(6));
		this.m_nameTxt_5 = <fgui.GTextField>(this.getChildAt(8));
	}
}