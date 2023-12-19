/** This is an automatically generated class by FairyGUI. Please do not modify it. **/

export default class UI_PVP_Info extends fgui.GComponent {

	public m_nameTxt:fgui.GTextField;
	public m_nameTxt_2:fgui.GTextField;
	public m_nameTxt_3:fgui.GTextField;
	public static URL:string = "ui://n6cyu7a6ja10n";

	public static createInstance():UI_PVP_Info {
		return <UI_PVP_Info>(fgui.UIPackage.createObject("PVPUI", "PVP Info"));
	}

	protected onConstruct():void {
		this.m_nameTxt = <fgui.GTextField>(this.getChildAt(1));
		this.m_nameTxt_2 = <fgui.GTextField>(this.getChildAt(2));
		this.m_nameTxt_3 = <fgui.GTextField>(this.getChildAt(3));
	}
}