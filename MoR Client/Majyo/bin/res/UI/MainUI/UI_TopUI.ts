/** This is an automatically generated class by FairyGUI. Please do not modify it. **/

export default class UI_TopUI extends fgui.GComponent {

	public m_nameTxt:fgui.GTextField;
	public static URL:string = "ui://89m4e0wjddsn2g";

	public static createInstance():UI_TopUI {
		return <UI_TopUI>(fgui.UIPackage.createObject("MainUI", "TopUI"));
	}

	protected onConstruct():void {
		this.m_nameTxt = <fgui.GTextField>(this.getChildAt(6));
	}
}