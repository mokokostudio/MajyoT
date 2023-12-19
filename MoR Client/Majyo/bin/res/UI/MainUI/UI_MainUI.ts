/** This is an automatically generated class by FairyGUI. Please do not modify it. **/

export default class UI_MainUI extends fgui.GComponent {

	public m_nameTxt:fgui.GTextField;
	public m_nameTxt_2:fgui.GTextField;
	public m_nameTxt_3:fgui.GTextField;
	public m_nameTxt_4:fgui.GTextField;
	public m_nameTxt_5:fgui.GTextField;
	public m_nameTxt_6:fgui.GTextField;
	public m_nameTxt_7:fgui.GTextField;
	public m_nameTxt_8:fgui.GTextField;
	public m_PropInfo:fgui.GGroup;
	public m_nameTxt_9:fgui.GTextField;
	public m_nameTxt_10:fgui.GTextField;
	public m_BossTip:fgui.GGroup;
	public m_ItemTip:fgui.GGroup;
	public static URL:string = "ui://89m4e0wjun3u0";

	public static createInstance():UI_MainUI {
		return <UI_MainUI>(fgui.UIPackage.createObject("MainUI", "MainUI"));
	}

	protected onConstruct():void {
		this.m_nameTxt = <fgui.GTextField>(this.getChildAt(56));
		this.m_nameTxt_2 = <fgui.GTextField>(this.getChildAt(57));
		this.m_nameTxt_3 = <fgui.GTextField>(this.getChildAt(59));
		this.m_nameTxt_4 = <fgui.GTextField>(this.getChildAt(60));
		this.m_nameTxt_5 = <fgui.GTextField>(this.getChildAt(62));
		this.m_nameTxt_6 = <fgui.GTextField>(this.getChildAt(63));
		this.m_nameTxt_7 = <fgui.GTextField>(this.getChildAt(64));
		this.m_nameTxt_8 = <fgui.GTextField>(this.getChildAt(65));
		this.m_PropInfo = <fgui.GGroup>(this.getChildAt(66));
		this.m_nameTxt_9 = <fgui.GTextField>(this.getChildAt(69));
		this.m_nameTxt_10 = <fgui.GTextField>(this.getChildAt(72));
		this.m_BossTip = <fgui.GGroup>(this.getChildAt(73));
		this.m_ItemTip = <fgui.GGroup>(this.getChildAt(87));
	}
}