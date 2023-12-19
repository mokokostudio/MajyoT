/** This is an automatically generated class by FairyGUI. Please do not modify it. **/

export default class UI_NFTUI extends fgui.GComponent {

	public m_c1:fgui.Controller;
	public m_MintBtn:fgui.GButton;
	public static URL:string = "ui://o2qwdak911fs21";

	public static createInstance():UI_NFTUI {
		return <UI_NFTUI>(fgui.UIPackage.createObject("NFTUI", "NFTUI"));
	}

	protected onConstruct():void {
		this.m_c1 = this.getControllerAt(0);
		this.m_MintBtn = <fgui.GButton>(this.getChildAt(4));
	}
}