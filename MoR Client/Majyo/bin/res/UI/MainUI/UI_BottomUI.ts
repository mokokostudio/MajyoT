/** This is an automatically generated class by FairyGUI. Please do not modify it. **/

export default class UI_BottomUI extends fgui.GComponent {

	public m_c1:fgui.Controller;
	public m_HomeBtn:fgui.GButton;
	public m_NftBtn:fgui.GButton;
	public m_MarketBtn:fgui.GButton;
	public m_PvpBtn:fgui.GButton;
	public static URL:string = "ui://89m4e0wjddsn2h";

	public static createInstance():UI_BottomUI {
		return <UI_BottomUI>(fgui.UIPackage.createObject("MainUI", "BottomUI"));
	}

	protected onConstruct():void {
		this.m_c1 = this.getControllerAt(0);
		this.m_HomeBtn = <fgui.GButton>(this.getChildAt(1));
		this.m_NftBtn = <fgui.GButton>(this.getChildAt(2));
		this.m_MarketBtn = <fgui.GButton>(this.getChildAt(3));
		this.m_PvpBtn = <fgui.GButton>(this.getChildAt(4));
	}
}