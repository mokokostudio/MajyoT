/** This is an automatically generated class by FairyGUI. Please do not modify it. **/

export default class UI_PVPUI extends fgui.GComponent {

	public m_HistoryBtn:fgui.GButton;
	public m_nameTxt:fgui.GTextField;
	public m_nameTxt_2:fgui.GTextField;
	public m_nameTxt_3:fgui.GTextField;
	public m_nameTxt_4:fgui.GTextField;
	public m_nameTxt_5:fgui.GTextField;
	public m_nameTxt_6:fgui.GTextField;
	public m_nameTxt_7:fgui.GTextField;
	public m_nameTxt_8:fgui.GTextField;
	public m_nameTxt_9:fgui.GTextField;
	public m_nameTxt_10:fgui.GTextField;
	public m_nameTxt_11:fgui.GTextField;
	public static URL:string = "ui://n6cyu7a6ja100";

	public static createInstance():UI_PVPUI {
		return <UI_PVPUI>(fgui.UIPackage.createObject("PVPUI", "PVPUI"));
	}

	protected onConstruct():void {
		this.m_HistoryBtn = <fgui.GButton>(this.getChildAt(2));
		this.m_nameTxt = <fgui.GTextField>(this.getChildAt(3));
		this.m_nameTxt_2 = <fgui.GTextField>(this.getChildAt(8));
		this.m_nameTxt_3 = <fgui.GTextField>(this.getChildAt(9));
		this.m_nameTxt_4 = <fgui.GTextField>(this.getChildAt(11));
		this.m_nameTxt_5 = <fgui.GTextField>(this.getChildAt(12));
		this.m_nameTxt_6 = <fgui.GTextField>(this.getChildAt(14));
		this.m_nameTxt_7 = <fgui.GTextField>(this.getChildAt(15));
		this.m_nameTxt_8 = <fgui.GTextField>(this.getChildAt(16));
		this.m_nameTxt_9 = <fgui.GTextField>(this.getChildAt(17));
		this.m_nameTxt_10 = <fgui.GTextField>(this.getChildAt(18));
		this.m_nameTxt_11 = <fgui.GTextField>(this.getChildAt(19));
	}
}